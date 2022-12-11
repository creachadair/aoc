package sub

import (
	"fmt"
	"strconv"
	"strings"
)

type Marine struct {
	HPos  int
	Depth int
	Aim   int // approximately: declination
}

type Action struct {
	Tag string
	Arg int
}

type Program []Action

func (p Program) ApplyAbs(m *Marine) {
	for _, act := range p {
		switch act.Tag {
		case "forward":
			m.HPos += act.Arg
		case "down":
			m.Depth += act.Arg
		case "up":
			m.Depth -= act.Arg
			if m.Depth < 0 {
				panic("you're out of your depth")
			}
		default:
			panic("invalid tag: " + act.Tag)
		}
	}
}

func (p Program) ApplyAimed(m *Marine) {
	for _, act := range p {
		switch act.Tag {
		case "forward":
			m.HPos += act.Arg
			m.Depth += act.Arg * m.Aim
			if m.Depth < 0 {
				panic("you're out of your depth")
			}
		case "down":
			m.Aim += act.Arg
		case "up":
			m.Aim -= act.Arg
		default:
			panic("invalid tag: " + act.Tag)
		}
	}
}

func ParseProgram(s string) (Program, error) {
	var out []Action
	for i, line := range strings.Split(strings.TrimSpace(s), "\n") {
		fs := strings.Fields(line)
		if len(fs) != 2 {
			return nil, fmt.Errorf("line %d: invalid format %q", i+1, line)
		}
		v, err := strconv.Atoi(fs[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid argument: %w", i+1, err)
		}
		out = append(out, Action{Tag: fs[0], Arg: v})
	}
	return out, nil
}
