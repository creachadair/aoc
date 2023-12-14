package aoc

import (
	"errors"
	"fmt"
	"strings"
)

// A Map represents a rectangular grid of one-byte glyphs.
// Each row is encoded on one line.
// This is a very common format for AoC puzzle inputs.
type Map struct {
	nr, nc int
	data   []byte
}

func (m *Map) Rows() int            { return m.nr }
func (m *Map) Cols() int            { return m.nc }
func (m *Map) At(r, c int) byte     { return m.data[r*m.nc+c] }
func (m *Map) Set(r, c int, b byte) { m.data[r*m.nc+c] = b }

// Transpose returns a copy of m in which the rows and columns of m are
// exchanged (transposed) around the primary diagonal.
func (m *Map) Transpose() *Map {
	out := &Map{nr: m.nc, nc: m.nr, data: make([]byte, len(m.data))}
	for r := 0; r < m.nr; r++ {
		for c := 0; c < m.nc; c++ {
			out.data[c*out.nc+r] = m.At(r, c)
		}
	}
	return out
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

// ParseMaps parses an input comprising one or more maps separated by blank
// lines.
func ParseMaps(input []byte) ([]*Map, error) {
	lines := SplitLines(input)
	var out []*Map

	pos := 0
	for pos < len(lines) {
		i := pos
		for i < len(lines) && lines[i] != "" {
			i++
		}
		m, err := ParseMap(lines[pos:i])
		if err != nil {
			return nil, err
		}
		out = append(out, m)
		pos = i + 1
	}
	return out, nil
}

// ParseMap parses a single map from a sequence of input lines.
func ParseMap(lines []string) (*Map, error) {
	var nr, nc int
	var buf []byte
	for i, line := range lines {
		if i == 0 {
			nc = len(line)
		} else if len(line) != nc {
			return nil, fmt.Errorf("line %d: got %d columns, want %d", i+1, len(line), nc)
		}
		buf = append(buf, line...)
		nr++
	}
	if nr == 0 || nc == 0 {
		return nil, errors.New("invalid input map")
	}
	return &Map{nr: nr, nc: nc, data: buf}, nil
}
