package lib

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/creachadair/aoc/aoc"
)

type Insn struct {
	Label string
	L, R  string
}

var insnRE = regexp.MustCompile(`^(\w+) = \((\w+), (\w+)\)$`)

func parseInsn(line string) (Insn, error) {
	m := insnRE.FindStringSubmatch(line)
	if m == nil {
		return Insn{}, errors.New("invalid instruction")
	}
	return Insn{
		Label: m[1],
		L:     m[2],
		R:     m[3],
	}, nil
}

type Pgm struct {
	Scheme string
	Insn   []Insn

	next map[string]*Insn
}

func (p *Pgm) Steps(from, to string) int {
	cur := p.next[from]
	if cur == nil {
		panic("origin not found: " + from)
	}
	pos, ns := 0, 0
	for cur.Label != to {
		var next string
		switch p.Scheme[pos] {
		case 'L':
			next = cur.L
		case 'R':
			next = cur.R
		default:
			panic("invalid step: " + string(p.Scheme[pos]))
		}
		log.Printf("At %q go %c to %q", cur.Label, p.Scheme[pos], next)
		cur = p.next[next]
		if cur == nil {
			panic("step not found: " + next)
		}
		ns++
		pos = (pos + 1) % len(p.Scheme)
	}
	return ns
}

func ParseProgram(input []byte) (*Pgm, error) {
	lines := aoc.SplitLines(input)
	if len(lines) < 3 || lines[1] != "" {
		return nil, errors.New("invalid program format")
	}
	var insn []Insn
	for i, line := range lines[2:] {
		p, err := parseInsn(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+3)
		}
		insn = append(insn, p)
	}
	next := make(map[string]*Insn)
	for i, p := range insn {
		next[p.Label] = &insn[i]
	}
	return &Pgm{Scheme: lines[0], Insn: insn, next: next}, nil
}
