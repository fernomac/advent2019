package grid

import (
	"strings"
)

// A Point on a grid.
type Point struct {
	X, Y int
}

// Neighbors returns the neighbors of this point.
func (p Point) Neighbors() []Point {
	return []Point{
		{p.X + 1, p.Y},
		{p.X - 1, p.Y},
		{p.X, p.Y + 1},
		{p.X, p.Y - 1},
	}
}

// A Node on the grid.
type Node interface {
	// ID returns the ID of this node.
	ID() string
}

// A Grid of points.
type Grid struct {
	grid  map[Point]Node
	index map[string]Point
}

// New directly creates a new grid.
func New(grid map[Point]Node) *Grid {
	index := map[string]Point{}
	for p, n := range grid {
		if n.ID() != "" {
			index[n.ID()] = p
		}
	}
	return &Grid{grid, index}
}

// Parse parses a grid from basic text form.
func Parse(in string, nf func(Point, rune) Node) *Grid {
	grid := map[Point]Node{}
	index := map[string]Point{}

	lines := strings.Split(strings.TrimSpace(in), "\n")
	for y, l := range lines {
		for x, r := range l {
			node := nf(Point{x, y}, r)
			grid[Point{x, y}] = node
			if node.ID() != "" {
				index[node.ID()] = Point{x, y}
			}
		}
	}

	return &Grid{grid, index}
}

// At returns the node at the given point.
func (g *Grid) At(x, y int) Node {
	return g.AtPt(Point{x, y})
}

// AtPt returns the node at the given point.
func (g *Grid) AtPt(p Point) Node {
	return g.grid[p]
}

// Point returns the point with the given ID.
func (g *Grid) Point(id string) (Point, bool) {
	p, ok := g.index[id]
	return p, ok
}

// NodesWhere returns the set of nodes that match the given predicate.
func (g *Grid) NodesWhere(matches func(Node) bool) []Node {
	ret := []Node{}
	for _, n := range g.grid {
		if matches(n) {
			ret = append(ret, n)
		}
	}
	return ret
}
