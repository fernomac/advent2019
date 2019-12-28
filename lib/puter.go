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

type cpustate uint8

const (
	initial cpustate = iota
	running
	blocked
	done
)

// A Puter is an intcode computer.
type Puter struct {
	mem map[int]int
	ip  int
	rbo int

	state    cpustate
	blocking bool

	stdin     <-chan int
	stdout    chan<- int
	stdoutbuf []int
}

// NewPuter creates a new Puter.
func NewPuter(prog Program) *Puter {
	mem := map[int]int{}
	for i, b := range prog {
		mem[i] = b
	}

	return &Puter{
		mem:   mem,
		state: initial,
	}
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
	ch := make(chan int, len(stdin))
	for _, i := range stdin {
		ch <- i
	}
	p.stdin = ch
}

// Stdout returns the computer's output stream.
func (p *Puter) Stdout() []int {
	if p.stdout != nil {
		panic("stdout not buffered")
	}
	return p.stdoutbuf
}

// DropStdout drops any buffered stdout
func (p *Puter) DropStdout() {
	p.stdoutbuf = nil
}

// StdinCh sets the computer's input stream as a channel.
func (p *Puter) StdinCh(ch <-chan int) {
	p.stdin = ch
}

// StdoutCh sets this puter's stdout channel.
func (p *Puter) StdoutCh(ch chan<- int) {
	p.stdout = ch
}

// Run runs the intcode program until it quits.
func (p *Puter) Run() {
	if !p.run(true) {
		panic("run(true) returned false")
	}
}

// RunNB runs the intcode program until it either quits or needs more input.
func (p *Puter) RunNB() bool {
	return p.run(false)
}

func (p *Puter) run(blocking bool) bool {
	p.blocking = blocking
	p.state = running

	for p.state == running {
		p.step()
	}

	return p.state == done
}

func (p *Puter) step() {
	opcode, modes := decode(p.mem[p.ip])

	op, ok := opcodes[opcode]
	if !ok {
		panic(fmt.Sprintf("invalid instruction %v at position %v", opcode, p.ip))
	}
	op(p, modes)
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
	relativeMode  mode = 2
)

type modeset int

func (m modeset) get(n int) mode {
	temp := int(m)
	for n > 0 {
		temp /= 10
		n--
	}
	return mode(temp % 10)
}
