package lib

import (
	"errors"
	"fmt"
	"strings"

	"github.com/creachadair/aoc/aoc"
	"github.com/creachadair/mds/mapset"
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

func (g Grid) Mark(r, c int) { g.data[r*g.nc+c] = '*' }

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
	Start      Cell
	Max        int
	StartShape byte
	Path       mapset.Set[Cell]
}

func (g *Grid) FindLoop(r, c int) Loop {
	a, b, startShape := g.seeds(r, c)
	seen := mapset.New[Cell]()
	seen.Add(Cell{r, c}, a, b)

	dist := 1
	for a != b {
		dist++

		for _, x := range g.exits(a) {
			if !seen.Has(x) {
				a = x
			}
		}
		for _, x := range g.exits(b) {
			if !seen.Has(x) {
				b = x
			}
		}
		seen.Add(a, b)
	}
	return Loop{Start: Cell{r, c}, Max: dist, StartShape: startShape, Path: seen}
}

func (g *Grid) IsInside(loop Loop, r, c int) bool {
	// All the paths from a fully-enclosed position to any edge of the map must
	// cross the boundary an odd number of times.
	//
	// However, "crossing" is slightly subtle: Say we are currently outside the
	// figure, traverse a section of pipe, say:
	//
	//     .F---7. >> direction of travel
	//
	// then this should not count as a "crossing" because we did not enter the
	// interior of the figure. However:
	//
	//     .F---J. >> direction of travel
	//
	// should count as a "crossing" because we exited into the interior as we
	// departed from J. So F-J and L-7 transitions cross, F-7 and L-J do not.
	if loop.Path.Has(Cell{r, c}) {
		return false // path elements are not contained by the path
	}

	var cross int
	var pathStart byte
	for h := c; h < g.nc; h++ {
		cur := byte('.')
		if loop.Path.Has(Cell{r, h}) {
			cur = g.At(r, h)

			// To simplify the logic below, treat "S" as whatever shape we
			// inferred when constructing the loop.
			if cur == 'S' {
				cur = loop.StartShape
			}
		}

		// Case 1: Traversing a horizontal edge.
		if pathStart != 0 {
			// Check whether the path ends here.
			if isEnd := strings.IndexByte("LJF7", cur) >= 0; isEnd {
				// If so, check whether the path twisted and count a crossing if so.
				if (pathStart == 'L' && cur == '7') || (pathStart == 'F' && cur == 'J') {
					cross++
				}
				pathStart = 0
			} else if cur != '-' {
				panic(fmt.Sprintf("(%d,%d:%d) impossible in value: %c", r, c, h, cur))
			}
			continue
		}

		// Case 2: Not traversing a horizontal edge.
		if cur == '|' {
			cross++
		} else if cur == 'L' || cur == 'F' {
			pathStart = cur // entering a horizontal edge
		} else if cur != '.' && cur != 'S' {
			panic(fmt.Sprintf("(%d,%d:%d) impossible out value: %c", r, c, h, cur))
		}
	}
	return cross%2 == 1
}

type Cell [2]int

const (
	up byte = 1 << iota
	down
	left
	right
)

var startShape = []byte{up | right: 'L', up | down: '|', up | left: 'J', down | left: '7', down | right: 'F', right | left: '-'}

func (g *Grid) seeds(r, c int) (a, b Cell, start byte) {
	var out []Cell
	var dir byte
	if p := g.At(r-1, c); p == '|' || p == 'F' || p == '7' {
		dir |= up
		out = append(out, Cell{r - 1, c})
	}
	if p := g.At(r+1, c); p == '|' || p == 'J' || p == 'L' {
		dir |= down
		out = append(out, Cell{r + 1, c})
	}
	if p := g.At(r, c-1); p == '-' || p == 'F' || p == 'L' {
		dir |= left
		out = append(out, Cell{r, c - 1})
	}
	if p := g.At(r, c+1); p == '-' || p == '7' || p == 'J' {
		dir |= right
		out = append(out, Cell{r, c + 1})
	}
	// N.B. This will panic if the grid does not have a coherent loop.
	return out[0], out[1], startShape[dir]
}

func (g *Grid) exits(cell Cell) []Cell {
	r, c := cell[0], cell[1]
	switch g.At(r, c) {
	case '-':
		return []Cell{{r, c - 1}, {r, c + 1}}
	case '|':
		return []Cell{{r - 1, c}, {r + 1, c}}
	case 'F':
		return []Cell{{r, c + 1}, {r + 1, c}}
	case 'L':
		return []Cell{{r - 1, c}, {r, c + 1}}
	case '7':
		return []Cell{{r, c - 1}, {r + 1, c}}
	case 'J':
		return []Cell{{r - 1, c}, {r, c - 1}}
	default:
		return nil
	}
}

func (g *Grid) CleanString(loop Loop) string {
	var buf strings.Builder
	for r := 0; r < g.nr; r++ {
		for c := 0; c < g.nc; c++ {
			cur := g.At(r, c)
			if loop.Path.Has(Cell{r, c}) || cur == '*' {
				buf.WriteByte(cur)
			} else {
				buf.WriteByte('.')
			}
		}
		buf.WriteByte('\n')
	}
	return buf.String()
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
