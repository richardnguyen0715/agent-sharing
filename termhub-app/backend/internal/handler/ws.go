package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/termhub/backend/internal/domain"
	"github.com/termhub/backend/internal/policy"
	ptymanager "github.com/termhub/backend/internal/pty"
	"github.com/termhub/backend/internal/recording"
	"github.com/termhub/backend/internal/service"
	wshub "github.com/termhub/backend/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(_ *http.Request) bool {
		// Allow all origins for development.
		return true
	},
}

// WSHandler handles WebSocket connections for terminal I/O.
type WSHandler struct {
	hub             *wshub.Hub
	ptyManager      *ptymanager.Manager
	terminalService *service.TerminalService
	authService     *service.AuthService
	auditService    *service.AuditService
	policyEngine    *policy.Engine
	recordingMgr    *recording.Manager
}

// NewWSHandler creates a new WSHandler.
func NewWSHandler(
	hub *wshub.Hub,
	ptyManager *ptymanager.Manager,
	termSvc *service.TerminalService,
	authSvc *service.AuthService,
	auditSvc *service.AuditService,
	policyEngine *policy.Engine,
	recordingMgr *recording.Manager,
) *WSHandler {
	return &WSHandler{
		hub:             hub,
		ptyManager:      ptyManager,
		terminalService: termSvc,
		authService:     authSvc,
		auditService:    auditSvc,
		policyEngine:    policyEngine,
		recordingMgr:    recordingMgr,
	}
}

// Routes returns a chi.Router with WebSocket routes mounted.
// Note: NO auth middleware — authentication is done inside the handler via query param.
func (h *WSHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/terminal/{id}", h.HandleTerminal)
	return r
}

// HandleTerminal upgrades the HTTP connection to a WebSocket and bridges terminal I/O.
func (h *WSHandler) HandleTerminal(w http.ResponseWriter, r *http.Request) {
	// 1. Extract terminal ID from URL.
	terminalIDStr := chi.URLParam(r, "id")
	terminalID, err := uuid.Parse(terminalIDStr)
	if err != nil {
		http.Error(w, "invalid terminal ID", http.StatusBadRequest)
		return
	}

	// 2. Extract JWT token from query param.
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "missing token query parameter", http.StatusUnauthorized)
		return
	}

	// 3. Validate JWT token.
	claims, err := h.authService.ValidateToken(tokenStr)
	if err != nil {
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Reject refresh tokens.
	if claims.Type == "refresh" {
		http.Error(w, "refresh token cannot be used for WebSocket auth", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	// 4. Verify terminal ownership.
	terminal, err := h.terminalService.Get(r.Context(), terminalID, userID)
	if err != nil {
		log.Error().Err(err).
			Str("terminal_id", terminalIDStr).
			Str("user_id", claims.UserID).
			Msg("terminal access denied or not found")
		http.Error(w, "terminal not found or access denied", http.StatusForbidden)
		return
	}
	_ = terminal

	// 5. Get PTY session.
	ptySession, ok := h.ptyManager.Get(terminalID)
	if !ok {
		http.Error(w, "terminal PTY session not found", http.StatusNotFound)
		return
	}

	// 6. Upgrade to WebSocket.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to upgrade websocket connection")
		return
	}

	// 7. Create Client and register with Hub.
	client := wshub.NewClient(terminalID, claims.UserID, conn, h.hub)
	h.hub.Register(client)

	// Audit log the WebSocket connection.
	if auditErr := h.auditService.Log(context.Background(), &userID, &terminalID, domain.AuditActionTerminalConnect, map[string]interface{}{
		"client_id":  client.ID,
		"ip_address": r.RemoteAddr,
	}); auditErr != nil {
		log.Error().Err(auditErr).Str("terminal_id", terminalIDStr).Msg("failed to write connect audit log")
	}

	log.Info().
		Str("client_id", client.ID).
		Str("terminal_id", terminalIDStr).
		Str("user_id", claims.UserID).
		Msg("WebSocket connection established")

	// 8. Start goroutine to read PTY output and broadcast to all clients of this terminal.
	go func() {
		buf := make([]byte, 4096)
		for {
			n, readErr := ptySession.Reader().Read(buf)
			if readErr != nil {
				// PTY closed or process exited.
				errMsg := wshub.WSMessage{Type: "error", Data: "terminal process exited"}
				jsonData, _ := json.Marshal(errMsg)
				h.hub.SendToTerminal(terminalID, jsonData)
				log.Info().
					Str("terminal_id", terminalIDStr).
					Msg("PTY read loop ended")
				return
			}
			if n > 0 {
				data := buf[:n]
				msg := wshub.WSMessage{Type: "output", Data: string(data)}
				jsonData, marshalErr := json.Marshal(msg)
				if marshalErr != nil {
					continue
				}
				h.hub.SendToTerminal(terminalID, jsonData)

				// Record output (non-blocking, best-effort).
				h.recordingMgr.RecordOutput(terminalIDStr, data)
			}
		}
	}()

	// 9. Define PTY write function with policy enforcement and recording.
	ptyWriter := func(data []byte) error {
		input := string(data)

		// Evaluate input against the policy engine.
		result := h.policyEngine.Evaluate(input)

		switch result.Action {
		case policy.ActionBlock:
			// Send error back to client, don't write to PTY.
			errMsg := wshub.WSMessage{
				Type: "error",
				Data: fmt.Sprintf("Command blocked by policy: %s", result.Message),
			}
			jsonData, _ := json.Marshal(errMsg)
			select {
			case client.Send <- jsonData:
			default:
			}
			log.Warn().
				Str("client_id", client.ID).
				Str("terminal_id", terminalIDStr).
				Str("rule", result.Rule).
				Msg("command blocked by policy")
			return nil

		case policy.ActionWarn:
			// Write to PTY but log a warning audit event.
			if auditErr := h.auditService.Log(context.Background(), &userID, &terminalID, domain.AuditActionPolicyViolation, map[string]interface{}{
				"rule":    result.Rule,
				"message": result.Message,
				"input":   input,
			}); auditErr != nil {
				log.Error().Err(auditErr).Msg("failed to write policy warning audit log")
			}
		}

		// Record input (non-blocking, best-effort).
		h.recordingMgr.RecordInput(terminalIDStr, input)

		_, writeErr := ptySession.Writer().Write(data)
		if writeErr != nil {
			return fmt.Errorf("failed to write to PTY: %w", writeErr)
		}
		return nil
	}

	resizeFn := func(cols, rows uint16) error {
		return h.ptyManager.Resize(terminalID, cols, rows)
	}

	// 10. Start WritePump in goroutine, ReadPump blocks.
	go client.WritePump()
	client.ReadPump(ptyWriter, resizeFn)

	// ReadPump returned — connection is closed, client already unregistered.
	// Audit log the WebSocket disconnection.
	if auditErr := h.auditService.Log(context.Background(), &userID, &terminalID, domain.AuditActionTerminalDisconnect, map[string]interface{}{
		"client_id": client.ID,
	}); auditErr != nil {
		log.Error().Err(auditErr).Str("terminal_id", terminalIDStr).Msg("failed to write disconnect audit log")
	}

	log.Info().
		Str("client_id", client.ID).
		Str("terminal_id", terminalIDStr).
		Msg("WebSocket connection closed")
}
