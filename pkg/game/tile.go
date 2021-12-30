package game

type Tile struct {
	x         int
	y         int
	width     float64
	height    float64
	cellCount int
}

func (t *Tile) ResetCellCount() {
	t.cellCount = 0
}

func (t *Tile) IncCellCount() {
	t.cellCount++
}

func (t *Tile) CellCount() int {
	return t.cellCount
}

func (t *Tile) DecCellCount() {
	if t.cellCount > 0 {
		t.cellCount--
	}
}
