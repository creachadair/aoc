package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	inputFile  = flag.String("input", "", "Input file path")
	traceTicks = flag.String("trace", "20,60,100,140,180,220",
		"Comma-separated instruction ticks to trace")
)

type Insn struct {
	Op  string
	Arg int
}

func eval(pgm []Insn, f func(tick, x int)) {
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

func mustParseInsn(text string) []Insn {
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

func mustParseTraces(csv string) map[int]bool {
	if csv == "" {
		return nil
	}
	m := make(map[int]bool)
	for _, tick := range strings.Split(csv, ",") {
		v, err := strconv.Atoi(tick)
		if err != nil {
			log.Fatalf("Invalid trace tick %q: %v", tick, err)
		}
		m[v] = true
	}
	return m
}

func main() {
	flag.Parse()
	input, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}

	program := mustParseInsn(string(input))
	trace := mustParseTraces(*traceTicks)
	var total int
	eval(program, func(tick, x int) {
		if trace[tick] {
			total += tick * x
			fmt.Printf("tick %d x=%d [%d]\n", tick, x, tick*x)
		}
	})
	fmt.Printf("sum: %d\n", total)
}
