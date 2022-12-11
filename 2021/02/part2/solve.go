package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"aoc/2021/02/sub"
)

var inputFile = flag.String("input", "input.txt", "Input file")

func main() {
	flag.Parse()
	input, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}

	pgm, err := sub.ParseProgram(string(input))
	if err != nil {
		log.Fatalf("Parse input: %v", err)
	}

	var m sub.Marine
	pgm.ApplyAimed(&m)
	fmt.Println(m.HPos * m.Depth)
}
