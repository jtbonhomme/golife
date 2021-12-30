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
	cells        []*cell.Cell
	ScreenWidth  int
	ScreenHeight int
	debug        bool
}

func New(w, h int) *Game {
	g := &Game{counter: 0, ScreenWidth: w, ScreenHeight: h, debug: true}
	g.cells = []*cell.Cell{}
	for i := 0; i < 10; i++ {
		c := cell.New(vector.Vector2D{
			X: float64(rand.Int31n(int32(w))),
			Y: float64(rand.Int31n(int32(h))),
		}, g.ScreenWidth, g.ScreenHeight)
		c.Debug(g.debug)
		g.cells = append(g.cells, c)
	}

	return g
}

func (g *Game) Update() error {
	g.counter++
	for _, c := range g.cells {
		c.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	if g.debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nCounter: %d", ebiten.CurrentTPS(), ebiten.CurrentFPS(), g.counter))
		g.LinkAgents(screen, 250.0)
	}
	// Draw elements on top of debug informations
	for _, c := range g.cells {
		c.Draw(screen, g.counter, emptySubImage)
	}
}

// LinkAgents draws a line between two close agents
func (g *Game) LinkAgents(screen *ebiten.Image, radius float64) {
	for i, ci := range g.cells {
		for j := i; j < len(g.cells); j++ {
			if ci.Position().Distance(g.cells[j].Position()) < radius {
				// Draw line between agents
				ebitenutil.DrawLine(
					screen,
					ci.Position().X, ci.Position().Y,
					g.cells[j].Position().X, g.cells[j].Position().Y,
					color.Gray16{0xdddd},
				)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return g.ScreenWidth, g.ScreenHeight
}
