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

	c := grid.NewCursor(m)
	for _, act := range m.Spec {
		switch act.Op {
		case "R":
			c.R()
		case "L":
			c.L()
		case "":
			c.F(act.N)
		}
		c.Plot()
	}

	if *doVerbose {
		fmt.Println(m)
	}
	score := 1000*c.CR + 4*c.CC + int(c.Head)
	fmt.Println("Score:", score)
}
