package main

import (
	"fmt"

	"github.com/fernomac/advent2019/lib"
)

type point struct {
	x, y int
}

var leftturn = map[byte]byte{
	'^': '<',
	'<': 'v',
	'v': '>',
	'>': '^',
}

var rightturn = map[byte]byte{
	'^': '>',
	'>': 'v',
	'v': '<',
	'<': '^',
}

func turn(dir byte, lr int) byte {
	switch lr {
	case 0:
		return leftturn[dir]
	case 1:
		return rightturn[dir]
	default:
		panic(fmt.Sprintf("bad lr: %v", lr))
	}
}

func move(pos point, dir byte) point {
	switch dir {
	case '<':
		pos.x--
	case '>':
		pos.x++
	case '^':
		pos.y--
	case 'v':
		pos.y++
	default:
		panic(fmt.Sprintf("bad dir: %v", dir))
	}
	return pos
}

func paint(init int) map[point]int {
	in := make(chan int, 1)
	out := make(chan int, 2)

	p := lib.NewPuter(lib.LoadProgram("input.txt"))
	p.StdinCh(in)
	p.StdoutCh(out)

	pos := point{0, 0}
	dir := byte('^')

	grid := map[point]int{
		point{0, 0}: init,
	}

	for {
		in <- grid[pos]

		if p.RunNB() {
			return grid
		}

		select {
		case paint := <-out:
			grid[pos] = paint

			lr := <-out

			dir = turn(dir, lr)
			pos = move(pos, dir)

		default:
		}
	}
}

func minmax(grid map[point]int) (point, point) {
	min := point{0, 0}
	max := point{0, 0}

	for pt := range grid {
		if pt.x < min.x {
			min.x = pt.x
		}
		if pt.x > max.x {
			max.x = pt.x
		}
		if pt.y < min.y {
			min.y = pt.y
		}
		if pt.y > max.y {
			max.y = pt.y
		}
	}

	return min, max
}

func draw(grid map[point]int) {
	min, max := minmax(grid)

	for y := min.y; y <= max.y; y++ {
		for x := min.x; x <= max.x; x++ {
			if grid[point{x, y}] == 1 {
				fmt.Print("X")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func main() {
	grid1 := paint(0)
	fmt.Println(len(grid1))

	grid2 := paint(1)
	draw(grid2)
}
