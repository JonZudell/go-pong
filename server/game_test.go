package server

import (
	"testing"
)

func TestNoCollision(t *testing.T) {
	game := &Game{
		Ball: &Ball{
			Position: Position{X: 250, Y: 200},
			Velocity: Velocity{VX: 0, VY: 0},
			Radius:   10,
		},
		PlayerA: &Paddle{
			Position: Position{X: 200, Y: 200},
			Velocity: Velocity{VX: 0, VY: 0},
			Width:    100,
			Height:   20,
		},
		PlayerB: &Paddle{
			Position: Position{X: 300, Y: 200},
			Velocity: Velocity{VX: 0, VY: 0},
			Width:    100,
			Height:   20,
		},
	}

	// Perform all types of collisions
	game.checkCollision()
	game.checkWallCollision()
	game.checkPaddleWallCollision()

	// Check if the ball is not colliding with any paddle
	if game.isBallCollidingWithPaddle(game.PlayerA) {
		t.Errorf("Ball is colliding with PlayerA paddle")
	}

	if game.isBallCollidingWithPaddle(game.PlayerB) {
		t.Errorf("Ball is colliding with PlayerB paddle")
	}
}
func TestPlayerACollision(t *testing.T) {
	game := &Game{
		Ball: &Ball{
			Position: Position{X: 250, Y: 200},
			Velocity: Velocity{VX: 0, VY: 0},
			Radius:   10,
		},
		PlayerA: &Paddle{
			Position: Position{X: 200, Y: 200},
			Velocity: Velocity{VX: 0, VY: 0},
			Width:    100,
			Height:   20,
		},
		PlayerB: &Paddle{
			Position: Position{X: 300, Y: 200},
			Velocity: Velocity{VX: 0, VY: 0},
			Width:    100,
			Height:   20,
		},
	}

	// Set the ball position to collide with PlayerA paddle
	game.Ball.Position.X = 200
	game.Ball.Position.Y = 200

	// Perform collision check
	game.checkCollision()

	// Check if the ball is colliding with PlayerA paddle
	if !game.isBallCollidingWithPaddle(game.PlayerA) {
		t.Errorf("Ball is not colliding with PlayerA paddle")
	}
}
func TestPlayerBCollision(t *testing.T) {
	game := &Game{
		Ball: &Ball{
			Position: Position{X: 250, Y: 200},
			Velocity: Velocity{VX: 0, VY: 0},
			Radius:   10,
		},
		PlayerA: &Paddle{
			Position: Position{X: 200, Y: 200},
			Velocity: Velocity{VX: 0, VY: 0},
			Width:    100,
			Height:   20,
		},
		PlayerB: &Paddle{
			Position: Position{X: 300, Y: 200},
			Velocity: Velocity{VX: 0, VY: 0},
			Width:    100,
			Height:   20,
		},
	}

	// Set the ball position to collide with PlayerB paddle
	game.Ball.Position.X = 300
	game.Ball.Position.Y = 200

	// Perform collision check
	game.checkCollision()

	// Check if the ball is colliding with PlayerB paddle
	if !game.isBallCollidingWithPaddle(game.PlayerB) {
		t.Errorf("Ball is not colliding with PlayerB paddle")
	}
}
