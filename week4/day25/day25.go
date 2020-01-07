package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fernomac/advent2019/lib"
)

func collect(p *lib.ASCIIPuter) {
	cmds := []string{
		"south",
		"take hologram",
		"north",
		"west",
		"take mutex",
		"south",
		"south",
		"south",
		"take polygon",
		"north",
		"east",
		"take weather machine",
		"west",
		"north",
		"north",
		"east",
		"north",
		"north",
		"north",
		"take semiconductor",
		"west",
		"take monolith",
		"east",
		"east",
		"take prime number",
		"west",
		"south",
		"west",
		"north",
		"take jam",
		"west",
	}
	for _, cmd := range cmds {
		p.Exec(cmd)
	}
}

func bruteforce(p *lib.ASCIIPuter) {
	items := []string{
		"mutex",
		"hologram",
		"polygon",
		"jam",
		"semiconductor",
		"prime number",
		"monolith",
		"weather machine",
	}
	inv := int(0xFF)

	for dinv := 1; dinv <= 0xFF; dinv++ {
		for i := 0; i < 8; i++ {
			bit := 1 << i
			if inv&bit != 0 && dinv&bit == 0 {
				p.Exec("drop " + items[i])
			} else if inv&bit == 0 && dinv&bit != 0 {
				p.Exec("take " + items[i])
			}
		}
		inv = dinv

		out := p.Exec("north")
		if strings.Contains(out, "You may proceed.") {
			fmt.Println(out)
			return
		}
	}
}

func interactive(p *lib.ASCIIPuter) {
	r := bufio.NewScanner(os.Stdin)
	for !p.Run() {
		fmt.Print(p.Output())
		if !r.Scan() {
			panic(r.Err())
		}
		p.Input(r.Text() + "\n")
	}

	fmt.Print(p.Output())
}

func main() {
	{
		p := lib.ASCIIfy(lib.NewPuter(lib.LoadProgram("input.txt")))
		collect(p)
		bruteforce(p)
		// interactive(p)
	}
}
