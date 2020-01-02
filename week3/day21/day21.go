package main

import (
	"fmt"

	"github.com/fernomac/advent2019/lib"
)

func main() {
	prog := lib.LoadProgram("input.txt")

	{
		p := lib.ASCIIfy(lib.NewPuter(prog))
		p.Run()
		// fmt.Print(p.Output())

		// Jump if:
		// - There's a hole immediately in front of us (may as well).
		// - There's a hole two or three spaces in front of us and ground 4 spaces in front.

		p.Input("NOT B J\n") // j = ~b
		p.Input("NOT C T\n") // t = ~c
		p.Input("OR J T\n")  // t = ~b || ~c
		p.Input("AND D T\n") // t = (~b || ~c) && d
		p.Input("NOT A J\n") // j = ~a
		p.Input("OR T J\n")  // j = ~a || ((~b || ~c) && d)
		p.Input("WALK\n")

		p.Run()
		fmt.Print(p.Output())
		fmt.Println(p.Result())
	}

	{
		p := lib.ASCIIfy(lib.NewPuter(prog))
		p.Run()
		fmt.Print(p.Output())

		// .@...............
		// #####.#.##.#.####
		//   ABCDEFGHI

		// ..@..............
		// #####.#.##.#.####
		//    ABCDEFGHI

		// ...@.............
		// #####.#.##.#.####
		//     ABCDEFGHI

		// ....@............
		// #####.#.##.#.####
		//      ABCDEFGHI

		// need to jump if: ~a | ~b | ~c
		// safe to jump if: d & (h | (e & (i | f))

		// .................
		// .................
		// ....@............
		// #####...####..###
		//      ABCDEFGHI

		// j := (~a | ~b | ~c) & (d & (h | (e & (i | f))))
		// j := (~a | ~b | ~c) & (d & (h | (e & ~(~i & ~f))))
		// j := (~a | ~b | ~c) & (d & (h | ~(~e | (~i & ~f))))
		// j := (~a | ~b | ~c) & (d & ~(~h & (~e | (~i & ~f))))
		// j := (~a | ~b | ~c) & ~(~d | (~h & (~e | (~i & ~f))))
		// j := ~(a & b & c) & ~(~d | (~h & (~e | (~i & ~f))))
		// j := ~((a & b & c) | (~d | (~h & (~e | (~i & ~f))))

		// t = (~i & ~f)
		p.Input("NOT I J\n")
		p.Input("NOT F T\n")
		p.Input("AND J T\n")

		// t = (~e | (~i & ~f)) = (~e | t)
		p.Input("NOT E J\n")
		p.Input("OR J T\n")

		// t = (~h & (~e | (~i & ~f))) = (~h & t)
		p.Input("NOT H J\n")
		p.Input("AND J T\n")

		// t = (~d | (~h & (~e | (~i & ~f)))) = (~d | t)
		p.Input("NOT D J\n")
		p.Input("OR J T\n")

		// j = a & b & c
		p.Input("NOT A J\n")
		p.Input("NOT J J\n")
		p.Input("AND B J\n")
		p.Input("AND C J\n")

		// j = ~((a & b & c) | (~d | (~h & (~e | (~i & ~f)))) = ~(j | t)
		p.Input("OR J T\n")
		p.Input("NOT T J\n")
		p.Input("RUN\n")

		p.Run()
		fmt.Print(p.Output())
		fmt.Println(p.Result())
	}
}
