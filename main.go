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

// func upgrade(w http.ResponseWriter, r *http.Request) {
// 	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
// 	// upgrade this connection to a WebSocket connection
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("upgrade error %s", err)
// 		return
// 	}
// 	for {
// 		// Read message from browser
// 		msgType, msg, err := ws.ReadMessage()
// 		if err != nil {
// 			return
// 		}

// 		// Print the message to the console
// 		fmt.Printf("%s sent: %s\n", ws.RemoteAddr(), string(msg))

// 		// Write message back to browser
// 		if err = ws.WriteMessage(msgType, msg); err != nil {
// 			return
// 		}
// 	}
// }

func main() {
	c := make(chan os.Signal)
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
		server.Shutdown(ctx)
	}()

	log.Println("ListenAndServe")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Println("Clean Shutdown")
}
