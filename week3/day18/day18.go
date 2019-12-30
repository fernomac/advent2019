package main

import (
	"container/heap"
	"container/list"
	"fmt"
	"math"
	"strings"

	"github.com/fernomac/advent2019/lib"
)

type point struct {
	x, y int
}

type tile uint8

const (
	empty tile = iota
	wall
	door
	key
)

type space struct {
	tile tile
	id   int
}

type grid struct {
	spaces map[point]space
	keyset int
	start  point
}

func newgrid(input string) *grid {
	spaces := map[point]space{}
	keyset := 0
	start := point{-1, -1}

	lines := strings.Split(strings.TrimSpace(input), "\n")
	for y, line := range lines {
		for x, c := range line {
			switch {
			case c == '@':
				start = point{x, y}

			case c == '.':
				spaces[point{x, y}] = space{tile: empty}

			case c == '#':
				spaces[point{x, y}] = space{tile: wall}

			case c >= 'a' && c <= 'z':
				id := 1 << (c - 'a')
				keyset |= id
				spaces[point{x, y}] = space{
					tile: key,
					id:   id,
				}

			case c >= 'A' && c <= 'Z':
				id := 1 << (c - 'A')
				spaces[point{x, y}] = space{
					tile: door,
					id:   id,
				}

			default:
				panic(fmt.Sprintf("bogus char %v at (%v, %v)", c, x, y))
			}
		}
	}

	return &grid{spaces, keyset, start}
}

type move struct {
	loc   point
	dist  int
	key   int
	doors int
}

type movequeue struct {
	q *list.List
}

func (q *movequeue) init() {
	q.q = list.New()
}

func (q *movequeue) len() int {
	return q.q.Len()
}

func (q *movequeue) push(m move) {
	q.q.PushBack(m)
}

func (q *movequeue) pop() move {
	f := q.q.Front()
	q.q.Remove(f)
	return f.Value.(move)
}

func (g *grid) moves(from point) []move {
	ret := []move{}
	visited := map[point]bool{}

	queue := movequeue{}
	queue.init()
	queue.push(move{loc: from})

	for queue.len() > 0 {
		ret = g.visit(&queue, visited, ret)
	}

	return ret
}

var dirs = []point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

func (g *grid) visit(queue *movequeue, visited map[point]bool, ret []move) []move {
	cur := queue.pop()

	// If we've found a move that lands on a key, add it to ret.
	if space := g.spaces[cur.loc]; space.tile == key {
		cur.key = space.id
		ret = append(ret, cur)
	}

	// Add all the new spaces that are reachable from here to the queue.
	for _, dir := range dirs {
		next := move{
			loc:   point{cur.loc.x + dir.x, cur.loc.y + dir.y},
			dist:  cur.dist + 1,
			doors: cur.doors,
		}

		if visited[next.loc] {
			// We've already been here, don't bother.
			continue
		}

		space := g.spaces[next.loc]
		if space.tile == wall {
			// Can't move here, it's a wall.
			continue
		}

		if space.tile == door {
			// Can move here iff we have the matching key.
			next.doors |= space.id
		}

		queue.push(next)
		visited[next.loc] = true
	}

	return ret
}

type path struct {
	loc    point
	dist   int
	keyset int
}

type pathqueue struct {
	items []path
}

func (q *pathqueue) Len() int {
	return len(q.items)
}

func (q *pathqueue) Less(i, j int) bool {
	return q.items[i].dist < q.items[j].dist
}

func (q *pathqueue) Swap(i, j int) {
	q.items[i], q.items[j] = q.items[j], q.items[i]
}

func (q *pathqueue) Push(x interface{}) {
	q.items = append(q.items, x.(path))
}

func (q *pathqueue) Pop() interface{} {
	n := len(q.items)
	p := q.items[n-1]
	q.items = q.items[:n-1]
	return p
}

type state struct {
	loc    point
	keyset int
}

func (g *grid) traverse() int {
	// Do a first pass to sort out what moves are possible.
	allmoves := map[point][]move{}
	allmoves[g.start] = g.moves(g.start)
	for loc, spc := range g.spaces {
		if spc.tile == key {
			allmoves[loc] = g.moves(loc)
		}
	}

	// Now find the optimal path using those moves.
	pq := pathqueue{}
	heap.Init(&pq)
	heap.Push(&pq, path{loc: g.start})

	visited := map[state]bool{}

	for {
		cur := heap.Pop(&pq).(path)
		if cur.keyset == g.keyset {
			// We've found all the keys!
			return cur.dist
		}

		// We've already been here but more efficiently, save some time.
		st := state{cur.loc, cur.keyset}
		if visited[st] {
			continue
		}
		visited[st] = true

		// Try all of the potential moves from this position.
		for _, move := range allmoves[cur.loc] {
			if move.key&cur.keyset != 0 {
				// We already have that key, not worth it.
				continue
			}
			if move.doors&^cur.keyset != 0 {
				// We need a key we don't have yet to get there, can't do it yet.
				continue
			}

			next := path{
				loc:    move.loc,
				dist:   cur.dist + move.dist,
				keyset: cur.keyset | move.key,
			}
			heap.Push(&pq, next)
		}
	}
}

func (g *grid) subdivide() []*grid {
	g.spaces[g.start] = space{tile: wall}
	for _, dir := range dirs {
		g.spaces[point{g.start.x + dir.x, g.start.y + dir.y}] = space{tile: wall}
	}

	ul := g.subgrid(point{0, 0}, g.start, point{-1, -1})
	ur := g.subgrid(point{g.start.x, 0}, point{math.MaxInt32, g.start.y}, point{1, -1})
	bl := g.subgrid(point{0, g.start.y}, point{g.start.x, math.MaxInt32}, point{-1, 1})
	br := g.subgrid(g.start, point{math.MaxInt32, math.MaxInt32}, point{1, 1})

	return []*grid{ul, ur, bl, br}
}

func (g *grid) subgrid(min, max, ds point) *grid {
	spaces := map[point]space{}
	keyset := 0
	start := point{g.start.x + ds.x, g.start.y + ds.y}

	// Filter out only stuff in the given quadrant.
	for loc, spc := range g.spaces {
		if loc.x >= min.x && loc.y >= min.y && loc.x <= max.x && loc.y <= max.y {
			spaces[loc] = spc
			if spc.tile == key {
				keyset |= spc.id
			}
		}
	}

	// Ignore doors that we can't get the key for.
	for loc, spc := range spaces {
		if spc.tile == door && spc.id&keyset == 0 {
			spaces[loc] = space{tile: empty}
		}
	}

	return &grid{spaces, keyset, start}
}

func main() {
	input := lib.Load("input.txt")
	{
		g := newgrid(input)
		path := g.traverse()
		fmt.Println(path)
	}

	{
		g := newgrid(input)
		subgrids := g.subdivide()

		total := 0
		for _, sg := range subgrids {
			total += sg.traverse()
		}

		fmt.Println(total)
	}
}
