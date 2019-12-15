package lib

import (
	"fmt"
	"testing"
)

func TestMath(t *testing.T) {
	test := func(mem []int, eval []int) {
		t.Run(fmt.Sprintf("%v", mem), func(t *testing.T) {
			p := NewPuter(mem)
			p.Run()

			if !IntEq(p.mem, eval) {
				t.Fatalf("expected %v, got %v", eval, p.mem)
			}
		})
	}

	test([]int{1, 0, 0, 0, 99}, []int{2, 0, 0, 0, 99})
	test([]int{2, 3, 0, 3, 99}, []int{2, 3, 0, 6, 99})
	test([]int{2, 4, 4, 5, 99, 0}, []int{2, 4, 4, 5, 99, 9801})
	test([]int{1, 1, 1, 4, 99, 5, 6, 0, 99}, []int{30, 1, 1, 4, 2, 5, 6, 0, 99})
}

func TestIO(t *testing.T) {
	testPuter(t, []int{3, 0, 4, 0, 99}, []int{42}, []int{42})
	testPuter(t, []int{104, 420, 99}, nil, []int{420})
	testPuter(t, []int{1002, 6, 3, 6, 4, 6, 33}, nil, []int{99})
	testPuter(t, []int{1101, 100, -1, 4, 0}, nil, nil)
}

func TestNBIO(t *testing.T) {
	p := NewPuter([]int{3, 0, 4, 0, 99})
	in := make(chan int, 1)
	p.StdinCh(in)

	if done := p.RunNB(); done {
		t.Fatalf("expected to block, but didn't")
	}

	in <- 666

	if done := p.RunNB(); !done {
		t.Fatalf("expected to finish, but didn't")
	}

	out := p.Stdout()
	if !IntEq(out, []int{666}) {
		t.Fatalf("expected [666], got %v", out)
	}
}

func TestLTEQ(t *testing.T) {
	testPuter(t, []int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8}, []int{7}, []int{1})
	testPuter(t, []int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8}, []int{8}, []int{0})
	testPuter(t, []int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8}, []int{9}, []int{0})

	testPuter(t, []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8}, []int{7}, []int{0})
	testPuter(t, []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8}, []int{8}, []int{1})
	testPuter(t, []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8}, []int{9}, []int{0})

	testPuter(t, []int{3, 3, 1107, -1, 8, 3, 4, 3, 99}, []int{7}, []int{1})
	testPuter(t, []int{3, 3, 1107, -1, 8, 3, 4, 3, 99}, []int{8}, []int{0})
	testPuter(t, []int{3, 3, 1107, -1, 8, 3, 4, 3, 99}, []int{9}, []int{0})

	testPuter(t, []int{3, 3, 1108, -1, 8, 3, 4, 3, 99}, []int{7}, []int{0})
	testPuter(t, []int{3, 3, 1108, -1, 8, 3, 4, 3, 99}, []int{8}, []int{1})
	testPuter(t, []int{3, 3, 1108, -1, 8, 3, 4, 3, 99}, []int{9}, []int{0})
}

func TestJumps(t *testing.T) {
	testPuter(t, []int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9}, []int{0}, []int{0})
	testPuter(t, []int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9}, []int{42}, []int{1})

	testPuter(t, []int{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1}, []int{0}, []int{0})
	testPuter(t, []int{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1}, []int{42}, []int{1})

	prog := []int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31, 1106, 0,
		36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104, 999, 1105, 1, 46, 1101, 1000,
		1, 20, 4, 20, 1105, 1, 46, 98, 99}

	testPuter(t, prog, []int{7}, []int{999})
	testPuter(t, prog, []int{8}, []int{1000})
	testPuter(t, prog, []int{9}, []int{1001})
}

func testPuter(t *testing.T, mem []int, stdin []int, eval []int) {
	t.Run(fmt.Sprintf("%v", mem), func(t *testing.T) {
		p := NewPuter(mem)
		p.Stdin(stdin)
		p.Run()

		val := p.Stdout()
		if !IntEq(val, eval) {
			t.Fatalf("expected %v, got %v", eval, val)
		}
	})
}
