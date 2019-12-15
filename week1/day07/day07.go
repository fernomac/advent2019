package main

import (
	"fmt"

	"github.com/fernomac/advent2019/lib"
)

type circuit struct {
	prog lib.Program
}

func (c *circuit) run(phases []int) int {
	chans := []chan int{}
	for i := 0; i < len(phases)+1; i++ {
		chans = append(chans, make(chan int, 2))
	}

	// Initialize the 'puters.
	puters := []*lib.Puter{}
	for i, phase := range phases {
		p := lib.NewPuter(c.prog)
		puters = append(puters, p)

		in := chans[i]
		out := chans[i+1]

		in <- phase
		p.StdinCh(in)
		p.StdoutCh(out)
	}

	// Provide the initial input.
	in := chans[0]
	in <- 0

	out := chans[len(phases)]

	for {
		done := false

		// Run all the 'puters forward one iteration. They should consume one
		// input, produce one output, and either terminate or block waiting
		// for another input.
		for _, p := range puters {
			done = p.RunNB()
		}

		val := <-out
		if done {
			return val
		}

		in <- val
	}
}

func permute(a []int, f func([]int)) {
	doPermute(len(a), a, f)
}

func doPermute(k int, a []int, f func([]int)) {
	if k == 1 {
		f(a)
		return
	}

	doPermute(k-1, a, f)

	for i := 0; i < k-1; i++ {
		if k&1 == 0 {
			a[i], a[k-1] = a[k-1], a[i]
		} else {
			a[0], a[k-1] = a[k-1], a[0]
		}

		doPermute(k-1, a, f)
	}
}

func (c *circuit) optimize(options []int) (int, []int) {
	max := 0
	maxphases := []int{}

	permute(options, func(phases []int) {
		val := c.run(phases)
		if val > max {
			max = val
			maxphases = append([]int{}, phases...)
		}
	})

	return max, maxphases
}

func (c *circuit) optimize1() (int, []int) {
	return c.optimize([]int{0, 1, 2, 3, 4})
}

func (c *circuit) optimize2() (int, []int) {
	return c.optimize([]int{5, 6, 7, 8, 9})
}

func main() {
	c := circuit{lib.LoadProgram("input.txt")}
	{
		max, phases := c.optimize1()
		fmt.Println(max, phases)
	}

	{
		max, phases := c.optimize2()
		fmt.Println(max, phases)
	}
}
