package http

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (s *Server) JSONResponse(w http.ResponseWriter, r *http.Request, value interface{}) {
	body, err := json.Marshal(value)
	if err != nil {
		s.logger.Error("faild to Json", "error", err)
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(JSONIndent(body))
}

func JSONIndent(data []byte) []byte {
	var out bytes.Buffer
	// Learn what are the args
	json.Indent(&out, data, "", "  ")

	return out.Bytes()
}
