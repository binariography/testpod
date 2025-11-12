package http

import (
	"log/slog"
	"os"

	"github.com/gorilla/mux"
)

const (
	LevelTrace     = slog.Level(-8)
	LevelDebug     = slog.LevelDebug
	LevelInfo      = slog.LevelInfo
	LevelNotice    = slog.Level(2)
	LevelWarning   = slog.LevelWarn
	LevelError     = slog.LevelError
	LevelEmergency = slog.Level(12)
)

func NewMockServer() (*Server, error) {
	config := &Config{
		Port: 9898,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// Set a custom level to show all log output. The default value is
		// LevelInfo, which would drop Debug and Trace logs.
		Level: LevelTrace,

		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time from the output for predictable test output.
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))

	return &Server{
		router: mux.NewRouter(),
		config: config,
		logger: logger,
	}, nil
}
