package lib

import (
	"errors"
	"fmt"
	"strings"

	"github.com/creachadair/aoc/aoc"
	"golang.org/x/exp/slices"
)

type Map struct {
	nr, nc int
	data   []byte
}

func (m *Map) Rows() int        { return m.nr }
func (m *Map) Cols() int        { return m.nc }
func (m *Map) At(r, c int) byte { return m.data[r*m.nc+c] }

func (m *Map) Expand() *Map {
	erows, ecols := m.EmptySpace()
	if len(erows) == 0 && len(ecols) == 0 {
		return m
	}
	var buf []byte
	for r := 0; r < m.nr; r++ {
		var row []byte
		for c := 0; c < m.nc; c++ {
			row = append(row, m.At(r, c))
			if slices.Contains(ecols, c) {
				row = append(row, m.At(r, c))
			}
		}
		buf = append(buf, row...)
		if slices.Contains(erows, r) {
			buf = append(buf, row...)
		}
	}
	return &Map{nr: m.nr + len(erows), nc: m.nc + len(ecols), data: buf}
}

func (m *Map) EmptySpace() (erows, ecols []int) {
nextRow:
	for r := 0; r < m.nr; r++ {
		for c := 0; c < m.nc; c++ {
			if m.At(r, c) != '.' {
				continue nextRow
			}
		}
		erows = append(erows, r)
	}
nextCol:
	for c := 0; c < m.nc; c++ {
		for r := 0; r < m.nr; r++ {
			if m.At(r, c) != '.' {
				continue nextCol
			}
		}
		ecols = append(ecols, c)
	}
	return
}

func (m *Map) String() string {
	var buf strings.Builder
	for r := 0; r < m.nr; r++ {
		base := r * m.nc
		buf.Write(m.data[base : base+m.nc])
		buf.WriteByte('\n')
	}
	return buf.String()
}

func ParseMap(input []byte) (*Map, error) {
	var nr, nc int
	var buf []byte
	for i, line := range aoc.SplitLines(aoc.MustReadInput()) {
		if i == 0 {
			nc = len(line)
		} else if len(line) != nc {
			return nil, fmt.Errorf("line %d: got %d columns, want %d", i+1, len(line), nc)
		}
		buf = append(buf, line...)
		nr++
	}
	if nr == 0 || nc == 0 {
		return nil, errors.New("invalid input map")
	}
	return &Map{nr: nr, nc: nc, data: buf}, nil
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func MDist(r1, c1, r2, c2 int) int { return abs(r2-r1) + abs(c2-c1) }
