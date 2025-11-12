package http

import (
	"log/slog"
	"net/http"
)

type LoggingMiddleware struct {
	logger *slog.Logger
}

func NewLoggingMiddleware(logger *slog.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

func (lm *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		lm.logger.Info("request recieved",
			"addr", r.RemoteAddr,
			"uri", r.RequestURI,
			"method", r.Method,
			"Length", r.ContentLength,
			"user-agent", r.UserAgent(),
		)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
