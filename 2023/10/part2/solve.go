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
	fmt.Println(g.CleanString(loop))
	for r := 0; r < g.Rows(); r++ {
		for c := 0; c < g.Cols(); c++ {
			if g.IsInside(loop, r, c) {
				numIn++
				log.Printf("(%d,%d)", r, c)
			}
		}
	}
	fmt.Println(numIn)
}
