package gui

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"

	"github.com/ebitenui/ebitenui"
	ebimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/name212/mines/pkg/game"
	customwidget "github.com/name212/mines/pkg/gui/widget"
	"github.com/name212/mines/pkg/view"
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
		widget.ContainerOpts.BackgroundImage(ebimage.NewNineSliceColor(color.RGBA{0x13, 0x1a, 0x22, 0xff})),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(15),
			),
		),
	)

	statsContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(ebimage.NewNineSliceColor(color.RGBA{0x13, 0x1a, 0x22, 0xff})),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
				widget.RowLayoutOpts.Spacing(15),
			),
		),
	)

	root.AddChild(statsContainer)

	fieldContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(ebimage.NewNineSliceColor(color.RGBA{0x13, 0x1a, 0x22, 0xff})),
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
				Left:   13,
				Right:  13,
				Top:    5,
				Bottom: 5,
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

	// construct a button

	root.AddChild(fieldContainer)

	window = widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(root),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(500, 200),
		widget.WindowOpts.MaxSize(700, 400),
	)

	wWidth, wHeight := window.GetContainer().PreferredSize()

	r := image.Rect(0, 0, wWidth, wHeight)
	r = r.Add(image.Point{windowWidth/2 - wWidth/2, windowHeight/2 - wHeight/2})
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
			if bombs := cell.BombsAround(); bombs > 0 && bombs < 9 {
				fmt.Println("aaaaaaaa")
				button.Image.Disabled = colorsCell[bombs]
			}
		}
	}
}
