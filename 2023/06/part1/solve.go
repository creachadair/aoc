package main

import (
	"flag"
	"fmt"
	"log"

	lib "github.com/creachadair/aoc/2023/06"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	races, err := lib.ParseRaces(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Parse races: %v", err)
	}
	prod := 1
	for i, r := range races {
		min := lib.MinGreater(r.Time, r.Dist)
		max := r.Time - min
		nd := (max + 1) - min
		log.Printf("Race %d: %d ≤ i ≤ %d (%d)", i+1, min, max, nd)
		prod *= nd
	}
	fmt.Println(prod)
}
