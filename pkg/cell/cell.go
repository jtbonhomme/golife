package cell

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	evector "github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jtbonhomme/golife/internal/fonts"
	"github.com/jtbonhomme/golife/internal/vector"
	colorful "github.com/lucasb-eyer/go-colorful"
	log "github.com/sirupsen/logrus"
)

const (
	cellMaxVelocity float64 = 0.9
)

type Cell struct {
	size        float64
	rnd10       int32
	id          uuid.UUID
	orientation float64 // theta (radian)

	position     vector.Vector2D
	velocity     vector.Vector2D
	maxVelocity  float64
	acceleration vector.Vector2D

	screenWidth  float64
	screenHeight float64

	debug bool
}

func maxVelocity(size float64) float64 {
	vel := 15 / size
	if vel > cellMaxVelocity {
		vel = cellMaxVelocity
	}
	return vel
}

func New(position vector.Vector2D, w, h int) *Cell {
	size := 5.0 + rand.Float64()*45.0
	c := &Cell{
		position:     position,
		orientation:  rand.Float64() * 2 * math.Pi,
		size:         size,
		rnd10:        rand.Int31n(10),
		id:           uuid.New(),
		screenWidth:  float64(w),
		screenHeight: float64(h),
		maxVelocity:  maxVelocity(size),
	}
	return c
}

func (c *Cell) Debug(state bool) {
	c.debug = state
}

func maxCounter(index, rnd10 int) int {
	return 50 + rnd10 + (25*index+rnd10)%64
}

func angularToCartesian(dist, orientation float64) (x, y float32) {
	return float32(dist * math.Cos(orientation)), float32(dist * math.Sin(orientation))
}

func addVector(position vector.Vector2D, dist, orientation float64) (x, y float32) {
	acX, acY := angularToCartesian(dist, orientation)
	return float32(position.X) + acX, float32(position.Y) + acY
}

func (c *Cell) drawCellBody(screen *ebiten.Image, counter int, emptySubImage *ebiten.Image) {
	var path evector.Path
	npoints := 16

	indexToDirection := func(i int) float64 {
		return c.orientation - float64(2*i+1)*math.Pi/float64(npoints)
	}
	indexToDist := func(i, counter int) float64 {
		return c.size + c.size*0.1*math.Sin(float64(counter)*2*math.Pi/float64(maxCounter(i, int(c.rnd10))))
	}

	for i := 0; i <= npoints; i++ {
		if i == 0 {
			path.MoveTo(addVector(c.position, indexToDist(i, counter), indexToDirection(i)))
			continue
		}
		cpx0, cpy0 := addVector(c.position, indexToDist(i, counter), indexToDirection(i-1)-math.Pi/16)
		cpx1, cpy1 := addVector(c.position, indexToDist(i, counter), indexToDirection(i)+math.Pi/16)
		cpx2, cpy2 := addVector(c.position, indexToDist(i, counter), indexToDirection(i))
		path.CubicTo(cpx0, cpy0, cpx1, cpy1, cpx2, cpy2)
	}

	// Get the color (120° is green, 0° is red)
	cellColor := colorful.HSLuv(c.size*360/50, 1, 0.5)

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(cellColor.R)
		vs[i].ColorG = float32(cellColor.G)
		vs[i].ColorB = float32(cellColor.B)
	}
	screen.DrawTriangles(vs, is, emptySubImage, op)
}

func (c *Cell) drawEyes(screen *ebiten.Image, dist, side, size, bg float64, emptySubImage *ebiten.Image) {
	var path evector.Path

	randomizedFloat64 := func(in float64) float64 {
		return in + rand.Float64()*2
	}

	cpx0, cpy0 := addVector(c.position, dist-randomizedFloat64(size), c.orientation+side*math.Pi/randomizedFloat64(12))

	path.Arc(cpx0, cpy0, float32(size), float32(0), float32(2*math.Pi), evector.Clockwise)

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(bg) / float32(0xff)
		vs[i].ColorG = float32(bg) / float32(0xff)
		vs[i].ColorB = float32(bg) / float32(0xff)
	}
	screen.DrawTriangles(vs, is, emptySubImage, op)
}

func (c *Cell) Draw(screen *ebiten.Image, counter int, emptySubImage *ebiten.Image) {
	c.drawCellBody(screen, counter, emptySubImage)
	c.drawEyes(screen, c.size*0.9, -1, c.size*0.1, 0xff, emptySubImage)
	c.drawEyes(screen, c.size*0.9, -1, c.size*0.05, 0x00, emptySubImage)
	c.drawEyes(screen, c.size*0.9, 1, c.size*0.1, 0xff, emptySubImage)
	c.drawEyes(screen, c.size*0.9, 1, c.size*0.05, 0x00, emptySubImage)
	if c.debug {
		c.DrawBodyBoundaryBox(screen)
		msg := c.String()
		textDim := text.BoundString(fonts.MonoSansRegularFont, msg)
		textWidth := textDim.Max.X - textDim.Min.X
		text.Draw(screen,
			msg,
			fonts.MonoSansRegularFont,
			int(c.position.X)-textWidth/2,
			int(c.position.Y+c.size+5),
			color.Gray16{0x999f})
	}
}

// DrawBodyBoundaryBox draws a box around the body, based on its dimension.
func (c *Cell) DrawBodyBoundaryBox(screen *ebiten.Image) {
	// Top boundary
	ebitenutil.DrawLine(
		screen,
		c.position.X-c.size,
		c.position.Y-c.size,
		c.position.X+c.size,
		c.position.Y-c.size,
		color.Gray16{0xbbbb},
	)
	// Right boundary
	ebitenutil.DrawLine(
		screen,
		c.position.X+c.size,
		c.position.Y-c.size,
		c.position.X+c.size,
		c.position.Y+c.size,
		color.Gray16{0xbbbb},
	)
	// Bottom boundary
	ebitenutil.DrawLine(
		screen,
		c.position.X-c.size,
		c.position.Y+c.size,
		c.position.X+c.size,
		c.position.Y+c.size,
		color.Gray16{0xbbbb},
	)
	// Left boundary
	ebitenutil.DrawLine(
		screen,
		c.position.X-c.size,
		c.position.Y-c.size,
		c.position.X-c.size,
		c.position.Y+c.size,
		color.Gray16{0xbbbb},
	)
}

// String displays physic body information as a string.
func (c *Cell) String() string {
	return fmt.Sprintf("pos [%d, %d]\nsize [%d] orient %0.2f rad (%0.0f °)\nvel {%0.2f %0.2f} acc {%0.2f %0.2f}",
		int(c.position.X),
		int(c.position.Y),
		int(c.size),
		c.orientation,
		c.orientation*180/math.Pi,
		c.velocity.X,
		c.velocity.Y,
		c.acceleration.X,
		c.acceleration.Y)
}

// ID displays physic body unique ID.
func (c *Cell) ID() string {
	return c.id.String()
}

// Intersect returns true if the physical body collide another one.
// Collision is computed based on Axis-Aligned Bounding Boxes.
// https://developer.mozilla.org/en-US/docs/Games/Techniques/2D_collision_detection
func (c *Cell) Intersect(c2 *Cell) bool {
	ax, ay := c.position.X, c.position.Y
	aw, ah := c.size, c.size

	bx, by := c2.position.X, c2.position.Y
	bw, bh := c2.size, c2.size

	return (ax < bx+bw && ay < by+bh) && (ax+aw > bx && ay+ah > by)
}

// IntersectMultiple checks if multiple physical bodies are colliding with the first
func (c *Cell) IntersectMultiple(cells map[string]*Cell) (string, bool) {
	for _, c2 := range cells {
		if c.Intersect(c2) {
			log.Warnf("%s [%d , %d] (%dx%d) intersect with %s [%d , %d] (%dx%d)",
				c.ID(),
				int(c.position.X), int(c.position.Y),
				int(c.size), int(c.size),
				c2.ID(),
				int(c2.position.X), int(c2.position.Y),
				int(c2.size), int(c2.size))
			return c2.ID(), true
		}
	}
	return "", false
}

// Position returns physical body position.
func (c *Cell) Position() vector.Vector2D {
	return c.position
}
