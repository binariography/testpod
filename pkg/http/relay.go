package http

import (
	"bytes"
	"context"
	"fmt"
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
			backendReq, err := http.NewRequestWithContext(ctx, r.Method, backend, bytes.NewReader(body))
			if err != nil {
				s.logger.Error("backend call failed", "url", backend, "Error", err)
				return
			}

			// forward headers
			copyTracingHeaders(r, backendReq)

			resp, err := client.Do(backendReq)
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

func copyTracingHeaders(from *http.Request, to *http.Request) {
	fmt.Println("#################################COPYTRACINGHEADERS", from.Header)
	headers := []string{
		"x-request-id",
		"x-b3-traceid",
		"x-b3-spanid",
		"x-b3-parentspanid",
		"x-b3-sampled",
		"x-b3-flags",
		"x-ot-span-context",
	}

	for i := range headers {
		headerValue := from.Header.Get(headers[i])
		if len(headerValue) > 0 {
			fmt.Println("###################################################", headers[i], ":", headerValue)
			to.Header.Set(headers[i], headerValue)
		}
	}
}
