package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fernomac/advent2019/lib"
)

type moon1d struct {
	pos, vel int
}

type system1d struct {
	moons []moon1d
}

func (s *system1d) gravitize() {
	for i := range s.moons {
		m1 := &s.moons[i]

		for j := i + 1; j < len(s.moons); j++ {
			m2 := &s.moons[j]

			if m1.pos < m2.pos {
				m1.vel++
				m2.vel--
			} else if m1.pos > m2.pos {
				m1.vel--
				m2.vel++
			}
		}
	}
}

func (s *system1d) velocitize() {
	for i := range s.moons {
		m := &s.moons[i]
		m.pos += m.vel
	}
}

func (s *system1d) step() {
	s.gravitize()
	s.velocitize()
}

func (s *system1d) stringify() string {
	b := strings.Builder{}

	for i := range s.moons {
		b.WriteString(strconv.Itoa(s.moons[i].pos))
		b.WriteRune(':')
		b.WriteString(strconv.Itoa(s.moons[i].vel))
		b.WriteRune(',')
	}

	return b.String()
}

func (s *system1d) findloop() int {
	n := 0
	init := s.stringify()

	for {
		s.step()
		n++

		state := s.stringify()
		if state == init {
			return n
		}
	}
}

type system struct {
	x, y, z system1d
}

func (s *system) step() {
	s.x.step()
	s.y.step()
	s.z.step()
}

func (s *system) stepN(n int) {
	for i := 0; i < n; i++ {
		s.step()
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func absum(a, b, c int) int {
	return abs(a) + abs(b) + abs(c)
}

func (s *system) energy() int {
	total := 0

	for i := range s.x.moons {
		pot := absum(s.x.moons[i].pos, s.y.moons[i].pos, s.z.moons[i].pos)
		kin := absum(s.x.moons[i].vel, s.y.moons[i].vel, s.z.moons[i].vel)
		total += pot * kin
	}

	return total
}

func (s *system) findloop() int {
	xp := s.x.findloop()
	yp := s.y.findloop()
	zp := s.z.findloop()

	return lib.Lcm(xp, yp, zp)
}

type moon struct {
	x, y, z int
}

func newsys(moons ...moon) *system {
	sys := &system{}

	for _, moon := range moons {
		sys.x.moons = append(sys.x.moons, moon1d{pos: moon.x})
		sys.y.moons = append(sys.y.moons, moon1d{pos: moon.y})
		sys.z.moons = append(sys.z.moons, moon1d{pos: moon.z})
	}

	return sys
}

func main() {
	input := []moon{
		moon{8, 0, 8},     // <x=8, y=0, z=8>
		moon{0, -5, -10},  // <x=0, y=-5, z=-10>
		moon{16, 10, -5},  // <x=16, y=10, z=-5>
		moon{19, -10, -7}, // <x=19, y=-10, z=-7>
	}

	{
		s := newsys(input...)
		s.stepN(1000)
		fmt.Println("energy:", s.energy())
	}

	{
		s := newsys(input...)
		fmt.Println("period:", s.findloop())
	}
}
