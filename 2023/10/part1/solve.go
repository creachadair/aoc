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
	fmt.Println(g.FindLoop(g.Start()))
}