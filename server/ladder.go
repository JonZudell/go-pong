package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Ladder struct {
	conns      map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

func NewLadder() *Ladder {
	return &Ladder{
		conns:      make(map[*websocket.Conn]bool),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}
func (l *Ladder) run() {
	write := time.NewTicker(time.Second)

	for {
		select {
		case conn := <-l.register:
			log.Println("Registering Connection")
			go func() {
				for {
					conn.SetReadDeadline(time.Now().Add(time.Second * 60))
					_, _, err := conn.ReadMessage()
					if err != nil {
						log.Println("Failed to read message:", err)
						l.unregister <- conn
						conn.Close()
						break
					}
				}
			}()
			go func() {
				for range write.C {
					conn.SetWriteDeadline(time.Now().Add(time.Second * 60))
					log.Println("Sending Ping to", conn.RemoteAddr())
					err := conn.WriteMessage(websocket.PingMessage, []byte{})
					if err != nil {
						log.Println("Failed to write message:", err)
						l.unregister <- conn
						conn.Close()
						break
					}
				}
			}()
			l.conns[conn] = true
		case conn := <-l.unregister:
			log.Println("Unregistering Connection")
			delete(l.conns, conn)
		}
	}
}

func (l *Ladder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	conn.SetReadLimit(1024)
	conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	conn.SetPongHandler(func(string) error {
		log.Println("recieved pong from", conn.RemoteAddr())
		conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		return nil
	})
	if err != nil {
		log.Println("Failed to upgrade websocket:", err)
		return
	}
	l.register <- conn
}
