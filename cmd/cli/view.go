package cli

import (
	"fmt"
	"github.com/name212/mines/pkg/game"
)

const bombSymb = "*"

type Game interface {
	Field() *game.Field
}

type View struct {
	game      Game
	width     int
	showBombs bool
}

func NewView(game Game) *View {
	return &View{
		game: game,
	}
}

func (v *View) Render(showBombs bool) {
	f := v.game.Field()
	v.width = f.Width()
	v.showBombs = showBombs

	fmt.Printf("  | ")

	for i := 0; i < f.Width(); i++ {
		fmt.Printf("%d | ", i)
	}

	fmt.Printf("\n---")

	for i := 0; i < f.Width(); i++ {
		fmt.Printf("----")
	}

	fmt.Printf("\n")

	f.Walk(v)
}

func (v *View) printLineDelim() {
	fmt.Println()
}

func (v *View) HandleCell(cell *game.Cell) bool {
	if cell.X() == 0 {
		fmt.Printf("%d | ", cell.Y())
	}

	sym := "#"
	if cell.Opened() {
		bombs := cell.BombsAround()
		sym = " "
		if bombs > 0 {
			sym = fmt.Sprintf("%d", bombs)
		} else if cell.HasBomb() {
			sym = bombSymb
		}
	} else {
		if cell.MarkedAsBomb() {
			sym = "!"
		}
		if cell.HasBomb() && v.showBombs {
			sym = bombSymb
		}
	}
	fmt.Printf("%s | ", sym)

	if cell.X() == v.width-1 {
		v.printLineDelim()
	}

	return false
}
