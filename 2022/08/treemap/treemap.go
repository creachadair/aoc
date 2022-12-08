package treemap

import (
	"strings"
)

type Map struct {
	elts   string
	nr, nc int
}

func New(input string) Map {
	nc := strings.Index(input, "\n")
	nr := strings.Count(input, "\n")
	return Map{
		elts: strings.ReplaceAll(input, "\n", ""),
		nr:   nr,
		nc:   nc,
	}
}

func (m Map) Rows() int        { return m.nr }
func (m Map) Cols() int        { return m.nc }
func (m Map) Row(r int) string { p := (r - 1) * m.nc; return m.elts[p : p+m.nc] }

func (m Map) Col(c int) string {
	b := make([]byte, m.nr)
	for p, i := 0, c-1; i < len(m.elts); i += m.nc {
		b[p] = m.elts[i]
		p++
	}
	return string(b)
}

func (m Map) Visible(r, c int) bool {
	return isVisible(m.Row(r), c) || isVisible(m.Col(c), r)
}

func (m Map) ViewScore(r, c int) int {
	lhs, rhs := visDistance(m.Row(r), c)
	top, bot := visDistance(m.Col(c), r)
	return lhs * rhs * top * bot
}

func (m Map) Each(f func(r, c int)) {
	for r := 1; r <= m.Rows(); r++ {
		for c := 1; c <= m.Cols(); c++ {
			f(r, c)
		}
	}
}

func isVisible(rc string, pos int) bool {
	return allLess(rc[:pos-1], rc[pos-1]) || allLess(rc[pos:], rc[pos-1])
}

func visDistance(rc string, pos int) (fore, aft int) {
	v := rc[pos-1]
	for i := pos - 2; i >= 0; i-- {
		fore++
		if rc[i] >= v {
			break
		}
	}
	for i := pos; i < len(rc); i++ {
		aft++
		if rc[i] >= v {
			break
		}
	}
	return
}

func allLess(s string, b byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= b {
			return false
		}
	}
	return true
}
