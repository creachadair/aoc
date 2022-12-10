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
	inputFile = flag.String("input", "input.txt", "Input file path")
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

func main() {
	flag.Parse()
	input, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}

	program := mustParseInsn(string(input))

	var crt CRT
	eval(program, crt.paint)
	fmt.Println(crt)
}

const (
	crtWidth = 40 // pixels
	crtRows  = 6
)

type CRT [crtRows][crtWidth]bool

func (c *CRT) paint(tick, x int) {
	tick = (tick - 1) % 240 // pin to the refresh cycle
	row, col := tick/crtWidth, tick%crtWidth
	lit := col >= x-1 && col <= x+1
	c[row][col] = lit
}

var pix = map[bool]byte{false: '.', true: '#'}

func (c CRT) String() string {
	rows := make([]string, crtRows)
	for i, row := range c {
		cols := make([]byte, crtWidth)
		for j, col := range row {
			cols[j] = pix[col]
		}
		rows[i] = string(cols)
	}
	return strings.Join(rows, "\n")
}
