package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/termhub/backend/internal/domain"
	"github.com/termhub/backend/internal/middleware"
	ptymanager "github.com/termhub/backend/internal/pty"
	"github.com/termhub/backend/internal/repository"
	"github.com/termhub/backend/internal/service"
	wshub "github.com/termhub/backend/internal/ws"
)

// AdminHandler handles administrative HTTP requests.
type AdminHandler struct {
	auditService    *service.AuditService
	terminalService *service.TerminalService
	userRepo        *repository.UserRepository
	ptyManager      *ptymanager.Manager
	hub             *wshub.Hub
	jwtSecret       string
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(
	auditService *service.AuditService,
	terminalService *service.TerminalService,
	userRepo *repository.UserRepository,
	ptyManager *ptymanager.Manager,
	hub *wshub.Hub,
	jwtSecret string,
) *AdminHandler {
	return &AdminHandler{
		auditService:    auditService,
		terminalService: terminalService,
		userRepo:        userRepo,
		ptyManager:      ptyManager,
		hub:             hub,
		jwtSecret:       jwtSecret,
	}
}

// Routes returns a chi.Router with all admin routes mounted.
func (h *AdminHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Auth(h.jwtSecret))
	r.Use(middleware.RequireRole(domain.RoleAdmin))

	r.Get("/users", h.ListUsers)
	r.Get("/audit-logs", h.ListAuditLogs)
	r.Get("/terminals", h.ListAllTerminals)
	r.Get("/metrics", h.Metrics)

	return r
}

// ListUsers returns all registered users (admin only).
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 100
	offset := 0
	if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
		limit = v
	}
	if v, err := strconv.Atoi(offsetStr); err == nil && v >= 0 {
		offset = v
	}

	users, err := h.userRepo.List(r.Context(), limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to list users")
		writeError(w, http.StatusInternalServerError, "failed to list users")
		return
	}

	if users == nil {
		users = []*domain.User{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"users": users,
		"count": len(users),
	})
}

// ListAuditLogs returns audit log entries with pagination (admin only).
func (h *AdminHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 100
	offset := 0
	if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
		limit = v
	}
	if v, err := strconv.Atoi(offsetStr); err == nil && v >= 0 {
		offset = v
	}

	logs, err := h.auditService.List(r.Context(), limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to list audit logs")
		writeError(w, http.StatusInternalServerError, "failed to list audit logs")
		return
	}

	if logs == nil {
		logs = []*domain.AuditLog{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"audit_logs": logs,
		"count":      len(logs),
	})
}

// ListAllTerminals returns all terminals across all users (admin only).
func (h *AdminHandler) ListAllTerminals(w http.ResponseWriter, r *http.Request) {
	terminals, err := h.terminalService.ListAll(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to list all terminals")
		writeError(w, http.StatusInternalServerError, "failed to list terminals")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"terminals": terminals,
		"count":     len(terminals),
	})
}

// Metrics returns basic system metrics (admin only).
func (h *AdminHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	userCount, err := h.userRepo.Count(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to get user count for metrics")
		userCount = -1
	}

	terminalRepo := h.terminalService.GetRepo()
	totalTerminals, err := terminalRepo.Count(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to get terminal count for metrics")
		totalTerminals = -1
	}

	activeTerminals, err := terminalRepo.CountByStatus(r.Context(), domain.TerminalStatusActive)
	if err != nil {
		log.Error().Err(err).Msg("failed to get active terminal count for metrics")
		activeTerminals = -1
	}

	activePTYSessions := h.ptyManager.ActiveCount()
	wsConnections := h.hub.ClientCount()

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"users":               userCount,
		"total_terminals":     totalTerminals,
		"active_terminals":    activeTerminals,
		"active_pty_sessions": activePTYSessions,
		"ws_connections":      wsConnections,
	})
}
