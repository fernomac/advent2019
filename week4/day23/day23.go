package main

import (
	"fmt"

	"github.com/fernomac/advent2019/lib"
)

type network struct {
	ps   []*lib.Puter
	ins  []chan int
	outs []chan int

	natx, naty int
	lasty      int
}

func internetwork(size int) *network {
	prog := lib.LoadProgram("input.txt")

	ps := make([]*lib.Puter, size)
	ins := make([]chan int, size)
	outs := make([]chan int, size)

	for i := 0; i < size; i++ {
		p := lib.NewPuter(prog)
		in := make(chan int, 1024)
		out := make(chan int, 1024)

		p.StdinCh(in)
		p.StdoutCh(out)

		// Send the 'puter its ID.
		in <- i
		p.RunNB()

		ps[i] = p
		ins[i] = in
		outs[i] = out
	}

	return &network{
		ps:   ps,
		ins:  ins,
		outs: outs,
		natx: -1,
	}
}

func (n *network) send() (bool, int) {
	idle := true
	for i := range n.ps {
		if n.sendfrom(i) {
			idle = false
		}
	}

	if idle && n.natx != -1 {
		if n.naty == n.lasty {
			// Part 2: the first duplicate Y value sent FROM the NAT.
			return true, n.naty
		}

		n.lasty = n.naty
		n.ins[0] <- n.natx
		n.ins[0] <- n.naty
	}

	return false, 0
}

func (n *network) sendfrom(src int) bool {
	sent := false
	out := n.outs[src]

	for {
		select {
		case dst := <-out:
			// A packet is queued for sending, process it.
			x := <-out
			y := <-out
			sent = true

			if dst == 255 {
				if n.natx == -1 {
					// Part 1: the Y value of the first packet sent to the NAT.
					fmt.Println(y)
				}
				n.natx = x
				n.naty = y
			} else {
				n.ins[dst] <- x
				n.ins[dst] <- y
			}

		default:
			// This 'puter has nothing more to say right now, that's fine.
			return sent
		}
	}
}

func (n *network) receive() {
	for i := range n.ps {
		in := n.ins[i]

		if len(in) == 0 {
			// No mail for you!
			in <- -1
		}

		n.ps[i].RunNB()
	}
}

func main() {
	n := internetwork(50)

	for {
		done, ret := n.send()
		if done {
			fmt.Println(ret)
			break
		}

		n.receive()
	}
}
