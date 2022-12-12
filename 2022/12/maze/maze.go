package maze

import (
	"strings"
)

type M struct {
	Grid   []byte // packed rows
	R, C   int    // number of rows, columns
	SR, SC int    // start row, column
	ER, EC int    // end row, column
}

func (m *M) rcToPos(r, c int) int { return r*m.C + c }

func (m *M) posToRC(pos int) (r, c int) { return pos / m.C, pos % m.C }

func (m *M) inBounds(r, c int) bool { return r >= 0 && r < m.R && c >= 0 && c < m.C }

func (m *M) String() string {
	var buf strings.Builder
	for i, ch := range m.Grid {
		switch i {
		case m.rcToPos(m.SR, m.SC):
			buf.WriteByte('S')
		case m.rcToPos(m.ER, m.EC):
			buf.WriteByte('E')
		default:
			buf.WriteByte(ch + 'a')
		}
		if (i+1)%m.C == 0 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func (m *M) GetRC(r, c int) byte { return m.Grid[m.rcToPos(r, c)] }

func (m *M) BFS(sr, sc int, adj func(r, c, or, oc int) bool, goal func(r, c int) bool) ([]int, bool) {
	seen := make(map[int]bool)

	q := [][]int{{m.rcToPos(sr, sc)}} // inner slice is reverse path to origin
	addIfNew := func(r, c int, tail []int) {
		pos := m.rcToPos(r, c)
		if !seen[pos] {
			seen[pos] = true
			q = append(q, concat(pos, tail))
		}
	}

	for len(q) != 0 {
		next := q[0]
		q = q[:copy(q, q[1:])]

		r, c := m.posToRC(next[0])
		if goal(r, c) {
			return next, true
		}
		if nr, nc := r-1, c; m.inBounds(nr, nc) && adj(nr, nc, r, c) {
			addIfNew(nr, nc, next)
		}
		if nr, nc := r+1, c; m.inBounds(nr, nc) && adj(nr, nc, r, c) {
			addIfNew(nr, nc, next)
		}
		if nr, nc := r, c-1; m.inBounds(nr, nc) && adj(nr, nc, r, c) {
			addIfNew(nr, nc, next)
		}
		if nr, nc := r, c+1; m.inBounds(nr, nc) && adj(nr, nc, r, c) {
			addIfNew(nr, nc, next)
		}
	}
	return nil, false
}

func (m *M) NotTooHigh(r, c, or, oc int) bool {
	return m.GetRC(r, c) <= m.GetRC(or, oc)+1
}

func (m *M) IsEnd(r, c int) bool { return r == m.ER && c == m.EC }

func concat(v int, path []int) []int { return append([]int{v}, path...) }

func Parse(input string) *M {
	rows := strings.Fields(input)
	nr, nc := len(rows), len(rows[0])
	m := &M{
		Grid: make([]byte, nr*nc),
		R:    nr,
		C:    nc,
	}
	pos := 0
	for i, row := range rows {
		for j, ch := range row {
			if ch == 'S' {
				m.SR, m.SC = i, j
				ch = 'a' // S is at altitude "a"
			} else if ch == 'E' {
				m.ER, m.EC = i, j
				ch = 'z' // E is at altitude "z"
			}
			m.Grid[pos] = byte(ch - 'a') // normalize to zero
			pos++
		}
	}

	return m
}
