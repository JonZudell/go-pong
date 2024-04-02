package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"zudell.io/go-pong/server"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	ctx := context.Background()
	log.Println("Initializing Server with address:", os.Getenv("SERVER_URL"))
	server := server.NewServer(os.Getenv("SERVER_URL"))
	go func() {
		<-c
		log.Println("Caught SIGINT")
		server.Shutdown(ctx)
	}()

	log.Println("ListenAndServe")

	server.ListenAndServe()

	log.Println("Clean Shutdown")
}
