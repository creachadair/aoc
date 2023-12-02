package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/creachadair/aoc/2022/08/treemap"
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

	var nv int
	m := treemap.New(string(input))
	m.Each(func(r, c int) {
		if m.Visible(r, c) {
			nv++
		}
	})
	fmt.Println(nv)
}
