package main

import (
	"flag"
	"fmt"

	"aoc/2022/18/cubes"
)

var inputFile = flag.String("input", "input.txt", "Input file path")

func main() {
	flag.Parse()

	m := cubes.MustParse(*inputFile)

	// Flood fill everything reachable from "outside" the bounding region of the
	// points that were discovered.
	marked := m.Flood()
	min, max := m.Bounds()

	// Sweep the entire bounded region, filling in any free but unmarked cube.
	// This has the effect of removing the interior regions.
	for x := min[0]; x <= max[0]; x++ {
		for y := min[1]; y <= max[1]; y++ {
			for z := min[2]; z <= max[2]; z++ {
				cur := cubes.Pos{x, y, z}
				if _, ok := m[cur]; !ok && !marked[cur] {
					m.Add(cur)
				}
			}
		}
	}

	fmt.Println(m.Sum())
}
