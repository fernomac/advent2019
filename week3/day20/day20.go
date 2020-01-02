package main

import (
	"fmt"
	"strings"

	"github.com/fernomac/advent2019/lib"
	"github.com/fernomac/advent2019/lib/grid"
	"github.com/fernomac/advent2019/lib/queue"
)

type tile uint8

const (
	empty tile = iota
	wall
	warp
)

type node struct {
	tile  tile
	dst   grid.Point
	loop  grid.Point
	inner bool
	id    string
}

func (n *node) ID() string {
	return n.id
}

type portal struct {
	src, dst grid.Point
}

func parse(in string) *grid.Grid {
	lines := strings.Split(in, "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	size := grid.Point{X: len(lines[0]), Y: len(lines)}

	// Parse out all the underlying runes.
	runes := map[grid.Point]rune{}
	for y, l := range lines {
		for x, r := range l {
			if r != ' ' {
				runes[grid.Point{X: x, Y: y}] = r
			}
		}
	}

	// Find all the portals.
	nodes := map[grid.Point]grid.Node{}
	portals := map[string][]portal{}

	for p, r := range runes {
		switch {
		case r == '.':
			nodes[p] = &node{tile: empty}

		case r == '#':
			nodes[p] = &node{tile: wall}

		case r >= 'A' && r <= 'Z':
			name, portal := portalize(runes, p, r)
			nodes[portal.src] = &node{
				tile:  warp,
				inner: inner(size, portal.src),
			}
			portals[name] = append(portals[name], portal)
		}
	}

	// Tag the start and end tiles with an ID.
	identify(nodes, portals, "AA")
	identify(nodes, portals, "ZZ")

	// Link up the remaining portals.
	for _, portal := range portals {
		if len(portal) != 2 {
			panic(fmt.Sprintf("mismatched portals: %v", portal))
		}
		{
			n := nodes[portal[0].src].(*node)
			n.dst = portal[1].dst
			n.loop = portal[0].dst
		}
		{
			n := nodes[portal[1].src].(*node)
			n.dst = portal[0].dst
			n.loop = portal[1].dst
		}
	}

	return grid.New(nodes)
}

func portalize(runes map[grid.Point]rune, p grid.Point, r rune) (string, portal) {
	r2 := rune(0)
	r2p := grid.Point{}
	src := grid.Point{}
	dst := grid.Point{}

	foundr2 := false
	founddst := false

	for _, n := range p.Neighbors() {
		nr := runes[n]
		switch {
		case nr == '.':
			// Found this portal's destination.
			founddst = true
			src = p
			dst = n

		case nr >= 'A' && nr <= 'Z':
			// Found the other half of the warp label.
			foundr2 = true
			r2 = nr
			r2p = n
			delete(runes, r2p)
		}
	}

	if !foundr2 {
		panic(fmt.Sprintf("no neighboring label rune at %v", p))
	}

	if !founddst {
		src = r2p
		for _, n := range r2p.Neighbors() {
			if runes[n] == '.' {
				founddst = true
				dst = n
				break
			}
		}
		if !founddst {
			panic(fmt.Sprintf("neither %v nor %v is next to a destination", p, r2p))
		}
	}

	// Labels are read left-to-right or top-to-bottom.
	if r2p.X < p.X || r2p.Y < p.Y {
		r, r2 = r2, r
	}

	return string([]rune{r2, r}), portal{src, dst}
}

func inner(gridsize grid.Point, loc grid.Point) bool {
	if loc.X <= 1 || loc.Y <= 1 {
		return false
	}
	if gridsize.X-loc.X <= 2 || gridsize.Y-loc.Y <= 2 {
		return false
	}
	return true
}

func identify(nodes map[grid.Point]grid.Node, portals map[string][]portal, id string) {
	portal := portals[id]
	if len(portal) == 0 {
		panic(fmt.Sprintf("no portal with id %v", id))
	}
	if len(portal) > 1 {
		panic(fmt.Sprintf("more than one portal with id %v", id))
	}

	delete(portals, id)
	delete(nodes, portal[0].src)

	nodes[portal[0].dst].(*node).id = id
}

type move struct {
	to    grid.Point
	steps int
}

func movesFrom(g *grid.Grid, from grid.Point) []move {
	ret := []move{}
	visited := map[grid.Point]bool{}

	q := queue.New()
	q.Push(move{to: from})

	for q.Len() > 0 {
		cur := q.Pop().(move)

		curnode := g.AtPt(cur.to).(*node)
		if curnode.tile == warp || curnode.id == "ZZ" {
			ret = append(ret, cur)
		}

		for _, n := range cur.to.Neighbors() {
			if visited[n] {
				continue
			}

			nn := g.AtPt(n)
			if nn == nil {
				continue
			}
			nextnode := nn.(*node)

			if nextnode.tile == wall {
				continue
			}
			if nextnode.tile == warp && nextnode.loop == from {
				// Don't step right back into the portal that got us here.
				continue
			}

			q.Push(move{n, cur.steps + 1})
			visited[n] = true
		}
	}

	return ret
}

type state struct {
	loc   grid.Point
	steps int
}

func (s state) Priority() int {
	return s.steps
}

func walk(g *grid.Grid) int {
	start, ok := g.Point("AA")
	if !ok {
		panic("no starting point")
	}
	end, ok := g.Point("ZZ")
	if !ok {
		panic("no ending point")
	}

	movesfrom := map[grid.Point][]move{}
	movesfrom[start] = movesFrom(g, start)
	for _, n := range g.NodesWhere(func(n grid.Node) bool { return n.(*node).tile == warp }) {
		dst := n.(*node).dst
		movesfrom[dst] = movesFrom(g, dst)
	}

	q := queue.NewPriority()
	q.Push(state{start, 0})
	visited := map[grid.Point]bool{}

	for q.Len() > 0 {
		s := q.Pop().(state)

		if s.loc == end {
			return s.steps
		}

		if visited[s.loc] {
			continue
		}
		visited[s.loc] = true

		for _, move := range movesfrom[s.loc] {
			dst := move.to
			if n := g.AtPt(dst).(*node); n.tile == warp {
				dst = n.dst
			}
			q.Push(state{dst, s.steps + move.steps})
		}
	}

	panic("maze is unsolveable")
}

type state2 struct {
	loc   grid.Point
	depth int
	steps int
}

func (s state2) Priority() int {
	return s.steps
}

type pd struct {
	point grid.Point
	depth int
}

func walk2(g *grid.Grid) int {
	start, ok := g.Point("AA")
	if !ok {
		panic("no starting point")
	}
	end, ok := g.Point("ZZ")
	if !ok {
		panic("no ending point")
	}

	movesfrom := map[grid.Point][]move{}
	movesfrom[start] = movesFrom(g, start)
	for _, n := range g.NodesWhere(func(n grid.Node) bool { return n.(*node).tile == warp }) {
		dst := n.(*node).dst
		movesfrom[dst] = movesFrom(g, dst)
	}

	q := queue.NewPriority()
	q.Push(state2{start, 0, 0})
	visited := map[pd]bool{}

	for q.Len() > 0 {
		s := q.Pop().(state2)

		if s.loc == end && s.depth == 0 {
			return s.steps
		}

		if visited[pd{s.loc, s.depth}] {
			continue
		}
		visited[pd{s.loc, s.depth}] = true

		// fmt.Println("visiting", s.loc, "@", s.depth, "after", s.steps, "steps")

		for _, move := range movesfrom[s.loc] {
			dst := move.to
			depth := s.depth

			n := g.AtPt(dst).(*node)
			if n.tile == warp {
				dst = n.dst
				if n.inner {
					depth++
				} else {
					if depth == 0 {
						continue
					}
					depth--
				}
			}

			q.Push(state2{dst, depth, s.steps + move.steps})
		}
	}

	panic("maze is unsolveable")
}

func main() {
	g := parse(lib.Load("input.txt"))
	fmt.Println("steps:", walk(g))
	fmt.Println("steps2:", walk2(g))
}
