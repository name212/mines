package game

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type RandNumGenerator interface {
	Next(max int) (int, error)
}

type CryptoRandNumberGenerator struct{}

func (c *CryptoRandNumberGenerator) Next(max int) (int, error) {
	maxx := big.NewInt(int64(max))
	num, err := rand.Int(rand.Reader, maxx)
	if err != nil {
		return 0, err
	}

	return int(num.Int64()), nil
}

func randomBombsPositions(f *Field, bombs int, reader RandNumGenerator, excludePos []int) ([]int, error) {
	maxPos := f.size()
	if bombs > maxPos {
		return nil, fmt.Errorf("incorrect max bombs")
	}

	posMap := make(map[int]struct{}, len(excludePos)+bombs)
	for _, p := range excludePos {
		posMap[p] = struct{}{}
	}

	i := 0
	for i < bombs {
		p, err := reader.Next(maxPos)
		if err != nil {
			return nil, err
		}

		if _, ok := posMap[p]; !ok {
			posMap[p] = struct{}{}
			i++
		}
	}

	for _, p := range excludePos {
		delete(posMap, p)
	}

	res := make([]int, len(posMap))
	j := 0
	for pos := range posMap {
		res[j] = pos
		j++
	}

	return res, nil
}

type initWalkCellAround struct {
	bombs int
}

func (i *initWalkCellAround) HandleCell(c *Cell) (stop bool) {
	if c.HasBomb() {
		i.bombs++
	}
	return false
}

func (i *initWalkCellAround) Reset() {
	i.bombs = 0
}

func initField(f *Field, bombs int, gen RandNumGenerator, startX, startY int) error {
	size := f.size()
	if bombs > size {
		return fmt.Errorf("incorrect bombs %d, maximum is %d", bombs, f.size())
	}

	linearStartPos := f.linearFromPos(startX, startY)
	excludedPos := make([]int, 0, 9)
	excludedPos = append(excludedPos, linearStartPos)
	f.WalkAroundCell(startX, startY, NewFuncCellWalker(func(c *Cell) (stop bool) {
		excludedPos = append(excludedPos, f.linearFromPos(c.X(), c.Y()))
		return false
	}))

	bombsInLinear, err := randomBombsPositions(f, bombs, gen, excludedPos)
	if err != nil {
		return err
	}

	for _, pos := range bombsInLinear {
		cell := f.cell(f.posFromLinear(pos))
		cell.hasBomb = true
	}

	walker := &initWalkCellAround{}
	f.Walk(NewFuncCellWalker(func(c *Cell) (stop bool) {
		if c.HasBomb() {
			return false
		}

		walker.Reset()
		f.WalkAroundCell(c.X(), c.Y(), walker)
		c.bombsAround = walker.bombs

		return false
	}))

	return nil
}
