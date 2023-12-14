package lib

import (
	"github.com/creachadair/aoc/aoc"
)

func Toggle(m *aoc.Map, r, c int) {
	if m.At(r, c) == '.' {
		m.Set(r, c, '#')
	} else {
		m.Set(r, c, '.')
	}
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func Mirrors(m *aoc.Map, f func(byte, int)) {
	FindMirror(m, func(c int) { f('V', c) })
	FindMirror(m.Clone().Transpose(), func(r int) { f('H', r) })
}

// FindMirror reports the columns of m that have a mirror split.
func FindMirror(m *aoc.Map, f func(int)) {
	for c := 0; c < m.Cols(); c++ {
		if IsMirror(m, c) {
			f(c)
		}
	}
}

func IsMirror(m *aoc.Map, c int) bool {
	w := min(c, m.Cols()-c)
	if w == 0 {
		return false
	}
	lo, hi := c-w, c+w-1
	for r := 0; r < m.Rows(); r++ {
		for i := 0; lo+i < hi-i; i++ {
			if m.At(r, lo+i) != m.At(r, hi-i) {
				return false
			}
		}
	}
	return true
}
