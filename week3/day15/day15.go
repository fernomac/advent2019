package main

import (
	"container/list"
	"fmt"

	"github.com/fernomac/advent2019/lib"
)

type point struct {
	x, y int
}

type kind uint8

const (
	wallkind kind = iota
	openkind
	oxykind
)

type node struct {
	kind  kind
	place point
}

var wall = &node{kind: wallkind}

type game struct {
	nodes map[point]*node
	oxy   *node

	in    chan int
	out   chan int
	puter *lib.Puter
}

func newgame() *game {
	nodes := map[point]*node{
		point{0, 0}: &node{
			kind: openkind,
		},
	}

	in := make(chan int, 1)
	out := make(chan int, 1)

	puter := lib.NewPuter(lib.LoadProgram("input.txt"))
	puter.StdinCh(in)
	puter.StdoutCh(out)

	return &game{
		nodes: nodes,
		in:    in,
		out:   out,
		puter: puter,
	}
}

// explore performs a depth-first search of the environment to build the
// full `nodes` graph and find the oxygen.
func (g *game) explore() {
	g.exploreFrom(g.nodes[point{0, 0}])
}

func (g *game) exploreFrom(n *node) {
	for dir := 1; dir <= 4; dir++ {
		pt := move(n.place, dir)
		if g.nodes[pt] == nil {
			g.exploreFromHelper(n, pt, dir)
		}
	}
}

func (g *game) exploreFromHelper(n *node, pt point, dir int) {
	res := g.move(dir)
	switch res {
	case 0:
		// There was a wall, we didn't move.
		g.nodes[pt] = wall

	case 1, 2:
		// We moved, record and recursively explore the new space.
		nn := &node{
			kind:  openkind,
			place: pt,
		}

		if res == 2 {
			nn.kind = oxykind
			g.oxy = nn
		}

		g.nodes[pt] = nn

		g.exploreFrom(nn)

		if res := g.move(reverse(dir)); res != 1 {
			panic("error backtracking")
		}

	default:
		panic(res)
	}
}

func (g *game) move(dir int) int {
	g.in <- dir
	if g.puter.RunNB() {
		panic("puter exited")
	}
	return <-g.out
}

func (g *game) neighbor(n *node, dir int) *node {
	pt := move(n.place, dir)
	return g.nodes[pt]
}

func move(pt point, dir int) point {
	switch dir {
	case 1:
		return point{pt.x, pt.y - 1}

	case 2:
		return point{pt.x, pt.y + 1}

	case 3:
		return point{pt.x - 1, pt.y}

	case 4:
		return point{pt.x + 1, pt.y}

	default:
		panic(dir)
	}
}

func reverse(dir int) int {
	switch dir {
	case 1:
		return 2
	case 2:
		return 1
	case 3:
		return 4
	case 4:
		return 3
	default:
		panic(dir)
	}
}

func (g *game) draw() {
	min, max := g.dims()

	for y := min.y; y <= max.y; y++ {
		for x := min.x; x <= max.x; x++ {

			if x == 0 && y == 0 {
				fmt.Print("X")
				continue
			}

			n := g.nodes[point{x, y}]
			if n == nil {
				fmt.Print(" ")
				continue
			}

			switch n.kind {
			case openkind:
				fmt.Print(".")

			case wallkind:
				fmt.Print("#")

			case oxykind:
				fmt.Print("O")

			default:
				panic(n)
			}
		}
		fmt.Println()
	}
}

func (g *game) dims() (point, point) {
	min, max := point{}, point{}

	for pt := range g.nodes {
		if pt.x < min.x {
			min.x = pt.x
		}
		if pt.x > max.x {
			max.x = pt.x
		}
		if pt.y < min.y {
			min.y = pt.y
		}
		if pt.y > max.y {
			max.y = pt.y
		}
	}

	return min, max
}

type bfs struct {
	node *node
	path []int
}

func (g *game) pathfind() []int {
	q := list.New()
	q.PushBack(bfs{
		node: g.nodes[point{0, 0}],
	})

	visited := map[point]bool{
		point{0, 0}: true,
	}

	for q.Len() > 0 {
		front := q.Front()
		q.Remove(front)
		cur := front.Value.(bfs)

		if cur.node.kind == oxykind {
			return cur.path
		}

		for dir := 1; dir <= 4; dir++ {
			neighbor := g.neighbor(cur.node, dir)
			if neighbor.kind == wallkind {
				continue
			}
			if visited[neighbor.place] {
				continue
			}

			q.PushBack(bfs{
				node: neighbor,
				path: append(cur.path, dir),
			})
			visited[neighbor.place] = true
		}
	}

	panic("finished our search without finding the oxygen")
}

func (g *game) oxify() int {
	q := list.New()
	q.PushBack(bfs{
		node: g.oxy,
	})

	visited := map[point]bool{
		g.oxy.place: true,
	}

	max := 0

	for q.Len() > 0 {
		front := q.Front()
		q.Remove(front)
		cur := front.Value.(bfs)

		if len(cur.path) > max {
			max = len(cur.path)
		}

		for dir := 1; dir <= 4; dir++ {
			neighbor := g.neighbor(cur.node, dir)
			if neighbor.kind == wallkind {
				continue
			}
			if visited[neighbor.place] {
				continue
			}

			q.PushBack(bfs{
				node: neighbor,
				path: append(cur.path, dir),
			})
			visited[neighbor.place] = true
		}
	}

	return max
}

func main() {
	g := newgame()
	g.explore()
	path := g.pathfind()
	fmt.Println("steps:", len(path))

	time := g.oxify()
	fmt.Println("time to oxify:", time)
}
