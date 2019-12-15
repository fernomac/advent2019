package main

import (
	"fmt"
	"testing"
)

func TestOrbits(t *testing.T) {
	s := newSystem(
		"COM)B",
		"B)C",
		"C)D",
		"D)E",
		"E)F",
		"B)G",
		"G)H",
		"D)I",
		"E)J",
		"J)K",
		"K)L")

	test := func(body string, eorbs int) {
		t.Run(body, func(t *testing.T) {
			orbs := s.orbitsFor(body)
			if orbs != eorbs {
				t.Fatalf("expected %v, got %v", eorbs, orbs)
			}
		})
	}

	test("D", 3)
	test("L", 7)

	total := s.orbits()
	if total != 42 {
		t.Fatalf("expected 42, got %v", total)
	}
}

func TestLineage(t *testing.T) {
	s := newSystem("COM)B",
		"B)C",
		"C)D",
		"D)E",
		"E)F",
		"B)G",
		"G)H",
		"D)I",
		"E)J",
		"J)K",
		"K)L",
		"K)YOU",
		"I)SAN")

	test := func(body string, elin []string) {
		t.Run(body, func(t *testing.T) {
			lin := s.lineageFor(body)
			if !eq(lin, elin) {
				t.Fatalf("expected %v, got %v", elin, lin)
			}
		})
	}

	test("SAN", []string{"COM", "B", "C", "D", "I"})
	test("YOU", []string{"COM", "B", "C", "D", "E", "J", "K"})
}

func TestPath(t *testing.T) {
	s := newSystem("COM)B",
		"B)C",
		"C)D",
		"D)E",
		"E)F",
		"B)G",
		"G)H",
		"D)I",
		"E)J",
		"J)K",
		"K)L",
		"K)YOU",
		"I)SAN")

	test := func(from, to string, epath []string) {
		t.Run(fmt.Sprintf("%v->%v", from, to), func(t *testing.T) {
			path := s.path(from, to)
			if !eq(path, epath) {
				t.Fatalf("expected %v, got %v", epath, path)
			}
		})
	}

	test("YOU", "SAN", []string{
		"K->J",
		"J->E",
		"E->D",
		"D->I",
	})

	test("E", "SAN", []string{
		"D->I",
	})

	test("F", "YOU", []string{
		"E->J",
		"J->K",
	})
}

func eq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
