package cli

import (
	"flag"
	"fmt"

	"github.com/name212/mines/pkg/game"
)

func GetSettings() (width, height, bombs int) {
	widthFlag := flag.Int("width", -1, "Field width")
	heightFlag := flag.Int("height", -1, "Field height")
	bombsFlag := flag.Int("bombs", -1, "Bombs on field")

	flag.Parse()

	if *widthFlag > 0 && *heightFlag > 0 || *bombsFlag > 0 {
		width = *widthFlag
		height = *heightFlag
		bombs = *bombsFlag

		return width, height, bombs
	}

	var w, h, b int
	for {
		fmt.Print("Enter field width: ")
		fmt.Scanf("%d", &w)
		if w > 0 {
			break
		}
	}

	for {
		fmt.Print("Enter field height: ")
		fmt.Scanf("%d", &h)
		if h > 0 {
			break
		}
	}

	for {
		fmt.Print("Enter bombs count: ")
		fmt.Scanf("%d", &b)
		if b > 0 {
			break
		}
	}

	return w, h, b
}

func Loop() {
	width, height, bombs := GetSettings()
	mines := game.NewGame(width, height, bombs, &game.CryptoRandNumberGenerator{})
	view := NewView(mines)
	for {
		status := mines.Status()
		switch status {
		case game.Error:
			fmt.Println("Game in error state")
			return
		case game.Lose:
			view.Render(true)
			fmt.Println("You've lost :-(")
			return
		case game.Win:
			view.Render(false)
			fmt.Printf("You've WIN! :-) Your time: %s\n", mines.FinishedAt().Sub(mines.StartedAt()).String())
			return
		}

		var x, y int
		var mode string
		fmt.Printf("Enter x and y for open cell [%d/%d]: ", mines.BombsMarked(), mines.Bombs())
		fmt.Scanf("%d %d %s", &x, &y, &mode)

		if x < 0 || x >= width || y < 0 || y >= height {
			continue
		}

		if mode != "" && mode != "b" {
			continue
		}

		if mode == "b" {
			mines.SwitchMarkAsBomb(x, y)
		} else {
			mines.Open(x, y)
		}

		view.Render(false)
	}
}
