package lib

import (
	"fmt"
	"strconv"
	"strings"
)

// An Intcode program.
type Intcode []int

// LoadIntcode loads an intcode program from a file.
func LoadIntcode(filename string) Intcode {
	file := strings.TrimSpace(Load(filename))
	nums := strings.Split(file, ",")

	mem := make([]int, len(nums))
	for i, num := range nums {
		n, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		mem[i] = n
	}

	return mem
}

// A Puter is an intcode computer.
type Puter struct {
	mem []int
	ip  int
}

// NewPuter creates a new Puter.
func NewPuter(program Intcode) *Puter {
	return &Puter{mem: append([]int{}, program...)}
}

// Read reads a value from memory.
func (p *Puter) Read(n int) int {
	return p.mem[n]
}

// Write writes a value into memory.
func (p *Puter) Write(n, val int) {
	p.mem[n] = val
}

// Run runs the intcode program.
func (p *Puter) Run() {
	for p.Step() {
		// keep going.
	}
}

var opcodes = func() map[int]func(*Puter) bool {
	return map[int]func(p *Puter) bool{
		1:  (*Puter).add,
		2:  (*Puter).mult,
		99: func(p *Puter) bool { return false },
	}
}()

// Step steps the intcode program forward on instruction.
func (p *Puter) Step() bool {
	opcode := p.mem[p.ip]
	op, ok := opcodes[opcode]
	if !ok {
		panic(fmt.Sprintf("invalid instruction %v at position %v", opcode, p.ip))
	}

	return op(p)
}

func (p *Puter) add() bool {
	return p.math(func(a, b int) int { return a + b })
}

func (p *Puter) mult() bool {
	return p.math(func(a, b int) int { return a * b })
}

func (p *Puter) math(f func(int, int) int) bool {
	arg1 := p.mem[p.ip+1]
	arg2 := p.mem[p.ip+2]
	arg3 := p.mem[p.ip+3]

	lhs := p.mem[arg1]
	rhs := p.mem[arg2]

	res := f(lhs, rhs)

	p.mem[arg3] = res
	p.ip += 4

	return true
}
