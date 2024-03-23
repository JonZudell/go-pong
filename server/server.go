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
	ladder := newLadder()
	go ladder.run()
	r := mux.NewRouter()
	r.HandleFunc("/", logging(index))
	r.HandleFunc("/upgrade", logging(func(w http.ResponseWriter, r *http.Request) {
		serveWs(ladder, w, r)
	}))
	server.Handler = r
	return &Server{server: server, ladder: ladder}
}

func (s *Server) Shutdown(ctx context.Context) {
	s.server.Shutdown(ctx)
	s.ladder.Shutdown(ctx)
}

func (s *Server) ListenAndServe() {
	go s.ladder.run()
	err := s.server.ListenAndServe()
	if err != nil {
		log.Println("ListenAndServe: ", err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.server.Handler.ServeHTTP(w, req)
}

//go:embed static/index.html
var s string

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(s))
}

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		f(w, r)
	}
}
