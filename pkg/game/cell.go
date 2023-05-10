package game

type Cell struct {
	x, y         int
	bombsAround  int
	hasBomb      bool
	opened       bool
	markedAsBomb bool
}

func newCell(x, y int) *Cell {
	return &Cell{
		x: x,
		y: y,
	}
}

func (c *Cell) BombsAround() int {
	return c.bombsAround
}

func (c *Cell) HasBomb() bool {
	return c.hasBomb
}

func (c *Cell) Opened() bool {
	return c.opened
}

func (c *Cell) MarkedAsBomb() bool {
	return c.markedAsBomb
}

func (c *Cell) setAsOpened() {
	c.opened = true
}

func (c *Cell) markAsBomb() {
	c.markedAsBomb = true
}

func (c *Cell) X() int {
	return c.x
}

func (c *Cell) Y() int {
	return c.y
}

func (c *Cell) clone() *Cell {
	cc := newCell(c.x, c.y)
	cc.hasBomb = c.hasBomb
	cc.bombsAround = c.bombsAround
	cc.markedAsBomb = c.markedAsBomb
	cc.opened = c.opened

	return cc
}
