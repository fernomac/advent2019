package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fernomac/advent2019/lib"
)

type point struct {
	x, y int
}

type tile uint8

const (
	empty tile = iota
	grid
	intersection
	robot
)

func parse(in []int) map[point]tile {
	x := 0
	y := 0

	out := map[point]tile{}

	for _, i := range in {
		switch i {
		case '^':
			out[point{x, y}] = robot
			x++

		case '#':
			out[point{x, y}] = grid
			x++

		case '.':
			x++

		case '\n':
			x = 0
			y++
		}
	}

	return out
}

func move(pt, dir point) point {
	return point{pt.x + dir.x, pt.y + dir.y}
}

func notempty(m map[point]tile, pt point) bool {
	return m[pt] != empty
}

func intersect(m map[point]tile) {
	for pt := range m {
		if notempty(m, move(pt, point{1, 0})) &&
			notempty(m, move(pt, point{-1, 0})) &&
			notempty(m, move(pt, point{0, 1})) &&
			notempty(m, move(pt, point{0, -1})) {

			m[pt] = intersection
		}
	}
}

func align(m map[point]tile) int {
	sum := 0
	for pt, v := range m {
		if v == intersection {
			sum += pt.x * pt.y
		}
	}
	return sum
}

func draw(m map[point]tile) {
	size := size(m)
	for y := 0; y <= size.y; y++ {
		for x := 0; x <= size.x; x++ {
			switch m[point{x, y}] {
			case empty:
				fmt.Print(".")
			case grid:
				fmt.Print("#")
			case intersection:
				fmt.Print("O")
			case robot:
				fmt.Print("^")
			}
		}
		fmt.Println()
	}
}

func size(m map[point]tile) point {
	x, y := 0, 0
	for pt := range m {
		if pt.x > x {
			x = pt.x
		}
		if pt.y > y {
			y = pt.y
		}
	}
	return point{x, y}
}

func solve(m map[point]tile) []string {
	cmds := []string{}

	pos := findrobot(m)
	dir := point{0, -1}
	delete(m, pos)

	moves := 0

	for len(m) > 0 {
		// If we can move forward, do it.
		next := move(pos, dir)
		if notempty(m, next) {
			switch m[next] {
			case grid:
				delete(m, next)
			case intersection:
				m[next] = grid
			default:
				panic("wat?")
			}

			moves++
			pos = next
			continue
		}

		// Done moving forward.
		if moves > 0 {
			cmds = append(cmds, strconv.Itoa(moves))
			moves = 0
		}

		// Time to turn.
		if notempty(m, move(pos, left(dir))) {
			cmds = append(cmds, "L")
			dir = left(dir)
			continue
		}

		if notempty(m, move(pos, right(dir))) {
			cmds = append(cmds, "R")
			dir = right(dir)
			continue
		}

		panic("nowhere to turn")
	}

	if moves > 0 {
		cmds = append(cmds, strconv.Itoa(moves))
	}

	return cmds
}

func findrobot(m map[point]tile) point {
	for pt, v := range m {
		if v == robot {
			return pt
		}
	}
	panic("no robot")
}

func left(pt point) point {
	return point{pt.y, -pt.x}
}

func right(pt point) point {
	return point{-pt.y, pt.x}
}

func fold(cmds []string) string {
	b := strings.Builder{}

	for _, cmd := range cmds {
		if b.Len() > 0 {
			b.WriteByte(',')
		}
		b.WriteString(cmd)
	}

	return b.String()
}

func asciiprint(out []int) int {
	for _, c := range out {
		if c > 0xFF {
			return c
		}
		fmt.Printf("%c", byte(c))
	}
	return 0
}

func run(p *lib.Puter) {
	p.RunNB()
	out := p.Stdout()
	_ = out
	// asciiprint(out)
	p.DropStdout()
}

func send(ch chan int, str string) {
	for _, b := range str {
		ch <- int(b)
	}
}

func main() {
	var m map[point]tile

	{
		p := lib.NewPuter(lib.LoadProgram("input.txt"))
		p.Run()

		out := p.Stdout()

		m = parse(out)
		intersect(m)
		sum := align(m)
		fmt.Println("Sum:", sum)
	}

	{
		cmds := solve(m)
		fmt.Println("Solution:", fold(cmds))

		// I tried to do this programmatically and failed, but
		// got close enough to solve for my specific case manually.

		// ....................
		// A,B,A,B,A,C,B,C,A,C
		// L,10,L,12,R,6
		// R,10,L,4,L,4,L,12
		// L,10,R,10,R,6,L,4

		in := make(chan int, 20)

		prog := lib.LoadProgram("input.txt")
		prog[0] = 2
		p := lib.NewPuter(prog)
		p.StdinCh(in)

		run(p)
		send(in, "A,B,A,B,A,C,B,C,A,C\n")

		run(p)
		send(in, "L,10,L,12,R,6\n")

		run(p)
		send(in, "R,10,L,4,L,4,L,12\n")

		run(p)
		send(in, "L,10,R,10,R,6,L,4\n")

		run(p)
		send(in, "n\n")

		p.RunNB()
		out := p.Stdout()

		fmt.Println("Dust:", out[len(out)-1])
	}
}
