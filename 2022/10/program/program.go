package program

import (
	"strconv"
	"strings"
)

type Insn struct {
	Op  string
	Arg int
}

func Eval(pgm []Insn, f func(tick, x int)) {
	tick, x := 1, 1
	for _, insn := range pgm {
		switch insn.Op {
		case "noop":
			tick++
			f(tick, x)

		case "addx":
			tick++
			f(tick, x)

			x += insn.Arg

			tick++
			f(tick, x)

		default:
			panic("unknown instruction: " + insn.Op)
		}
	}
}

func Parse(text string) []Insn {
	var insns []Insn
	for _, line := range strings.Split(strings.TrimSpace(text), "\n") {
		op, arg, ok := strings.Cut(line, " ")
		insn := Insn{Op: op}
		if ok {
			insn.Arg, _ = strconv.Atoi(arg)
		}
		insns = append(insns, insn)
	}
	return insns
}
