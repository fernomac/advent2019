package main

import (
	"fmt"
	"strings"

	"github.com/fernomac/advent2019/lib"
)

type body struct {
	name   string
	parent string
}

func parseBody(b string) *body {
	parts := strings.Split(b, ")")
	if len(parts) != 2 {
		panic(b)
	}
	return &body{
		name:   parts[1],
		parent: parts[0],
	}
}

type lineage []string

func (l lineage) indexOf(s string) int {
	for i, b := range l {
		if b == s {
			return i
		}
	}
	return -1
}

type system struct {
	bodies map[string]*body
}

func newSystem(bodies ...string) *system {
	sys := &system{map[string]*body{}}
	for _, body := range bodies {
		sys.add(parseBody(body))
	}
	return sys
}

func (s *system) add(b *body) {
	s.bodies[b.name] = b
}

func (s *system) orbitsFor(body string) int {
	if body == "COM" {
		return 0
	}
	return s.orbitsFor(s.bodies[body].parent) + 1
}

func (s *system) orbits() int {
	total := 0
	for body := range s.bodies {
		total += s.orbitsFor(body)
	}
	return total
}

func (s *system) lineageFor(body string) lineage {
	if body == "COM" {
		return []string{}
	}

	parent := s.bodies[body].parent
	return append(s.lineageFor(parent), parent)
}

func (s *system) path(from, to string) []string {
	path := []string{}
	if from == to {
		return path
	}

	cur := s.bodies[from].parent
	lineage := s.lineageFor(to)
	index := lineage.indexOf(cur)

	// Up the nearest ancestor.
	for index == -1 {
		parent := s.bodies[cur].parent
		path = append(path, fmt.Sprintf("%v->%v", cur, parent))
		cur = parent

		index = lineage.indexOf(cur)
	}

	// Back down to the target.
	for i := index + 1; i < len(lineage); i++ {
		next := lineage[i]
		path = append(path, fmt.Sprintf("%v->%v", cur, next))
		cur = next
	}

	return path
}

func main() {
	s := newSystem()
	lib.ForLines("input.txt", func(l string) {
		s.add(parseBody(l))
	})

	fmt.Println("checksum:", s.orbits())

	path := s.path("YOU", "SAN")
	fmt.Println("len(path): ", len(path))
}
