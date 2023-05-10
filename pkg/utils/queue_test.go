package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueue(t *testing.T) {
	q := NewQueue[int]()
	a := 1
	b := 2

	require.True(t, q.IsEmpty())

	q.Add(&a)
	require.Len(t, q.queue, 1)

	require.False(t, q.IsEmpty())
	require.Len(t, q.queue, 1)

	q.Add(&b)
	require.Len(t, q.queue, 2)

	require.False(t, q.IsEmpty())
	require.Len(t, q.queue, 2)

	aa := q.Dequeue()
	require.Equal(t, a, *aa)
	require.Len(t, q.queue, 1)
	require.False(t, q.IsEmpty())

	bb := q.Dequeue()
	require.Equal(t, b, *bb)
	require.Len(t, q.queue, 0)
	require.True(t, q.IsEmpty())

	require.Nil(t, q.Dequeue())
	require.Len(t, q.queue, 0)
	require.True(t, q.IsEmpty())
}
