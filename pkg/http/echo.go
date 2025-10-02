package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) EchoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "echo: %v\n", vars["text"])
}
