package lib

import "fmt"

type op func(*Puter, modeset)

// opcodes is a map from opcode to opcode handler function.
var opcodes = func() map[int]op {
	return map[int]op{
		1: (*Puter).add,
		2: (*Puter).mult,

		3: (*Puter).in,
		4: (*Puter).out,

		5: (*Puter).jnz,
		6: (*Puter).jz,
		7: (*Puter).lt,
		8: (*Puter).eq,
	}
}()

//
// Opcodes.
//

func (p *Puter) add(modes modeset) {
	lhs := p.readarg(0, modes)
	rhs := p.readarg(1, modes)

	res := lhs + rhs

	p.writearg(2, modes, res)
	p.ip += 4
}

func (p *Puter) mult(modes modeset) {
	lhs := p.readarg(0, modes)
	rhs := p.readarg(1, modes)

	res := lhs * rhs

	p.writearg(2, modes, res)
	p.ip += 4
}

func (p *Puter) in(modes modeset) {
	in := p.stdin[0]
	p.stdin = p.stdin[1:]

	p.writearg(0, modes, in)
	p.ip += 2
}

func (p *Puter) out(modes modeset) {
	out := p.readarg(0, modes)

	p.stdout = append(p.stdout, out)
	p.ip += 2
}

func (p *Puter) jnz(modes modeset) {
	test := p.readarg(0, modes)

	if test != 0 {
		p.ip = p.readarg(1, modes)
	} else {
		p.ip += 3
	}
}

func (p *Puter) jz(modes modeset) {
	test := p.readarg(0, modes)

	if test == 0 {
		p.ip = p.readarg(1, modes)
	} else {
		p.ip += 3
	}
}

func (p *Puter) lt(modes modeset) {
	lhs := p.readarg(0, modes)
	rhs := p.readarg(1, modes)

	res := 0
	if lhs < rhs {
		res = 1
	}

	p.writearg(2, modes, res)
	p.ip += 4
}

func (p *Puter) eq(modes modeset) {
	lhs := p.readarg(0, modes)
	rhs := p.readarg(1, modes)

	res := 0
	if lhs == rhs {
		res = 1
	}

	p.writearg(2, modes, res)
	p.ip += 4
}

//
// Helpers
//

func (p *Puter) readarg(n int, modes modeset) int {
	return p.read(p.arg(n), modes.get(n))
}

func (p *Puter) writearg(n int, modes modeset, val int) {
	p.write(p.arg(n), modes.get(n), val)
}

func (p *Puter) arg(n int) int {
	return p.mem[p.ip+n+1]
}

func (p *Puter) read(arg int, mode mode) int {
	switch mode {
	case immediateMode:
		return arg

	case positionMode:
		return p.mem[arg]

	default:
		panic(fmt.Sprintf("bad read mode %v", mode))
	}
}

func (p *Puter) write(arg int, mode mode, val int) {
	switch mode {
	case positionMode:
		p.mem[arg] = val

	default:
		panic(fmt.Sprintf("bad write mode %v", mode))
	}
}
