package main

import (
	"flag"
	"fmt"
	"log"

	lib "github.com/creachadair/aoc/2023/05"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	a, err := lib.ParseAlmanac(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Parse almanac: %v", err)
	}
	minLoc := -1
	for _, s := range a.Seeds {
		data := a.Track(s)
		if minLoc < 0 || data["location"] < minLoc {
			minLoc = data["location"]
		}
		log.Printf("Seed %d: %v", s, data)
	}
	fmt.Println(minLoc)
}
