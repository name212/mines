package widget

import (
	img "image"
	"image/color"

	"github.com/ebitenui/ebitenui/event"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	. "github.com/ebitenui/ebitenui/widget"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type CustomButton struct {
	Image             *CustomButtonImage
	KeepPressedOnExit bool
	ToggleMode        bool
	GraphicImage      *CustomButtonImageImage
	TextColor         *CustomButtonTextColor

	PressedEvent       *event.Event
	ReleasedEvent      *event.Event
	ClickedEvent       *event.Event
	CursorEnteredEvent *event.Event
	CursorExitedEvent  *event.Event
	StateChangedEvent  *event.Event

	widgetOpts               []WidgetOpt
	autoUpdateTextAndGraphic bool
	textPadding              Insets
	graphicPadding           Insets

	init      *MultiOnce
	widget    *Widget
	container *Container
	graphic   *Graphic
	text      *Text
	hovering  bool
	pressing  bool
	state     WidgetState

	tabOrder      int
	focused       bool
	justSubmitted bool
}

type CustomButtonOpt func(b *CustomButton)

type CustomButtonImage struct {
	Idle         *image.NineSlice
	Hover        *image.NineSlice
	Pressed      *image.NineSlice
	PressedHover *image.NineSlice
	Disabled     *image.NineSlice
}

func (i *CustomButtonImage) Clone() *CustomButtonImage {
	return &CustomButtonImage{
		Idle:         i.Idle,
		Hover:        i.Hover,
		Pressed:      i.Pressed,
		PressedHover: i.PressedHover,
		Disabled:     i.Disabled,
	}
}

type CustomButtonImageImage struct {
	Idle     *ebiten.Image
	Disabled *ebiten.Image
}

type CustomButtonTextColor struct {
	Idle     color.Color
	Disabled color.Color
}

type CustomButtonPressedEventArgs struct {
	MouseButton ebiten.MouseButton
	Button      *CustomButton
	OffsetX     int
	OffsetY     int
}

type CustomButtonReleasedEventArgs struct {
	MouseButton ebiten.MouseButton
	Button      *CustomButton
	Inside      bool
	OffsetX     int
	OffsetY     int
}

type CustomButtonClickedEventArgs struct {
	Button *CustomButton
}
type CustomButtonHoverEventArgs struct {
	Button  *CustomButton
	Entered bool
}
type CustomButtonChangedEventArgs struct {
	Button *CustomButton
	State  WidgetState
}
type CustomButtonPressedHandlerFunc func(args *CustomButtonPressedEventArgs)

type CustomButtonReleasedHandlerFunc func(args *CustomButtonReleasedEventArgs)

type CustomButtonClickedHandlerFunc func(args *CustomButtonClickedEventArgs)

type CustomButtonCursorHoverHandlerFunc func(args *CustomButtonHoverEventArgs)

type CustomButtonChangedHandlerFunc func(args *CustomButtonChangedEventArgs)

type CustomButtonOptions struct {
}

var CustomButtonOpts CustomButtonOptions

func NewCustomButton(opts ...CustomButtonOpt) *CustomButton {
	b := &CustomButton{
		PressedEvent:       &event.Event{},
		ReleasedEvent:      &event.Event{},
		ClickedEvent:       &event.Event{},
		CursorEnteredEvent: &event.Event{},
		CursorExitedEvent:  &event.Event{},
		StateChangedEvent:  &event.Event{},

		init: &MultiOnce{},
	}

	b.init.Append(b.createWidget)

	for _, o := range opts {
		o(b)
	}

	return b
}

func (o CustomButtonOptions) WidgetOpts(opts ...WidgetOpt) CustomButtonOpt {
	return func(b *CustomButton) {
		b.widgetOpts = append(b.widgetOpts, opts...)
	}
}

func (o CustomButtonOptions) Image(i *CustomButtonImage) CustomButtonOpt {
	return func(b *CustomButton) {
		b.Image = i
	}
}

func (o CustomButtonOptions) TextSimpleLeft(label string, face font.Face, color *CustomButtonTextColor, padding Insets) CustomButtonOpt {
	return func(b *CustomButton) {
		b.init.Append(func() {
			b.container = NewContainer(
				ContainerOpts.Layout(NewAnchorLayout(AnchorLayoutOpts.Padding(padding))),
				ContainerOpts.AutoDisableChildren(),
			)

			b.text = NewText(
				TextOpts.WidgetOpts(WidgetOpts.LayoutData(AnchorLayoutData{
					HorizontalPosition: AnchorLayoutPositionStart,
					VerticalPosition:   AnchorLayoutPositionCenter,
				})),
				TextOpts.Text(label, face, color.Idle),
				TextOpts.Position(TextPositionStart, TextPositionCenter),
			)
			b.container.AddChild(b.text)

			b.autoUpdateTextAndGraphic = true
			b.TextColor = color
		})
	}
}

func (o CustomButtonOptions) Text(label string, face font.Face, color *CustomButtonTextColor) CustomButtonOpt {
	return func(b *CustomButton) {
		b.init.Append(func() {
			b.container = NewContainer(
				ContainerOpts.Layout(NewAnchorLayout(AnchorLayoutOpts.Padding(b.textPadding))),
				ContainerOpts.AutoDisableChildren(),
			)

			b.text = NewText(
				TextOpts.WidgetOpts(WidgetOpts.LayoutData(AnchorLayoutData{
					HorizontalPosition: AnchorLayoutPositionCenter,
					VerticalPosition:   AnchorLayoutPositionCenter,
				})),
				TextOpts.Text(label, face, color.Idle),
				TextOpts.Position(TextPositionCenter, TextPositionCenter),
			)
			b.container.AddChild(b.text)

			b.autoUpdateTextAndGraphic = true
			b.TextColor = color
		})
	}
}

// TODO: add parameter for image position (start/end)
func (o CustomButtonOptions) TextAndImage(label string, face font.Face, image *CustomButtonImageImage, color *CustomButtonTextColor) CustomButtonOpt {
	return func(b *CustomButton) {
		b.init.Append(func() {
			b.container = NewContainer(
				ContainerOpts.Layout(NewAnchorLayout(AnchorLayoutOpts.Padding(b.textPadding))),
				ContainerOpts.AutoDisableChildren(),
			)

			c := NewContainer(
				ContainerOpts.WidgetOpts(WidgetOpts.LayoutData(AnchorLayoutData{
					HorizontalPosition: AnchorLayoutPositionCenter,
					VerticalPosition:   AnchorLayoutPositionCenter,
				})),
				ContainerOpts.Layout(NewRowLayout(RowLayoutOpts.Spacing(10))),
				ContainerOpts.AutoDisableChildren(),
			)
			b.container.AddChild(c)

			b.text = NewText(
				TextOpts.WidgetOpts(WidgetOpts.LayoutData(RowLayoutData{
					Stretch: true,
				})),
				TextOpts.Text(label, face, color.Idle))
			c.AddChild(b.text)

			b.graphic = NewGraphic(
				GraphicOpts.WidgetOpts(WidgetOpts.LayoutData(RowLayoutData{
					Stretch: true,
				})),
				GraphicOpts.Image(image.Idle))
			c.AddChild(b.graphic)

			b.autoUpdateTextAndGraphic = true
			b.GraphicImage = image
			b.TextColor = color
		})
	}
}

func (o CustomButtonOptions) TextPadding(p Insets) CustomButtonOpt {
	return func(b *CustomButton) {
		b.textPadding = p
	}
}

func (o CustomButtonOptions) Graphic(i *ebiten.Image) CustomButtonOpt {
	return o.withGraphic(GraphicOpts.Image(i))
}

func (o CustomButtonOptions) GraphicNineSlice(i *image.NineSlice) CustomButtonOpt {
	return o.withGraphic(GraphicOpts.ImageNineSlice(i))
}

func (o CustomButtonOptions) withGraphic(opt GraphicOpt) CustomButtonOpt {
	return func(b *CustomButton) {
		b.init.Append(func() {
			b.container = NewContainer(
				ContainerOpts.Layout(NewAnchorLayout(AnchorLayoutOpts.Padding(b.graphicPadding))),
				ContainerOpts.AutoDisableChildren())

			b.graphic = NewGraphic(
				opt,
				GraphicOpts.WidgetOpts(WidgetOpts.LayoutData(AnchorLayoutData{
					HorizontalPosition: AnchorLayoutPositionCenter,
					VerticalPosition:   AnchorLayoutPositionCenter,
				})),
			)
			b.container.AddChild(b.graphic)

			b.autoUpdateTextAndGraphic = true
		})
	}
}

func (o CustomButtonOptions) GraphicPadding(i Insets) CustomButtonOpt {
	return func(b *CustomButton) {
		b.graphicPadding = i
	}
}

func (o CustomButtonOptions) KeepPressedOnExit() CustomButtonOpt {
	return func(b *CustomButton) {
		b.KeepPressedOnExit = true
	}
}

func (o CustomButtonOptions) ToggleMode() CustomButtonOpt {
	return func(b *CustomButton) {
		b.ToggleMode = true
	}
}

func (o CustomButtonOptions) PressedHandler(f CustomButtonPressedHandlerFunc) CustomButtonOpt {
	return func(b *CustomButton) {
		b.PressedEvent.AddHandler(func(args interface{}) {
			f(args.(*CustomButtonPressedEventArgs))
		})
	}
}

func (o CustomButtonOptions) ReleasedHandler(f CustomButtonReleasedHandlerFunc) CustomButtonOpt {
	return func(b *CustomButton) {
		b.ReleasedEvent.AddHandler(func(args interface{}) {
			f(args.(*CustomButtonReleasedEventArgs))
		})
	}
}

func (o CustomButtonOptions) ClickedHandler(f CustomButtonClickedHandlerFunc) CustomButtonOpt {
	return func(b *CustomButton) {
		b.ClickedEvent.AddHandler(func(args interface{}) {
			f(args.(*CustomButtonClickedEventArgs))
		})
	}
}

func (o CustomButtonOptions) CursorEnteredHandler(f CustomButtonCursorHoverHandlerFunc) CustomButtonOpt {
	return func(b *CustomButton) {
		b.CursorEnteredEvent.AddHandler(func(args interface{}) {
			f(args.(*CustomButtonHoverEventArgs))
		})
	}
}

func (o CustomButtonOptions) CursorExitedHandler(f CustomButtonCursorHoverHandlerFunc) CustomButtonOpt {
	return func(b *CustomButton) {
		b.CursorExitedEvent.AddHandler(func(args interface{}) {
			f(args.(*CustomButtonHoverEventArgs))
		})
	}
}

func (o CustomButtonOptions) StateChangedHandler(f CustomButtonChangedHandlerFunc) CustomButtonOpt {
	return func(b *CustomButton) {
		b.StateChangedEvent.AddHandler(func(args interface{}) {
			f(args.(*CustomButtonChangedEventArgs))
		})
	}
}

func (o CustomButtonOptions) TabOrder(tabOrder int) CustomButtonOpt {
	return func(b *CustomButton) {
		b.tabOrder = tabOrder
	}
}

func (tw *CustomButton) State() WidgetState {
	return tw.state
}

func (tw *CustomButton) SetState(state WidgetState) {
	if state != tw.state {
		tw.state = state

		tw.StateChangedEvent.Fire(&CustomButtonChangedEventArgs{
			Button: tw,
			State:  tw.state,
		})
	}
}

func (tw *CustomButton) getStateChangedEvent() *event.Event {
	return tw.StateChangedEvent
}

func (b *CustomButton) Configure(opts ...CustomButtonOpt) {
	for _, o := range opts {
		o(b)
	}
}

func (b *CustomButton) Focus(focused bool) {
	b.init.Do()
	b.GetWidget().FireFocusEvent(b, focused, img.Point{-1, -1})
	b.focused = focused
}

func (b *CustomButton) TabOrder() int {
	return b.tabOrder
}

func (b *CustomButton) GetWidget() *Widget {
	b.init.Do()
	return b.widget
}

func (b *CustomButton) PreferredSize() (int, int) {
	b.init.Do()

	w, h := 50, 50

	if b.container != nil && len(b.container.Children()) > 0 {
		w, h = b.container.PreferredSize()
	}

	if b.widget != nil && h < b.widget.MinHeight {
		h = b.widget.MinHeight
	}
	if b.widget != nil && w < b.widget.MinWidth {
		w = b.widget.MinWidth
	}

	iw, ih := b.Image.Idle.MinSize()
	if w < iw {
		w = iw
	}
	if h < ih {
		h = ih
	}

	return w, h
}

func (b *CustomButton) SetLocation(rect img.Rectangle) {
	b.init.Do()
	b.widget.Rect = rect
}

func (b *CustomButton) RequestRelayout() {
	b.init.Do()

	if b.container != nil {
		b.container.RequestRelayout()
	}
}

func (b *CustomButton) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	b.init.Do()

	if b.container != nil {
		b.container.SetupInputLayer(def)
	}
}

func (b *CustomButton) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	b.init.Do()

	if b.container != nil {
		w := b.container.GetWidget()
		w.Rect = b.widget.Rect
		w.Disabled = b.widget.Disabled
		b.container.RequestRelayout()
	}

	b.widget.Render(screen, def)
	b.handleSubmit()
	b.draw(screen)

	if b.autoUpdateTextAndGraphic {
		if b.graphic != nil {
			if b.widget.Disabled {
				b.graphic.Image = b.GraphicImage.Disabled
			} else {
				b.graphic.Image = b.GraphicImage.Idle
			}
		}

		if b.text != nil {
			if b.widget.Disabled {
				b.text.Color = b.TextColor.Disabled
			} else {
				b.text.Color = b.TextColor.Idle
			}
		}
	}

	if b.container != nil {
		b.container.Render(screen, def)
	}
}

func (b *CustomButton) draw(screen *ebiten.Image) {
	i := b.Image.Idle
	switch {
	case b.widget.Disabled:
		if b.Image.Disabled != nil {
			i = b.Image.Disabled
		}
	case b.focused, b.hovering:
		if b.ToggleMode && b.state == WidgetChecked || b.pressing && (b.hovering || b.KeepPressedOnExit) {
			if b.Image.PressedHover != nil {
				i = b.Image.PressedHover
			} else {
				i = b.Image.Pressed
			}
		} else {
			if b.Image.Hover != nil {
				i = b.Image.Hover
			}
		}
	case b.pressing && (b.hovering || b.KeepPressedOnExit) || (b.ToggleMode && b.state == WidgetChecked):
		if b.Image.Pressed != nil {
			i = b.Image.Pressed
		}

	}

	if i != nil {
		i.Draw(screen, b.widget.Rect.Dx(), b.widget.Rect.Dy(), func(opts *ebiten.DrawImageOptions) {
			opts.GeoM.Translate(float64(b.widget.Rect.Min.X), float64(b.widget.Rect.Min.Y))
			b.drawImageOptions(opts)
		})
	}
}

func (b *CustomButton) handleSubmit() {
	if input.KeyPressed(ebiten.KeyEnter) || input.KeyPressed(ebiten.KeySpace) {
		if !b.justSubmitted && b.focused {
			b.ClickedEvent.Fire(&CustomButtonClickedEventArgs{
				Button: b,
			})
			if b.ToggleMode {
				if b.state == WidgetUnchecked {
					b.state = WidgetChecked
				} else {
					b.state = WidgetUnchecked
				}
				b.StateChangedEvent.Fire(&CustomButtonChangedEventArgs{
					Button: b,
					State:  b.state,
				})
			}
			b.justSubmitted = true
		}
	} else {
		b.justSubmitted = false
	}
}

func (b *CustomButton) drawImageOptions(opts *ebiten.DrawImageOptions) {
	if b.widget.Disabled && b.Image.Disabled == nil {
		opts.ColorM.Scale(1, 1, 1, 0.35)
	}
}

func (b *CustomButton) Text() *Text {
	b.init.Do()
	return b.text
}

func (b *CustomButton) createWidget() {
	b.widget = NewWidget(append(b.widgetOpts, []WidgetOpt{
		WidgetOpts.CursorEnterHandler(func(_ *WidgetCursorEnterEventArgs) {
			if !b.widget.Disabled {
				b.hovering = true
			}
			b.CursorEnteredEvent.Fire(&CustomButtonHoverEventArgs{
				Button:  b,
				Entered: true,
			})
		}),

		WidgetOpts.CursorExitHandler(func(_ *WidgetCursorExitEventArgs) {
			b.hovering = false
			b.CursorExitedEvent.Fire(&CustomButtonHoverEventArgs{
				Button:  b,
				Entered: false,
			})
		}),

		WidgetOpts.MouseButtonPressedHandler(func(args *WidgetMouseButtonPressedEventArgs) {
			b.pressing = true

			b.PressedEvent.Fire(&CustomButtonPressedEventArgs{
				MouseButton: args.Button,
				Button:      b,
				OffsetX:     args.OffsetX,
				OffsetY:     args.OffsetY,
			})
		}),

		WidgetOpts.MouseButtonReleasedHandler(func(args *WidgetMouseButtonReleasedEventArgs) {
			b.pressing = false

			if !b.widget.Disabled && args.Button == ebiten.MouseButtonLeft {
				b.ReleasedEvent.Fire(&CustomButtonReleasedEventArgs{
					Button:  b,
					Inside:  args.Inside,
					OffsetX: args.OffsetX,
					OffsetY: args.OffsetY,
				})

				if args.Inside {
					b.ClickedEvent.Fire(&CustomButtonClickedEventArgs{
						Button: b,
					})
					if b.ToggleMode {
						if b.state == WidgetUnchecked {
							b.state = WidgetChecked
						} else {
							b.state = WidgetUnchecked
						}
						b.StateChangedEvent.Fire(&CustomButtonChangedEventArgs{
							Button: b,
							State:  b.state,
						})
					}
				}
			}
		}),
	}...)...)
	b.widgetOpts = nil
}
