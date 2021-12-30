package game

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/jtbonhomme/golife/internal/fonts"
	"github.com/jtbonhomme/golife/internal/vector"
	"github.com/jtbonhomme/golife/pkg/cell"
)

const (
	nCells int = 40
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
	counter       int
	cells         map[string]*cell.Cell
	tiles         [][]*Tile
	TileDimension int
	ScreenWidth   int
	ScreenHeight  int
	debug         bool
	startTime     time.Time
	gameDuration  time.Duration
}

func New(w, h, t int) *Game {
	g := &Game{
		counter:       0,
		ScreenWidth:   w,
		ScreenHeight:  h,
		TileDimension: t,
		startTime:     time.Now(),
		gameDuration:  0,
		debug:         true,
		cells:         make(map[string]*cell.Cell),
	}
	g.tiles = [][]*Tile{}
	for i := 0; i < nCells; i++ {
		c := cell.New(vector.Vector2D{
			X: float64(rand.Int31n(int32(w))),
			Y: float64(rand.Int31n(int32(h))),
		}, g.ScreenWidth, g.ScreenHeight, g.Detect)
		c.Debug(g.debug)
		g.cells[c.ID()] = c
	}
	g.tiles = make([][]*Tile, g.ScreenWidth/g.TileDimension)
	for i := 0; i < g.ScreenWidth/g.TileDimension; i++ {
		g.tiles[i] = make([]*Tile, g.ScreenHeight/g.TileDimension)
		for j := 0; j < g.ScreenHeight/g.TileDimension; j++ {
			g.tiles[i][j] = &Tile{x: i, y: j, width: float64(g.TileDimension), height: float64(g.TileDimension), cells: []string{}}
		}
	}
	return g
}

func (g *Game) removeCell(c *cell.Cell) {
	delete(g.cells, c.ID())
}

func (g *Game) resetTiles() {
	for i := 0; i < g.ScreenWidth/g.TileDimension; i++ {
		for j := 0; j < g.ScreenHeight/g.TileDimension; j++ {
			g.tiles[i][j].ResetCellCount()
		}
	}
}

func (g *Game) Update() error {
	g.counter++
	g.resetTiles()
	if len(g.cells) == 0 {
		return fmt.Errorf("all cells are dead")
	}

	for _, c := range g.cells {
		if c.IsDead() {
			g.removeCell(c)
			continue
		}
		// update tile count
		x := int(math.Floor(c.Position().X / float64(g.TileDimension)))
		y := int(math.Floor(c.Position().Y / float64(g.TileDimension)))
		g.tiles[x][y].AddCell(c.ID())

		// update cell state
		c.Update(g.counter)
	}
	g.gameDuration = time.Since(g.startTime).Round(time.Second)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// blank screen
	screen.Fill(color.White)
	// draw first debug information
	if g.debug {
		for i := 0; i < g.ScreenWidth/g.TileDimension; i++ {
			for j := 0; j < g.ScreenHeight/g.TileDimension; j++ {
				if g.tiles[i][j].CellCount() > 0 {
					ebitenutil.DrawRect(
						screen,
						float64(i*g.TileDimension),
						float64(j*g.TileDimension),
						float64(g.TileDimension),
						float64(g.TileDimension),
						color.Gray16{0xeeee},
					)
				}
			}
		}

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
	g.drawTimeElapsed(screen)
}

// LinkAgents draws a line between two close agents
func (g *Game) linkAgents(screen *ebiten.Image, radius float64) {
	for _, ci := range g.cells {
		for _, cj := range g.cells {
			if ci.ID() != cj.ID() && ci.Position().Distance(cj.Position()) < radius && !cj.IsDead() {
				// Draw line between agents
				ebitenutil.DrawLine(
					screen,
					ci.Position().X, ci.Position().Y,
					cj.Position().X, cj.Position().Y,
					color.Gray16{0xcccc},
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

func (g *Game) drawTimeElapsed(screen *ebiten.Image) {
	// Time elapsed
	elapsed := "Time elapsed " + g.gameDuration.String()
	elapsedTextDim := text.BoundString(fonts.MonoSansRegularFont, elapsed)
	elapsedTextHeight := elapsedTextDim.Max.Y - elapsedTextDim.Min.Y
	text.Draw(
		screen,
		elapsed,
		fonts.MonoSansRegularFont,
		100,
		elapsedTextHeight+10,
		color.Black,
	)
}

// Detect returns all cells located in a radius from (x,y)
func (g *Game) Detect(pos vector.Vector2D, radius float64) []*cell.Cell {
	nearestCells := []*cell.Cell{}

	for _, c := range g.cells {
		if pos.SquareDistance(c.Position()) < radius*radius {
			nearestCells = append(nearestCells, c)
		}
	}

	return nearestCells
}
