package view

import (
	"fmt"

	"github.com/name212/mines/pkg/game"
)

type Symbols struct {
	Closed     string
	Bomb       string
	MarkedBomb string
	Empty      string
}

func SymbolForCell(cell *game.Cell, syms *Symbols, showsBomb bool) string {
	sym := syms.Empty
	if cell.Opened() {
		bombs := cell.BombsAround()
		sym = " "
		if bombs > 0 {
			sym = fmt.Sprintf("%d", bombs)
		} else if cell.HasBomb() {
			sym = syms.Bomb
		}
	} else {
		sym = syms.Closed
		if cell.MarkedAsBomb() {
			sym = syms.MarkedBomb
		}
		if cell.HasBomb() && showsBomb {
			sym = syms.Bomb
		}
	}

	return sym
}
