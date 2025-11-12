package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/binariography/testpod/pkg/http"

	// From Individuals
	"github.com/alexflint/go-arg"
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

func initLog() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		// Set a custom level to show all log output. The default value is
		// LevelInfo, which would drop Debug and Trace logs.
		Level: LevelDebug,

		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time from the output for predictable test output.
			if a.Key == slog.TimeKey {
				a.Key = "ts"
			}

			// Customize the name of the level key and the output string, including
			// custom level values.
			if a.Key == slog.LevelKey {
				// Handle custom level values.
				level := a.Value.Any().(slog.Level)

				switch {
				case level < LevelDebug:
					a.Value = slog.StringValue("TRACE")
				case level < LevelInfo:
					a.Value = slog.StringValue("DEBUG")
				case level < LevelNotice:
					a.Value = slog.StringValue("INFO")
				case level < LevelWarning:
					a.Value = slog.StringValue("NOTICE")
				case level < LevelError:
					a.Value = slog.StringValue("WARNING")
				case level < LevelEmergency:
					a.Value = slog.StringValue("ERROR")
				default:
					a.Value = slog.StringValue("EMERGENCY")
				}
			}

			return a
		},
	}))

	logger.Info("Logging Initialized...")

	return logger
}

func main() {

	logger := initLog()
	var srvConf http.Config

	arg.MustParse(&srvConf)
	hostname, err := os.Hostname()

	// At this stage Hostname is not important for running the server
	// only A warning log
	if err != nil {
		logger.Warn("Cannot Retrieve Hostname", "Error", err)
	}
	srvConf.Hostname = hostname
	srv, _ := http.NewServer(&srvConf, logger)

	//Start HTTP server
	httpServer := srv.ListenAndServe()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	s := <-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	logger.Info("msg", "context", ctx)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	logger.Info("Shutting down, Bye...", "Signal", s)
	httpServer.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	os.Exit(0)
}
