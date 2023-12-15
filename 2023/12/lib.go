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

// Solve counts the number of distinct assignments of wildcards in the pattern
// of r satisfying the specified groups.
func Solve(r Record) int {
	// Memoize solutions to avoid recomputing them when backtracking.
	type skey struct {
		pat  string // the pattern suffix
		v    int    // the number of "#" marks seen so far
		rest int    // the number of incomplete groups remaining
	}
	done := make(map[skey]int) // value is solution count

	var solve func(string, int, []int) int
	solve = func(pat string, seen int, want []int) (out int) {
		key := skey{pat, seen, len(want)}
		v, ok := done[key]
		if !ok {
			// Base case when the pattern runs out.
			if pat == "" {
				// There are no more groups to consider, or the last group is full.
				if len(want) == 0 && seen == 0 || len(want) == 1 && seen == want[0] {
					return 1
				}
				// The pattern ended with unfinished groups.
				return 0
			}

			// Check for some boundary conditions to simplify the logic below.
			var isFull, isPartial bool
			if len(want) != 0 {
				if seen == want[0] {
					// The current group is full. Record this fact so we can check
					// the validity of advancing in the pattern, and update to the
					// next desired group (if any).
					isFull = true
					seen, want = 0, want[1:]
				} else if seen > 0 {
					// A "partial" group is one where we have seen at least one
					// mark, but it is not yet full.
					isPartial = true
				}
			}

			switch pat[0] {
			case '#':
				if !isFull {
					v = solve(pat[1:], seen+1, want)
				}
				// Otherwise: The group is overfull.

			case '.':
				if !isPartial {
					v = solve(pat[1:], seen, want)
				}
				// Otherwise: The group is incomplete.

			case '?':
				var dot, hash int
				if !isFull {
					// A "#" is possible, check for solutions beginning there.
					hash = solve("#"+pat[1:], seen, want)
				}
				if !isPartial {
					// A "." is possible, check for solutions beginning there.
					dot = solve("."+pat[1:], seen, want)
				}
				v = dot + hash

			default:
				panic(fmt.Sprintf("invalid pattern word %c", pat[0]))
			}
			done[key] = v
		}
		return v
	}
	return solve(r.Pattern, 0, r.Groups)
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
