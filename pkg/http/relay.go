package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptrace"
	"sync"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func (s *Server) RelayHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := s.tracer.Start(r.Context(), "relayhandler")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.Error("Request failed", "Error", err)
		return
	}
	defer r.Body.Close()

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	var result string
	var wg sync.WaitGroup
	wg.Add(1)

	if s.config.BackendURL != "" {
		backend := s.config.BackendURL
		go func(timeout time.Duration) {
			defer wg.Done()
			ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, r.Method, backend, bytes.NewReader(body))
			if err != nil {
				s.logger.Error("backend call failed", "url", backend, "Error", err)
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				s.logger.Error("Request failed", "Error", err)
				return
			}
			defer resp.Body.Close()

			rbody, err := io.ReadAll(resp.Body)
			if err != nil {
				s.logger.Info("Response is not readable", "status", resp.Status)
			}
			result = string(rbody)
		}(s.config.ReadTimeout)
		wg.Wait()
		s.JSONResponse(w, r, result)

	} else {
		w.WriteHeader(http.StatusAccepted)
		w.Write(body)
	}
}
