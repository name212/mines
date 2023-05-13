package gui

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"github.com/ebitenui/ebitenui"
	ebimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

type settings struct {
	width  int
	height int
	bombs  int
}

func openSettingsDialog(ui *ebitenui.UI, defaultSettings *settings, onStart func(s settings)) {
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

	settingsContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(ebimage.NewNineSliceColor(color.RGBA{0x13, 0x1a, 0x22, 0xff})),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Stretch([]bool{false, false}, []bool{false, false, false}),
				widget.GridLayoutOpts.Spacing(10, 15),
			),
		),
	)

	tOpts := []widget.TextInputOpt{
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextInputOpts.Image(imageTextInput),
		widget.TextInputOpts.Color(colorTextInput),
		widget.TextInputOpts.Padding(insetsTextInput),
		widget.TextInputOpts.Face(fontLabel),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(fontLabel, 2),
		),
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.MinSize(150, -1)),
		widget.TextInputOpts.Validation(func(newInputText string) (bool, *string) {
			i, err := strconv.Atoi(newInputText)
			if err != nil {
				fmt.Println(err)
				return false, nil
			}

			return i > 0, nil
		}),
	}

	settingsContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Width", fontLabel, colorLabelIdle),
		widget.TextOpts.Insets(insetsLabel),
	))

	widthText := widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})))...,
	)

	if defaultSettings != nil {
		widthText.InputText = fmt.Sprintf("%d", defaultSettings.width)
	}

	settingsContainer.AddChild(widthText)

	settingsContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Heights", fontLabel, colorLabelIdle),
		widget.TextOpts.Insets(insetsLabel),
	))

	heightText := widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})))...,
	)

	if defaultSettings != nil {
		heightText.InputText = fmt.Sprintf("%d", defaultSettings.height)
	}

	settingsContainer.AddChild(heightText)

	settingsContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Bombs", fontLabel, colorLabelIdle),
		widget.TextOpts.Insets(insetsLabel),
	))

	bombsText := widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})))...,
	)

	if defaultSettings != nil {
		bombsText.InputText = fmt.Sprintf("%d", defaultSettings.bombs)
	}

	settingsContainer.AddChild(bombsText)

	var rw widget.RemoveWindowFunc

	// construct a button
	button := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
			}),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Start!", fontLabel, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			widthField, err := strconv.Atoi(widthText.InputText)
			if err != nil {
				return
			}

			heightField, err := strconv.Atoi(heightText.InputText)
			if err != nil {
				return
			}

			bombs, err := strconv.Atoi(bombsText.InputText)
			if err != nil {
				return
			}

			rw()

			onStart(settings{
				width:  widthField,
				height: heightField,
				bombs:  bombs,
			})
		}),
	)

	buttonContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(ebimage.NewNineSliceColor(color.RGBA{0x13, 0x1a, 0x22, 0xff})),
		widget.ContainerOpts.Layout(
			widget.NewAnchorLayout(),
		),
		widget.ContainerOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
			}),
		),
	)

	buttonContainer.AddChild(button)

	root.AddChild(settingsContainer)
	root.AddChild(buttonContainer)

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

	rw = ui.AddWindow(window)
}
