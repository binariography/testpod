package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	//logger
	//config
}

func NewServer() (*Server, error) {
	srv := &Server{
		router: mux.NewRouter(),
	}

	return srv, nil
}

func (s *Server) ListenAndServe() *http.Server {
	s.registerHandlers()
	srv := s.startServer()

	return srv
}

func (s *Server) registerHandlers() {
	s.router.HandleFunc("/info", s.InfoHandler)
	s.router.HandleFunc("/echo/{text}", s.EchoHandler)
}

func (s *Server) startServer() *http.Server {
	srv := &http.Server{
		Addr:    "localhost:9090",
		Handler: s.router,
	}
	log.Println("Sending to a function to start the server")

	go func() {
		log.Println("Starting HTTP Server....")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("HTTP server crashed", err)
		}
	}()
	log.Println("Moved on from the function that started the server")

	return srv
}
