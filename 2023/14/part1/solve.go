package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	m, err := aoc.ParseMap(aoc.SplitLines(aoc.MustReadInput()))
	if err != nil {
		log.Fatalf("Parse map: %v", err)
	}
	var sum int
	for c := 0; c < m.Cols(); c++ {
		for r := 0; r < m.Rows(); r++ {
			if m.At(r, c) == 'O' {
				// Slide this rock up till it hits something.
				u := r
				for u > 0 {
					if m.At(u-1, c) != '.' {
						break
					}
					m.Set(u-1, c, 'O')
					m.Set(u, c, '.')
					u--
				}
				sum += m.Rows() - u // attribute mass
			}
		}
	}
	fmt.Println(sum)
}
