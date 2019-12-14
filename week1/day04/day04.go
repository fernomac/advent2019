package main

import "fmt"

const (
	min = 197487
	max = 673251
)

type num struct {
	digits [6]int
}

func newnum(n int) *num {
	if n > 999999 {
		panic("too big")
	}

	ret := &num{}
	for i := 5; i >= 0; i-- {
		ret.digits[i] = n % 10
		n /= 10
	}

	return ret
}

func (n *num) test1() bool {
	pair := false

	for i := 1; i < 6; i++ {
		if n.digits[i-1] == n.digits[i] {
			pair = true
		} else if n.digits[i-1] > n.digits[i] {
			return false
		}
	}

	return pair
}

type state uint8

const (
	base state = iota
	maybepair
	toomany
	pair
)

func (s state) update(match bool) state {
	switch s {
	case base:
		if match {
			// It might be a pair.
			return maybepair
		}
		// Nothing to see here.
		return base

	case maybepair:
		if match {
			// It's at least a triple.
			return toomany
		}
		// It's a real pair.
		return pair

	case toomany:
		if match {
			// Still too many.
			return toomany
		}
		// Reset.
		return base

	case pair:
		// Once we've found a pair we're chillin'.
		return pair

	default:
		panic(fmt.Sprintf("bad state %v", s))
	}
}

func (n *num) test2() bool {
	state := base

	for i := 1; i < 6; i++ {
		if n.digits[i-1] > n.digits[i] {
			return false
		}

		match := n.digits[i-1] == n.digits[i]
		state = state.update(match)
	}

	return (state == pair || state == maybepair)
}

func main() {
	count1 := 0
	count2 := 0

	for n := min; n <= max; n++ {
		num := newnum(n)
		if num.test1() {
			count1++
		}
		if num.test2() {
			count2++
		}
	}

	fmt.Println(count1)
	fmt.Println(count2)
}
