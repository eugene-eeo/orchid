package metronome

import "errors"

var QueueEmpty error = errors.New("queue empty")
var QueueFull error = errors.New("queue full")

func mod(r, m int) int {
	t := r % m
	if t < 0 {
		t += m
	}
	return t
}

type Queue struct {
	data []float64
	head int
	tail int
}

func NewQueue(maxSize int) *Queue {
	return &Queue{
		data: make([]float64, maxSize+1),
		head: 0,
		tail: 0,
	}
}

func (q *Queue) Empty() bool {
	return q.head == q.tail
}

func (q *Queue) Size() int {
	return mod(q.tail-q.head, len(q.data))
}

func (q *Queue) Put(x float64) error {
	if q.Size() == len(q.data)-1 {
		return QueueFull
	}
	q.data[q.tail] = x
	q.tail = (q.tail + 1) % len(q.data)
	return nil
}

func (q *Queue) Dequeue() (float64, error) {
	if q.Size() == 0 {
		return 0, QueueEmpty
	}
	x := q.data[q.head]
	q.head = (q.head + 1) % len(q.data)
	return x, nil
}

func (q *Queue) Each(f func(x float64) bool) {
	N := len(q.data)
	for i := q.head; i < q.tail; i = (i + 1) % N {
		if !f(q.data[i]) {
			break
		}
	}
}
