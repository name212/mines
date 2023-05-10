package game

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type testRandomNumGenerator struct {
	pos     int
	numbers []int
}

func (c *testRandomNumGenerator) Next(max int) (int, error) {
	n := c.numbers[c.pos]
	c.pos++

	return n, nil
}

func TestRandomBombs(t *testing.T) {
	field := newEmptyField(3, 3)

	gen := &testRandomNumGenerator{
		numbers: []int{3, 1, 0, 0, 4},
	}

	const bombs = 3

	linearPositions, err := randomBombsPositions(field, bombs, gen, []int{3})

	require.NoError(t, err)
	require.Len(t, linearPositions, bombs)

	for _, p := range []int{1, 0, 4} {
		require.Contains(t, linearPositions, p)
	}

	randGen := &CryptoRandNumberGenerator{}

	linearPositions, err = randomBombsPositions(field, bombs, randGen, []int{2})

	require.NoError(t, err)
	require.Len(t, linearPositions, bombs)
	require.NotContains(t, linearPositions, 2)
	for _, p := range linearPositions {
		require.Less(t, p, field.size())
	}
}
