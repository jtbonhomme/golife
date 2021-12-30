package game

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jtbonhomme/golife/internal/vector"
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
	ScreenWidth  int
	ScreenHeight int
}

func New(w, h int) *Game {
	g := &Game{counter: 0, ScreenWidth: w, ScreenHeight: h}
	g.c = []*cell.Cell{}
	for i := 0; i < 10; i++ {
		c := cell.New(vector.Vector2D{
			X: float64(rand.Int31n(int32(w))),
			Y: float64(rand.Int31n(int32(h))),
		}, g.ScreenWidth, g.ScreenHeight)
		g.c = append(g.c, c)
	}

	return g
}

func (g *Game) Update() error {
	g.counter++
	for _, c := range g.c {
		c.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	for _, c := range g.c {
		c.Draw(screen, g.counter, emptySubImage)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nCounter: %d", ebiten.CurrentTPS(), ebiten.CurrentFPS(), g.counter))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return g.ScreenWidth, g.ScreenHeight
}
