package http

import (
	"net/http"
	"os"
	"strings"
)

type Info struct {
	Title   string
	Headers map[string]string
	EnvVars []string
}

func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	_, span := s.tracer.Start(r.Context(), "infohandler")
	defer span.End()

	httpBody := Info{}
	httpBody.Title = s.config.Hostname
	httpBody.Headers = make(map[string]string)
	httpBody.EnvVars = os.Environ()

	for k, v := range r.Header {
		httpBody.Headers[k] = strings.Join(v, ", ")
	}
	if httpBody.Headers["Host"] == "" {
		httpBody.Headers["Host"] = r.Host
	}
	s.JSONResponse(w, r, httpBody)
}
