package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/fernomac/advent2019/lib"
)

type point struct {
	x, y int
}

type visitor interface {
	visit(x, y int)
}

type mapper struct {
	points map[point]int
	dist   int
}

func newMapper() *mapper {
	return &mapper{
		points: map[point]int{},
		dist:   0,
	}
}

func (m *mapper) visit(x, y int) {
	m.dist++
	if _, ok := m.points[point{x, y}]; !ok {
		m.points[point{x, y}] = m.dist
	}
}

type collider struct {
	points map[point]int
	dist   int

	nydist   int
	linedist int
}

func newCollider(points map[point]int) *collider {
	return &collider{
		points:   points,
		dist:     0,
		nydist:   math.MaxInt32,
		linedist: math.MaxInt32,
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (c *collider) visit(x, y int) {
	c.dist++
	if dist, ok := c.points[point{x, y}]; ok {
		// There's a collision at (x,y)

		nydist := abs(x) + abs(y)
		linedist := dist + c.dist

		fmt.Printf("collision at (%v, %v): ny=%v, line=%v\n", x, y, nydist, linedist)

		if nydist < c.nydist {
			c.nydist = nydist
		}
		if linedist < c.linedist {
			c.linedist = linedist
		}
	}
}

func trace(str string, v visitor) {
	steps := strings.Split(str, ",")

	x := 0
	y := 0

	for _, step := range steps {
		dir := step[0]
		len, err := strconv.Atoi(step[1:])
		if err != nil {
			panic(err)
		}

		dx := 0
		dy := 0

		switch dir {
		case 'L':
			dx = -1
		case 'R':
			dx = 1
		case 'U':
			dy = -1
		case 'D':
			dy = 1
		default:
			panic(dir)
		}

		for i := 0; i < len; i++ {
			x += dx
			y += dy
			v.visit(x, y)
		}
	}
}

func collide(l1, l2 string) (int, int) {
	m := newMapper()
	trace(l1, m)

	c := newCollider(m.points)
	trace(l2, c)

	return c.nydist, c.linedist
}

func main() {
	file := lib.Load("input.txt")
	lines := strings.Split(file, "\n")
	if len(lines) != 2 {
		panic(len(lines))
	}

	nydist, linedist := collide(lines[0], lines[1])

	fmt.Println("ny:", nydist, ", line:", linedist)
}
