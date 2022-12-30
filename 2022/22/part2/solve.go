package main

import (
	"flag"
	"fmt"

	"aoc/2022/22/grid"
)

var (
	inputFile = flag.String("input", "input.txt", "Input file path")
	doVerbose = flag.Bool("v", false, "Verbose output")
)

func main() {
	flag.Parse()

	m := grid.MustParseInput(*inputFile)
	fmt.Println(m.Tiles)

	fmt.Println(m)
	m.Tiles.Normalize(m)
	fmt.Println(m)
}
