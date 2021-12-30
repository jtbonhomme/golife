package cell

import (
	"math"

	"github.com/jtbonhomme/golife/internal/vector"
)

// Accelerate set physical body acceleration.
func (c *Cell) Accelerate(acceleration vector.Vector2D) {
	c.acceleration = acceleration
}

func (c *Cell) Update(counter int) {
	c.neighbors = c.detect(c.position, 250)

	if counter > c.lastEnergyBurn+150 {
		c.energy -= 2
		c.lastEnergyBurn = counter
	}

	if counter > c.lastGrowth+1000 {
		c.energy -= 5
		c.size += 5
		c.lastGrowth = counter
	}

	if c.energy <= 0 {
		c.Kill()
		return
	}
	predators := []vector.Vector2D{}
	preyDistance := c.detectionRadius
	preyPosition := vector.Vector2D{}

	for _, c1 := range c.neighbors {
		// Don't compare to myself
		if c1.ID() == c.ID() {
			continue
		}
		// Eat smaller cells in the neighborood
		if c.Intersect(c1) && c.Size() > c1.Size()*1.1 {
			c.Eat(c1)
		}
		dist := c.Position().Distance(c1.Position())

		// find the nearest prey in neighborood
		if c.Size() > c1.Size() && dist < preyDistance {
			preyDistance = dist
			preyPosition = c1.Position()
		}
		// record all the predators in neighborood
		if c1.Size() > c.Size() {
			predators = append(predators, c1.Position())
		}
	}

	acceleration := vector.Vector2D{
		X: math.Cos(c.orientation),
		Y: math.Sin(c.orientation),
	}

	if len(predators) > 0 {
		// if there is a predator in the neighborood, flee !
		flee := c.avoid(predators)
		acceleration.Add(flee)
	} else if preyDistance < c.detectionRadius {
		// else pursuit prey
		chase := vector.Vector2D{
			X: 0,
			Y: 0,
		}
		d := c.Position().Distance(preyPosition)
		diff := c.Position()
		diff.Subtract(preyPosition)
		diff.Normalize()
		diff.Divide(d)
		chase.Add(diff)

		acceleration.Add(chase)
	}
	// else continue in the same direction

	c.Accelerate(acceleration)
	c.UpdateVelocity()
	c.UpdateOrientation()
	c.UpdatePosition()
}

func (c *Cell) avoid(predators []vector.Vector2D) vector.Vector2D {
	result := vector.Vector2D{
		X: 0,
		Y: 0,
	}
	cells := 0.0
	for _, p := range predators {
		cells++
		d := c.Position().Distance(p)
		diff := c.Position()
		diff.Subtract(p)
		diff.Normalize()
		diff.Divide(d)
		result.Add(diff)
	}
	if cells > 0 {
		result.Divide(cells)
		result.Normalize()
		result.Multiply(cellMaxVelocity)
		result.Subtract(c.Velocity())
		result.Limit(cellMaxForce)
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

	if c.position.X >= c.screenWidth-1.0 {
		c.position.X = 0
	} else if c.position.X < 0 {
		c.position.X = c.screenWidth - 1.0
	}
	if c.position.Y >= c.screenHeight-1.0 {
		c.position.Y = 0
	} else if c.position.Y < 0 {
		c.position.Y = c.screenHeight - 1.0
	}
}
