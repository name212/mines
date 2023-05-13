package gui

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/name212/mines/pkg/game"
	"image/color"
)

type ebitenGame struct {
	ui             *ebitenui.UI
	mines          *game.Mines
	gameController *controller
}

func newEbitenGame() *ebitenGame {
	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	ui := &ebitenui.UI{
		Container: rootContainer,
	}

	return &ebitenGame{
		ui:             ui,
		gameController: newController(),
	}
}

func (g *ebitenGame) openSettingsDialog() {
	defaultSettings := settings{
		bombs:  10,
		height: 8,
		width:  8,
	}

	openSettingsDialog(g.ui, &defaultSettings, func(s settings) {
		g.mines = game.NewGame(s.width, s.height, s.bombs, &game.CryptoRandNumberGenerator{})
		g.gameController.setGame(g.mines)
		fieldVied := newField(g.ui, g.mines.Field(), g.gameController)
		g.gameController.setView(fieldVied)
	})
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *ebitenGame) Update() error {
	g.ui.Update()
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *ebitenGame) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *ebitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
