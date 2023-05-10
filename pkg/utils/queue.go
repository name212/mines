package utils

type Queue[C any] struct {
	queue []*C
}

func NewQueue[C any]() *Queue[C] {
	return &Queue[C]{}
}

func (q *Queue[C]) Add(c *C) {
	q.queue = append(q.queue, c)
}

func (q *Queue[C]) Dequeue() *C {
	if len(q.queue) == 0 {
		return nil
	}

	ret := q.queue[0]
	q.queue = q.queue[1:]

	return ret
}

func (q *Queue[C]) IsEmpty() bool {
	return len(q.queue) == 0
}
