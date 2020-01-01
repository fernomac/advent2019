package main

import (
	"fmt"

	"github.com/fernomac/advent2019/lib"
)

type point struct {
	x, y int
}

type tbeam struct {
	prog lib.Program
}

func (t *tbeam) on(at point) bool {
	if at.x < 0 || at.y < 0 {
		return false
	}

	p := lib.NewPuter(t.prog)
	p.Stdin([]int{at.x, at.y})
	p.Run()
	val := (p.Stdout()[0] == 1)

	return val
}

func (t *tbeam) count(max point) int {
	count := 0

	for y := 0; y < max.y; y++ {
		for x := 0; x < max.x; x++ {
			if t.on(point{x, y}) {
				count++
			}
		}
	}

	return count
}

func (t *tbeam) fit(size point) point {
	left := 0
	right := 0

	for y := size.y; ; y++ {
		// Find the left edge of the beam.
		for !t.on(point{left, y}) {
			left++
		}

		// Find the right edge of the beam.
		if right < left {
			right = left
		}
		for t.on(point{right + 1, y}) {
			right++
		}

		// Check if we've got a hit.
		bottom := y + (size.y - 1)
		x := right - (size.x - 1)

		if t.on(point{x, bottom}) {
			// Check if we can slide it left at all.
			for x > left && t.on(point{x - 1, bottom}) {
				x--
			}
			return point{x, y}
		}
	}
}

func main() {
	t := tbeam{
		prog: lib.LoadProgram("input.txt"),
	}
	fmt.Println("count:", t.count(point{50, 50}))
	fmt.Println("point:", t.fit(point{100, 100}))
}
