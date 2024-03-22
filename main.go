package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
    "github.com/gorilla/websocket"
    "github.com/gorilla/mux"
)
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<!-- websockets.html -->
	<input id="input" type="text" />
	<button onclick="connect()">Connect</button>
	<button onclick="send()">Send</button>
	<pre id="output"></pre>
	<script>
		var input = document.getElementById("input");
		var output = document.getElementById("output");
		var socket = null;
	
	
		function send() {
			socket.send(input.value);
			input.value = "";
		}

		function connect() {
			socket = new WebSocket("ws://localhost:3000/echo");
			socket.onopen = function () {
				output.innerHTML += "Status: Connected\n";
			};
		
			socket.onmessage = function (e) {
				output.innerHTML += "Server: " + e.data + "\n";
			};
		}
	</script>`)
}

func upgrade(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade error %s", err)
		return
	}
	for {
		// Read message from browser
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", ws.RemoteAddr(), string(msg))

		// Write message back to browser
		if err = ws.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}

func main() {
	ctx := context.Background()
	fmt.Println("Initializing Server with address:", os.Getenv("SERVER_URL") )
	server := http.Server{
		Addr:         os.Getenv("SERVER_URL"),
		Handler:      nil,
		ReadTimeout:  1000,
		WriteTimeout: 1000,
	}
	r := mux.NewRouter()
	server.Handler = r
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	r.HandleFunc("/", index)
	r.HandleFunc("/echo", upgrade)
	
	go func() {
		<-c
		fmt.Println("Caught SIGINT")
		server.Shutdown(ctx)
	}()
	fmt.Println("Listening and Serving")
	server.ListenAndServe()
	fmt.Println("Clean Shutdown")
}
