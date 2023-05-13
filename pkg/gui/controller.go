package gui

import (
	"github.com/name212/mines/pkg/game"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func Loop() {
	g := newEbitenGame()

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Your game's title")
	g.openSettingsDialog()
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type fieldView interface {
	SetField(field2 *game.Field)
}

type controller struct {
	mines *game.Mines
	view  fieldView
}

func newController() *controller {
	return &controller{}
}

func (c *controller) OnCellLeftClick(x, y int) {
	if c.mines != nil {
		c.mines.Open(x, y)
		if c.view != nil {
			c.view.SetField(c.mines.Field())
		}
	}
}

func (c *controller) OnCellRightClick(x, y int) {
	if c.mines != nil {
		c.mines.SwitchMarkAsBomb(x, y)
		if c.view != nil {
			c.view.SetField(c.mines.Field())
		}
	}
}

func (c *controller) setGame(mines *game.Mines) {
	c.mines = mines
}

func (c *controller) setView(v fieldView) {
	c.view = v
}
