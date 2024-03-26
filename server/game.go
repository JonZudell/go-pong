package server

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Ball struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	VX     float64 `json:"vx"`
	VY     float64 `json:"vy"`
	Radius float64 `json:"radius"`
}
type Paddle struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	VX     float64 `json:"vx"`
	VY     float64 `json:"vy"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Game struct {
	clientA        *Client
	clientB        *Client
	ladder         *Ladder
	PlayerA        *Paddle
	PlayerB        *Paddle
	Ball           *Ball
	ScoreA         int
	ScoreB         int
	lastUpdateTime time.Time
	Started        bool
}

// Encode game as JSON
type GameStateMessage struct {
	Type string `json:"type"`
	Game *Game  `json:"game"`
}

func (g *Game) update() {
	// Compute deltaTime
	if !g.Started {
		g.lastUpdateTime = time.Now()
		g.Started = true
	}
	currentTime := time.Now()

	deltaTime := currentTime.Sub(g.lastUpdateTime).Seconds()
	if deltaTime > float64(time.Second)/64 {
		deltaTime = float64(time.Second) / 64
	}

	g.Ball.X += g.Ball.VX * deltaTime
	g.Ball.Y += g.Ball.VY * deltaTime
	g.PlayerA.Y += g.PlayerA.VY * deltaTime
	g.PlayerB.Y += g.PlayerB.VY * deltaTime
	// Update ball position
	// AABB intersection test for paddles
	if g.Ball.X-g.Ball.Radius <= g.PlayerA.X+g.PlayerA.Width &&
		g.Ball.X+g.Ball.Radius >= g.PlayerA.X &&
		g.Ball.Y-g.Ball.Radius <= g.PlayerA.Y+g.PlayerA.Height &&
		g.Ball.Y+g.Ball.Radius >= g.PlayerA.Y {
		// Collision with player A paddle
		g.Ball.VX = -g.Ball.VX
		g.Ball.X = g.PlayerA.X + g.PlayerA.Width + (g.Ball.Radius + 1)
	}

	if g.Ball.X-g.Ball.Radius <= g.PlayerB.X+g.PlayerB.Width &&
		g.Ball.X+g.Ball.Radius >= g.PlayerB.X &&
		g.Ball.Y-g.Ball.Radius <= g.PlayerB.Y+g.PlayerB.Height &&
		g.Ball.Y+g.Ball.Radius >= g.PlayerB.Y {
		// Collision with player B paddle
		g.Ball.VX = -g.Ball.VX
		g.Ball.X = g.PlayerB.X - (g.Ball.Radius + 1)
	}
	// AABB intersection test for ball
	if g.Ball.X-g.Ball.Radius <= 0 ||
		g.Ball.X+g.Ball.Radius >= 1000 {
		// Collision with left or right wall
		// Handle collision logic here
		if g.Ball.X-g.Ball.Radius <= 0 {
			// Collision with left wall
			// Increment player B's score
			g.ScoreB++
			g.Ball.X = 500
			g.Ball.Y = 375
		} else if g.Ball.X+g.Ball.Radius >= 1000 {
			// Collision with right wall
			// Increment player A's score
			g.ScoreA++
			g.Ball.X = 500
			g.Ball.Y = 375
		}
	}

	if g.Ball.Y-g.Ball.Radius <= 0 ||
		g.Ball.Y+g.Ball.Radius >= 750 {
		// Collision with top or bottom wall
		// Handle collision logic here
		// Collision with top wall
		if g.Ball.Y+g.Ball.Radius >= 750 {
			g.Ball.VY = -g.Ball.VY
			g.Ball.Y = g.Ball.Radius
		}
		// Collision with bottom wall
		if g.Ball.Y-g.Ball.Radius <= 0 {
			g.Ball.VY = -g.Ball.VY
			g.Ball.Y = 750 - g.Ball.Radius
		}
	}
	gameJSON, err := json.Marshal(GameStateMessage{Type: "gamestate", Game: g})
	if err != nil {
		log.Printf("Error encoding game state: %v", err)
	}

	// Send game JSON to client A
	err = g.clientA.conn.WriteMessage(websocket.TextMessage, gameJSON)
	if err != nil {
		log.Printf("Error sending game state to client A: %v", err)
		g.ladder.unregister <- g.clientA
		panic("Couldn't talk to clientA")
	}

	// Send game JSON to client B
	err = g.clientB.conn.WriteMessage(websocket.TextMessage, gameJSON)
	if err != nil {
		log.Printf("Error sending game state to client B: %v", err)
		g.ladder.unregister <- g.clientB
		panic("Couldn't talk to clientB")
	}
	g.lastUpdateTime = time.Now()
}
func (g *Game) run() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Killing game in a Panic", r)
			g.ladder.RemoveGame(g)
		}
	}()

	// Your update logic here
	for range time.Tick(time.Second / 64) {
		g.update()
	}
}
