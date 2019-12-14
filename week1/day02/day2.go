package main

import (
	"fmt"

	"github.com/fernomac/advent2019/lib"
)

var code = lib.LoadProgram("input.txt")

func compute(noun, verb int) int {
	p := lib.NewPuter(code)

	p.Write(1, noun)
	p.Write(2, verb)

	p.Run()

	return p.Read(0)
}

func _main() {
	n := compute(12, 2)
	fmt.Println(n)
}

func main() {
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			if compute(noun, verb) == 19690720 {
				fmt.Println("found it: noun =", noun, ", verb =", verb)
				return
			}
		}
	}
	panic("not found")
}
