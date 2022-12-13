package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"aoc/2022/08/treemap"
)

var (
	inputFile = flag.String("input", "input.txt", "Input file (required)")
)

func main() {
	flag.Parse()
	input, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Reading input: %v", err)
	}

	var max int
	m := treemap.New(string(input))
	m.Each(func(r, c int) {
		if vs := m.ViewScore(r, c); vs > max {
			max = vs
		}
	})
	fmt.Println(max)
}
