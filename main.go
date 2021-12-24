package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

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
}

const (
	screenWidth  = 640
	screenHeight = 480
)

func maxCounter(index int) int {
	return 128 + (17*index+32)%64
}

type point struct {
	X float64
	Y float64
}

const PI float64 = math.Pi

func drawCriterBody(screen *ebiten.Image, counter int) {
	var path vector.Path
	npoints := 8
	center := point{X: 150, Y: 150}
	direction := -PI / 2
	size := float64(100.0)

	indexToPoint := func(i int, counter int) (float32, float32) {
		change := float64(30)
		return float32(center.X + (size+change*math.Sin(float64(counter)*2*math.Pi/float64(maxCounter(i))))*math.Cos(direction-float64(2*i+1)*PI/8)),
			float32(center.Y + (size+change*math.Sin(float64(counter)*2*math.Pi/float64(maxCounter(i))))*math.Sin(direction-float64(2*i+1)*PI/8))
	}

	for i := 0; i < npoints; i++ {
		if i == 0 {
			path.MoveTo(indexToPoint(i, 0))
			continue
		}
		cpx0, cpy0 := indexToPoint(i-1, counter)
		cpx1, cpy1 := indexToPoint(i, counter)
		cpx2, cpy2 := indexToPoint(i, 0)
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

func drawEyes(screen *ebiten.Image, side, size, bg float64, counter int) {
	var path vector.Path
	center := point{X: 150, Y: 150}
	dist := 90.0
	direction := -PI / 2

	path.Arc(float32(center.X-(size-dist)*math.Cos(direction+side*PI/12)), float32(center.Y-(size-dist)*math.Sin(direction+side*PI/12)), float32(size), float32(0), float32(2*PI), vector.Clockwise)

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
	drawCriterBody(screen, counter)
	drawEyes(screen, -1, 15, 0xff, counter)
	drawEyes(screen, -1, 5, 0x00, counter)
	drawEyes(screen, 1, 15, 0xff, counter)
	drawEyes(screen, 1, 5, 0x00, counter)
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
