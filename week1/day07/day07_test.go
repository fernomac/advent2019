package main

import (
	"fmt"
	"testing"

	"github.com/fernomac/advent2019/lib"
)

var (
	test1 = []int{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0}
	test2 = []int{3, 23, 3, 24, 1002, 24, 10, 24, 1002, 23, -1, 23, 101, 5, 23, 23, 1, 24, 23, 23, 4, 23, 99, 0, 0}
	test3 = []int{3, 31, 3, 32, 1002, 32, 10, 32, 1001, 31, -2, 31, 1007, 31, 0, 33, 1002, 33, 7, 33, 1, 33, 31, 31, 1, 32, 31, 31, 4, 31, 99, 0, 0, 0}

	test4 = []int{3, 26, 1001, 26, -4, 26, 3, 27, 1002, 27, 2, 27, 1, 27, 26, 27, 4, 27, 1001, 28, -1, 28, 1005, 28, 6, 99, 0, 0, 5}
	test5 = []int{3, 52, 1001, 52, -5, 52, 3, 53, 1, 52, 56, 54, 1007, 54, 5, 55, 1005, 55, 26, 1001, 54, -5, 54, 1105, 1, 12, 1, 53, 54, 53, 1008, 54, 0, 55, 1001, 55, 1, 55, 2, 53, 55, 53, 4, 53, 1001, 56, -1, 56, 1005, 56, 6, 99, 0, 0, 0, 0, 10}
)

func TestRun(t *testing.T) {
	test := func(prog []int, phases []int, eval int) {
		t.Run(fmt.Sprintf("%v", prog), func(t *testing.T) {
			c := circuit{prog}
			val := c.run(phases)
			if val != eval {
				t.Fatalf("expected %v, got %v", eval, val)
			}
		})
	}

	test(test1, []int{4, 3, 2, 1, 0}, 43210)
	test(test2, []int{0, 1, 2, 3, 4}, 54321)
	test(test3, []int{1, 0, 4, 3, 2}, 65210)
}

func TestOptimize1(t *testing.T) {
	test := func(prog []int, eval int, ephases []int) {
		t.Run(fmt.Sprintf("%v", prog), func(t *testing.T) {
			c := circuit{prog}
			val, phases := c.optimize1()
			if val != eval {
				t.Errorf("expected %v, got %v", eval, val)
			}
			if !lib.IntEq(phases, ephases) {
				t.Errorf("expected %v, got %v", ephases, phases)
			}
		})
	}

	test(test1, 43210, []int{4, 3, 2, 1, 0})
	test(test2, 54321, []int{0, 1, 2, 3, 4})
	test(test3, 65210, []int{1, 0, 4, 3, 2})
}

func TestOptimize2(t *testing.T) {
	test := func(prog []int, eval int, ephases []int) {
		t.Run(fmt.Sprintf("%v", prog), func(t *testing.T) {
			c := circuit{prog}
			val, phases := c.optimize2()
			if val != eval {
				t.Errorf("expected %v, got %v", eval, val)
			}
			if !lib.IntEq(phases, ephases) {
				t.Errorf("expected %v, got %v", ephases, phases)
			}
		})
	}

	test(test4, 139629729, []int{9, 8, 7, 6, 5})
	test(test5, 18216, []int{9, 7, 8, 5, 6})
}

func TestPermute(t *testing.T) {
	n := 0

	permute([]int{0, 1, 2, 3, 4}, func(a []int) {
		n++
		// fmt.Println(a)
	})

	// factorial(5)
	if n != 120 {
		t.Fatalf("expected %v, got %v", 120, n)
	}
}
