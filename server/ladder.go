package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Ladder struct {
	clients        map[*Client]bool
	games          map[*Game]bool
	broadcast      chan []byte
	register       chan *Client
	unregister     chan *Client
	gameRegister   chan *Game
	gameUnregister chan *Game
	log            [][]byte
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
		clients:        make(map[*Client]bool),
		games:          make(map[*Game]bool),
		broadcast:      make(chan []byte),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		gameRegister:   make(chan *Game),
		gameUnregister: make(chan *Game),
		log:            [][]byte{},
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
			log.Println("Unregistering Client")
			if _, ok := l.clients[client]; ok {
				if !client.closed {
					close(client.send)
				}
				delete(l.clients, client)
			}
		case game := <-l.gameRegister:
			log.Println("Registering Game")
			l.games[game] = true
		case game := <-l.gameUnregister:
			log.Println("Unregistering Game")
			if _, ok := l.games[game]; ok {
				delete(l.games, game)
			}
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
		if !contains(clientsInGame, client) && client.ready {
			unpairedClients = append(unpairedClients, client)
		}
	}
	paired := PairList(unpairedClients)
	for _, pair := range paired {

		game := &Game{
			clientA:        pair.First.(*Client),
			clientB:        pair.Second.(*Client),
			ladder:         l,
			PlayerA:        &Paddle{Position: Position{X: 37.5, Y: 325}, Velocity: Velocity{VX: 0, VY: 0}, Width: 25, Height: 100},
			PlayerB:        &Paddle{Position: Position{X: 962.5 - 25, Y: 325}, Velocity: Velocity{VX: 0, VY: 0}, Width: 25, Height: 100},
			Ball:           &Ball{Position: Position{X: 500, Y: 375}, Velocity: Velocity{VX: 400, VY: 0}, Radius: 10},
			ScoreA:         0,
			ScoreB:         0,
			lastUpdateTime: time.Time{},
			Started:        false,
			Paused:         false,
			Ended:          false,
		}
		pair.First.(*Client).game = game
		pair.Second.(*Client).game = game
		l.games[game] = true
		go game.run()

	}
	log.Printf("Number of clients: %d", len(l.clients))
	log.Printf("Number of games: %d", len(l.games))
	log.Printf("Number of clients in game: %d", len(clientsInGame))
	log.Printf("Number of unpaired clients: %d", len(unpairedClients))
	log.Printf("Number of paired clients: %d", len(paired))
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade websocket:", err)
		return
	}
	client := &Client{ladder: l, conn: conn, send: make(chan []byte, 256), ready: false, closed: false, Name: ""}
	client.ladder.register <- client

	go client.writePump()
	go client.readPump()
}

func (l *Ladder) RemoveGame(g *Game) {
	// implement the method
	delete(l.games, g)
}
