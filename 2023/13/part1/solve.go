package main

import (
	"flag"
	"fmt"
	"log"

	lib "github.com/creachadair/aoc/2023/13"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	ms, err := lib.ParseMaps(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Parse maps: %v", err)
	}
	var sum int
	for i, m := range ms {
		log.Printf("Map %d", i+1)
		for c := 0; c < m.Cols(); c++ {
			if m.IsMirror(c) {
				sum += c
				log.Printf("- V split col %d", c)
			}
		}
		mt := m.Transpose()
		for c := 0; c < mt.Cols(); c++ {
			if mt.IsMirror(c) {
				sum += 100 * c
				log.Printf("- H split at row %d", c)
			}
		}
	}
	fmt.Println(sum)
}
