package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging returns middleware that logs each HTTP request using zerolog.
func Logging(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := newResponseWriter(w)

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			var event *zerolog.Event
			switch {
			case wrapped.statusCode >= 500:
				event = logger.Error()
			case wrapped.statusCode >= 400:
				event = logger.Warn()
			default:
				event = logger.Info()
			}

			event.
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", wrapped.statusCode).
				Dur("duration", duration).
				Str("remote_addr", r.RemoteAddr).
				Msg("request completed")
		})
	}
}

func init() {
	// Suppress unused import warning — log is used indirectly via zerolog's global logger.
	_ = log.Logger
}
