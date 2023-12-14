package aoc

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/creachadair/mds/slice"
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

// Clone returns a copy of m with the same content.
func (m *Map) Clone() *Map {
	return &Map{nr: m.nr, nc: m.nc, data: bytes.Clone(m.data)}
}

// Equal reports whether the two maps are equal.
func (m *Map) Equal(o *Map) bool {
	return m.nr == o.nr && m.nc == o.nc && bytes.Equal(m.data, o.data)
}

// Transpose transposes the rows and columns of m in-place around its primary
// diagonal. It returns m to permit chaining.
func (m *Map) Transpose() *Map {
	// We could do this in-place by blocks, but it's simpler to write the copy.
	data := make([]byte, len(m.data))
	for r := 0; r < m.nr; r++ {
		for c := 0; c < m.nc; c++ {
			data[c*m.nr+r] = m.At(r, c)
		}
	}
	m.nr, m.nc, m.data = m.nc, m.nr, data
	return m
}

// FlipH flips m horizontally in-place. It returns m to permit chaining.
func (m *Map) FlipH() *Map {
	for r := 0; r < m.nr; r++ {
		base := r * m.nc
		slice.Reverse(m.data[base : base+m.nc])
	}
	return m
}

// FlipV flips m vertically in-place. It returns m to permit chaining.
func (m *Map) FlipV() *Map { return m.Transpose().FlipH().Transpose() }

// Rotate rotates m one quarter turn clockwise in-place.  It returns m to
// permit chaining.
func (m *Map) Rotate() *Map { return m.Transpose().FlipH() }

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
