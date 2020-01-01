package main

import (
	"fmt"
	"testing"
)

func TestStepFrom(t *testing.T) {
	try := try{
		size:    point{100, 100},
		station: point{4, 4},
	}
	test := func(x, y, edx, edy int) {
		t.Run(fmt.Sprintf("(%v,%v)", x, y), func(t *testing.T) {
			dx, dy := try.stepFrom(x, y)
			if dx != edx {
				t.Errorf("expected dx=%v, got %v", edx, dx)
			}
			if dy != edy {
				t.Errorf("expected dy=%v, got %v", edy, dy)
			}
		})
	}

	test(0, 0, 1, 1)
	test(1, 0, 3, 4)
	test(2, 0, 1, 2)
	test(3, 0, 1, 4)
	test(4, 0, 0, 1)
	test(5, 0, -1, 4)
	test(6, 0, -1, 2)
	test(7, 0, -3, 4)
	test(8, 0, -1, 1)
}

func TestTrace(t *testing.T) {
	try := try{
		roids: map[point]bool{
			point{1, 0}: true,
			point{2, 2}: true,
			point{3, 4}: true,
			point{4, 4}: true,
		},
		size:    point{5, 5},
		station: point{3, 4},
	}

	if val := try.trace(1, 0); !val {
		t.Error("expected trace(1,0)=true")
	}
	if val := try.trace(2, 2); val {
		t.Error("expected trace(2,2)=false")
	}
	if val := try.trace(3, 4); val {
		t.Error("expected trace(3,4)=false")
	}
	if val := try.trace(4, 4); !val {
		t.Error("expected trace(4,4)=true")
	}
}

func TestTry(t *testing.T) {
	test := func(str string, x, y int, ecount int) {
		t.Run(fmt.Sprintf("%v:(%v,%v)", str, x, y), func(t *testing.T) {
			space := parse(str)
			count := space.try(x, y)
			if count != ecount {
				t.Errorf("expected %v, got %v", ecount, count)
			}
		})
	}

	test(`
.#..#
.....
#####
....#
...##`, 3, 4, 8)

	test(`
......#.#.
#..#.#....
..#######.
.#.#.###..
.#..#.....
..#....#.#
#..#....#.
.##.#..###
##...#..#.
.#....####`, 5, 8, 33)

	test(`
#.#...#.#.
.###....#.
.#....#...
##.#.#.#.#
....#.#.#.
.##..###.#
..#...##..
..##....##
......#...
.####.###.`, 1, 2, 35)

	test(`
.#..#..###
####.###.#
....###.#.
..###.##.#
##.##.#.#.
....###..#
..#.#..#.#
#..#.#.###
.##...##.#
.....#.#..`, 6, 3, 41)

	test(`
.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##`, 11, 13, 210)
}

func TestOptimize(t *testing.T) {
	test := func(str string, ecount int, ewhere point) {
		t.Run(str, func(t *testing.T) {
			space := parse(str)
			count, where := space.optimize()
			if count != ecount {
				t.Errorf("expected %v, got %v", ecount, count)
			}
			if where != ewhere {
				t.Errorf("expected %v, got %v", ewhere, where)
			}
		})
	}

	test(`
.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##`, 210, point{11, 13})
}
