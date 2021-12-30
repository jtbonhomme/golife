package cell

import (
	"math"
	"math/rand"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	evector "github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jtbonhomme/golife/internal/vector"
	color "github.com/lucasb-eyer/go-colorful"
)

const (
	cellMaxVelocity float64 = 1.1
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
	vel := 18 / size
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
	cellColor := color.HSLuv(c.size*360/50, 1, 0.5)

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
}
