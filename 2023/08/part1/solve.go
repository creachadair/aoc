package main

import (
	"flag"
	"fmt"
	"log"

	lib "github.com/creachadair/aoc/2023/08"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	pgm, err := lib.ParseProgram(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Parse program: %v", err)
	}
	fmt.Println(pgm.Steps("AAA", "ZZZ"))
}
