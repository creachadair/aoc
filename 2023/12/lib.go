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

func Solve(r Record) int {
	type skey struct {
		pat  string
		v    int
		rest int
	}
	done := make(map[skey]int)
	var solve func(string, int, []int) int
	solve = func(pat string, seen int, want []int) (out int) {
		key := skey{pat, seen, len(want)}
		v, ok := done[key]
		if !ok {
			if pat == "" {
				if len(want) == 0 && seen == 0 || len(want) == 1 && seen == want[0] {
					return 1 // pattern complete
				}
				return 0 // pattern ended with unfinished groups
			}
			var isFull, isPartial bool
			if len(want) != 0 {
				if seen == want[0] {
					isFull = true
					seen, want = 0, want[1:]
				} else if seen > 0 {
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
				if isPartial {
					return 0 // incomplete group
				}
				v = solve(pat[1:], seen, want)

			case '?':
				var dot, hash int
				if !isFull {
					hash = solve("#"+pat[1:], seen, want) // substitute "#"
				}
				if !isPartial {
					dot = solve("."+pat[1:], seen, want) // substitute "."
				}
				v = dot + hash // cached by caller

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
