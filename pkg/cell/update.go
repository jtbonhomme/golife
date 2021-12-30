package cell

import (
	"math"

	"github.com/jtbonhomme/golife/internal/vector"
)

// Accelerate set physical body acceleration.
func (c *Cell) Accelerate(acceleration vector.Vector2D) {
	c.acceleration = acceleration
}

func (c *Cell) Update() {
	acceleration := vector.Vector2D{}
	acceleration.Add(c.move())

	c.Accelerate(acceleration)
	c.UpdateVelocity()
	c.UpdateOrientation()
	c.UpdatePosition()
}

func (c *Cell) move() vector.Vector2D {
	result := vector.Vector2D{
		X: 1,
		Y: 1,
	}

	return result
}

// UpdateVelocity computes new velocity.
func (c *Cell) UpdateVelocity() {
	// update velocity from acceleration
	c.velocity.Add(c.acceleration)

	// limit velocity to max value
	c.velocity.Limit(c.maxVelocity)
}

func (c *Cell) normalizeOrientation() {
	if c.orientation > 2*math.Pi {
		c.orientation -= 2 * math.Pi
	}
	if c.orientation < 0 {
		c.orientation += 2 * math.Pi
	}
}

// UpdateOrientation computes orientation from velocity.
func (c *Cell) UpdateOrientation() {
	// Update orientation from velocity
	if !c.velocity.IsNil() {
		c.orientation = c.velocity.Theta()
	}
	c.normalizeOrientation()
}

// UpdatePosition compute new position.
func (c *Cell) UpdatePosition() {
	c.position.Add(c.velocity)

	if c.position.X > c.ScreenWidth {
		c.position.X = 0
	} else if c.position.X < 0 {
		c.position.X = c.ScreenWidth
	}
	if c.position.Y > c.ScreenHeight {
		c.position.Y = 0
	} else if c.position.Y < 0 {
		c.position.Y = c.ScreenHeight
	}
}
