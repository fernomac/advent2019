package queue

import (
	"container/heap"
)

// PQI is an item in a PriorityQueue.
type PQI interface {
	Priority() int
}

// Internal implementation of heap.Heap.
type pq struct {
	items []PQI
}

func (p *pq) Len() int {
	return len(p.items)
}

func (p *pq) Less(i, j int) bool {
	return p.items[i].Priority() < p.items[j].Priority()
}

func (p *pq) Swap(i, j int) {
	p.items[i], p.items[j] = p.items[j], p.items[i]
}

func (p *pq) Push(x interface{}) {
	p.items = append(p.items, x.(PQI))
}

func (p *pq) Pop() interface{} {
	n := len(p.items)
	i := p.items[n-1]
	p.items = p.items[:n-1]
	return i
}

// A PriorityQueue of items.
type PriorityQueue struct {
	p *pq
}

// NewPriority creates a new PriorityQueue.
func NewPriority() PriorityQueue {
	p := &pq{}
	heap.Init(p)
	return PriorityQueue{p}
}

// Push pushes an item onto the pqueue.
func (p PriorityQueue) Push(i PQI) {
	heap.Push(p.p, i)
}

// Pop pops an item off of the pqueue.
func (p PriorityQueue) Pop() PQI {
	if p.p.Len() == 0 {
		return nil
	}
	return heap.Pop(p.p).(PQI)
}

// Len returns the length of the pqueue.
func (p PriorityQueue) Len() int {
	return p.p.Len()
}
