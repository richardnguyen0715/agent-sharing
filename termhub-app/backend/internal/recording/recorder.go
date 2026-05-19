// Package recording provides session recording for terminal I/O.
// Recordings are stored as JSONL files (one JSON event per line) for
// efficient append-only writing and streaming playback.
package recording

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Event represents a single recorded I/O event.
type Event struct {
	Time float64 `json:"time"`
	Type string  `json:"type"` // "o" for output, "i" for input
	Data string  `json:"data"`
}

// RecordingInfo contains metadata about a stored recording file.
type RecordingInfo struct {
	TerminalID string    `json:"terminal_id"`
	Filename   string    `json:"filename"`
	StartedAt  time.Time `json:"started_at"`
	Size       int64     `json:"size_bytes"`
}

// Recording tracks an active recording session for a single terminal.
type Recording struct {
	TerminalID string
	UserID     string
	StartedAt  time.Time
	Events     []Event
	mu         sync.Mutex
	file       *os.File
	writer     *bufio.Writer
}

// Manager manages recording sessions across all terminals.
type Manager struct {
	recordings map[string]*Recording // key: terminal ID
	mu         sync.RWMutex
	baseDir    string
	logger     zerolog.Logger
}

// NewManager creates a new recording Manager.
// If baseDir is empty, it defaults to "./recordings".
func NewManager(baseDir string) *Manager {
	if baseDir == "" {
		baseDir = "./recordings"
	}
	return &Manager{
		recordings: make(map[string]*Recording),
		baseDir:    baseDir,
		logger:     log.With().Str("component", "recording-manager").Logger(),
	}
}

// StartRecording begins recording I/O for the given terminal.
func (m *Manager) StartRecording(terminalID, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.recordings[terminalID]; exists {
		return fmt.Errorf("recording: terminal %s is already being recorded", terminalID)
	}

	// Create directory for this terminal's recordings.
	dir := filepath.Join(m.baseDir, terminalID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("recording: failed to create directory %s: %w", dir, err)
	}

	now := time.Now().UTC()
	filename := now.Format("20060102-150405") + ".jsonl"
	filePath := filepath.Join(dir, filename)

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("recording: failed to create file %s: %w", filePath, err)
	}

	rec := &Recording{
		TerminalID: terminalID,
		UserID:     userID,
		StartedAt:  now,
		Events:     make([]Event, 0),
		file:       f,
		writer:     bufio.NewWriter(f),
	}

	m.recordings[terminalID] = rec

	m.logger.Info().
		Str("terminal_id", terminalID).
		Str("user_id", userID).
		Str("file", filePath).
		Msg("recording started")

	return nil
}

// StopRecording flushes and closes the recording for the given terminal.
func (m *Manager) StopRecording(terminalID string) error {
	m.mu.Lock()
	rec, exists := m.recordings[terminalID]
	if !exists {
		m.mu.Unlock()
		return fmt.Errorf("recording: no active recording for terminal %s", terminalID)
	}
	delete(m.recordings, terminalID)
	m.mu.Unlock()

	rec.mu.Lock()
	defer rec.mu.Unlock()

	if err := rec.writer.Flush(); err != nil {
		m.logger.Error().Err(err).Str("terminal_id", terminalID).Msg("failed to flush recording")
	}

	if err := rec.file.Close(); err != nil {
		return fmt.Errorf("recording: failed to close file for terminal %s: %w", terminalID, err)
	}

	m.logger.Info().
		Str("terminal_id", terminalID).
		Int("event_count", len(rec.Events)).
		Msg("recording stopped")

	return nil
}

// RecordOutput appends an output event to the recording. Non-blocking, best-effort.
func (m *Manager) RecordOutput(terminalID string, data []byte) {
	m.mu.RLock()
	rec, exists := m.recordings[terminalID]
	m.mu.RUnlock()

	if !exists {
		return
	}

	rec.mu.Lock()
	defer rec.mu.Unlock()

	elapsed := time.Since(rec.StartedAt).Seconds()
	event := Event{
		Time: elapsed,
		Type: "o",
		Data: string(data),
	}
	rec.Events = append(rec.Events, event)

	if line, err := json.Marshal(event); err == nil {
		rec.writer.Write(line)
		rec.writer.WriteByte('\n')
		// Flush periodically — every 100 events or so.
		if len(rec.Events)%100 == 0 {
			rec.writer.Flush()
		}
	}
}

// RecordInput appends an input event to the recording. Non-blocking, best-effort.
func (m *Manager) RecordInput(terminalID string, data string) {
	m.mu.RLock()
	rec, exists := m.recordings[terminalID]
	m.mu.RUnlock()

	if !exists {
		return
	}

	rec.mu.Lock()
	defer rec.mu.Unlock()

	elapsed := time.Since(rec.StartedAt).Seconds()
	event := Event{
		Time: elapsed,
		Type: "i",
		Data: data,
	}
	rec.Events = append(rec.Events, event)

	if line, err := json.Marshal(event); err == nil {
		rec.writer.Write(line)
		rec.writer.WriteByte('\n')
		if len(rec.Events)%100 == 0 {
			rec.writer.Flush()
		}
	}
}

// IsRecording reports whether the given terminal is currently being recorded.
func (m *Manager) IsRecording(terminalID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.recordings[terminalID]
	return exists
}

// ListRecordings returns metadata for all recording files found in baseDir.
func (m *Manager) ListRecordings() ([]RecordingInfo, error) {
	var infos []RecordingInfo

	termDirs, err := os.ReadDir(m.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return infos, nil
		}
		return nil, fmt.Errorf("recording: failed to read base dir: %w", err)
	}

	for _, termDir := range termDirs {
		if !termDir.IsDir() {
			continue
		}
		termID := termDir.Name()
		dirPath := filepath.Join(m.baseDir, termID)

		files, err := os.ReadDir(dirPath)
		if err != nil {
			m.logger.Warn().Err(err).Str("dir", dirPath).Msg("failed to read terminal recording dir")
			continue
		}

		for _, f := range files {
			if f.IsDir() {
				continue
			}
			info, err := f.Info()
			if err != nil {
				continue
			}

			infos = append(infos, RecordingInfo{
				TerminalID: termID,
				Filename:   f.Name(),
				StartedAt:  info.ModTime(),
				Size:       info.Size(),
			})
		}
	}

	return infos, nil
}

// GetRecording reads and parses a specific recording file.
func (m *Manager) GetRecording(terminalID, filename string) ([]Event, error) {
	filePath := filepath.Join(m.baseDir, terminalID, filename)

	// Sanitize to prevent directory traversal.
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("recording: invalid path: %w", err)
	}
	absBase, err := filepath.Abs(m.baseDir)
	if err != nil {
		return nil, fmt.Errorf("recording: invalid base path: %w", err)
	}
	if len(absPath) < len(absBase) || absPath[:len(absBase)] != absBase {
		return nil, fmt.Errorf("recording: path traversal not allowed")
	}

	f, err := os.Open(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("recording: file not found: %s/%s", terminalID, filename)
		}
		return nil, fmt.Errorf("recording: failed to open file: %w", err)
	}
	defer f.Close()

	var events []Event
	scanner := bufio.NewScanner(f)
	// Increase buffer for potentially large output lines.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		var event Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			m.logger.Warn().Err(err).Str("file", filePath).Msg("skipping malformed recording line")
			continue
		}
		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("recording: error reading file: %w", err)
	}

	return events, nil
}

// Cleanup stops all active recordings. Called during graceful shutdown.
func (m *Manager) Cleanup() {
	m.mu.Lock()
	ids := make([]string, 0, len(m.recordings))
	for id := range m.recordings {
		ids = append(ids, id)
	}
	m.mu.Unlock()

	for _, id := range ids {
		if err := m.StopRecording(id); err != nil {
			m.logger.Error().Err(err).Str("terminal_id", id).Msg("error stopping recording during cleanup")
		}
	}

	m.logger.Info().Msg("recording cleanup complete")
}
