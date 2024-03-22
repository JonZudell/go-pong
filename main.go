package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func handleFunction(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprint(w, "Hello World")
	case "/other":
		fmt.Fprint(w, "Hello other")
	default:
		fmt.Fprint(w, "No")
	}
}
func main() {
	ctx := context.Background()
	fmt.Println("Initializing Server")
	server := http.Server{
		Addr:         "",
		Handler:      nil,
		ReadTimeout:  1000,
		WriteTimeout: 1000,
	}
	var mux http.ServeMux
	server.Handler = &mux
	mux.HandleFunc("/", handleFunction)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Caught SIGINT")
		server.Shutdown(ctx)
	}()
	fmt.Println("Listening and Serving")
	server.ListenAndServe()
	fmt.Println("Clean Shutdown")
}
