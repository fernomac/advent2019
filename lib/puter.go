package lib

import (
	"fmt"
	"strconv"
	"strings"
)

// A Program written in Intcode.
type Program []int

// LoadProgram loads an intcode program from a file.
func LoadProgram(filename string) Program {
	file := strings.TrimSpace(Load(filename))
	nums := strings.Split(file, ",")

	prog := make([]int, len(nums))
	for i, num := range nums {
		n, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		prog[i] = n
	}

	return prog
}

// A Puter is an intcode computer.
type Puter struct {
	mem []int
	ip  int

	stdin  []int
	stdout []int
}

// NewPuter creates a new Puter.
func NewPuter(prog Program) *Puter {
	return &Puter{mem: append([]int{}, prog...)}
}

// Read reads a value from memory.
func (p *Puter) Read(n int) int {
	return p.mem[n]
}

// Write writes a value into memory.
func (p *Puter) Write(n, val int) {
	p.mem[n] = val
}

// Stdin sets the computer's input stream.
func (p *Puter) Stdin(stdin []int) {
	p.stdin = stdin
}

// Stdout returns the computer's output stream.
func (p *Puter) Stdout() []int {
	return p.stdout
}

// Run runs the intcode program.
func (p *Puter) Run() {
	for p.Step() {
		// keep going.
	}
}

// Step steps the intcode program forward one instruction.
func (p *Puter) Step() bool {
	opcode, modes := decode(p.mem[p.ip])
	if opcode == 99 {
		// We're done.
		return false
	}

	op, ok := opcodes[opcode]
	if !ok {
		panic(fmt.Sprintf("invalid instruction %v at position %v", opcode, p.ip))
	}
	op(p, modes)

	return true
}

func decode(raw int) (int, modeset) {
	high := raw / 100
	low := raw % 100
	return low, modeset(high)
}

type mode uint8

const (
	immediateMode mode = 1
	positionMode  mode = 0
)

type modeset int

func (m modeset) get(n int) mode {
	temp := int(m)
	if n > 0 {
		temp /= (n * 10)
	}
	return mode(temp % 10)
}
