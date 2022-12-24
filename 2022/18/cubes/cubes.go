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
