package aoc

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/creachadair/mds/slice"
)

// NewMap constructs a new map from the specified data.
func NewMap(nr, nc int, data []byte) *Map {
	if len(data) != nr*nc {
		panic(fmt.Sprintf("invalid map size: got %d, want %d", len(data), nr*nc))
	}
	return &Map{nr: nr, nc: nc, data: data}
}

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
func (m *Map) Data() string         { return string(m.data) }

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
	// In-place transposition:
	//
	// The item at position i moves to min(R, C)*i % (RC-1).
	// This results in max(R, C) permutation cycles to chase.
	// The elements at i=0 and i=(RC-1) do not move.
	rcm := len(m.data) - 1 // == (nr * nc) - 1
	min := min(m.nr, m.nc)

	// We need to keep track of which positions are already permuted, so that we
	// can find the next cycle anchor. Keep one bit per byte of data.
	bs := make([]uint64, (rcm+1+63)/64) // round up
	isSet := func(i int) bool { return bs[i/64]&(1<<(i%64)) != 0 }
	set := func(i int) { bs[i/64] |= 1 << (i % 64) }
	next := func(i int) int { return (min * i) % rcm }

	// Skip 0 and (rc-1) since those items never move.
	for i := 1; i < rcm; i++ {
		if isSet(i) {
			continue // this position was part of an earlier cycle
		}

		// Lift the item at the start of the cycle and rotate it through the
		// positions of the permutation. When we get back to here, drop the last
		// item we picked up back into place at the start.
		cur := m.data[i]
		for j := next(i); j != i; j = next(j) {
			m.data[j], cur = cur, m.data[j]
			set(j)
		}
		m.data[i] = cur
		set(i)
	}
	m.nr, m.nc = m.nc, m.nr
	return m
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
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
