package main

import (
	"flag"
	"fmt"
	"log"

	lib "github.com/creachadair/aoc/2023/11"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	m, err := lib.ParseMap(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Parse map: %v", err)
	}
	fmt.Printf("input\n%s\n", m)
	exp := m.Expand()
	fmt.Printf("expand\n%s\n", exp)

	type gxy struct {
		r, c int
		dist []int
	}
	var gx []*gxy
	for r := 0; r < exp.Rows(); r++ {
		for c := 0; c < exp.Cols(); c++ {
			if exp.At(r, c) == '#' {
				gx = append(gx, &gxy{r, c, nil})
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
