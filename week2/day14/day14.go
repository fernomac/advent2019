package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fernomac/advent2019/lib"
)

type chem struct {
	name  string
	count int
}

func parsechem(str string) chem {
	parts := strings.Split(str, " ")
	if len(parts) != 2 {
		panic(str)
	}

	count, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}

	return chem{parts[1], count}
}

type reaction struct {
	inputs []chem
	output chem
}

func parserx(str string) *reaction {
	halves := strings.Split(str, " => ")
	if len(halves) != 2 {
		panic(str)
	}

	parts := strings.Split(halves[0], ", ")
	inputs := make([]chem, len(parts))
	for i, part := range parts {
		inputs[i] = parsechem(part)
	}

	output := parsechem(halves[1])

	return &reaction{inputs, output}
}

type nanofac struct {
	reactions map[string]*reaction
	leftovers map[string]int
}

func parsenanofac(str string) *nanofac {
	lines := strings.Split(strings.TrimSpace(str), "\n")
	rxs := map[string]*reaction{}

	for _, line := range lines {
		rx := parserx(line)
		rxs[rx.output.name] = rx
	}

	return &nanofac{
		reactions: rxs,
		leftovers: map[string]int{},
	}
}

func (f *nanofac) cost(name string, n int) int {
	if name == "ORE" {
		// ORE can always be acquired directly.
		return n
	}

	left := f.leftovers[name]
	if left > n {
		// We've already got enough left over from something else we did.
		f.leftovers[name] = left - n
		return 0
	}
	if left > 0 {
		// We have some left, but not enough.
		delete(f.leftovers, name)
		n -= left
	}

	rx, ok := f.reactions[name]
	if !ok {
		panic(fmt.Sprintf("no reaction found to create '%v'", name))
	}

	// Figure out how many times we need to run the reaction.
	mult := n / rx.output.count
	if (n % rx.output.count) > 0 {
		mult++
	}

	// Recursively sort out how much this is going to cost.
	cost := 0
	for i := range rx.inputs {
		cost += f.cost(rx.inputs[i].name, rx.inputs[i].count*mult)
	}

	// Save any left-overs.
	left = (mult * rx.output.count) - n
	if left > 0 {
		f.leftovers[rx.output.name] = left
	}

	return cost
}

func max(name string, budget int) int {
	input := lib.Load("input.txt")

	n := 1300000
	for {
		nf := parsenanofac(input)
		cost := nf.cost("FUEL", n+1)
		if cost > budget {
			return n
		}
		n++
	}
}

func main() {
	{
		nf := parsenanofac(lib.Load("input.txt"))
		cost := nf.cost("FUEL", 1)
		fmt.Println("cost:", cost)
	}
	{
		fuel := max("FUEL", 1000000000000)
		fmt.Println("max:", fuel)
	}
}
