package lib

import (
	"github.com/creachadair/aoc/aoc"
	"golang.org/x/exp/slices"
)

func Expand(m *aoc.Map) *aoc.Map {
	erows, ecols := EmptySpace(m)
	if len(erows) == 0 && len(ecols) == 0 {
		return m
	}
	var buf []byte
	for r := 0; r < m.Rows(); r++ {
		var row []byte
		for c := 0; c < m.Cols(); c++ {
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
	return aoc.NewMap(m.Rows()+len(erows), m.Cols()+len(ecols), buf)
}

func EmptySpace(m *aoc.Map) (erows, ecols []int) {
nextRow:
	for r := 0; r < m.Rows(); r++ {
		for c := 0; c < m.Cols(); c++ {
			if m.At(r, c) != '.' {
				continue nextRow
			}
		}
		erows = append(erows, r)
	}
nextCol:
	for c := 0; c < m.Cols(); c++ {
		for r := 0; r < m.Rows(); r++ {
			if m.At(r, c) != '.' {
				continue nextCol
			}
		}
		ecols = append(ecols, c)
	}
	return
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func MDist(r1, c1, r2, c2 int) int { return abs(r2-r1) + abs(c2-c1) }
