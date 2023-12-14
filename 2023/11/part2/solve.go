package main

import (
	"flag"
	"fmt"
	"log"

	lib "github.com/creachadair/aoc/2023/11"
	"github.com/creachadair/aoc/aoc"
)

// To solve part1, set --expand=2
// For the example, try --expand=10 and --expand=100
var expandAmt = flag.Int("expand", 1000000, "Expansion factor")

func main() {
	flag.Parse()

	m, err := aoc.ParseMap(aoc.SplitLines(aoc.MustReadInput()))
	if err != nil {
		log.Fatalf("Parse map: %v", err)
	}
	fmt.Printf("input\n%s\n", m)

	// With the much higher expansion factor, instead of explicitly expanding
	// the map we'll just apply the effect to the positions of the nodes.
	erows, ecols := lib.EmptySpace(m)

	type gxy struct {
		r, c int
		dist []int
	}
	var gx []*gxy
	for r := 0; r < m.Rows(); r++ {
		for c := 0; c < m.Cols(); c++ {
			if m.At(r, c) == '#' {
				// Adjust the effective location by the empty space above and to
				// the left of the original location.
				er, ec := r, c
				for _, erow := range erows {
					if erow < r {
						er += *expandAmt - 1
					}
				}
				for _, ecol := range ecols {
					if ecol < c {
						ec += *expandAmt - 1
					}
				}
				gx = append(gx, &gxy{er, ec, nil})
			}
		}
	}
	var sum int
	for i, g := range gx {
		for _, o := range gx[i+1:] {
			d := lib.MDist(g.r, g.c, o.r, o.c)
			g.dist = append(g.dist, d)
			log.Printf("From #%d (%d, %d) to (%d, %d) d=%d", i+1, g.r, g.c, o.r, o.c, d)
			sum += d
		}
	}
	fmt.Println(sum)
}
