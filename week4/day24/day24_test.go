package main

import "testing"

const step1 = `
....#
#..#.
#..##
..#..
#....`

const step2 = `
#..#.
####.
###.#
##.##
.##..`

const step3 = `
#####
....#
....#
...#.
#.###`

const stepN = `
.....
.....
.....
#....
.#...`

func TestParse(t *testing.T) {
	s := parse(step1)
	e := state(0b00001_00100_11001_01001_10000)
	if s != e {
		t.Errorf("expected %v, got %v", int(s), int(e))
	}

	if s.get(0, 0) {
		t.Error("expected empty at (0,0)")
	}
	if !s.get(0, 4) {
		t.Error("expected bug at (4,0)")
	}
	if !s.get(2, 3) {
		t.Error("expected bug at (2,3)")
	}
	if s.get(3, 3) {
		t.Error("expected empty at (3,3)")
	}

	if s.get(-1, -1) {
		t.Error("expected empty at (-1, -1)")
	}
	if s.get(5, 5) {
		t.Error("expected empty at (5,5)")
	}
}

func TestStep(t *testing.T) {
	s1 := parse(step1)
	s2 := parse(step2)
	s3 := parse(step3)

	test := s1.step()
	if test != s2 {
		t.Errorf("expected %v, got %v", s2, test)
	}

	test = test.step()
	if test != s3 {
		t.Errorf("expected %v, got %v", s3, test)
	}
}

func TestBiodiversity(t *testing.T) {
	sn := parse(stepN)
	bio := sn.biodiversity()
	if bio != 2129920 {
		t.Errorf("expected 2129920, got %v", bio)
	}
}

func TestFindRepeat(t *testing.T) {
	rep := findRepeat(parse(step1))
	sn := parse(stepN)

	if rep != sn {
		t.Errorf("expected %v, got %v", sn, rep)
	}
}

func TestMultiverse(t *testing.T) {
	m := multiversify(parse(step1))
	for i := 0; i < 10; i++ {
		m = m.step()
	}

	if c := m.count(); c != 99 {
		t.Errorf("expected 99, got %v", c)
	}
}
