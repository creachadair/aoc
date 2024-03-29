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

	ms, err := aoc.ParseMaps(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Parse maps: %v", err)
	}
	var sum int
	for i, m := range ms {
		lib.Mirrors(m, func(kind byte, pos int) {
			log.Printf("Map %d has %c split at %d", i+1, kind, pos)
			switch kind {
			case 'V':
				sum += pos
			case 'H':
				sum += 100 * pos
			default:
				panic("invalid mirror kind")
			}
		})
	}
	fmt.Println(sum)
}
