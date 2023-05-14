package gui

import (
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	customwidget "github.com/name212/mines/pkg/gui/widget"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"strconv"
)

var (
	windowWidth, windowHeight = 700, 480

	fontLabel font.Face
	fontCell  font.Face

	colorLabelIdle = hexToColor("ffffff")

	colorTextInput = &widget.TextInputColor{
		Idle:          hexToColor("ffffff"),
		Disabled:      hexToColor("fafafa"),
		Caret:         hexToColor("ffffff"),
		DisabledCaret: hexToColor("fafafa"),
	}

	imageTextInput = &widget.TextInputImage{
		Idle:     image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
		Disabled: image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
	}

	insetsTextInput = widget.Insets{
		Left:   13,
		Right:  13,
		Top:    0,
		Bottom: 0,
	}

	insetsLabel = widget.Insets{
		Left:   13,
		Right:  13,
		Top:    12,
		Bottom: 12,
	}

	buttonImage = &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255}),
		Hover:   image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255}),
		Pressed: image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255}),
	}

	imageCustomButton = &customwidget.CustomButtonImage{
		Idle:    image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255}),
		Hover:   image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255}),
		Pressed: image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255}),
	}

	colorsCell = []*image.NineSlice{
		nil,
		image.NewNineSliceColor(hexToColor("ddfac3")),
		image.NewNineSliceColor(hexToColor("ecedbf")),
		image.NewNineSliceColor(hexToColor("eddab4")),
		image.NewNineSliceColor(hexToColor("edc38a")),
		image.NewNineSliceColor(hexToColor("f7a1a2")),
		image.NewNineSliceColor(hexToColor("fea785")),
		image.NewNineSliceColor(hexToColor("ff7d60")),
		image.NewNineSliceColor(hexToColor("ff323c")),
	}

	colorsCellTxt = &customwidget.CustomButtonTextColor{
		Idle:     hexToColor("2e3436"),
		Disabled: hexToColor("2e3436"),
	}

	containerBackground = image.NewNineSliceColor(color.RGBA{0x13, 0x1a, 0x22, 0xff})
)

func init() {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}

	fontLabel = truetype.NewFace(ttfFont, &truetype.Options{
		Size:    28,
		DPI:     144,
		Hinting: font.HintingFull,
	})

	fontCell = truetype.NewFace(ttfFont, &truetype.Options{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

func hexToColor(h string) color.Color {
	u, err := strconv.ParseUint(h, 16, 0)
	if err != nil {
		panic(err)
	}

	return color.NRGBA{
		R: uint8(u & 0xff0000 >> 16),
		G: uint8(u & 0xff00 >> 8),
		B: uint8(u & 0xff),
		A: 255,
	}
}
