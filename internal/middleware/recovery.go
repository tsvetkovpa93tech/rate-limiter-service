package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// RecoveryMiddleware returns a middleware that recovers from panics
func RecoveryMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					logger.Error("Panic recovered",
						"error", rvr,
						"stack", string(debug.Stack()),
						"path", r.URL.Path,
						"method", r.Method,
						"remote_addr", r.RemoteAddr,
					)

					if rvr == http.ErrAbortHandler {
						panic(rvr)
					}

					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, `{"error":"Internal server error"}`)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

