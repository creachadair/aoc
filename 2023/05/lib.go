package lib

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/creachadair/aoc/aoc"
)

type Range struct {
	Dst int
	Src int
	N   int
}

func (r Range) Map(in int) (int, bool) {
	if in < r.Src || in > r.Src+r.N-1 {
		return in, false
	}
	return r.Dst + (in - r.Src), true
}

type Ranges []Range

func (rs Ranges) Map(in int) int {
	for _, r := range rs {
		if out, ok := r.Map(in); ok {
			return out
		}
	}
	return in // unmapped
}

func parseRange(line string) (Range, error) {
	dsn, err := parseInts(strings.Fields(line))
	if err != nil {
		return Range{}, fmt.Errorf("invalid spec: %w", err)
	} else if len(dsn) != 3 {
		return Range{}, fmt.Errorf("got %d fields, want 3", len(dsn))
	}
	return Range{Dst: dsn[0], Src: dsn[1], N: dsn[2]}, nil
}

type Map struct {
	In, Out string
	Ranges
}

type Almanac struct {
	Seeds []int
	Maps  map[string]Map
}

func (a Almanac) Track(seed int) map[string]int {
	out := make(map[string]int)

	kind, cur := "seed", seed
	for {
		out[kind] = cur
		next, ok := a.Maps[kind]
		if !ok {
			break
		}
		kind, cur = next.Out, next.Map(cur)
	}
	return out
}

func ParseAlmanac(input []byte) (Almanac, error) {
	lines := aoc.SplitLines(input)
	if len(lines) == 0 {
		return Almanac{}, errors.New("empty input")
	}
	seeds, ok := strings.CutPrefix(lines[0], "seeds:")
	if !ok {
		return Almanac{}, errors.New("missing seeds list")
	}
	sv, err := parseInts(strings.Fields(seeds))
	if err != nil {
		return Almanac{}, fmt.Errorf("invalid seeds: %w", err)
	}
	out := Almanac{
		Seeds: sv,
		Maps:  make(map[string]Map),
	}
	var curMap Map
	push := func() {
		if curMap.In != "" {
			out.Maps[curMap.In] = curMap
			curMap = Map{}
		}
	}
	for i, line := range lines[1:] {
		if line == "" {
			push()
			continue
		} else if curMap.In == "" {
			name, ok := strings.CutSuffix(line, " map:")
			if !ok {
				return Almanac{}, fmt.Errorf("line %d: invalid map label %q", i+2, line)
			}
			ps := strings.SplitN(name, "-", 3)
			curMap.In, curMap.Out = ps[0], ps[2]
			continue
		}
		r, err := parseRange(line)
		if err != nil {
			return Almanac{}, fmt.Errorf("line %d: invalid range: %w", i+2, err)
		}
		curMap.Ranges = append(curMap.Ranges, r)
	}
	push()
	return out, nil
}

func parseInts(ss []string) ([]int, error) {
	out := make([]int, len(ss))
	for i, s := range ss {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("offset %d: %w", i, err)
		}
		out[i] = v
	}
	return out, nil
}
