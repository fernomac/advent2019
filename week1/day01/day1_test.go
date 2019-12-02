package main

import (
	"fmt"
	"testing"
)

func TestFuelFor(t *testing.T) {
	test := func(m, ef int) {
		t.Run(fmt.Sprintf("%v", m), func(t *testing.T) {
			f := fuelFor(m)
			if f != ef {
				t.Errorf("expected %v, got %v", ef, f)
			}
		})
	}

	test(-1, 0)
	test(0, 0)
	test(1, 0)
	test(8, 0)
	test(9, 1)
	test(14, 2)
	test(1969, 966)
	test(100756, 50346)
}
