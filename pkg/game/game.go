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

func (g *Game) removeCell(i int) error {
	if i < 0 || i >= len(g.cells) {
		return fmt.Errorf("incorrect index: %d", i)
	}
	g.cells = append(g.cells[:i], g.cells[i:]...)
	return nil
}

func (g *Game) Update() error {
	g.counter++
	for i, c := range g.cells {
		if c.IsDead() {
			err := g.removeCell(i)
			if err != nil {
				return fmt.Errorf("cannot remove cell index: %d: %w", i, err)
			}
			continue
		}
		c.Update(g.counter)
	}

	g.DetectCollision()

	return nil
}

func (g *Game) DetectCollision() {
	for i, c1 := range g.cells {
		for j, c2 := range g.cells {
			if c1.ID() != c2.ID() && c1.Intersect(c2) {
				if c1.Size() > c2.Size()*1.1 {
					c1.Eat(c2)
					g.removeCell(j)
				} else if c1.Size()*1.1 < c2.Size() {
					c2.Eat(c1)
					g.removeCell(i)
					continue
				}
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	// draw first debug information
	if g.debug {
		ebitenutil.DebugPrint(
			screen,
			fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nCounter: %d\nCreatures: %d",
				ebiten.CurrentTPS(),
				ebiten.CurrentFPS(),
				g.counter,
				len(g.cells),
			),
		)
		g.linkAgents(screen, 250.0)
	}
	// Draw elements on top of debug information
	for _, c := range g.cells {
		if !c.IsDead() {
			c.Draw(screen, g.counter, emptySubImage)
		}
	}
}

// LinkAgents draws a line between two close agents
func (g *Game) linkAgents(screen *ebiten.Image, radius float64) {
	for i, ci := range g.cells {
		for j := i; j < len(g.cells); j++ {
			if ci.Position().Distance(g.cells[j].Position()) < radius && !g.cells[j].IsDead() {
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

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.ScreenWidth, g.ScreenHeight
}
