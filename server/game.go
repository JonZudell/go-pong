package server

import (
	"time"
)

type Ball struct {
	x      float64
	y      float64
	vx     float64
	vy     float64
	radius float64
}

type Paddle struct {
	x      float64
	y      float64
	vx     float64
	vy     float64
	width  float64
	height float64
}

type Game struct {
	clientA        *Client
	clientB        *Client
	playerA        *Paddle
	playerB        *Paddle
	ball           *Ball
	scoreA         int
	scoreB         int
	lastUpdateTime time.Time
}

func (g *Game) update() {
	// Compute deltaTime
	currentTime := time.Now()
	deltaTime := currentTime.Sub(g.lastUpdateTime).Seconds()
	g.lastUpdateTime = currentTime
	// Update paddle positions
	g.playerA.y += g.playerA.vy * deltaTime
	g.playerB.y += g.playerB.vy * deltaTime
	// Update ball position
	// AABB intersection test for paddles
	if g.ball.x-g.ball.radius <= g.playerA.x+g.playerA.width &&
		g.ball.x+g.ball.radius >= g.playerA.x &&
		g.ball.y-g.ball.radius <= g.playerA.y+g.playerA.height &&
		g.ball.y+g.ball.radius >= g.playerA.y {
		// Collision with player A paddle
		g.ball.vx = -g.ball.vx
		g.ball.x = g.playerA.x + g.playerA.width + g.ball.radius
	}

	if g.ball.x-g.ball.radius <= g.playerB.x+g.playerB.width &&
		g.ball.x+g.ball.radius >= g.playerB.x &&
		g.ball.y-g.ball.radius <= g.playerB.y+g.playerB.height &&
		g.ball.y+g.ball.radius >= g.playerB.y {
		// Collision with player B paddle
		g.ball.vx = -g.ball.vx
		g.ball.x = g.playerB.x + g.playerB.width - g.ball.radius
	}
	// AABB intersection test for ball
	if g.ball.x-g.ball.radius <= 0 ||
		g.ball.x+g.ball.radius >= 1000 {
		// Collision with left or right wall
		// Handle collision logic here
		if g.ball.x-g.ball.radius <= 0 {
			// Collision with left wall
			// Increment player B's score
			g.scoreB++
		} else if g.ball.x+g.ball.radius >= 1000 {
			// Collision with right wall
			// Increment player A's score
			g.scoreA++
		}
	}

	if g.ball.y-g.ball.radius <= 0 ||
		g.ball.y+g.ball.radius >= 750 {
		// Collision with top or bottom wall
		// Handle collision logic here
		// Collision with top wall
		if g.ball.y+g.ball.radius >= 750 {
			g.ball.vy = -g.ball.vy
			g.ball.y = g.ball.radius
		}
		// Collision with bottom wall
		if g.ball.y-g.ball.radius <= 0 {
			g.ball.vy = -g.ball.vy
			g.ball.y = 750 - g.ball.radius
		}
	}
}
