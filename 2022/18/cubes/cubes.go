package cubes

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Pos [3]int

type Map map[Pos]int

func (m Map) Add(cube Pos) {
	m[cube] = 6
	for _, n := range neighborsOf(cube) {
		if _, ok := m[n]; ok {
			m[cube]--
			m[n]--
		}
	}
}

func (m Map) Sum() int {
	var sum int
	for _, v := range m {
		sum += v
	}
	return sum
}

func (m Map) Bounds() (min, max Pos) {
	var maxx, maxy, maxz int
	for pos := range m {
		if pos[0] > maxx {
			maxx = pos[0]
		}
		if pos[1] > maxy {
			maxy = pos[1]
		}
		if pos[2] > maxz {
			maxz = pos[2]
		}
	}
	minx, miny, minz := maxx, maxy, maxz
	for pos := range m {
		if pos[0] < minx {
			minx = pos[0]
		}
		if pos[1] < miny {
			miny = pos[1]
		}
		if pos[2] < minz {
			minz = pos[2]
		}
	}
	return Pos{minx, miny, minz}, Pos{maxx, maxy, maxz}
}

func (m Map) Flood() map[Pos]bool {
	min, max := m.Bounds()
	seed := Pos{min[0] - 1, min[1] - 1, min[2] - 1}

	q := []Pos{seed}
	marked := map[Pos]bool{seed: true}
	for len(q) != 0 {
		next := q[len(q)-1]
		q = q[:len(q)-1]

		for _, adj := range neighborsOf(next) {
			// Filter out out-of-bound points. But note we allow one slice beyond
			// the bounding cubes so we can't get trapped in a cul-de-sac.
			if adj[0] < min[0]-1 || adj[1] < min[1]-1 || adj[2] < min[2]-2 {
				continue // past the inner bounding point
			} else if adj[0] > max[0]+1 || adj[1] > max[1]+1 || adj[2] > max[2]+2 {
				continue // past the outer bounding point
			} else if _, ok := m[adj]; ok || marked[adj] {
				continue // skip marked cubes and cubes already in the map
			}

			marked[adj] = true
			q = append(q, adj)
		}
	}
	return marked
}

func neighborsOf(cube Pos) []Pos {
	return []Pos{
		{cube[0] - 1, cube[1], cube[2]},
		{cube[0] + 1, cube[1], cube[2]},
		{cube[0], cube[1] - 1, cube[2]},
		{cube[0], cube[1] + 1, cube[2]},
		{cube[0], cube[1], cube[2] - 1},
		{cube[0], cube[1], cube[2] + 1},
	}
}

func MustParse(path string) Map {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}
	text := strings.TrimSpace(string(data))

	m := make(Map)
	for i, line := range strings.Split(text, "\n") {
		pos, err := parseCube(line)
		if err != nil {
			log.Fatalf("Line %d: %v", i+1, err)
		}
		m.Add(pos)
	}
	return m
}

func parseCube(text string) (Pos, error) {
	fs := strings.Split(text, ",")
	if len(fs) != 3 {
		return Pos{}, fmt.Errorf("wrong number of fields %d", len(fs))
	}
	x, _ := strconv.Atoi(fs[0])
	y, _ := strconv.Atoi(fs[1])
	z, _ := strconv.Atoi(fs[2])
	return Pos{x, y, z}, nil
}
