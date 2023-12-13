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
		if v := m.FindMirror(); v >= 0 {
			sum += v
			log.Printf("- V split col %d", v)
		}
		if v := m.Transpose().FindMirror(); v >= 0 {
			sum += 100 * v
			log.Printf("- H split at row %d", v)
		}
	}
	fmt.Println(sum)
}
