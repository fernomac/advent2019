package queue

import "testing"

func TestQueue(t *testing.T) {
	q := New()
	q.Push(1)
	q.Push(2)
	q.Push(3)

	if val := q.Len(); val != 3 {
		t.Errorf("expected 3, got %v", val)
	}

	if val := q.Pop().(int); val != 1 {
		t.Fatalf("expected 1, got %v", val)
	}
	if val := q.Pop().(int); val != 2 {
		t.Fatalf("expected 2, got %v", val)
	}

	if val := q.Len(); val != 1 {
		t.Errorf("expected 1, got %v", val)
	}

	q.Push(4)

	if val := q.Len(); val != 2 {
		t.Errorf("expected 2, got %v", val)
	}

	if val := q.Pop().(int); val != 3 {
		t.Fatalf("expected 3, got %v", val)
	}
	if val := q.Pop().(int); val != 4 {
		t.Fatalf("expected 4, got %v", val)
	}

	if val := q.Len(); val != 0 {
		t.Errorf("expected 0, got %v", val)
	}
	if val := q.Pop(); val != nil {
		t.Fatalf("expected nil, got %v", val)
	}
}
