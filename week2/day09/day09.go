package main

import (
	"fmt"

	"github.com/fernomac/advent2019/lib"
)

func main() {
	{
		p := lib.NewPuter(lib.LoadProgram("input.txt"))
		p.Stdin([]int{1})
		p.Run()
		fmt.Println(p.Stdout())
	}

	{
		p := lib.NewPuter(lib.LoadProgram("input.txt"))
		p.Stdin([]int{2})
		p.Run()
		fmt.Println(p.Stdout())
	}
}
