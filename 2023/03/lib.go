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

func (s Schematic) NearSymbol(row, col int) bool {
	for r := -1; r <= 1; r++ {
		for c := -1; c <= 1; c++ {
			if r == 0 && c == 0 {
				continue
			}
			if isSymbol(s.At(row+r, col+c)) {
				return true
			}
		}
	}
	return false
}

func (s Schematic) IsLabel(row, lo, hi int) bool {
	if row < 0 || row >= len(s) {
		return false
	}
	for col := lo; col < hi; col++ {
		if s.NearSymbol(row, col) {
			return true
		}
	}
	return false
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
