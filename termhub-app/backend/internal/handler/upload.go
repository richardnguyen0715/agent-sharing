package handler

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/termhub/backend/internal/middleware"
	"github.com/termhub/backend/internal/service"
)

// maxUploadSize caps a single pasted/dropped file (25 MiB).
const maxUploadSize = 25 << 20

// UploadHandler handles files pasted or dropped into a terminal (e.g. images
// copied from the clipboard). A terminal can only accept text, so uploaded
// files are written to disk and the resulting path is handed back to the
// frontend, which types it into the terminal in place of the raw paste.
type UploadHandler struct {
	terminalService *service.TerminalService
	jwtSecret       string
	baseDir         string
}

// NewUploadHandler creates a new UploadHandler. Files are stored under
// baseDir/{terminalID}/. If baseDir is empty, it defaults to "./uploads".
func NewUploadHandler(terminalService *service.TerminalService, jwtSecret string, baseDir string) *UploadHandler {
	if baseDir == "" {
		baseDir = "./uploads"
	}
	return &UploadHandler{terminalService: terminalService, jwtSecret: jwtSecret, baseDir: baseDir}
}

// TerminalUploadRoutes returns routes for terminal-scoped file uploads.
// Mounted under /api/terminals/{id}/upload.
func (h *UploadHandler) TerminalUploadRoutes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Auth(h.jwtSecret))
	r.Post("/", h.Upload)
	return r
}

// Upload accepts a multipart file upload, saves it to disk, and returns its
// absolute path so it can be referenced from the terminal's shell.
func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
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

	// Verify ownership before writing anything to disk.
	if _, err := h.terminalService.Get(r.Context(), terminalID, ownerID); err != nil {
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

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeError(w, http.StatusBadRequest, "file too large or invalid multipart form (max 25MB)")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing file field")
		return
	}
	defer file.Close()

	dir := filepath.Join(h.baseDir, terminalID.String())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Error().Err(err).Str("terminal_id", terminalID.String()).Msg("failed to create upload directory")
		writeError(w, http.StatusInternalServerError, "failed to store file")
		return
	}

	// Sanitize: discard the client-supplied basename entirely (avoids path
	// traversal and collisions) and keep only its extension.
	safeName := uuid.New().String() + filepath.Ext(header.Filename)
	destPath := filepath.Join(dir, safeName)

	dest, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		log.Error().Err(err).Str("terminal_id", terminalID.String()).Msg("failed to create upload file")
		writeError(w, http.StatusInternalServerError, "failed to store file")
		return
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		log.Error().Err(err).Str("terminal_id", terminalID.String()).Msg("failed to write upload file")
		writeError(w, http.StatusInternalServerError, "failed to store file")
		return
	}

	absPath, err := filepath.Abs(destPath)
	if err != nil {
		absPath = destPath
	}

	log.Info().
		Str("terminal_id", terminalID.String()).
		Str("user_id", claims.UserID).
		Str("original_filename", header.Filename).
		Str("path", absPath).
		Msg("file uploaded to terminal")

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"path":     absPath,
		"filename": header.Filename,
	})
}
