package lib

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/creachadair/aoc/aoc"
)

var doDebug = flag.Bool("debug", false, "Enable debug logging")

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

func (p *Pgm) Find(label string) Insn {
	if c := p.next[label]; c != nil {
		return *c
	}
	panic("missing label: " + label)
}

func dprintf(msg string, args ...any) {
	if *doDebug {
		log.Printf(msg, args...)
	}
}

func (p *Pgm) Steps(from, to string) (int, string) {
	var cur *Insn
	for _, c := range p.Insn {
		if strings.HasSuffix(c.Label, from) {
			cur = p.next[c.Label]
			break
		}
	}
	if cur == nil {
		panic("origin not found: " + from)
	}
	pos, ns := 0, 0
	for !strings.HasSuffix(cur.Label, to) {
		var next string
		switch p.Scheme[pos] {
		case 'L':
			next = cur.L
		case 'R':
			next = cur.R
		default:
			panic("invalid step: " + string(p.Scheme[pos]))
		}
		dprintf("At %q go %c to %q", cur.Label, p.Scheme[pos], next)
		cur = p.next[next]
		if cur == nil {
			panic("step not found: " + next)
		}
		ns++
		pos = (pos + 1) % len(p.Scheme)
	}
	return ns, cur.Label
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
