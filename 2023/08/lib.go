package lib

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/creachadair/aoc/aoc"
)

type Insn struct {
	Label string
	L, R  string
}

var insnRE = regexp.MustCompile(`^(\w+) = \((\w+), (\w+)\)$`)

func parseInsn(line string) (Insn, error) {
	var out Insn
	if err := aoc.Scanx(insnRE, line, &out.Label, &out.L, &out.R); err != nil {
		return out, fmt.Errorf("invalid instruction: %w", err)
	}
	return out, nil
}

type Pgm struct {
	Scheme string
	Insn   []Insn

	next map[string]*Insn
}

func (p *Pgm) Find(label string) Insn {
	if c := p.next[label]; c != nil {
		return *c
	}
	panic("missing label: " + label)
}

type Iter struct {
	pgm *Pgm
	cur *Insn
	ns  int
	pos int
}

func (it *Iter) Insn() Insn           { return *it.cur }
func (it *Iter) State() (int, string) { return it.ns, it.cur.Label }

func (it *Iter) Next() (int, string) {
	var next string
	dir := it.pgm.Scheme[it.pos]
	switch dir {
	case 'L':
		next = it.cur.L
	case 'R':
		next = it.cur.R
	default:
		panic("invalid step: " + string(dir))
	}
	aoc.Dprintf("step %d: from %q go %c to %q", it.ns, it.cur.Label, dir, next)
	it.cur = it.pgm.next[next]
	if it.cur == nil {
		panic("node not found: " + next)
	}
	it.ns++
	it.pos = (it.pos + 1) % len(it.pgm.Scheme)
	return it.ns, next
}

func (p *Pgm) Iter(from string) *Iter {
	cur := p.next[from]
	if cur == nil {
		panic("missing label: " + from)
	}
	return &Iter{pgm: p, cur: cur}
}

func (p *Pgm) Steps(from, to string) *Iter {
	it := p.Iter(from)
	for !strings.HasSuffix(it.Insn().Label, to) {
		it.Next()
	}
	return it
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
			return nil, fmt.Errorf("line %d: %w", i+3, err)
		}
		insn = append(insn, p)
	}
	next := make(map[string]*Insn)
	for i, p := range insn {
		next[p.Label] = &insn[i]
	}
	return &Pgm{Scheme: lines[0], Insn: insn, next: next}, nil
}
