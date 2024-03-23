package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Ladder struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	log        [][]byte
}

func newLadder() *Ladder {
	return &Ladder{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (l *Ladder) run() {
	ticker := time.NewTicker(time.Second / 128)
	for {
		select {
		case client := <-l.register:
			log.Println("Registering Client")
			l.clients[client] = true
			for line := range l.log {
				client.send <- l.log[line]
			}
		case client := <-l.unregister:
			log.Println("Unregistering Client")
			if _, ok := l.clients[client]; ok {
				delete(l.clients, client)
				close(client.send)
			}
		case message := <-l.broadcast:
			log.Println("Broadcasting Message to all Clients")
			for client := range l.clients {
				select {
				case client.send <- message:
					l.log = append(l.log, message)
				default:
					close(client.send)
					delete(l.clients, client)
				}
			}
		case time := <-ticker.C:
			message := fmt.Sprintf("Tick %s", time)
			log.Println(message)
			for client := range l.clients {
				select {
				case client.send <- []byte(message):
				default:
					close(client.send)
					delete(l.clients, client)
				}
			}
		}
	}
}
func (l *Ladder) Shutdown(ctx context.Context) {
	log.Println("Shutdown Ladder")
	for client := range l.clients {
		client.ladder.unregister <- client
	}
}
