package cell

import (
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jtbonhomme/golife/internal/point"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func maxCounter(index, rnd10 int) int {
	return 50 + rnd10 + (25*index+rnd10)%64
}

func angularToCartesian(dist, direction float64) (x, y float32) {
	return float32(dist * math.Cos(direction)), float32(dist * math.Sin(direction))
}

func addVector(center point.Point, dist, direction float64) (x, y float32) {
	acX, acY := angularToCartesian(dist, direction)
	return float32(center.X) + acX, float32(center.Y) + acY
}

func (c *Cell) drawCellBody(screen *ebiten.Image, center point.Point, direction, size float64, counter int, emptySubImage *ebiten.Image) {
	var path vector.Path
	npoints := 16

	indexToDirection := func(i int) float64 {
		return direction - float64(2*i+1)*math.Pi/float64(npoints)
	}
	indexToDist := func(i, counter int) float64 {
		return size + size*0.1*math.Sin(float64(counter)*2*math.Pi/float64(maxCounter(i, int(c.rnd10))))
	}

	for i := 0; i <= npoints; i++ {
		if i == 0 {
			path.MoveTo(addVector(center, indexToDist(i, counter), indexToDirection(i)))
			continue
		}
		cpx0, cpy0 := addVector(center, indexToDist(i, counter), indexToDirection(i-1)-math.Pi/16)
		cpx1, cpy1 := addVector(center, indexToDist(i, counter), indexToDirection(i)+math.Pi/16)
		cpx2, cpy2 := addVector(center, indexToDist(i, counter), indexToDirection(i))
		path.CubicTo(cpx0, cpy0, cpx1, cpy1, cpx2, cpy2)
	}

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = 0x22 / float32(0xff)
		vs[i].ColorG = 0x33 / float32(0xff)
		vs[i].ColorB = 0x66 / float32(0xff)
	}
	screen.DrawTriangles(vs, is, emptySubImage, op)
}

func (c *Cell) drawEyes(screen *ebiten.Image, center point.Point, direction, dist, side, size, bg float64, counter int, emptySubImage *ebiten.Image) {
	var path vector.Path

	randomizedFloat64 := func(in float64) float64 {
		return in + rand.Float64()*2
	}

	cpx0, cpy0 := addVector(center, dist-randomizedFloat64(size), direction+side*math.Pi/randomizedFloat64(12))

	path.Arc(cpx0, cpy0, float32(size), float32(0), float32(2*math.Pi), vector.Clockwise)

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

type Cell struct {
	center    point.Point
	direction float64
	size      float64
	rnd10     int32
}

func NewCell(center point.Point) *Cell {
	c := &Cell{
		center:    center,
		direction: rand.Float64() * 2 * math.Pi,
		size:      20.0 + rand.Float64()*20.0,
		rnd10:     rand.Int31n(10),
	}
	return c
}

func (c *Cell) DrawCell(screen *ebiten.Image, counter int, emptySubImage *ebiten.Image) {
	c.drawCellBody(screen, c.center, c.direction, c.size, counter, emptySubImage)
	c.drawEyes(screen, c.center, c.direction, c.size*0.9, -1, c.size*0.1, 0xff, counter, emptySubImage)
	c.drawEyes(screen, c.center, c.direction, c.size*0.9, -1, c.size*0.05, 0x00, counter, emptySubImage)
	c.drawEyes(screen, c.center, c.direction, c.size*0.9, 1, c.size*0.1, 0xff, counter, emptySubImage)
	c.drawEyes(screen, c.center, c.direction, c.size*0.9, 1, c.size*0.05, 0x00, counter, emptySubImage)
}
