package lib

import (
	"fmt"
	"strings"

	"github.com/creachadair/aoc/aoc"
)

type Seq []int

func (s Seq) Project() int {
	if s.IsConst() {
		return s[len(s)-1]
	}
	ds := s.Diff().Project()
	return s[len(s)-1] + ds
}

func (s Seq) Reject() int {
	if s.IsConst() {
		return s[0]
	}
	ds := s.Diff().Reject()
	return s[0] - ds
}

func (s Seq) Diff() Seq {
	out := make(Seq, len(s)-1)
	for i := 0; i < len(out); i++ {
		out[i] = s[i+1] - s[i]
	}
	return out
}

func (s Seq) IsConst() bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[i-1] {
			return false
		}
	}
	return true
}

func ParseSeq(input []byte) ([]Seq, error) {
	var seqs []Seq
	for i, line := range aoc.SplitLines(input) {
		vs, err := aoc.ParseInts(strings.Fields(line))
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}
		seqs = append(seqs, Seq(vs))
	}
	return seqs, nil
}
