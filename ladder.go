package main

import (
	"context"
	"log"
)

type Ladder struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
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
	for {
		select {
		case client := <-l.register:
			log.Println("Registering Client")
			l.clients[client] = true
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
		close(client.send)
		delete(l.clients, client)
	}
}
