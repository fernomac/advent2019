package main

import (
	"fmt"
	"testing"
)

func TestPhysics(t *testing.T) {
	s := newsys(
		moon{-1, 0, 2},
		moon{2, -10, -7},
		moon{4, -8, 8},
		moon{3, 5, -1},
	)

	s.stepN(10)

	test := func(i int, pos []int, vel []int) {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			x, y, z := s.x.moons[i].pos, s.y.moons[i].pos, s.z.moons[i].pos
			dx, dy, dz := s.x.moons[i].vel, s.y.moons[i].vel, s.z.moons[i].vel

			if x != pos[0] || y != pos[1] || z != pos[2] {
				t.Errorf("expected %v, got [%v %v %v]", pos, x, y, z)
			}
			if dx != vel[0] || dy != vel[1] || dz != vel[2] {
				t.Errorf("expected %v, got [%v %v %v]", vel, dx, dy, dz)
			}
		})
	}

	test(0, []int{2, 1, -3}, []int{-3, -2, 1})
	test(1, []int{1, -8, 0}, []int{-1, 1, 3})
	test(2, []int{3, -6, 1}, []int{3, 2, -3})
	test(3, []int{2, 0, 4}, []int{1, -1, -1})
}

func TestEnergy(t *testing.T) {
	s := newsys(
		moon{-8, -10, 0},
		moon{5, 5, 10},
		moon{2, -7, 3},
		moon{9, -8, -3},
	)

	s.stepN(100)

	if e := s.energy(); e != 1940 {
		t.Errorf("expected system=1940, got %v", e)
	}
}

func TestFindLoop(t *testing.T) {
	{
		s := newsys(
			moon{-1, 0, 2},
			moon{2, -10, -7},
			moon{4, -8, 8},
			moon{3, 5, -1},
		)

		if n := s.findloop(); n != 2772 {
			t.Errorf("expected 2772, got %v", n)
		}
	}
	{
		s := newsys(
			moon{-8, -10, 0},
			moon{5, 5, 10},
			moon{2, -7, 3},
			moon{9, -8, -3},
		)

		if n := s.findloop(); n != 4686774924 {
			t.Errorf("expected 4686774924, got %v", n)
		}
	}
}
