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
		r, c, kind, pos := findSmudge(m)
		log.Printf("Map %d has smudge at (%d, %d)", i+1, r, c)
		log.Printf(" - now has %c split at %d", kind, pos)
		switch kind {
		case 'V':
			sum += pos
		case 'H':
			sum += 100 * pos
		default:
			panic("invalid mirror")
		}
	}
	fmt.Println(sum)
}

func findSmudge(m *aoc.Map) (r, c int, kind byte, pos int) {
	// Find the initial split.
	lib.Mirrors(m, func(k byte, p int) { kind, pos = k, p })

	ok, op := kind, pos
	for r = 0; r < m.Rows(); r++ {
		for c = 0; c < m.Cols(); c++ {
			found := false
			lib.Toggle(m, r, c)
			lib.Mirrors(m, func(k byte, p int) {
				if k != ok || p != op {
					kind, pos = k, p
					found = true
				}
			})
			if found {
				return
			}
			lib.Toggle(m, r, c) // put it back
		}
	}
	return -1, -1, 0, 0
}
