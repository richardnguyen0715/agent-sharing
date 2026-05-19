// Package handler provides HTTP handlers and JSON response utilities.
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ErrorResponse is the standard error response body.
type ErrorResponse struct {
	Error string `json:"error"`
}

// writeJSON writes a JSON response with the given status code and data.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// At this point headers are already sent, so we can only log.
			// The caller should ensure data is serializable.
			_ = err
		}
	}
}

// writeError writes a JSON error response with the given status code and message.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

// readJSON decodes the request body into the destination struct.
// It enforces a 1MB body size limit and rejects unknown fields.
func readJSON(r *http.Request, dst interface{}) error {
	const maxBodySize = 1 << 20 // 1 MB

	r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("handler: invalid JSON body: %w", err)
	}

	// Ensure there's no additional JSON in the body.
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("handler: body must contain a single JSON object")
	}

	return nil
}
