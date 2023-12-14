package main

import (
	"flag"
	"log"

	"github.com/creachadair/aoc/aoc"
)

var numCycles = flag.Int("cycles", 1000000000, "Number of cycles")

func main() {
	flag.Parse()

	m, err := aoc.ParseMap(aoc.SplitLines(aoc.MustReadInput()))
	if err != nil {
		log.Fatalf("Parse map: %v", err)
	}

	var cyc []*aoc.Map
	seen := make(map[string]int)

	for i := 0; i < *numCycles; i++ {
		n := m.Clone()
		shiftUp(n)          // N
		shiftUp(n.Rotate()) // W
		shiftUp(n.Rotate()) // S
		shiftUp(n.Rotate()) // E
		n.Rotate()          // restore orientation

		s := n.String()
		if p, ok := seen[s]; ok {
			rem := (*numCycles - i - 1) % (len(cyc) - p)
			hit := cyc[p+rem]
			log.Printf("At i=%d found cycle rem=%d\n%s\nload=%d",
				i, rem, hit, load(hit))
			break
		} else {
			seen[s] = i
			cyc = append(cyc, n)
		}
		m = n
	}
}

func load(m *aoc.Map) int {
	var sum int
	for r := 0; r < m.Rows(); r++ {
		for c := 0; c < m.Cols(); c++ {
			if m.At(r, c) == 'O' {
				sum += m.Rows() - r
			}
		}
	}
	return sum
}

func shiftUp(m *aoc.Map) *aoc.Map {
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
			}
		}
	}
	return m
}
