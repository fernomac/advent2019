package main

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/fernomac/advent2019/lib"
)

type point struct {
	x, y int
}

type try struct {
	roids   map[point]bool
	size    point
	station point
}

func (t *try) stepFrom(x, y int) (int, int) {
	dx := t.station.x - x
	dy := t.station.y - y

	div := lib.Gcd(dx, dy)
	if div == 0 {
		return 0, 0
	}

	return dx / div, dy / div
}

func (t *try) trace(x, y int) bool {
	ret := false

	pt := point{x, y}
	dx, dy := t.stepFrom(x, y)

	for pt.x != t.station.x || pt.y != t.station.y {
		if t.roids[pt] {
			ret = true
			delete(t.roids, pt)
		}

		pt.x += dx
		pt.y += dy
	}

	return ret
}

func (t *try) traceline(y int) int {
	count := 0
	for x := 0; x < t.size.x; x++ {
		if t.trace(x, y) {
			count++
		}
	}
	return count
}

func (t *try) count() int {
	count := 0

	for y := 0; y < t.station.y; y++ {
		count += t.traceline(y)
	}
	for y := t.size.y - 1; y > t.station.y; y-- {
		count += t.traceline(y)
	}

	if t.trace(0, t.station.y) {
		count++
	}
	if t.trace(t.size.x-1, t.station.y) {
		count++
	}

	return count
}

type space struct {
	roids map[point]bool
	size  point
}

func parse(str string) space {
	lines := strings.Split(strings.TrimSpace(str), "\n")

	roids := map[point]bool{}
	size := point{len(lines[0]), len(lines)}

	for y, line := range lines {
		for x, char := range line {
			switch char {
			case '#':
				roids[point{x, y}] = true
			}
		}
	}

	return space{roids, size}
}

func (s *space) try(x, y int) int {
	roids := map[point]bool{}
	for pt := range s.roids {
		roids[pt] = true
	}

	t := try{
		roids:   roids,
		size:    s.size,
		station: point{x, y},
	}

	return t.count()
}

func (s *space) optimize() (int, point) {
	max := -1
	maxpt := point{}

	for pt := range s.roids {
		val := s.try(pt.x, pt.y)
		if val > max {
			max = val
			maxpt = pt
		}
	}

	return max, maxpt
}

type target struct {
	point    point
	angle    float64
	distance int
}

func (s *space) targets(sx, sy int) []target {
	ret := []target{}

	for pt := range s.roids {
		dx := pt.x - sx
		dy := sy - pt.y

		angle := math.Atan2(float64(dx), float64(dy))
		if angle < 0 {
			angle += (math.Pi * 2)
		}

		ret = append(ret, target{pt, angle, dx + dy})
	}

	return ret
}

func (s *space) asplode(n int, targets []target) target {
	sort.Slice(targets, func(i, j int) bool {
		return targets[i].angle < targets[j].angle
	})

	i := 0
	q := []target{}

	for {
		t := targets[0]
		i++

		fmt.Printf("%v: asploding %v\n", i, t.point)
		if i == n {
			return t
		}

		j := 1
		for ; targets[j].angle == t.angle; j++ {
			// Keep going.
		}

		if j > 1 {
			q = append(q, targets[1:j]...)
		}
		targets = targets[j:]

		if len(targets) == 0 {
			targets = q
			q = []target{}
		}
	}
}

func main() {
	space := parse(lib.Load("input.txt"))
	score, where := space.optimize()
	fmt.Println(score, ":", where)

	targets := space.targets(where.x, where.y)
	space.asplode(200, targets)
}
