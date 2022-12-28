package cipher

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/creachadair/mds/mlink"
)

type State []*mlink.Ring[int]

func MustParse(path string) State {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	values := make([]int, len(lines))
	for i, line := range lines {
		v, err := strconv.Atoi(line)
		if err != nil {
			log.Fatalf("Line %d: invalid input: %v", i+1, err)
		}
		values[i] = v
	}
	return NewState(values)
}

func NewState(values []int) State {
	r := mlink.RingOf(values...)
	out := make(State, len(values))
	for i, cur := 0, r; i < len(values); i++ {
		out[i] = cur
		cur = cur.Next()
	}
	return out
}

func (s State) Mix() {
	for i, cur := range s {
		s.Move(i, cur.Value)
	}
}

func (s State) findValue(v int) *mlink.Ring[int] {
	for _, cur := range s {
		if cur.Value == v {
			return cur
		}
	}
	return nil
}

func (s State) Values() []int {
	vs := make([]int, 0, len(s))
	s.findValue(0).Each(func(v int) bool {
		vs = append(vs, v)
		return true
	})
	return vs
}

func (s State) Move(pos, delta int) {
	if delta == 0 {
		return
	}
	delta %= len(s) - 1

	cur := s[pos]
	for i := delta; i > 0; i-- {
		cur = cur.Next()
	}
	for i := delta; i <= 0; i++ { // one extra to skip the start value itself
		cur = cur.Prev()
	}
	s[pos].Pop()
	cur.Join(s[pos].Pop())
}
