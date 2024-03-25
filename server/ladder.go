package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Ladder struct {
	clients    map[*Client]bool
	games      map[*Game]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	log        [][]byte
}
type Pair struct {
	First, Second interface{}
}

func PairList(list []*Client) []Pair {
	var pairedList []Pair
	for i := 0; i < len(list)-1; i += 2 {
		pairedList = append(pairedList, Pair{list[i], list[i+1]})
	}
	return pairedList
}
func NewLadder() *Ladder {
	return &Ladder{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		games:      make(map[*Game]bool),
	}
}
func (l *Ladder) unregisterClients(client *Client) {
	log.Println("Unregistering Client")
	if _, ok := l.clients[client]; ok {
		delete(l.clients, client)
		close(client.send)
	}
}
func (l *Ladder) run() {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case client := <-l.register:
			log.Println("Registering Client")
			l.clients[client] = true
			for line := range l.log {
				client.send <- l.log[line]
			}
		case client := <-l.unregister:
			l.unregisterClients(client)
		case <-ticker.C:
			l.ladderTick()
		}
	}
}

func (l *Ladder) ladderTick() {
	clientsInGame := make([]*Client, 0)
	unpairedClients := make([]*Client, 0)
	for game := range l.games {
		clientsInGame = append(clientsInGame, game.clientA, game.clientB)
	}
	for client := range l.clients {
		if !contains(clientsInGame, client) {
			unpairedClients = append(unpairedClients, client)
		}
	}

	log.Println("Games:", len(l.games))
	log.Println("Clients in Game:", len(clientsInGame))
	log.Println("Clients not in game:", len(unpairedClients))
	paired := PairList(unpairedClients)
	for _, pair := range paired {
		game := &Game{
			clientA:        pair.First.(*Client),
			clientB:        pair.Second.(*Client),
			ladder:         l,
			PlayerA:        &Paddle{X: 37.5, Y: 325, Width: 25, Height: 100},
			PlayerB:        &Paddle{X: 962.5, Y: 325, Width: 25, Height: 100},
			Ball:           &Ball{X: 500, Y: 375, VX: 10, VY: 0, Radius: 10},
			ScoreA:         0,
			ScoreB:         0,
			lastUpdateTime: time.Time{},
		}
		pair.First.(*Client).game = game
		pair.Second.(*Client).game = game
		l.games[game] = true
		go game.run()
	}
}

func contains(clientsInGame []*Client, client *Client) bool {
	for _, c := range clientsInGame {
		if c == client {
			return true
		}
	}
	return false
}

func (l *Ladder) Shutdown(ctx context.Context) {
	log.Println("Shutdown Ladder")
	for client := range l.clients {
		client.ladder.unregister <- client
	}
}

func (l *Ladder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{ladder: l, conn: conn, send: make(chan []byte, 256)}
	client.ladder.register <- client

	go client.writePump()
	go client.readPump()
}

func (l *Ladder) RemoveGame(g *Game) {
	// implement the method
	delete(l.games, g)
}
