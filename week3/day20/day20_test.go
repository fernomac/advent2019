package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fernomac/advent2019/lib"
	"github.com/fernomac/advent2019/lib/grid"
)

var maze1 = strings.TrimSpace(`
_________A___________
_________A___________
__#######.#########__
__#######.........#__
__#######.#######.#__
__#######.#######.#__
__#######.#######.#__
__#####__B____###.#__
BC...##__C____###.#__
__##.##_______###.#__
__##...DE__F__###.#__
__#####____G__###.#__
__#########.#####.#__
DE..#######...###.#__
__#.#########.###.#__
FG..#########.....#__
__###########.#####__
_____________Z_______
_____________Z_______`)

var maze2 = strings.TrimSpace(`
_____________Z_L_X_W_______C_________________
_____________Z_P_Q_B_______K_________________
__###########.#.#.#.#######.###############__
__#...#.......#.#.......#.#.......#.#.#...#__
__###.#.#.#.#.#.#.#.###.#.#.#######.#.#.###__
__#.#...#.#.#...#.#.#...#...#...#.#.......#__
__#.###.#######.###.###.#.###.###.#.#######__
__#...#.......#.#...#...#.............#...#__
__#.#########.#######.#.#######.#######.###__
__#...#.#____F_______R_I_______Z____#.#.#.#__
__#.###.#____D_______E_C_______H____#.#.#.#__
__#.#...#___________________________#...#.#__
__#.###.#___________________________#.###.#__
__#.#....OA_______________________WB..#.#..ZH
__#.###.#___________________________#.#.#.#__
CJ......#___________________________#.....#__
__#######___________________________#######__
__#.#....CK_________________________#......IC
__#.###.#___________________________#.###.#__
__#.....#___________________________#...#.#__
__###.###___________________________#.#.#.#__
XF....#.#_________________________RF..#.#.#__
__#####.#___________________________#######__
__#......CJ_______________________NM..#...#__
__###.#.#___________________________#.###.#__
RE....#.#___________________________#......RF
__###.###________X___X_______L______#.#.#.#__
__#.....#________F___Q_______P______#.#.#.#__
__###.###########.###.#######.#########.###__
__#.....#...#.....#.......#...#.....#.#...#__
__#####.#.###.#######.#######.###.###.#.#.#__
__#.......#.......#.#.#.#.#...#...#...#.#.#__
__#####.###.#####.#.#.#.#.###.###.#.###.###__
__#.......#.....#.#...#...............#...#__
__#############.#.#.###.###################__
_______________A_O_F___N_____________________
_______________A_A_D___M_____________________`)

func TestParse(t *testing.T) {
	g := parse(maze1)

	if val := g.At(2, 2).(*node); val.tile != wall {
		t.Error("expected wall at (2,2), got", val)
	}
	if val := g.At(2, 8).(*node); val.tile != empty {
		t.Error("expected empty at (2,8), got", val)
	}
	if val := g.At(1, 8).(*node); val.tile != warp {
		t.Error("expected warp at (1,8), got ", val)
	}

	pt28 := grid.Point{X: 2, Y: 8}
	if val := g.At(9, 7).(*node); val.dst != pt28 {
		t.Error("expected", pt28, "got", val)
	}

	pt96 := grid.Point{X: 9, Y: 6}
	if val := g.At(1, 8).(*node); val.dst != pt96 {
		t.Error("expected", pt96, "got", val)
	}
}

func TestParse2(t *testing.T) {
	g := parse(maze2)

	{
		n := g.At(43, 13).(*node)
		if n.tile != warp {
			t.Errorf("expected warp, got %v", n.tile)
		}
		if n.inner {
			t.Errorf("expected inner=false, got true")
		}
	}
	{
		n := g.At(23, 35).(*node)
		if n.tile != warp {
			t.Errorf("expected warp, got %v", n.tile)
		}
		if n.inner {
			t.Errorf("expected inner=false, got true")
		}
	}
}

func TestParse3(t *testing.T) {
	g := parse(lib.Load("input.txt"))

	start, _ := g.Point("AA")
	if start.X != 37 || start.Y != 2 {
		t.Errorf("expected {37 2}, got %v", start)
	}

	end, _ := g.Point("ZZ")
	if end.X != 83 || end.Y != 120 {
		t.Errorf("expected {83 120}, got %v", end)
	}

	if n := g.At(39, 121).(*node); n.tile != warp || n.inner {
		t.Errorf("expected inner(39,121)=false, got true")
	}
	if n := g.At(127, 81).(*node); n.tile != warp || n.inner {
		t.Errorf("expected inner(127,81)=false, got true")
	}
}

func TestMovesFrom(t *testing.T) {
	g := parse(maze1)

	test := func(from grid.Point, ems []move) {
		t.Run(fmt.Sprintf("%v", from), func(t *testing.T) {
			ms := movesFrom(g, from)
			if !moveseq(ms, ems) {
				t.Errorf("expected %v, got %v", ems, ms)
			}
		})
	}

	test(grid.Point{X: 9, Y: 2}, []move{
		move{grid.Point{X: 9, Y: 7}, 5},
		move{grid.Point{X: 13, Y: 16}, 26},
		move{grid.Point{X: 11, Y: 11}, 31},
	})

	test(grid.Point{X: 2, Y: 8}, []move{
		move{grid.Point{X: 7, Y: 10}, 7},
	})
}

func moveseq(a, b []move) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestWalk(t *testing.T) {
	test := func(name, maze string, eval int) {
		t.Run(name, func(t *testing.T) {
			g := parse(maze)
			val := walk(g)
			if val != eval {
				t.Errorf("expected %v, got %v", eval, val)
			}
		})
	}

	test("maze1", maze1, 23)
}

func TestWalk2(t *testing.T) {
	test := func(name, maze string, eval int) {
		t.Run(name, func(t *testing.T) {
			g := parse(maze)
			val := walk2(g)
			if val != eval {
				t.Errorf("expected %v, got %v", eval, val)
			}
		})
	}

	test("maze1", maze1, 26)
	test("maze2", maze2, 396)
}
