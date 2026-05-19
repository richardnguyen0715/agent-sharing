package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/termhub/backend/internal/domain"
	"github.com/termhub/backend/internal/middleware"
	"github.com/termhub/backend/internal/service"
)

// CreateTerminalRequest represents the create terminal request body.
type CreateTerminalRequest struct {
	Title    string `json:"title"`
	CWD      string `json:"cwd"`
	Shell    string `json:"shell"`
	ColorTag string `json:"colorTag"`
}

// ResizeTerminalRequest represents the resize terminal request body.
type ResizeTerminalRequest struct {
	Cols uint16 `json:"cols"`
	Rows uint16 `json:"rows"`
}

// TerminalHandler handles terminal-related HTTP requests.
type TerminalHandler struct {
	terminalService *service.TerminalService
	auditService    *service.AuditService
	jwtSecret       string
}

// NewTerminalHandler creates a new TerminalHandler.
func NewTerminalHandler(terminalService *service.TerminalService, auditService *service.AuditService, jwtSecret string) *TerminalHandler {
	return &TerminalHandler{terminalService: terminalService, auditService: auditService, jwtSecret: jwtSecret}
}

// Routes returns a chi.Router with all terminal routes mounted.
func (h *TerminalHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Auth(h.jwtSecret))

	r.Post("/", h.Create)
	r.Get("/", h.List)
	r.Get("/{id}", h.GetByID)
	r.Delete("/{id}", h.Delete)
	r.Post("/{id}/resize", h.Resize)

	return r
}

// Create provisions a new terminal.
func (h *TerminalHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	ownerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var req CreateTerminalRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	terminal, err := h.terminalService.Create(r.Context(), ownerID, req.Title, req.CWD, req.Shell, req.ColorTag)
	if err != nil {
		log.Error().Err(err).Str("owner_id", claims.UserID).Msg("failed to create terminal")
		writeError(w, http.StatusInternalServerError, "failed to create terminal")
		return
	}

	// Audit log the terminal creation.
	termID := terminal.ID
	if auditErr := h.auditService.Log(r.Context(), &ownerID, &termID, domain.AuditActionTerminalCreate, map[string]interface{}{
		"title": terminal.Title,
		"shell": terminal.Shell,
		"cwd":   terminal.CWD,
	}); auditErr != nil {
		log.Error().Err(auditErr).Str("terminal_id", termID.String()).Msg("failed to write terminal create audit log")
	}

	writeJSON(w, http.StatusCreated, terminal)
}

// List returns all terminals for the authenticated user.
func (h *TerminalHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	ownerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	terminals, err := h.terminalService.List(r.Context(), ownerID)
	if err != nil {
		log.Error().Err(err).Str("owner_id", claims.UserID).Msg("failed to list terminals")
		writeError(w, http.StatusInternalServerError, "failed to list terminals")
		return
	}

	writeJSON(w, http.StatusOK, terminals)
}

// GetByID retrieves a specific terminal.
func (h *TerminalHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	ownerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	terminalID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid terminal ID")
		return
	}

	terminal, err := h.terminalService.Get(r.Context(), terminalID, ownerID)
	if err != nil {
		if errors.Is(err, service.ErrTerminalNotFound) {
			writeError(w, http.StatusNotFound, "terminal not found")
			return
		}
		if errors.Is(err, service.ErrNotOwner) {
			writeError(w, http.StatusForbidden, "access denied")
			return
		}
		log.Error().Err(err).Str("terminal_id", terminalID.String()).Msg("failed to get terminal")
		writeError(w, http.StatusInternalServerError, "failed to get terminal")
		return
	}

	writeJSON(w, http.StatusOK, terminal)
}

// Delete terminates and removes a terminal.
func (h *TerminalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	ownerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	terminalID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid terminal ID")
		return
	}

	// Audit log BEFORE delete — terminal_id FK must still exist.
	if auditErr := h.auditService.Log(r.Context(), &ownerID, &terminalID, domain.AuditActionTerminalDelete, map[string]interface{}{
		"ip_address": r.RemoteAddr,
	}); auditErr != nil {
		log.Error().Err(auditErr).Str("terminal_id", terminalID.String()).Msg("failed to write terminal delete audit log")
	}

	if err := h.terminalService.Delete(r.Context(), terminalID, ownerID); err != nil {
		if errors.Is(err, service.ErrTerminalNotFound) {
			writeError(w, http.StatusNotFound, "terminal not found")
			return
		}
		if errors.Is(err, service.ErrNotOwner) {
			writeError(w, http.StatusForbidden, "access denied")
			return
		}
		log.Error().Err(err).Str("terminal_id", terminalID.String()).Msg("failed to delete terminal")
		writeError(w, http.StatusInternalServerError, "failed to delete terminal")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "terminal deleted",
	})
}

// Resize updates the terminal dimensions.
func (h *TerminalHandler) Resize(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	ownerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	terminalID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid terminal ID")
		return
	}

	// Verify ownership.
	_, err = h.terminalService.Get(r.Context(), terminalID, ownerID)
	if err != nil {
		if errors.Is(err, service.ErrTerminalNotFound) {
			writeError(w, http.StatusNotFound, "terminal not found")
			return
		}
		if errors.Is(err, service.ErrNotOwner) {
			writeError(w, http.StatusForbidden, "access denied")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to verify terminal ownership")
		return
	}

	var req ResizeTerminalRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Cols == 0 || req.Rows == 0 {
		writeError(w, http.StatusBadRequest, "cols and rows must be greater than 0")
		return
	}

	if err := h.terminalService.Resize(r.Context(), terminalID, req.Cols, req.Rows); err != nil {
		log.Error().Err(err).Str("terminal_id", terminalID.String()).Msg("failed to resize terminal")
		writeError(w, http.StatusInternalServerError, "failed to resize terminal")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "terminal resized",
	})
}
