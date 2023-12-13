package lib

import (
	"errors"
	"fmt"
	"strings"

	"github.com/creachadair/aoc/aoc"
)

type Map struct {
	nr, nc int
	data   []byte
}

func (m *Map) Rows() int        { return m.nr }
func (m *Map) Cols() int        { return m.nc }
func (m *Map) At(r, c int) byte { return m.data[r*m.nc+c] }

func (m *Map) Transpose() *Map {
	out := &Map{nr: m.nc, nc: m.nr, data: make([]byte, len(m.data))}
	for r := 0; r < m.nr; r++ {
		for c := 0; c < m.nc; c++ {
			out.data[c*out.nc+r] = m.At(r, c)
		}
	}
	return out
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func (m *Map) IsMirror(c int) bool {
	w := min(c, m.nc-c)
	if w == 0 {
		return false
	}
	lo, hi := c-w, c+w-1
	for r := 0; r < m.nr; r++ {
		for i := 0; lo+i < hi-i; i++ {
			if m.At(r, lo+i) != m.At(r, hi-i) {
				return false
			}
		}
	}
	return true
}

func (m *Map) String() string {
	var buf strings.Builder
	for r := 0; r < m.nr; r++ {
		base := r * m.nc
		buf.Write(m.data[base : base+m.nc])
		buf.WriteByte('\n')
	}
	return buf.String()
}

func ParseMaps(input []byte) ([]*Map, error) {
	lines := aoc.SplitLines(input)
	var out []*Map

	pos := 0
	for pos < len(lines) {
		i := pos
		for i < len(lines) && lines[i] != "" {
			i++
		}
		m, err := parseMap(pos, lines[pos:i])
		if err != nil {
			return nil, err
		}
		out = append(out, m)
		pos = i + 1
	}
	return out, nil
}

func parseMap(start int, lines []string) (*Map, error) {
	var nr, nc int
	var buf []byte
	for i, line := range lines {
		if i == 0 {
			nc = len(line)
		} else if len(line) != nc {
			return nil, fmt.Errorf("line %d: got %d columns, want %d", start+i+1, len(line), nc)
		}
		buf = append(buf, line...)
		nr++
	}
	if nr == 0 || nc == 0 {
		return nil, errors.New("invalid input map")
	}
	return &Map{nr: nr, nc: nc, data: buf}, nil
}
