package server

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
	clientA *Client
	clientB *Client
	playerA *Paddle
	playerB *Paddle
	ball    *Ball
	scoreA  int
	scoreB  int
}

func (g *Game) update() {
	// Update the game state
}
