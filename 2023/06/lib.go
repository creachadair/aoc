package lib

import (
	"errors"
	"fmt"
	"strings"

	"github.com/creachadair/aoc/aoc"
)

type Race struct {
	Time int
	Dist int
}

// Dist computes the distance travelled in a race of max time with the given
// initial hold.
//
// The maximum possible occurs at the midpoint, and is symmetric about that.
//
//	0 (1*(max-1)) (2*(max-2)) (3*(max-3)) ... ((max-2)*2) ((max-1)*1 0
//
// So you can binary search in the interval (0..(max/2)) to find a breakpoint
// where the distance exceeds some threshold (if one exists).
func Dist(hold, max int) int {
	if hold == 0 || hold >= max {
		return 0
	}
	return hold * (max - hold)
}

func MinGreater(max, record int) int {
	lo, hi := 0, max/2
	for lo < hi {
		hold := (lo + hi) / 2
		dist := Dist(hold, max)
		if dist <= record {
			lo = hold + 1
		} else {
			hi = hold
		}
	}
	if Dist(lo, max) > record {
		return lo
	}
	return -1
}

func ParseRaces(input []byte) ([]Race, error) {
	lines := aoc.SplitLines(input)
	if len(lines) != 2 {
		return nil, fmt.Errorf("found %d lines, want 2", len(lines))
	}
	ts, ok := strings.CutPrefix(lines[0], "Time:")
	if !ok {
		return nil, errors.New("missing times")
	}
	ds, ok := strings.CutPrefix(lines[1], "Distance:")
	if !ok {
		return nil, errors.New("missing distances")
	}
	times, err := aoc.ParseInts(strings.Fields(ts))
	if err != nil {
		return nil, fmt.Errorf("times: %w", err)
	}
	dists, err := aoc.ParseInts(strings.Fields(ds))
	if err != nil {
		return nil, fmt.Errorf("distances: %w", err)
	}
	if len(times) != len(dists) {
		return nil, fmt.Errorf("got %d times but %d distances", len(times), len(dists))
	}
	var races []Race
	for i := 0; i < len(times); i++ {
		races = append(races, Race{
			Time: times[i],
			Dist: dists[i],
		})
	}
	return races, nil
}
