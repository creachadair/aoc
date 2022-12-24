package main

import (
	"flag"
	"fmt"

	"aoc/2022/18/cubes"
)

var inputFile = flag.String("input", "input.txt", "Input file path")

func main() {
	flag.Parse()

	m := cubes.MustParse(*inputFile)
	fmt.Println(m.Sum())
}
