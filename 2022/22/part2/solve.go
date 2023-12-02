package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/creachadair/aoc/2022/22/grid"
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
		log.Printf("POS NOW %v,%v %v in %v", c.CR, c.CC, c.Head, c.Tile)
		switch act.Op {
		case "R":
			log.Print("TURN RIGHT")
			c.R()
		case "L":
			log.Print("TURN LEFT")
			c.L()
		case "":
			log.Printf("FORWARD %d", act.N)
			c.CF(act.N)
		}
		c.Plot()
	}

	if *doVerbose {
		fmt.Println(m)
	}
	fmt.Println("Final:", c)

	score := 1000*c.CR + 4*c.CC + int(c.Head)
	fmt.Println("Score:", score)
}
