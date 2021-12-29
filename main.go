package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	emptyImage    = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	emptyImage.Fill(color.White)
	rand.Seed(time.Now().UnixNano())
}

const (
	screenWidth  int = 640
	screenHeight int = 480
)

func maxCounter(index int) int {
	return 128 + (17*index+32)%64
}

type point struct {
	X float64
	Y float64
}

const PI float64 = math.Pi

func angularToCartesian(dist, direction float64) (x, y float32) {
	return float32(dist * math.Cos(direction)), float32(dist * math.Sin(direction))
}

func addVector(center point, dist, direction float64) (x, y float32) {
	acX, acY := angularToCartesian(dist, direction)
	return float32(center.X) + acX, float32(center.Y) + acY
}

func drawCriterBody(screen *ebiten.Image, center point, direction, size float64, counter int) {
	var path vector.Path
	npoints := 16

	indexToDirection := func(i int) float64 {
		return direction - float64(2*i+1)*PI/float64(npoints)
	}
	indexToDist := func(i, counter int) float64 {
		return size + size*0.1*math.Sin(float64(counter)*2*PI/float64(maxCounter(i)))
	}

	for i := 0; i <= npoints; i++ {
		if i == 0 {
			path.MoveTo(addVector(center, indexToDist(i, counter), indexToDirection(i)))
			continue
		}
		cpx0, cpy0 := addVector(center, indexToDist(i, counter), indexToDirection(i-1)-PI/16)
		cpx1, cpy1 := addVector(center, indexToDist(i, counter), indexToDirection(i)+PI/16)
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

func drawEyes(screen *ebiten.Image, center point, direction, dist, side, size, bg float64, counter int) {
	var path vector.Path

	randomizedFloat64 := func(in float64) float64 {
		return in + rand.Float64()*2
	}

	cpx0, cpy0 := addVector(center, dist-randomizedFloat64(size), direction+side*PI/randomizedFloat64(12))

	path.Arc(cpx0, cpy0, float32(size), float32(0), float32(2*PI), vector.Clockwise)

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

func drawCriter(screen *ebiten.Image, counter int) {
	center := point{X: 150, Y: 150}
	direction := -PI / 2
	criterSize := 40.0
	drawCriterBody(screen, center, direction, criterSize, counter)
	drawEyes(screen, center, direction, criterSize*0.9, -1, criterSize*0.1, 0xff, counter)
	drawEyes(screen, center, direction, criterSize*0.9, -1, criterSize*0.05, 0x00, counter)
	drawEyes(screen, center, direction, criterSize*0.9, 1, criterSize*0.1, 0xff, counter)
	drawEyes(screen, center, direction, criterSize*0.9, 1, criterSize*0.05, 0x00, counter)
}

type Game struct {
	counter int
}

func (g *Game) Update() error {
	g.counter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	drawCriter(screen, g.counter)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nCounter: %d", ebiten.CurrentTPS(), ebiten.CurrentFPS(), g.counter))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{counter: 0}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Vector (Ebiten Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
