package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	inputPath = flag.String("input", "input.txt", "Input file path")
	scanRow   = flag.Int("row", 2000000, "Scan row")

	sensorLine = regexp.MustCompile(
		strings.ReplaceAll(sensorText, "%d", `(-?\d+)`) + "$",
	)
)

const sensorText = `Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d`

func main() {
	flag.Parse()
	input, err := os.ReadFile(*inputPath)
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}

	var ss []*sensor
	for _, line := range strings.Split(string(input), "\n") {
		if line == "" {
			continue
		}
		ss = append(ss, newSensor(line))
	}

	var invalid int
	xmin, xmax := xminmax(ss)
	fmt.Println(xmin, xmax)

	for i := xmin; i <= xmax; i++ {
		for _, s := range ss {
			if s.dist(i, *scanRow) <= s.minBeacon() && !s.isBeacon(i, *scanRow) {
				invalid++
				break
			}
		}
	}
	fmt.Println(invalid)
}

type sensor struct {
	x, y   int
	bx, by int
}

func (s *sensor) String() string {
	return fmt.Sprintf(sensorText+" {%d}", s.x, s.y, s.bx, s.by, s.minBeacon())
}

func (s *sensor) isBeacon(x, y int) bool { return x == s.bx && y == s.by }

func (s *sensor) dist(x, y int) int { return mdist(x, y, s.x, s.y) }

func (s *sensor) minBeacon() int { return mdist(s.x, s.y, s.bx, s.by) }

func xminmax(ss []*sensor) (xmin, xmax int) {
	xmin, xmax = ss[0].x, ss[0].x
	for _, s := range ss {
		xmin = min(xmin, s.x, s.bx, s.x-s.minBeacon())
		xmax = max(xmax, s.x, s.bx, s.x+s.minBeacon())
	}
	return
}

func max(a int, bs ...int) int {
	for _, b := range bs {
		if b > a {
			a = b
		}
	}
	return a
}

func min(a int, bs ...int) int {
	for _, b := range bs {
		if b < a {
			a = b
		}
	}
	return a
}

func abs(z int) int {
	if z < 0 {
		return -z
	}
	return z
}

func mdist(x1, y1, x2, y2 int) int { return abs(x1-x2) + abs(y1-y2) }

func newSensor(text string) *sensor {
	m := sensorLine.FindStringSubmatch(strings.TrimSpace(text))
	var s sensor
	s.x, _ = strconv.Atoi(m[1])
	s.y, _ = strconv.Atoi(m[2])
	s.bx, _ = strconv.Atoi(m[3])
	s.by, _ = strconv.Atoi(m[4])
	return &s
}
