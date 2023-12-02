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
	var minPath []int
	for r := 0; r < m.R; r++ {
		for c := 0; c < m.C; c++ {
			if m.GetRC(r, c) != 0 {
				continue
			}
			path, ok := m.BFS(r, c, m.NotTooHigh, m.IsEnd)
			if ok && (minPath == nil || len(path) < len(minPath)) {
				minPath = path
			}
		}
	}
	fmt.Println(len(minPath) - 1)

	m.Plot(minPath)
	fmt.Println(m)
}
