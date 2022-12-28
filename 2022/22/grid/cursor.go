package grid

import (
	"fmt"
)

type Cursor struct {
	Map    *Map
	Head   Heading
	CR, CC int
}

func NewCursor(m *Map) *Cursor {
	sr, sc := m.Start()
	return &Cursor{Map: m, CR: sr, CC: sc, Head: Right}
}

func (c *Cursor) R() { c.Head = c.Head.R() }
func (c *Cursor) L() { c.Head = c.Head.L() }

func (c *Cursor) F(n int) {
	dhv := hdelta[c.Head]

	cr, cc := c.CR, c.CC
	for n > 0 {
		cr += dhv[0]
		cc += dhv[1]
		q := c.Map.At(cr, cc)
		switch q {
		case '.':
			// clear, keep moving
			n--
		case ' ':
			// skip, but don't consume a step
			continue
		case '#':
			// wall, stop
			return
		}

		// Pin the location to positive coordinates
		nr, nc := c.Map.norm(cr, cc)
		c.CR, c.CC = nr+1, nc+1
	}
}

func (c *Cursor) Plot() { c.Map.Plot(c.CR, c.CC, c.Head.String()[0]) }

func (c *Cursor) String() string {
	return fmt.Sprintf("Cursor(%d,%d|%q)", c.CR, c.CC, c.Head)
}

type Heading byte

const (
	headings = ">V<^"
	nhead    = len(headings)
)

const (
	Right Heading = iota
	Down
	Left
	Up
)

var hdelta = [...][2]int{Right: {0, 1}, Down: {1, 0}, Left: {0, -1}, Up: {-1, 0}}

func (h Heading) String() string { return string(headings[h]) }
func (h Heading) R() Heading     { return Heading((int(h) + 1) % nhead) }
func (h Heading) L() Heading     { return Heading(((int(h)-1)%nhead + nhead) % nhead) }
