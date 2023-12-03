package lib

import (
	"fmt"

	"github.com/creachadair/aoc/aoc"
)

type Schematic []string

func (s Schematic) Rows() int { return len(s) }
func (s Schematic) Cols() int { return len(s[0]) }

func (s Schematic) At(row, col int) byte {
	if row < 0 || row >= len(s) || col < 0 || col >= len(s[0]) {
		return '.'
	}
	return s[row][col]
}

func isSymbol(b byte) bool {
	if b >= '0' && b <= '9' {
		return false
	}
	return b != '.'
}

func (s Schematic) adjSymbol(row, col int) (r, c int) {
	for r := -1; r <= 1; r++ {
		for c := -1; c <= 1; c++ {
			if r == 0 && c == 0 {
				continue
			}
			if isSymbol(s.At(row+r, col+c)) {
				return row + r, col + c
			}
		}
	}
	return -1, -1
}

func (s Schematic) LabelOf(row, lo, hi int) (r, c int) {
	if row < 0 || row >= len(s) {
		return -1, -1
	}
	for col := lo; col < hi; col++ {
		if r, c := s.adjSymbol(row, col); r >= 0 && c >= 0 {
			return r, c
		}
	}
	return -1, -1
}

func (s Schematic) IsLabel(row, lo, hi int) bool {
	r, c := s.LabelOf(row, lo, hi)
	return r >= 0 && c >= 0
}

func ParseSchematic(input []byte) (Schematic, error) {
	var cols int
	var out Schematic
	for i, line := range aoc.SplitLines(input) {
		if cols == 0 {
			cols = len(line)
		} else if len(line) != cols {
			return nil, fmt.Errorf("line %d: want %d cols, got %d", i+1, cols, len(line))
		}
		out = append(out, line)
	}
	return out, nil
}
