package lib

import (
	"errors"

	"github.com/creachadair/aoc/aoc"
	"github.com/creachadair/mds/mapset"
	"github.com/creachadair/mds/mlink"
)

type Grid struct {
	nr, nc int
	data   []byte // r1 ... rc r2 ... rc r3 ... rc ...
}

func (g *Grid) Rows() int { return g.nr }
func (g *Grid) Cols() int { return g.nc }

func (g *Grid) At(r, c int) byte {
	if r < 0 || r >= g.nr || c < 0 || c >= g.nc {
		return '.'
	}
	return g.data[r*g.nc+c]
}

func (g *Grid) Start() (r, c int) {
	for r = 0; r < g.nr; r++ {
		for c = 0; c < g.nc; c++ {
			if g.At(r, c) == 'S' {
				return
			}
		}
	}
	return -1, -1
}

type Loop struct {
	Start Cell
	Max   int
}

func (g *Grid) FindLoop(r, c int) Loop {
	type cdist struct {
		Cell
		D int
	}
	seen := mapset.New[Cell]()
	queue := mlink.NewQueue[cdist]()
	queue.Add(cdist{Cell: Cell{r, c}})
	seen.Add(Cell{r, c})
	var max int
	for !queue.IsEmpty() {
		next, _ := queue.Pop()
		if next.D > max {
			max = next.D
		}
		for _, exit := range g.exits(next.Cell[0], next.Cell[1]) {
			if !seen.Has(exit) {
				queue.Add(cdist{Cell: exit, D: next.D + 1})
				seen.Add(exit)
			}
		}
	}
	return Loop{Start: Cell{r, c}, Max: max}
}

type Cell [2]int

func (g *Grid) exits(r, c int) []Cell {
	var out []Cell
	if p := g.At(r, c-1); p == '-' || p == 'L' || p == 'F' {
		out = append(out, Cell{r, c - 1})
	}
	if p := g.At(r-1, c); p == '|' || p == '7' || p == 'F' {
		out = append(out, Cell{r - 1, c})
	}
	if p := g.At(r, c+1); p == '-' || p == '7' || p == 'J' {
		out = append(out, Cell{r, c + 1})
	}
	if p := g.At(r+1, c); p == '|' || p == 'J' || p == 'L' {
		out = append(out, Cell{r + 1, c})
	}
	return out
}

func ParseGrid(input []byte) (*Grid, error) {
	lines := aoc.SplitLines(input)
	if len(lines) == 0 {
		return nil, errors.New("empty grid")
	}
	out := &Grid{
		nr:   len(lines),
		nc:   len(lines[0]),
		data: make([]byte, 0, len(lines)*len(lines[0])),
	}
	for _, line := range lines {
		out.data = append(out.data, line...)
	}
	return out, nil
}
