package http

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Hostname     string
	Host         string        `arg:"--host,env:HOST" default:"localhost"`
	WriteTimeout time.Duration `arg:"--write-timeout,env:WRITE_TIMEOUT" default:"15s"`
	ReadTimeout  time.Duration `arg:"--read-timeout,env:READ_TIMEOUT" default:"15s"`
	IdleTimeout  time.Duration `arg:"--idle-timeout,env:IDLE_TIMEOUT" default:"60s"`
	Port         int           `arg:"--port,env:PORT" default:"8081" help:"Port that server is listening on"`
	PortMetrics  int           `arg:"--port-metrics,env:PORT_METRICS" default:"9090" help:"Port that Prometheus is listening on"`
	LogLevel     string        `arg:"--log-level,env:LOG_LEVEL" default:"info" help:"set log level"`
	BackendURL   string        `arg:"--backend-url,env:BACKEND_URL" help:"set backend service URL"`
}

var T = true

type Server struct {
	router         *mux.Router
	config         *Config
	logger         *slog.Logger
	tracer         trace.Tracer
	tracerProvider *sdktrace.TracerProvider
}

func NewServer(config *Config, logger *slog.Logger) (*Server, error) {
	srv := &Server{
		router: mux.NewRouter(),
		config: config,
		logger: logger,
	}

	return srv, nil
}

func (s *Server) ListenAndServe() *http.Server {
	ctx := context.Background()

	fmt.Println(ctx)

	s.initTracer(ctx)

	go s.startMetricServer()

	s.registerHandlers()
	s.registerMiddlewares()

	srv := s.startServer()

	return srv
}

func (s *Server) registerHandlers() {
	s.router.HandleFunc("/info", s.infoHandler)
	s.router.HandleFunc("/relay/{text}", s.RelayHandler).Methods("POST")
}

func (s *Server) registerMiddlewares() {
	prom := NewMetricMiddleware()
	s.router.Use(prom.Handler)
	httpTracer := NewOtelMiddleware()
	s.router.Use(httpTracer)
	httpLogger := NewLoggingMiddleware(s.logger)
	s.router.Use(httpLogger.Handler)
}

func (s *Server) startServer() *http.Server {

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler: s.router,
	}

	go func() {
		s.logger.Info(
			"starting pod....",
			"address", srv.Addr,
		)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Info("HTTP server crashed:", "error", err)
		}
	}()

	return srv
}

func (s *Server) startMetricServer() {
	mux := http.DefaultServeMux
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{ ok }"))
	})
	s.logger.Info("Starting metrics server",
		"port", s.config.PortMetrics)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.PortMetrics),
		Handler: mux,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Metric server crashed: ", err)
	}
}
