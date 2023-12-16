package main

import (
	"flag"
	"fmt"
	"log"

	lib "github.com/creachadair/aoc/2023/16"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()
	m, err := aoc.ParseMap(aoc.MustReadLines())
	if err != nil {
		log.Fatalf("Parse map: %v", err)
	}
	fmt.Println(m)

	var max int
	for r := 0; r < m.Rows(); r++ {
		// Left edge, facing right.
		if n := lib.CountEnergizedFrom(m, r, 0, 0, 1); n > max {
			max = n
		}
		// Right edge, facing left.
		if n := lib.CountEnergizedFrom(m, r, m.Cols()-1, 0, -1); n > max {
			max = n
		}
	}
	for c := 0; c < m.Cols(); c++ {
		// Top row, facing down.
		if n := lib.CountEnergizedFrom(m, 0, c, 1, 0); n > max {
			max = n
		}
		// Bottom row, facing up.
		if n := lib.CountEnergizedFrom(m, m.Rows()-1, c, -1, 0); n > max {
			max = n
		}
	}
	fmt.Println(max)
}
