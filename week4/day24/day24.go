package main

import (
	"fmt"
	"strings"

	"github.com/fernomac/advent2019/lib"
)

type state int

func parse(in string) state {
	lines := strings.Split(strings.TrimSpace(in), "\n")
	if len(lines) != 5 {
		panic(fmt.Sprint("too many lines of input:", lines))
	}

	ret := state(0)

	for y, line := range lines {
		if len(line) != 5 {
			panic(fmt.Sprint("line too long:", line))
		}

		for x, char := range line {
			switch char {
			case '#':
				ret = ret.set(x, y)

			case '.':
				// Nothing happens.

			default:
				panic(fmt.Sprint("line malformed:", line))
			}
		}
	}

	return ret
}

func (s state) get(x, y int) bool {
	if x < 0 || x >= 5 || y < 0 || y >= 5 {
		return false
	}
	return (s & bit(x, y)) != 0
}

func (s state) set(x, y int) state {
	if x < 0 || x >= 5 || y < 0 || y >= 5 {
		panic(fmt.Sprintf("invalid coordinates (%v, %v)", x, y))
	}
	return s | bit(x, y)
}

func bit(x, y int) state {
	return 1 << ((y * 5) + x)
}

func (s state) step() state {
	next := state(0)

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if stepone(s.get(x, y), s.neighbors(x, y)) {
				next = next.set(x, y)
			}
		}
	}

	return next
}

func stepone(cur bool, neighbors int) bool {
	if cur {
		// A bug dies (becoming an empty space) unless there is exactly one bug adjacent to it.
		return neighbors == 1
	}

	// An empty space becomes infested with a bug if exactly one or two bugs are adjacent to it.
	return neighbors == 1 || neighbors == 2
}

func (s state) neighbors(x, y int) int {
	count := 0
	if s.get(x, y-1) {
		count++
	}
	if s.get(x, y+1) {
		count++
	}
	if s.get(x-1, y) {
		count++
	}
	if s.get(x+1, y) {
		count++
	}
	return count
}

func (s state) biodiversity() int {
	return int(s)
}

func (s state) print() {
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if s.get(x, y) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func (s state) count() int {
	sum := 0

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if s.get(x, y) {
				sum++
			}
		}
	}

	return sum
}

func findRepeat(init state) state {
	visited := map[state]bool{
		init: true,
	}
	s := init.step()

	for {
		if visited[s] {
			return s
		}
		visited[s] = true

		s = s.step()
	}
}

type multiverse struct {
	levels   map[int]state
	min, max int
}

func multiversify(s state) multiverse {
	return multiverse{
		levels: map[int]state{0: s},
	}
}

func (m multiverse) get(x, y, z int) bool {
	return m.levels[z].get(x, y)
}

func (m multiverse) getint(x, y, z int) int {
	if m.get(x, y, z) {
		return 1
	}
	return 0
}

func (m multiverse) step() multiverse {
	next := multiverse{
		levels: map[int]state{},
	}

	for z := m.min - 1; z <= m.max+1; z++ {
		nz := m.steplevel(z)
		if nz == 0 {
			continue
		}

		next.levels[z] = nz
		if z < next.min {
			next.min = z
		}
		if z > next.max {
			next.max = z
		}
	}

	return next
}

func (m multiverse) steplevel(z int) state {
	next := state(0)

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if x == 2 && y == 2 {
				continue
			}
			if stepone(m.get(x, y, z), m.neighbors(x, y, z)) {
				next = next.set(x, y)
			}
		}
	}

	return next
}

func (m multiverse) neighbors(x, y, z int) int {
	return m.left(x, y, z) +
		m.right(x, y, z) +
		m.up(x, y, z) +
		m.down(x, y, z)
}

func (m multiverse) left(x, y, z int) int {
	if x == 0 {
		// Off the outer-left edge; tile (1,2) of the enclosing.
		return m.getint(1, 2, z-1)
	}
	if x == 3 && y == 2 {
		// Off the inner-left edge; tiles (4,0-4) of the enclosed.
		sum := 0
		for i := 0; i < 5; i++ {
			if m.get(4, i, z+1) {
				sum++
			}
		}
		return sum
	}
	// A regular neighbor on this level.
	return m.getint(x-1, y, z)
}

func (m multiverse) right(x, y, z int) int {
	if x == 4 {
		// Off the outer-right edge; tile (3,2) of the enclosing.
		return m.getint(3, 2, z-1)
	}
	if x == 1 && y == 2 {
		// Off the inner-right edge; tiles (0,0-4) of the enclosed.
		sum := 0
		for i := 0; i < 5; i++ {
			if m.get(0, i, z+1) {
				sum++
			}
		}
		return sum
	}
	// A regular neighbor on this level.
	return m.getint(x+1, y, z)
}

func (m multiverse) up(x, y, z int) int {
	if y == 0 {
		// Off the outer-top edge; tile (2, 1) of the enclosing.
		return m.getint(2, 1, z-1)
	}
	if x == 2 && y == 3 {
		// Off the inner-top edge; tiles (0-4, 4) of the enclosed.
		sum := 0
		for i := 0; i < 5; i++ {
			if m.get(i, 4, z+1) {
				sum++
			}
		}
		return sum
	}
	// A regular neighbor on this level.
	return m.getint(x, y-1, z)
}

func (m multiverse) down(x, y, z int) int {
	if y == 4 {
		// Off the outer-bottom edge; tile (2,3) of the enclosing.
		return m.getint(2, 3, z-1)
	}
	if x == 2 && y == 1 {
		// Off the inner-bottom edge; tiles (0-4, 0) of the enclosed.
		sum := 0
		for i := 0; i < 5; i++ {
			if m.get(i, 0, z+1) {
				sum++
			}
		}
		return sum
	}
	// A regular neighbor on this level.
	return m.getint(x, y+1, z)
}

func (m multiverse) count() int {
	sum := 0
	for i := m.min; i <= m.max; i++ {
		sum += m.levels[i].count()
	}
	return sum
}

func (m multiverse) print() {
	for i := m.min; i <= m.max; i++ {
		fmt.Println("Depth", i)
		m.levels[i].print()
		fmt.Println()
	}
}

func main() {
	init := parse(lib.Load("input.txt"))

	{
		repeat := findRepeat(init)
		fmt.Println(repeat.biodiversity())
	}

	{
		m := multiversify(init)
		for i := 0; i < 200; i++ {
			m = m.step()
		}
		fmt.Println(m.count())
	}
}
