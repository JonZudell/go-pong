package server

import (
	"encoding/json"
	"log"
	"math"
	"time"
)

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Velocity struct {
	VX float64 `json:"vx"`
	VY float64 `json:"vy"`
}

type Ball struct {
	Position
	Velocity
	Radius float64 `json:"radius"`
}

type Paddle struct {
	Position
	Velocity
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Game struct {
	clientA        *Client
	clientB        *Client
	ladder         *Ladder
	PlayerA        *Paddle
	PlayerB        *Paddle
	PlayerAHasBall bool
	PlayerBHasBall bool
	Ball           *Ball
	ScoreA         int `json:"scoreA"`
	ScoreB         int `json:"scoreB"`
	lastUpdateTime time.Time
	Started        bool
	Paused         bool
}

// Encode game as JSON
type GameStateMessage struct {
	Type string `json:"type"`
	Game *Game  `json:"game"`
}

func (g *Game) checkCollision() {
	g.checkPaddleCollision(g.PlayerA)
	g.checkPaddleCollision(g.PlayerB)
	g.checkWallCollision()
}

func (g *Game) checkPaddleCollision(paddle *Paddle) {
	if g.isBallCollidingWithPaddle(paddle) {
		g.resolvePaddleCollision(paddle)
	}
}

func (g *Game) isBallCollidingWithPaddle(paddle *Paddle) bool {
	return g.Ball.X-g.Ball.Radius <= paddle.X+paddle.Width &&
		g.Ball.X+g.Ball.Radius >= paddle.X &&
		g.Ball.Y-g.Ball.Radius <= paddle.Y+paddle.Height &&
		g.Ball.Y+g.Ball.Radius >= paddle.Y
}

func (g *Game) resolvePaddleCollision(paddle *Paddle) {
	// Determine which edge was intersected
	if g.Ball.Y+g.Ball.Radius >= paddle.Y+paddle.Height {
		g.Ball.VY, paddle.VY = paddle.VY, g.Ball.VY
		g.Ball.Y = paddle.Y + paddle.Height + g.Ball.Radius + 1
	} else if g.Ball.Y-g.Ball.Radius <= paddle.Y {
		g.Ball.VY, paddle.VY = paddle.VY, g.Ball.VY
		g.Ball.Y = paddle.Y - g.Ball.Radius - 1
	} else if g.Ball.X+g.Ball.Radius >= paddle.X+paddle.Width {
		g.Ball.VX = -g.Ball.VX
		g.Ball.X = paddle.X + paddle.Width + g.Ball.Radius + 1
	} else if g.Ball.X-g.Ball.Radius <= paddle.X {
		g.Ball.VX = -g.Ball.VX
		g.Ball.X = paddle.X - g.Ball.Radius - 1
	}
}

func (g *Game) checkWallCollision() {
	if g.isBallCollidingWithTopOrBottomWall() {
		g.resolveTopOrBottomWallCollision()
	}
	if g.isBallCollidingWithLeftOrRightWall() {
		g.resolveLeftOrRightWallCollision()
	}
}

func (g *Game) isBallCollidingWithTopOrBottomWall() bool {
	return g.Ball.Y-g.Ball.Radius <= 0 || g.Ball.Y+g.Ball.Radius >= 750
}

func (g *Game) resolveTopOrBottomWallCollision() {
	if g.Ball.Y-g.Ball.Radius <= 0 {
		g.Ball.VY = math.Abs(g.Ball.VY)
	} else {
		g.Ball.VY = -math.Abs(g.Ball.VY)
	}
}

func (g *Game) isBallCollidingWithLeftOrRightWall() bool {
	return g.Ball.X-g.Ball.Radius <= 0 || g.Ball.X+g.Ball.Radius >= 1000
}

func (g *Game) resolveLeftOrRightWallCollision() {
	if g.Ball.X-g.Ball.Radius <= 0 {
		g.ScoreB++
	} else {
		g.ScoreA++
	}
	g.Ball.VX = -g.Ball.VX
	g.Ball.X = 500
	g.Ball.Y = 375
}
func (g *Game) handleInput(deltaTime float64) {
	if g.clientA.up && g.clientA.down {
	} else if g.clientA.up {
		g.PlayerA.VY = g.PlayerA.VY - (10000 * deltaTime)
	} else if g.clientA.down {
		g.PlayerA.VY = g.PlayerA.VY + (10000 * deltaTime)
	}

	if g.PlayerA.VY > 100 {
		g.PlayerA.VY = 100
	} else if g.PlayerA.VY < -100 {
		g.PlayerA.VY = -100
	}
	if g.clientB.up && g.clientB.down {

	} else if g.clientB.up {
		g.PlayerB.VY = g.PlayerB.VY - (10000 * deltaTime)
	} else if g.clientB.down {
		g.PlayerB.VY = g.PlayerB.VY + (10000 * deltaTime)
	}
	if g.PlayerB.VY > 100 {
		g.PlayerB.VY = 100
	} else if g.PlayerB.VY < -100 {
		g.PlayerB.VY = -100
	}
	// Apply friction to player A's VY
	if g.PlayerA.VY > 0 {
		g.PlayerA.VY -= 100 * deltaTime
		if g.PlayerA.VY < 0 {
			g.PlayerA.VY = 0
		}
	} else if g.PlayerA.VY < 0 {
		g.PlayerA.VY += 100 * deltaTime
		if g.PlayerA.VY > 0 {
			g.PlayerA.VY = 0
		}
	}

	g.PlayerA.Y += g.PlayerA.VY * deltaTime
	g.PlayerB.Y += g.PlayerB.VY * deltaTime
	// Apply friction to player B's VY
	if g.PlayerB.VY > 0 {
		g.PlayerB.VY -= 100 * deltaTime
		if g.PlayerB.VY < 0 {
			g.PlayerB.VY = 0
		}
	} else if g.PlayerB.VY < 0 {
		g.PlayerB.VY += 100 * deltaTime
		if g.PlayerB.VY > 0 {
			g.PlayerB.VY = 0
		}
	}

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
	if !g.Paused {
		g.Ball.X += g.Ball.VX * deltaTime
		g.Ball.Y += g.Ball.VY * deltaTime
		g.handleInput(deltaTime)
		g.checkCollision()
	}

	gameJSON, err := json.Marshal(GameStateMessage{Type: "gamestate", Game: g})
	if err != nil {
		log.Printf("Error encoding game state: %v", err)
	}

	g.clientA.send <- gameJSON

	g.clientB.send <- gameJSON

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
