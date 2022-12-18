package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/creachadair/taskgroup"
)

var (
	inputPath = flag.String("input", "input.txt", "Input file path")
	maxCoord  = flag.Int("max", 4000000, "Bounding square size")
	startVal  = flag.Int("start", 0, "Starting diagonal")

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
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].minBeacon() > ss[j].minBeacon()
	})

	const numTasks = 256
	type pair struct{ x, y int }
	g, start := taskgroup.New(nil).Limit(numTasks)
	ch := make(chan pair, numTasks)

	done := make(chan struct{})
	go func() {
		defer close(done)

	candidate:
		for c := range ch {
			fmt.Println(c.x, c.y, "?")
			for _, s := range ss {
				if s.isBeacon(c.x, c.y) {
					continue candidate
				}
			}
			fmt.Println(c.x, c.y, 4000000*c.x+c.y)
		}
	}()

	for k := *startVal; k < *maxCoord*2; k++ {
		if k%1000 == 0 {
			log.Printf("k=%d", k)
		}
		k := k
		start(func() error {
			i, j := 0, k
			for j >= 0 {
				if !anyClose(ss, i, j) {
					ch <- pair{i, j}
				}
				if pi, pj := *maxCoord-i, *maxCoord-j; !anyClose(ss, pi, pj) {
					ch <- pair{pi, pj}
				}
				i, j = i+1, j-1
			}
			return nil
		})
	}
	g.Wait()
	close(ch)
	<-done
}

func anyClose(ss []*sensor, i, j int) bool {
	for _, s := range ss {
		if s.dist(i, j) <= s.minBeacon() && !s.isBeacon(i, j) {
			return true
		}
	}
	return false
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
