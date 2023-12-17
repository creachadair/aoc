package aoc

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// MDist reports the Manhattan distance between (r1, c1) and (r2, c2).
func MDist(r1, c1, r2, c2 int) int { return abs(r2-r1) + abs(c2-c1) }
