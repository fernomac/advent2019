package queue

import "testing"

type prioritizedInt int

func (p prioritizedInt) Priority() int {
	return int(p)
}

func TestPriorityQueue(t *testing.T) {
	q := NewPriority()
	q.Push(prioritizedInt(2))
	q.Push(prioritizedInt(1))
	q.Push(prioritizedInt(3))

	if val := q.Pop().(prioritizedInt); val != 1 {
		t.Fatalf("expected 1, got %v", val)
	}
	if val := q.Pop().(prioritizedInt); val != 2 {
		t.Fatalf("expected 2, got %v", val)
	}

	q.Push(prioritizedInt(0))
	q.Push(prioritizedInt(4))

	if val := q.Pop().(prioritizedInt); val != 0 {
		t.Fatalf("expected 0, got %v", val)
	}
	if val := q.Pop().(prioritizedInt); val != 3 {
		t.Fatalf("expected 3, got %v", val)
	}
	if val := q.Pop().(prioritizedInt); val != 4 {
		t.Fatalf("expected 4, got %v", val)
	}

	if val := q.Pop(); val != nil {
		t.Fatalf("expected nil, got %v", val)
	}
}
