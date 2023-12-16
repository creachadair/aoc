package lib

import (
	"fmt"

	"github.com/creachadair/aoc/aoc"
	"github.com/creachadair/mds/mapset"
	"github.com/creachadair/mds/mlink"
)

type Beam struct {
	R, C   int
	RV, CV int // velocity
}

func (b *Beam) Step() {
	b.R += b.RV
	b.C += b.CV
}

func (b *Beam) String() string {
	var c byte
	if b.CV < 0 {
		c = '<'
	} else if b.CV > 0 {
		c = '>'
	} else if b.RV < 0 {
		c = '^'
	} else {
		c = 'v'
	}
	return fmt.Sprintf("[%d,%d]%c", b.R, b.C, c)
}

func CountEnergized(m *aoc.Map) int {
	// Avoid repeating beams.
	seen := mapset.New[Beam]()
	beams := mlink.NewQueue[*Beam]()
	add := func(r, c, rv, cv int) {
		beams.Add(&Beam{r, c, rv, cv})
	}
	add(0, 0, 0, 1) // top-left, facing right

	type cell struct{ r, c int }
	lit := mapset.New[cell]()

	for !beams.IsEmpty() {
		beam, _ := beams.Pop()
		if seen.Has(*beam) {
			continue
		}
		seen.Add(*beam)
		aoc.Dprintf("MJF :: start %v", beam)
	nextBeam:
		for m.InBounds(beam.R, beam.C) {
			aoc.Dprintf("MJF :: at %v", beam)
			lit.Add(cell{beam.R, beam.C})
			switch c := m.At(beam.R, beam.C); c {
			case '|':
				if beam.CV != 0 { // moving across the splitter
					aoc.Dprintf("MJF :: %v vsplit", beam)
					add(beam.R-1, beam.C, -1, 0) // up
					add(beam.R+1, beam.C, 1, 0)  // down
					break nextBeam
				}
			case '-':
				if beam.RV != 0 { // moving across the splitter
					aoc.Dprintf("MJF :: %v hsplit", beam)
					add(beam.R, beam.C-1, 0, -1) // left
					add(beam.R, beam.C+1, 0, 1)  // right
					break nextBeam
				}
			case '/':
				beam.RV, beam.CV = -beam.CV, -beam.RV // reflect clockwise
			case '\\':
				beam.RV, beam.CV = beam.CV, beam.RV // reflect counterclockwise
			default:
				// ignore
			}
			beam.Step()
		}
		aoc.Dprintf("MJF :: %v done %d left", beam, beams.Len())
	}
	return lit.Len()
}
