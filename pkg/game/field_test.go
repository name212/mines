package game

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPosFromLinear(t *testing.T) {
	fSquare := newEmptyField(3, 3)
	fRect1 := newEmptyField(2, 3)
	fRect2 := newEmptyField(3, 2)
	fHosLine := newEmptyField(3, 1)
	fVerLine := newEmptyField(1, 3)

	cases := []struct {
		field *Field

		pos       int
		expectedX int
		expectedY int
	}{
		{
			field:     fSquare,
			pos:       0,
			expectedX: 0,
			expectedY: 0,
		},
		{
			field:     fSquare,
			pos:       1,
			expectedX: 1,
			expectedY: 0,
		},
		{
			field:     fSquare,
			pos:       2,
			expectedX: 2,
			expectedY: 0,
		},
		{
			field:     fSquare,
			pos:       3,
			expectedX: 0,
			expectedY: 1,
		},
		{
			field:     fSquare,
			pos:       4,
			expectedX: 1,
			expectedY: 1,
		},
		{
			field:     fSquare,
			pos:       5,
			expectedX: 2,
			expectedY: 1,
		},
		{
			field:     fSquare,
			pos:       6,
			expectedX: 0,
			expectedY: 2,
		},
		{
			field:     fSquare,
			pos:       7,
			expectedX: 1,
			expectedY: 2,
		},
		{
			field:     fSquare,
			pos:       8,
			expectedX: 2,
			expectedY: 2,
		},
		// rect 1
		{
			field:     fRect1,
			pos:       0,
			expectedX: 0,
			expectedY: 0,
		},
		{
			field:     fRect1,
			pos:       1,
			expectedX: 1,
			expectedY: 0,
		},
		{
			field:     fRect1,
			pos:       2,
			expectedX: 0,
			expectedY: 1,
		},
		{
			field:     fRect1,
			pos:       3,
			expectedX: 1,
			expectedY: 1,
		},
		{
			field:     fRect1,
			pos:       4,
			expectedX: 0,
			expectedY: 2,
		},
		{
			field:     fRect1,
			pos:       5,
			expectedX: 1,
			expectedY: 2,
		},

		// react 2
		{
			field:     fRect2,
			pos:       0,
			expectedX: 0,
			expectedY: 0,
		},
		{
			field:     fRect2,
			pos:       1,
			expectedX: 1,
			expectedY: 0,
		},
		{
			field:     fRect2,
			pos:       2,
			expectedX: 2,
			expectedY: 0,
		},
		{
			field:     fRect2,
			pos:       3,
			expectedX: 0,
			expectedY: 1,
		},
		{
			field:     fRect2,
			pos:       4,
			expectedX: 1,
			expectedY: 1,
		},
		{
			field:     fRect2,
			pos:       5,
			expectedX: 2,
			expectedY: 1,
		},

		// hos line
		{
			field:     fHosLine,
			pos:       0,
			expectedX: 0,
			expectedY: 0,
		},
		{
			field:     fHosLine,
			pos:       1,
			expectedX: 1,
			expectedY: 0,
		},
		{
			field:     fHosLine,
			pos:       2,
			expectedX: 2,
			expectedY: 0,
		},

		// vert line
		{
			field:     fVerLine,
			pos:       0,
			expectedX: 0,
			expectedY: 0,
		},
		{
			field:     fVerLine,
			pos:       1,
			expectedX: 0,
			expectedY: 1,
		},
		{
			field:     fVerLine,
			pos:       2,
			expectedX: 0,
			expectedY: 2,
		},
	}

	for _, c := range cases {
		x, y := c.field.posFromLinear(c.pos)
		require.Equal(t, x, c.expectedX, fmt.Sprintf("incorrect x %d for %d pos in field %dx%d, should be %d", x, c.pos, c.field.width, c.field.height, c.expectedX))
		require.Equal(t, y, c.expectedY, fmt.Sprintf("incorrect y for %d pos in field %dx%d, should be %d", y, c.field.width, c.field.height, c.expectedY))

		pos := c.field.linearFromPos(c.expectedX, c.expectedY)
		require.Equal(t, pos, c.pos, fmt.Sprintf("incorrect pos %d for x=%d y=%d in field %dx%d, should be %d", pos, c.expectedX, c.expectedY, c.field.width, c.field.height, c.pos))
	}
}

type testPos struct {
	x, y int
}

type testCellAroundWalker struct {
	cells []testPos
}

func (w *testCellAroundWalker) HandleCell(c *Cell) (stop bool) {
	w.cells = append(w.cells, testPos{
		x: c.X(),
		y: c.Y(),
	})

	return false
}

func TestWalkCellAround(t *testing.T) {
	f := newEmptyField(8, 8)

	assertWalkAround := func(t *testing.T, x, y int, expectPos []testPos) {
		w := &testCellAroundWalker{}
		f.WalkAroundCell(x, y, w)
		require.Len(t, w.cells, len(expectPos))
		for _, p := range expectPos {
			require.Contains(t, w.cells, p)
		}
	}

	// left up corner
	assertWalkAround(t, 0, 0, []testPos{{1, 0}, {1, 1}, {0, 1}})

	// left down corner
	assertWalkAround(t, 0, 7, []testPos{{1, 7}, {1, 6}, {0, 6}})

	// right down corner
	assertWalkAround(t, 7, 7, []testPos{{6, 7}, {6, 6}, {7, 6}})

	// right up corner
	assertWalkAround(t, 7, 0, []testPos{{6, 0}, {6, 1}, {7, 1}})

	// up line
	assertWalkAround(t, 4, 0, []testPos{{3, 0}, {5, 0}, {3, 1}, {4, 1}, {5, 1}})

	// down line
	assertWalkAround(t, 4, 7, []testPos{{3, 7}, {5, 7}, {3, 6}, {4, 6}, {5, 6}})

	// left line
	assertWalkAround(t, 0, 4, []testPos{{0, 3}, {0, 5}, {1, 3}, {1, 4}, {1, 5}})

	// right line
	assertWalkAround(t, 7, 4, []testPos{{7, 3}, {7, 5}, {6, 3}, {6, 4}, {6, 5}})

	// center
	assertWalkAround(t, 4, 4, []testPos{{3, 3}, {3, 4}, {3, 5}, {4, 3}, {4, 5}, {5, 3}, {5, 4}, {5, 5}})
}
