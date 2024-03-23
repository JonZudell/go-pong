package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		f(w, r)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ctx := context.Background()

	log.Println("Initializing Server with address:", os.Getenv("SERVER_URL"))

	server := http.Server{
		Addr:         os.Getenv("SERVER_URL"),
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

	go func() {
		<-c
		log.Println("Caught SIGINT")
		ladder.Shutdown(ctx)
		server.Shutdown(ctx)
	}()

	log.Println("ListenAndServe")

	err := server.ListenAndServe()
	if err != nil {
		log.Println("ListenAndServe: ", err)
	}
	log.Println("Clean Shutdown")
}
