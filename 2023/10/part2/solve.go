package main

import (
	"flag"
	"fmt"
	"log"

	lib "github.com/creachadair/aoc/2023/10"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	g, err := lib.ParseGrid(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Parse grid: %v", err)
	}
	var numIn int
	loop := g.FindLoop(g.Start())
	for r := 0; r < g.Rows(); r++ {
		for c := 0; c < g.Cols(); c++ {
			if g.IsInside(loop, r, c) {
				numIn++
				g.Mark(r, c)
			}
		}
	}
	fmt.Println(g.CleanString())
	fmt.Println(numIn)
}
