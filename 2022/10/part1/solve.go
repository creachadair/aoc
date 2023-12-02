package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/creachadair/aoc/2022/10/program"
)

var (
	inputFile  = flag.String("input", "input.txt", "Input file path")
	traceTicks = flag.String("trace", "20,60,100,140,180,220",
		"Comma-separated instruction ticks to trace")
)

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

	code := program.Parse(string(input))
	trace := mustParseTraces(*traceTicks)
	var total int
	program.Eval(code, func(tick, x int) {
		if trace[tick] {
			total += tick * x
			fmt.Printf("tick %d x=%d [%d]\n", tick, x, tick*x)
		}
	})
	fmt.Printf("sum: %d\n", total)
}
