package grid

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func MustParseInput(path string) *Map {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}
	m, err := ParseMap(string(data))
	if err != nil {
		log.Fatalf("Parsing map: %v", err)
	}
	return m
}

// A Map is a map grid read from the puzzle input.
type Map struct {
	Data   []byte // packed rows
	NR, NC int    // rows, columns
	Spec   []Action

	overlay map[int]byte
}

func (m *Map) norm(r, c int) (nr, nc int) {
	return ((r-1)%m.NR + m.NR) % m.NR, ((c-1)%m.NC + m.NC) % m.NC
}

func (m *Map) pos(r, c int) int {
	nr, nc := m.norm(r, c)
	return nr*m.NC + nc
}

// Start returns the starting row and column (1-indexed).
func (m *Map) Start() (r, c int) {
	for c := 1; c <= m.NC; c++ {
		if m.At(1, c) == '.' {
			return 1, c
		}
	}
	panic("unreached")
}

// At reports the contents of the map at r, c (1-indexed).
func (m *Map) At(r, c int) byte { return m.Data[m.pos(r, c)] }

// Plot adds b to the map overlay at r, c (1-indexed).
func (m *Map) Plot(r, c int, b byte) {
	if m.overlay == nil {
		m.overlay = make(map[int]byte)
	}
	m.overlay[m.pos(r, c)] = b
}

// Unplot removes any overlay markings at r, c (1-indexed).
func (m *Map) Unplot(r, c int) {
	if m.overlay != nil {
		delete(m.overlay, m.pos(r, c))
	}
}

func (m *Map) String() string {
	var buf strings.Builder

	for r := 0; r < m.NR; r++ {
		start := r * m.NC
		row := []byte(string(m.Data[start : start+m.NC]))
		for c := 0; c < m.NC; c++ {
			if v, ok := m.overlay[m.pos(r+1, c+1)]; ok {
				row[c] = v
			}
		}
		buf.Write(row)
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')
	buf.WriteString(actionString(m.Spec))
	return buf.String()
}

func ParseMap(text string) (*Map, error) {
	lines := strings.Split(text, "\n")
	rule := lines[len(lines)-2]
	lines = lines[:len(lines)-2]

	// First count the number of rows and columns, so we know how long to make
	// each block in the data array.
	var nc, nr int
	for _, line := range lines {
		if line == "" {
			break
		} else if len(line) > nc {
			nc = len(line)
		}
		nr++
	}
	data := make([]byte, 0, nr*nc)
	for _, line := range lines {
		if line == "" {
			break
		}
		data = append(data, []byte(line)...)
		if len(line) < nc {
			data = append(data, bytes.Repeat([]byte(" "), nc-len(line))...)
		}
	}
	spec, err := ParseActions(rule)
	if err != nil {
		return nil, err
	}
	return &Map{Data: data, NR: nr, NC: nc, Spec: spec}, nil
}

type Action struct {
	Op string
	N  int
}

var actionRE = regexp.MustCompile(`(?:(\d+)|([RL]))`)

func ParseActions(text string) ([]Action, error) {
	ms := actionRE.FindAllStringSubmatch(text, -1)
	if ms == nil {
		return nil, errors.New("invalid action spec")
	}
	out := make([]Action, len(ms))
	for i, m := range ms {
		if m[2] != "" {
			out[i] = Action{Op: m[2]}
		} else if v, err := strconv.Atoi(m[1]); err != nil {
			return nil, fmt.Errorf("invalid action spec: %w", err)
		} else {
			out[i] = Action{N: v}
		}
	}
	return out, nil
}

func actionString(as []Action) string {
	var buf strings.Builder
	for _, a := range as {
		if a.Op != "" {
			buf.WriteString(a.Op)
		} else {
			buf.WriteString(strconv.Itoa(a.N))
		}
	}
	return buf.String()
}
