package lib

import (
	"github.com/creachadair/aoc/aoc"
)

type Schematic struct{ *aoc.Map }

func (s Schematic) At(row, col int) byte {
	if s.Map.InBounds(row, col) {
		return s.Map.At(row, col)
	}
	return '.'
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
	if row < 0 || row >= s.Map.Rows() {
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
	m, err := aoc.ParseMap(aoc.SplitLines(input))
	if err != nil {
		return Schematic{}, err
	}
	return Schematic{Map: m}, nil
}
