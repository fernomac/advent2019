package main

import (
	"fmt"

	"github.com/fernomac/advent2019/lib"
)

func main() {
	test := lib.LoadProgram("input.txt")

	{
		p := lib.NewPuter(test)
		p.Stdin([]int{1})
		p.Run()
		fmt.Println(p.Stdout())
	}

	{
		p := lib.NewPuter(test)
		p.Stdin([]int{5})
		p.Run()
		fmt.Println(p.Stdout())
	}
}
