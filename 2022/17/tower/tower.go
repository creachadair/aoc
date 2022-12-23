package tower

import (
	"fmt"
	"sort"
	"strings"
)

type Tower struct {
	nc     int
	base   int
	jets   string
	step   int
	shapes []*Shape
}

// New returns a new empty towers of nc columns and using the given pattern of
// shift jets ("<" and ">" permitted).
func New(width int, jets string) *Tower { return &Tower{nc: width, jets: jets} }

func (t *Tower) blow(s *Shape) {
	switch t.jets[t.step%len(t.jets)] {
	case '<':
		s.left(t.At)
	case '>':
		s.right(t.nc, t.At)
	default:
		panic("invalid jet spec")
	}
	t.step++
}

// Drop moves s to the current starting location and then updates its position
// until it can no longer move, then adds it to the tower.
func (t *Tower) Drop(s *Shape) {
	s.move(t.Start())
	for {
		t.blow(s)
		if !s.down(t.base, t.At) {
			break
		}
	}
	t.shapes = append(t.shapes, s)
	for i := len(t.shapes) - 1; i > 0 && t.shapes[i].isBefore(t.shapes[i-1]); i-- {
		t.shapes[i], t.shapes[i-1] = t.shapes[i-1], t.shapes[i]
	}
}

// Step returns the current time step.
func (t *Tower) Step() int { return t.step }

// Phase returns the current time phase.
func (t *Tower) Phase() int { return t.step % len(t.jets) }

// At reports the contents of the tower at the given x, y coordinate.
// The origin (0, 0) is at the bottom left of the tower.
func (t *Tower) At(x, y int) byte {
	// Binary search for the highest candidate shape. We can't raise the lower
	// bound because there could be multiple shapes with the same height bounds.
	hi := len(t.shapes)
	for hi > 0 {
		mid := hi / 2
		if t.shapes[mid].cy <= y {
			break
		}
		hi = mid
	}
	for _, s := range t.shapes[:hi] {
		if s.At(x, y) == '#' {
			return '#'
		}
	}
	return '.'
}

// String renders the current contents of the tower as a string.
func (t *Tower) String() string {
	var buf strings.Builder
	for r := t.Height() - 1; r >= t.base; r-- {
		buf.WriteByte('|')
		for c := 0; c < t.nc; c++ {
			buf.WriteByte(t.At(c, r))
		}
		buf.WriteString("|\n")
	}
	if t.base > 0 {
		fmt.Fprintf(&buf, "~%s~ %d\n", strings.Repeat(" ", t.nc), t.base)
	}
	buf.WriteString("+" + strings.Repeat("-", t.nc) + "+")
	return buf.String()
}

// Height reports the number of occupied rows of the tower.
func (t *Tower) Height() int {
	var highest int
	for _, s := range t.shapes {
		if t := s.Top(); t > highest {
			highest = t
		}
	}
	return highest
}

// Start returns the coordinates of the next starting position.
func (t *Tower) Start() (int, int) { return 2, t.Height() + 3 }

// Trim finds the highest level in t at which each column is filled at some
// higher level, and discards all shapes ending below that level.
func (t *Tower) Trim() {
	row := (1 << t.nc) - 1 // a 1 bit for each column
	for i := t.Height() - 1; i >= t.base; i-- {
		for c := 0; c < t.nc; c++ {
			if row&(1<<c) != 0 && t.At(c, i) == '#' {
				row &^= (1 << c)
			}

			// If we have seen an entry in each column, we can discard anything
			// earlier in the sequence.
			if row == 0 {
				t.base = i
				n := sort.Search(len(t.shapes), func(i int) bool {
					return t.shapes[i].Top() >= t.base
				})
				if n < len(t.shapes) {
					nc := copy(t.shapes, t.shapes[n:])
					t.shapes = t.shapes[:nc]
				}
				return
			}
		}
	}
}

type Shape struct {
	nr, nc int    // bounding box
	cx, cy int    // lower-left corner
	fill   string // w*h contents
}

// Top returns the y coordinate of the row just above s.
func (s *Shape) Top() int { return s.cy + s.nr }

func (s *Shape) At(x, y int) byte {
	rx, ry := x-s.cx, y-s.cy
	if rx >= 0 && rx < s.nc && ry >= 0 && ry < s.nr {
		pos := ry*s.nc + rx
		return s.fill[pos]
	}
	return 0
}

func (s *Shape) move(x, y int) *Shape { s.cx = x; s.cy = y; return s }

func (s *Shape) left(at func(x, y int) byte) {
	if s.cx > 0 && s.canShift(-1, 0, at) {
		s.cx--
	}
}

func (s *Shape) right(max int, at func(x, y int) byte) {
	if s.cx+s.nc < max && s.canShift(1, 0, at) {
		s.cx++
	}
}

func (s *Shape) down(base int, at func(x, y int) byte) bool {
	if s.cy > base && s.canShift(0, -1, at) {
		s.cy--
		return true
	}
	return false
}

func (s *Shape) canShift(dx, dy int, at func(x, y int) byte) bool {
	for pos, v := range s.fill {
		nx, ny := pos%s.nc+s.cx+dx, pos/s.nc+s.cy+dy
		if v != '.' && at(nx, ny) == '#' {
			return false
		}
	}
	return true
}

func (s *Shape) isBefore(s2 *Shape) bool {
	return s.cy < s2.cy || (s.cy == s2.cy && s.Top() < s2.Top())
}

const (
	line  = "####"
	plus  = ".#.###.#."
	angle = "###..#..#"
	box   = "####"
)

func HLine() *Shape { return &Shape{1, 4, 0, 0, line} }
func VLine() *Shape { return &Shape{4, 1, 0, 0, line} }
func Plus() *Shape  { return &Shape{3, 3, 0, 0, plus} }
func Angle() *Shape { return &Shape{3, 3, 0, 0, angle} }
func Box() *Shape   { return &Shape{2, 2, 0, 0, box} }

// CycleLen is the number of shapes per cycle.
const CycleLen = 5

// NewSequence returns a function that constructs shapes in the standard
// sequence.  Ecah call constructs a fresh shape.
func NewSequence() func() *Shape {
	seq := []func() *Shape{HLine, Plus, Angle, VLine, Box}
	i := 0
	return func() *Shape {
		out := seq[i]()
		i = (i + 1) % len(seq)
		return out
	}
}
