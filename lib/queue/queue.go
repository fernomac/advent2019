package queue

import (
	"container/list"
)

// A Queue of items.
type Queue struct {
	q *list.List
}

// New creates a new queue.
func New() Queue {
	return Queue{list.New()}
}

// Push pushes an item onto the queue.
func (q Queue) Push(v interface{}) {
	q.q.PushBack(v)
}

// Pop pops an item off the queue.
func (q Queue) Pop() interface{} {
	f := q.q.Front()
	if f == nil {
		return nil
	}
	return q.q.Remove(f)
}

// Len returns the length of the queue.
func (q Queue) Len() int {
	return q.q.Len()
}
