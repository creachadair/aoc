package maze

import (
	"strings"

	"github.com/creachadair/mds/mlink"
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
			if ch > 127 {
				buf.WriteByte(' ')
			} else {
				buf.WriteByte((ch & 0x7f) + 'a')
			}
		}
		if (i+1)%m.C == 0 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

// Plot marks each grid cell in path so that the string output will highlight
// those cells.
func (m *M) Plot(path []int) {
	for _, pos := range path {
		m.Grid[pos] |= 0x80
	}
}

// GetRC returns the altitude at the specified coordinates.
func (m *M) GetRC(r, c int) byte { return m.Grid[m.rcToPos(r, c)] }

// BFS executes a breadth-first search starting at sr, sc.
// The adj function reports whether (or, oc) is a valid neighbor of (r, c).
// The goal function reports whether (r, c) is a valid goal state.
//
// It reports whether a path was found, and if so gives the sequence of grid
// positions leading from the goal back to the origin.
func (m *M) BFS(sr, sc int, adj func(r, c, or, oc int) bool, goal func(r, c int) bool) ([]int, bool) {
	seen := make(map[int]bool)

	q := mlink.NewQueue[[]int]()
	q.Add([]int{m.rcToPos(sr, sc)}) // slice is reverse path to origin
	addIfNew := func(r, c int, tail []int) {
		pos := m.rcToPos(r, c)
		if !seen[pos] {
			seen[pos] = true
			q.Add(concat(pos, tail))
		}
	}

	for !q.IsEmpty() {
		next, _ := q.Pop()

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

// NotTooHigh accepts (or, oc) as a neighbor of (r, c) if its altitude is not
// more than one unit higher than (r, c).
func (m *M) NotTooHigh(r, c, or, oc int) bool {
	return m.GetRC(r, c) <= m.GetRC(or, oc)+1
}

// IsEnd reports whether r, c is the end cell.
func (m *M) IsEnd(r, c int) bool { return r == m.ER && c == m.EC }

func concat(v int, path []int) []int { return append([]int{v}, path...) }

// Parse parses an altitude map.
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
