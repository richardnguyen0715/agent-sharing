package handler

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/termhub/backend/internal/middleware"
	ptymanager "github.com/termhub/backend/internal/pty"
	"github.com/termhub/backend/internal/service"
)

// maxEditableFileSize caps files that can be opened in the editor (5 MiB) —
// large enough for any source file, small enough to keep edits snappy.
const maxEditableFileSize = 5 << 20

// FilesHandler exposes a VSCode-style file browser/editor scoped to a
// terminal's working directory. It grants no access beyond what the
// terminal's own shell already has — same user, same filesystem root.
type FilesHandler struct {
	terminalService *service.TerminalService
	jwtSecret       string
}

// NewFilesHandler creates a new FilesHandler.
func NewFilesHandler(terminalService *service.TerminalService, jwtSecret string) *FilesHandler {
	return &FilesHandler{terminalService: terminalService, jwtSecret: jwtSecret}
}

// TerminalFilesRoutes returns routes for terminal-scoped file browsing.
// Mounted under /api/terminals/{id}/files.
func (h *FilesHandler) TerminalFilesRoutes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Auth(h.jwtSecret))

	r.Get("/tree", h.Tree)
	r.Get("/content", h.GetContent)
	r.Put("/content", h.SaveContent)
	r.Post("/", h.Create)
	r.Post("/rename", h.Rename)
	r.Delete("/", h.Delete)

	return r
}

// TreeEntry describes a single file or directory entry.
type TreeEntry struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"` // relative to the terminal root, forward-slash separated
	IsDir   bool      `json:"isDir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

// rootAndPath resolves the authenticated caller's terminal and validates the
// requested relative path stays within that terminal's root directory.
func (h *FilesHandler) rootAndPath(r *http.Request, rawRel string) (root string, full string, err error) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		return "", "", errUnauthorized
	}

	ownerID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return "", "", errBadRequest
	}

	terminalID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return "", "", errBadRequest
	}

	terminal, err := h.terminalService.Get(r.Context(), terminalID, ownerID)
	if err != nil {
		return "", "", err
	}

	root = ptymanager.ResolveCWD(terminal.CWD)

	full, err = resolveWithinRoot(root, rawRel)
	if err != nil {
		return "", "", errBadRequest
	}

	return root, full, nil
}

var (
	errUnauthorized = errors.New("authentication required")
	errBadRequest   = errors.New("invalid request")
)

// resolveWithinRoot joins rel onto root and guarantees the result cannot
// escape root, regardless of ".." segments in rel.
func resolveWithinRoot(root, rel string) (string, error) {
	if strings.Contains(rel, "\x00") {
		return "", errBadRequest
	}
	// Prepending "/" then Cleaning collapses any leading ".." back to "/",
	// so the joined result can never climb above root.
	cleaned := filepath.Clean("/" + rel)
	return filepath.Join(root, cleaned), nil
}

func (h *FilesHandler) writeOwnershipError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errUnauthorized):
		writeError(w, http.StatusUnauthorized, "authentication required")
	case errors.Is(err, errBadRequest):
		writeError(w, http.StatusBadRequest, "invalid request")
	case errors.Is(err, service.ErrTerminalNotFound):
		writeError(w, http.StatusNotFound, "terminal not found")
	case errors.Is(err, service.ErrNotOwner):
		writeError(w, http.StatusForbidden, "access denied")
	default:
		log.Error().Err(err).Msg("files: failed to resolve terminal")
		writeError(w, http.StatusInternalServerError, "failed to resolve terminal")
	}
}

// Tree lists the contents of a single directory (non-recursive — the
// frontend expands directories on demand, same as VSCode's explorer).
func (h *FilesHandler) Tree(w http.ResponseWriter, r *http.Request) {
	root, full, err := h.rootAndPath(r, r.URL.Query().Get("path"))
	if err != nil {
		h.writeOwnershipError(w, err)
		return
	}

	info, err := os.Stat(full)
	if err != nil {
		writeError(w, http.StatusNotFound, "directory not found")
		return
	}
	if !info.IsDir() {
		writeError(w, http.StatusBadRequest, "path is not a directory")
		return
	}

	dirEntries, err := os.ReadDir(full)
	if err != nil {
		log.Error().Err(err).Str("path", full).Msg("files: failed to read directory")
		writeError(w, http.StatusInternalServerError, "failed to read directory")
		return
	}

	entries := make([]TreeEntry, 0, len(dirEntries))
	for _, de := range dirEntries {
		// Skip VCS internals — noise no one wants in a file tree.
		if de.Name() == ".git" {
			continue
		}
		entryInfo, err := de.Info()
		if err != nil {
			continue
		}
		relPath, err := filepath.Rel(root, filepath.Join(full, de.Name()))
		if err != nil {
			continue
		}
		entries = append(entries, TreeEntry{
			Name:    de.Name(),
			Path:    filepath.ToSlash(relPath),
			IsDir:   de.IsDir(),
			Size:    entryInfo.Size(),
			ModTime: entryInfo.ModTime(),
		})
	}

	// Directories first, then alphabetical — matches VSCode's default sort.
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"path":    filepath.ToSlash(strings.TrimPrefix(strings.TrimPrefix(full, root), "/")),
		"entries": entries,
	})
}

// GetContent reads a file's contents for editing.
func (h *FilesHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	_, full, err := h.rootAndPath(r, r.URL.Query().Get("path"))
	if err != nil {
		h.writeOwnershipError(w, err)
		return
	}

	info, err := os.Stat(full)
	if err != nil {
		writeError(w, http.StatusNotFound, "file not found")
		return
	}
	if info.IsDir() {
		writeError(w, http.StatusBadRequest, "path is a directory")
		return
	}
	if info.Size() > maxEditableFileSize {
		writeError(w, http.StatusRequestEntityTooLarge, "file too large to edit (max 5MB)")
		return
	}

	data, err := os.ReadFile(full)
	if err != nil {
		log.Error().Err(err).Str("path", full).Msg("files: failed to read file")
		writeError(w, http.StatusInternalServerError, "failed to read file")
		return
	}

	if looksBinary(data) {
		writeError(w, http.StatusUnsupportedMediaType, "binary file cannot be edited")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"content": string(data),
		"size":    info.Size(),
	})
}

// looksBinary makes a best-effort guess at whether data is a text file, the
// same heuristic `file`/git use: a NUL byte in the first chunk means binary.
func looksBinary(data []byte) bool {
	probe := data
	if len(probe) > 8000 {
		probe = probe[:8000]
	}
	return bytes.IndexByte(probe, 0) != -1
}

// SaveContent overwrites an existing file with new content.
func (h *FilesHandler) SaveContent(w http.ResponseWriter, r *http.Request) {
	_, full, err := h.rootAndPath(r, r.URL.Query().Get("path"))
	if err != nil {
		h.writeOwnershipError(w, err)
		return
	}

	info, statErr := os.Stat(full)
	if statErr != nil {
		writeError(w, http.StatusNotFound, "file not found")
		return
	}
	if info.IsDir() {
		writeError(w, http.StatusBadRequest, "path is a directory")
		return
	}

	var body struct {
		Content string `json:"content"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := os.WriteFile(full, []byte(body.Content), 0o644); err != nil {
		log.Error().Err(err).Str("path", full).Msg("files: failed to save file")
		writeError(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "saved"})
}

// Create makes a new empty file or directory.
func (h *FilesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Path string `json:"path"`
		Type string `json:"type"` // "file" or "dir"
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Path == "" {
		writeError(w, http.StatusBadRequest, "path is required")
		return
	}

	_, full, err := h.rootAndPath(r, body.Path)
	if err != nil {
		h.writeOwnershipError(w, err)
		return
	}

	if _, statErr := os.Stat(full); statErr == nil {
		writeError(w, http.StatusConflict, "already exists")
		return
	}

	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		log.Error().Err(err).Str("path", full).Msg("files: failed to create parent directory")
		writeError(w, http.StatusInternalServerError, "failed to create")
		return
	}

	if body.Type == "dir" {
		if err := os.Mkdir(full, 0o755); err != nil {
			log.Error().Err(err).Str("path", full).Msg("files: failed to create directory")
			writeError(w, http.StatusInternalServerError, "failed to create directory")
			return
		}
	} else {
		f, err := os.OpenFile(full, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			log.Error().Err(err).Str("path", full).Msg("files: failed to create file")
			writeError(w, http.StatusInternalServerError, "failed to create file")
			return
		}
		_ = f.Close()
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "created", "path": body.Path})
}

// Rename moves/renames a file or directory within the terminal root.
func (h *FilesHandler) Rename(w http.ResponseWriter, r *http.Request) {
	var body struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.From == "" || body.To == "" {
		writeError(w, http.StatusBadRequest, "from and to are required")
		return
	}

	root, fromFull, err := h.rootAndPath(r, body.From)
	if err != nil {
		h.writeOwnershipError(w, err)
		return
	}
	toFull, err := resolveWithinRoot(root, body.To)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid destination path")
		return
	}

	if _, statErr := os.Stat(fromFull); statErr != nil {
		writeError(w, http.StatusNotFound, "source not found")
		return
	}
	if _, statErr := os.Stat(toFull); statErr == nil {
		writeError(w, http.StatusConflict, "destination already exists")
		return
	}

	if err := os.MkdirAll(filepath.Dir(toFull), 0o755); err != nil {
		log.Error().Err(err).Str("path", toFull).Msg("files: failed to create destination directory")
		writeError(w, http.StatusInternalServerError, "failed to rename")
		return
	}

	if err := os.Rename(fromFull, toFull); err != nil {
		log.Error().Err(err).Str("from", fromFull).Str("to", toFull).Msg("files: failed to rename")
		writeError(w, http.StatusInternalServerError, "failed to rename")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "renamed", "path": body.To})
}

// Delete removes a file or directory (recursively for directories).
func (h *FilesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	root, full, err := h.rootAndPath(r, r.URL.Query().Get("path"))
	if err != nil {
		h.writeOwnershipError(w, err)
		return
	}

	if full == root {
		writeError(w, http.StatusBadRequest, "cannot delete the root directory")
		return
	}

	if _, statErr := os.Stat(full); statErr != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	if err := os.RemoveAll(full); err != nil {
		log.Error().Err(err).Str("path", full).Msg("files: failed to delete")
		writeError(w, http.StatusInternalServerError, "failed to delete")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
