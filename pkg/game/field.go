package game

type CellWalker interface {
	HandleCell(c *Cell) (stop bool)
}

type FuncCellWalker struct {
	f func(c *Cell) bool
}

func NewFuncCellWalker(f func(c *Cell) (stop bool)) *FuncCellWalker {
	return &FuncCellWalker{
		f: f,
	}
}

func (w *FuncCellWalker) HandleCell(c *Cell) (stop bool) {
	return w.f(c)
}

type Field struct {
	width  int
	height int
	cells  [][]*Cell
}

func newField(width, height int, cellConstructor func(x, y int) *Cell) *Field {
	cells := make([][]*Cell, height)
	var i int
	for i = 0; i < height; i++ {
		cells[i] = make([]*Cell, width)
		var j int
		for j = 0; j < width; j++ {
			cells[i][j] = cellConstructor(j, i)
		}
	}

	return &Field{
		width:  width,
		height: height,
		cells:  cells,
	}
}

func newEmptyField(width, height int) *Field {
	return newField(width, height, func(x, y int) *Cell {
		return newCell(x, y)
	})
}

func (f *Field) Width() int {
	return f.width
}

func (f *Field) Height() int {
	return f.height
}

func (f *Field) Cell(x, y int) *Cell {
	return f.cells[y][x]
}

func (f *Field) Clone() *Field {
	return newField(f.Width(), f.Height(), func(x, y int) *Cell {
		return f.Cell(x, y).clone()
	})
}

func (f *Field) Walk(walker CellWalker) {
	for y := 0; y < f.Height(); y++ {
		for x := 0; x < f.Width(); x++ {
			if walker.HandleCell(f.Cell(x, y)) {
				return
			}
		}
	}
}

func (f *Field) WalkAroundCell(x, y int, walker CellWalker) {
	startX, endX := x-1, x+1
	if startX < 0 {
		startX = 0
	}
	if endX >= f.width {
		endX = f.width - 1
	}

	startY, endY := y-1, y+1
	if startY < 0 {
		startY = 0
	}
	if endY >= f.height {
		endY = f.height - 1
	}

	for i := startX; i <= endX; i++ {
		for j := startY; j <= endY; j++ {
			if i == x && j == y {
				continue
			}
			if walker.HandleCell(f.Cell(i, j)) {
				return
			}
		}
	}
}

func (f *Field) linearFromPos(x, y int) int {
	return y*f.width + x
}

func (f *Field) Size() int {
	return f.width * f.height
}

func (f *Field) posFromLinear(pos int) (x, y int) {
	y = pos / f.width
	x = pos % f.width

	return x, y
}
