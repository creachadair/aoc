package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/creachadair/aoc/2022/12/maze"
)

var inputFile = flag.String("input", "input.txt", "Input file path")

func main() {
	flag.Parse()
	input, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}

	m := maze.Parse(string(input))
	path, ok := m.BFS(m.SR, m.SC, m.NotTooHigh, m.IsEnd)
	if !ok {
		log.Fatal("no path")
	}
	fmt.Println(len(path) - 1) // -1 for steps, not positions
}
