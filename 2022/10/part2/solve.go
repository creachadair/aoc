package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/creachadair/aoc/2022/10/program"
)

var (
	inputFile = flag.String("input", "input.txt", "Input file path")
)

func main() {
	flag.Parse()
	input, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}

	code := program.Parse(string(input))

	var crt CRT
	program.Eval(code, crt.paint)
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
