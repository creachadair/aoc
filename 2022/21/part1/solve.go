package main

import (
	"flag"
	"fmt"

	"aoc/2022/21/rules"
)

var inputFile = flag.String("input", "input.txt", "Input file path")

func main() {
	flag.Parse()

	rs := rules.MustParse(*inputFile)
	g := rules.NewGraph(rs)
	g.Solve()

	fmt.Println(g.Values["root"])
}
