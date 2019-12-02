package lib

import (
	"fmt"
	"testing"
)

func TestPuter(t *testing.T) {
	test := func(mem []int, eval []int) {
		t.Run(fmt.Sprintf("%v", mem), func(t *testing.T) {
			p := Puter{mem: mem}
			p.Run()

			for i := range eval {
				if p.mem[i] != eval[i] {
					t.Fatalf("expected %v, got %v", eval, p.mem)
				}
			}
		})
	}

	test([]int{1, 0, 0, 0, 99}, []int{2, 0, 0, 0, 99})
	test([]int{2, 3, 0, 3, 99}, []int{2, 3, 0, 6, 99})
	test([]int{2, 4, 4, 5, 99, 0}, []int{2, 4, 4, 5, 99, 9801})
	test([]int{1, 1, 1, 4, 99, 5, 6, 0, 99}, []int{30, 1, 1, 4, 2, 5, 6, 0, 99})
}
