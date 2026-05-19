package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/termhub/backend/internal/middleware"
	"github.com/termhub/backend/internal/recording"
	"github.com/termhub/backend/internal/service"
)

// RecordingHandler handles recording-related HTTP requests.
type RecordingHandler struct {
	recordingMgr *recording.Manager
	authService  *service.AuthService
	jwtSecret    string
}

// NewRecordingHandler creates a new RecordingHandler.
func NewRecordingHandler(mgr *recording.Manager, authSvc *service.AuthService, jwtSecret string) *RecordingHandler {
	return &RecordingHandler{
		recordingMgr: mgr,
		authService:  authSvc,
		jwtSecret:    jwtSecret,
	}
}

// Routes returns a chi.Router with recording routes mounted.
func (h *RecordingHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Auth(h.jwtSecret))

	r.Get("/", h.ListRecordings)
	r.Get("/{terminalId}/{filename}", h.GetRecording)

	return r
}

// TerminalRecordingRoutes returns routes for terminal-scoped recording control.
// These should be mounted under /api/terminals/{id}/recording.
func (h *RecordingHandler) TerminalRecordingRoutes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Auth(h.jwtSecret))

	r.Post("/start", h.StartRecording)
	r.Post("/stop", h.StopRecording)

	return r
}

// ListRecordings returns all available recordings.
func (h *RecordingHandler) ListRecordings(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	recordings, err := h.recordingMgr.ListRecordings()
	if err != nil {
		log.Error().Err(err).Msg("failed to list recordings")
		writeError(w, http.StatusInternalServerError, "failed to list recordings")
		return
	}

	if recordings == nil {
		recordings = []recording.RecordingInfo{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"recordings": recordings,
		"count":      len(recordings),
	})
}

// GetRecording returns the events of a specific recording for playback.
func (h *RecordingHandler) GetRecording(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	terminalID := chi.URLParam(r, "terminalId")
	filename := chi.URLParam(r, "filename")

	if terminalID == "" || filename == "" {
		writeError(w, http.StatusBadRequest, "terminalId and filename are required")
		return
	}

	events, err := h.recordingMgr.GetRecording(terminalID, filename)
	if err != nil {
		log.Error().Err(err).
			Str("terminal_id", terminalID).
			Str("filename", filename).
			Msg("failed to get recording")
		writeError(w, http.StatusNotFound, "recording not found")
		return
	}

	if events == nil {
		events = []recording.Event{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"terminal_id": terminalID,
		"filename":    filename,
		"events":      events,
		"count":       len(events),
	})
}

// StartRecording begins recording a terminal session.
func (h *RecordingHandler) StartRecording(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	terminalID := chi.URLParam(r, "id")
	if terminalID == "" {
		writeError(w, http.StatusBadRequest, "terminal ID is required")
		return
	}

	if err := h.recordingMgr.StartRecording(terminalID, claims.UserID); err != nil {
		log.Error().Err(err).Str("terminal_id", terminalID).Msg("failed to start recording")
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message":     "recording started",
		"terminal_id": terminalID,
	})
}

// StopRecording stops recording a terminal session.
func (h *RecordingHandler) StopRecording(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	terminalID := chi.URLParam(r, "id")
	if terminalID == "" {
		writeError(w, http.StatusBadRequest, "terminal ID is required")
		return
	}

	if err := h.recordingMgr.StopRecording(terminalID); err != nil {
		log.Error().Err(err).Str("terminal_id", terminalID).Msg("failed to stop recording")
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message":     "recording stopped",
		"terminal_id": terminalID,
	})
}
