package main

import (
	"testing"
)

var maze1 = `
#########
#b.A.@.a#
#########`

var maze2 = `
########################
#f.D.E.e.C.b.A.@.a.B.c.#
######################.#
#d.....................#
########################`

var maze2_1 = `
########################
#f.D.E.e.............@.#
######################.#
#d.....................#
########################`

var maze3 = `
########################
#...............b.C.D.f#
#.######################
#.....@.a.B.c.d.A.e.F.g#
########################`

var maze4 = `
#################
#i.G..c...e..H.p#
########.########
#j.A..b...f..D.o#
########@########
#k.E..a...g..B.n#
########.########
#l.F..d...h..C.m#
#################`

var maze5 = `
########################
#@..............ac.GI.b#
###d#e#f################
###A#B#C################
###g#h#i################
########################`

func TestMoves(t *testing.T) {
	test := func(name, maze string, from point, eval []move) {
		t.Run(name, func(t *testing.T) {
			g := newgrid(maze)
			val := g.moves(from)
			if !moveeq(val, eval) {
				t.Fatalf("expected %v, got %v", eval, val)
			}
		})
	}

	test("maze1", maze1, point{5, 1}, []move{
		move{loc: point{7, 1}, dist: 2, key: 1, doors: 0},
		move{loc: point{1, 1}, dist: 4, key: 2, doors: 1},
	})

	test("maze2", maze2, point{15, 1}, []move{
		move{loc: point{17, 1}, dist: 2, key: 1, doors: 0},   // a: []
		move{loc: point{11, 1}, dist: 4, key: 2, doors: 1},   // b: [A]
		move{loc: point{21, 1}, dist: 6, key: 4, doors: 2},   // c: [B]
		move{loc: point{7, 1}, dist: 8, key: 16, doors: 5},   // e: [C B]
		move{loc: point{1, 1}, dist: 14, key: 32, doors: 29}, // f: [E D C B]
		move{loc: point{1, 3}, dist: 30, key: 8, doors: 2},   // d: [B]
	})
}

func moveeq(a, b []move) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestTraverse(t *testing.T) {
	test := func(name, maze string, eval int) {
		t.Run(name, func(t *testing.T) {
			g := newgrid(maze)
			val := g.traverse()
			if val != eval {
				t.Fatalf("expected %v, got %v", eval, val)
			}
		})
	}

	test("maze1", maze1, 8)
	test("maze2", maze2, 86)
	test("maze3", maze3, 132)
	test("maze4", maze4, 136)
	test("maze5", maze5, 81)
}
