package main

import (
	"fmt"
	"strconv"

	"github.com/fernomac/advent2019/lib"
)

func main() {
	total := 0

	lib.ForLines("input.txt", func(line string) {
		mass, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}

		fuel := fuelFor(mass)
		total += fuel
	})

	fmt.Println(total)
}

func fuelFor(mass int) int {
	fuel := (mass / 3) - 2
	if fuel <= 0 {
		return 0
	}
	return fuel + fuelFor(fuel)
}
