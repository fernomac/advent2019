package main

import (
	"fmt"
	"time"

	"github.com/fernomac/advent2019/lib"
)

type point struct {
	x, y int
}

type game struct {
	screen map[point]int
	puter  *lib.Puter

	in   chan int
	out  chan int
	done chan int

	ball   point
	paddle point

	score int
}

func newgame(coin int) *game {
	in := make(chan int, 1)
	out := make(chan int, 10240)
	done := make(chan int, 1)

	prog := lib.LoadProgram("input.txt")
	if coin > 0 {
		prog[0] = coin
	}

	puter := lib.NewPuter(prog)
	puter.StdinCh(in)
	puter.StdoutCh(out)

	g := &game{
		screen: make(map[point]int),
		puter:  puter,

		in:   in,
		out:  out,
		done: done,
	}

	return g
}

func (g *game) runonce() {
	go g.collect()

	g.puter.Run()

	close(g.out)
	<-g.done
}

func (g *game) play() {
	for !g.puter.RunNB() {
		g.collectNB()
		g.draw()

		dir := 0
		if g.ball.x < g.paddle.x {
			dir = -1
		} else if g.ball.x > g.paddle.x {
			dir = 1
		}

		time.Sleep(1 * time.Millisecond)
		g.in <- dir
	}

	g.collectNB()
	g.draw()
}

func (g *game) collect() {
	for {
		x, ok := <-g.out
		if !ok {
			close(g.done)
			return
		}

		y := <-g.out
		id := <-g.out

		// fmt.Println("(", x, ",", y, ")", ":", id)

		if x == -1 && y == 0 {
			g.score = id
		} else {
			g.screen[point{x, y}] = id
		}
	}
}

func (g *game) collectNB() {
	for {
		select {
		case x := <-g.out:
			y := <-g.out
			id := <-g.out

			// fmt.Println("(", x, ",", y, ")", ":", id)

			if x == -1 && y == 0 {
				g.score = id
			} else {
				g.screen[point{x, y}] = id
			}

			if id == 3 {
				g.paddle = point{x, y}
			} else if id == 4 {
				g.ball = point{x, y}
			}

		default:
			fmt.Println("done")
			return
		}
	}
}

func (g *game) draw() {
	for y := 0; y < 20; y++ {
		for x := 0; x < 37; x++ {
			id := g.screen[point{x, y}]

			var c string
			switch id {
			case 0:
				c = " "
			case 1:
				c = "W" // wall
			case 2:
				c = "b" // block
			case 3:
				c = "-" // paddle
			case 4:
				c = "o" // ball
			}

			fmt.Print(c)
		}
		fmt.Println()
	}
	fmt.Println("score:", g.score)
}

func (g *game) blocks() int {
	count := 0

	for _, id := range g.screen {
		if id == 2 {
			count++
		}
	}

	return count
}

func main() {
	{
		g := newgame(0)
		g.runonce()
		fmt.Println("initial blocks:", g.blocks())
	}

	{
		g := newgame(2)
		g.play()
		fmt.Println(g.score)
	}
}
