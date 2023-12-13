package lib

import (
	"fmt"
	"strings"

	"github.com/creachadair/aoc/aoc"
)

type Record struct {
	Pattern string
	Groups  []int
}

func (r Record) fillCount() int {
	var n int
	for _, g := range r.Groups {
		n += g
	}
	return len(r.Pattern) - n
}

// SumsOf calls f with each combination if k integers that sum to n.  The slice
// passed to f is only valid for the duration of the call.
func SumsOf(k, n int, f func([]int)) {
	var sums func(int, []int, func())
	sums = func(n int, buf []int, ready func()) {
		if len(buf) == 1 {
			buf[0] = n
			ready()
			return
		}
		for i := n; i >= 0; i-- {
			buf[0] = i
			sums(n-i, buf[1:], ready)
		}
	}
	buf := make([]int, k)
	sums(n, buf, func() { f(buf) })
}

func (r Record) at(i int, buf []int) byte {
	pos := 0
	for {
		if i < buf[pos] {
			return '.'
		}
		i -= buf[pos]
		if i < r.Groups[pos] {
			return '#'
		}
		i -= r.Groups[pos]
		pos++
	}
}

func (r Record) toString(buf []int) string {
	out := strings.Repeat(".", buf[0])
	for i, g := range r.Groups {
		out += strings.Repeat("#", g)
		out += strings.Repeat(".", buf[i+1])
	}
	return out
}

func (r Record) Solve(f func(soln string)) {
	// n+1 gaps, 1 + (n-1) + 1
	// spc is the amount of space we have to fill
	spc := r.fillCount()
	if spc < 0 {
		panic("impossible record")
	}

	// If we have enough spacers, it is possible there may be spaces before and
	// after all the required groups.
	SumsOf(len(r.Groups)+1, spc, func(buf []int) {
		// There must be at least one spacer between groups, so eliminate
		// combinations that have a zero in the intermediate position.
		for i := 1; i+1 < len(buf); i++ {
			if buf[i] == 0 {
				return
			}
		}

		// Check whether the proposed result is compatible with the pattern,
		// meaning at each position the solution matches the pattern, or at that
		// position is a wildcard.
		//
		//   pat: #.??.#
		//   ---
		//   ok:  #.#..#
		//   bad: ..##.#
		//
		for i := 0; i < len(r.Pattern); i++ {
			got, c := r.at(i, buf), r.Pattern[i]
			if got != c && c != '?' {
				return // incompatible
			}
		}
		f(r.toString(buf))
	})
}

func ParseRecords(input []byte) ([]Record, error) {
	var recs []Record
	for i, line := range aoc.SplitLines(input) {
		fs := strings.Fields(line)
		if len(fs) != 2 {
			return nil, fmt.Errorf("line %d: got %d fields, want 2", i+1, len(fs))
		}
		gs, err := aoc.ParseInts(strings.Split(fs[1], ","))
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid groups: %w", i+1, err)
		}
		recs = append(recs, Record{
			Pattern: fs[0],
			Groups:  gs,
		})
	}
	return recs, nil
}
