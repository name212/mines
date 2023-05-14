package gui

import (
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/name212/mines/pkg/game"
	customwidget "github.com/name212/mines/pkg/gui/widget"
	"github.com/name212/mines/pkg/view"
	"image"
)

type position struct {
	x, y int
}

type fieldEventsHandler interface {
	OnCellLeftClick(x, y int)
	OnCellRightClick(x, y int)
}

type field struct {
	closeFun  widget.RemoveWindowFunc
	gameField *game.Field
	handler   fieldEventsHandler
	buttons   map[*customwidget.CustomButton]position
}

func newField(ui *ebitenui.UI, f *game.Field, handler fieldEventsHandler) *field {
	buttons := make(map[*customwidget.CustomButton]position, f.Size())
	//var rw widget.RemoveWindowFunc
	var window *widget.Window

	root := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(containerBackground),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(15),
			),
		),
	)

	statsContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(containerBackground),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
				widget.RowLayoutOpts.Spacing(15),
			),
		),
	)

	root.AddChild(statsContainer)

	fieldContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(containerBackground),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(f.Width()),
				widget.GridLayoutOpts.Spacing(5, 5),
			),
		),
	)

	f.Walk(game.NewFuncCellWalker(func(c *game.Cell) (stop bool) {
		button := customwidget.NewCustomButton(
			// specify the images to use
			customwidget.CustomButtonOpts.Image(imageCustomButton.Clone()),

			customwidget.CustomButtonOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(30, 30),
			),

			customwidget.CustomButtonOpts.Text("", fontCell, colorsCellTxt),

			// specify that the button's text needs some padding for correct display
			customwidget.CustomButtonOpts.TextPadding(widget.Insets{
				Left:   4,
				Right:  4,
				Top:    4,
				Bottom: 4,
			}),

			// add a handler that reacts to clicking the button
			customwidget.CustomButtonOpts.PressedHandler(func(args *customwidget.CustomButtonPressedEventArgs) {
				pos, ok := buttons[args.Button]
				if !ok {
					panic("Cannot get button for click")
				}
				switch args.MouseButton {
				case ebiten.MouseButtonLeft:
					handler.OnCellLeftClick(pos.x, pos.y)
				case ebiten.MouseButtonRight:
					handler.OnCellRightClick(pos.x, pos.y)
				}
			}),
		)

		fieldContainer.AddChild(button)

		buttons[button] = position{
			x: c.X(),
			y: c.Y(),
		}

		return false
	}))

	scroll := widget.NewScrollContainer(
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle:     containerBackground,
			Disabled: containerBackground,
			Mask:     containerBackground,
		}),
		widget.ScrollContainerOpts.Content(fieldContainer),
	)

	root.AddChild(scroll)

	window = widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(root),
		widget.WindowOpts.MaxSize(windowWidth, windowHeight),
	)

	wWidth, wHeight := window.GetContainer().PreferredSize()

	r := image.Rect(0, 0, wWidth, wHeight)
	x := windowWidth/2 - wWidth/2
	if x < 0 {
		x = 0
	}
	y := windowHeight/2 - wHeight/2
	if y < 0 {
		y = 0
	}
	r = r.Add(image.Point{x, y})

	fmt.Println(r.String())
	window.SetLocation(r)

	rw := ui.AddWindow(window)

	fw := &field{
		closeFun: rw,
		handler:  handler,
		buttons:  buttons,
	}

	fw.SetField(f)

	return fw
}

var syms = &view.Symbols{
	MarkedBomb: "!",
	Bomb:       "*",
	Closed:     "#",
	Empty:      "",
}

func (f *field) SetField(nField *game.Field) {
	if f.gameField == nField {
		return
	}

	f.gameField = nField

	for button, pos := range f.buttons {
		cell := f.gameField.Cell(pos.x, pos.y)
		s := view.SymbolForCell(cell, syms, false)
		disabled := false
		switch s {
		case syms.Closed:
			s = ""
		case syms.Empty:
			fallthrough
		default:
			disabled = true
		}

		button.GetWidget().Disabled = disabled
		if txt := button.Text(); txt != nil {
			txt.Label = s
		}

		if button.Image != nil {
			if bombs, opened := cell.BombsAround(), cell.Opened(); opened && bombs > 0 && bombs < 9 {
				button.Image.Disabled = colorsCell[bombs]
			}
		}
	}
}
