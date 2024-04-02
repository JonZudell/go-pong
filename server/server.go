package server

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
	ladder *Ladder
}

func NewServer(addr string) *Server {
	server := &http.Server{
		Addr:         addr,
		Handler:      nil,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}
	ladder := NewLadder()
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.HandleFunc("/readiness/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	})
	r.HandleFunc("/upgrade/", logging(func(w http.ResponseWriter, r *http.Request) {
		ladder.ServeHTTP(w, r)
	}))
	server.Handler = r
	return &Server{server: server, ladder: ladder}
}

func (s *Server) Shutdown(ctx context.Context) {
	s.server.Shutdown(ctx)
	s.ladder.Shutdown(ctx)
}

func (s *Server) ListenAndServeTLS() {
	go s.ladder.run()
	err := s.server.ListenAndServeTLS("/etc/secrets/tls.crt", "/etc/secrets/tls.key")
	if err != nil {
		log.Println("ListenAndServe: ", err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	s.server.Handler.ServeHTTP(w, req)
}

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		f(w, r)
	}
}
