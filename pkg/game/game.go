package game

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jtbonhomme/golife/internal/point"
	"github.com/jtbonhomme/golife/pkg/cell"
)

var (
	emptyImage    = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	emptyImage.Fill(color.White)
	rand.Seed(time.Now().UnixNano())
}

type Game struct {
	counter      int
	c            []*cell.Cell
	screenWidth  int
	screenHeight int
}

func New(w, h int) *Game {
	g := &Game{counter: 0, screenWidth: w, screenHeight: h}
	g.c = []*cell.Cell{}
	for i := 0; i < 10; i++ {
		g.c = append(g.c, cell.New(point.Point{X: float64(rand.Int31n(int32(w))), Y: float64(rand.Int31n(int32(h)))}))
	}

	return g
}

func (g *Game) Update() error {
	g.counter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	for _, c := range g.c {
		c.DrawCell(screen, g.counter, emptySubImage)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nCounter: %d", ebiten.CurrentTPS(), ebiten.CurrentFPS(), g.counter))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return g.screenWidth, g.screenHeight
}
