package game

type Tile struct {
	x      int
	y      int
	width  float64
	height float64
	cells  []string
}

func (t *Tile) ResetCellCount() {
	t.cells = []string{}
}

func (t *Tile) AddCell(id string) {
	t.cells = append(t.cells, id)
}

func (t *Tile) CellCount() int {
	return len(t.cells)
}
