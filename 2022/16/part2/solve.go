package main

import (
	"bufio"
	"expvar"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"aoc/2022/16/graph"
)

var (
	inputFile = flag.String("input", "input.txt", "Input file path")
	maxSteps  = flag.Int("steps", 26, "Time steps to simulate")
	startNode = flag.String("start", "AA", "Name of the start node")
	doVerbose = flag.Bool("v", false, "Enable logging")
	port      = flag.Int("port", 0, "Service port")

	specRE = regexp.MustCompile(
		`Valve (\w+) has flow rate=(\d+); tunnels? leads? to valves? ((\w+, )*\w+)\s*$`)
)

func main() {
	flag.Parse()
	g, err := parseInput(*inputFile)
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}

	if *port > 0 {
		go func() {
			if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
				log.Fatalf("Start service: %v", err)
			}
		}()
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println(best(g, []string{*startNode}, []string{*startNode}, 0, 0, 0))
	log.Printf("optimizer calls=%d mhits=%d", ncalls.Value(), mhit.Value())
}

func logPrintf(msg string, args ...any) {
	if *doVerbose {
		log.Printf(msg, args...)
	}
}

type mem struct {
	node1, node2 string
	step1, step2 int
	val          int
}

var (
	ncalls = expvar.NewInt("optimizer_calls")
	mhit   = expvar.NewInt("memo_hits")
	memo   = make(map[mem]int)
)

func best(g *graph.G[string], path1, path2 []string, step1, step2, value int) int {
	ncalls.Add(1)
	cur1, cur2 := path1[len(path1)-1], path2[len(path2)-1]
	if v, ok := memo[mem{cur1, cur2, step1, step2, value}]; ok {
		mhit.Add(1)
		return v
	}

	bestc1, bestc2, bestv := "", "", value
	logPrintf("[%d/%d] start %v %v init %v", step1, step2, path1, path2, value)

	var candidates []string
	for _, next := range g.Nodes() {
		// Skip cycles. We're choosing a candidate to activate, so anything that
		// makes it into the path is already active and should not be chosen.
		if next == cur1 || next == cur2 || nameInPath(next, path1) || nameInPath(next, path2) {
			continue
		}
		candidates = append(candidates, next)
	}
	sort.Strings(candidates)

	// Choose a pair of candidates:
	//  - Each worker must go to a different location.
	//  - Both locations must be reachable.
	//  - Each worker has to have time to get there and turn the knob.
	//
	for _, c1 := range candidates {
		for _, c2 := range candidates {
			if c1 == c2 {
				continue // different location
			}

			d1, _ := g.Distance(cur1, c1)
			d2, _ := g.Distance(cur2, c2)

			// Don't choose a pair that runs either worker out of time.
			astep1, astep2 := step1+d1+1, step2+d2+1
			if astep1 > *maxSteps || astep2 > *maxSteps {
				continue
			}

			// We get the benefit of both improvements.
			incr := valueAtStep(g, c1, astep1) + valueAtStep(g, c2, astep2)

			nextv := best(g, append(path1, c1), append(path2, c2), astep1, astep2, value+incr)
			if nextv > bestv {
				bestc1, bestc2, bestv = c1, c2, nextv
			}
		}
	}

	memo[mem{cur1, cur2, step1, step2, value}] = bestv

	// It's possible no path we try from this point can improve our score.  That
	// basically means we could stop here for the rest of the time period
	// without affecting the outcome.
	if bestc1 != "" {
		logPrintf("[%d/%d] %v %v best (%v %v) total %v", step1, step2, path1, path2, bestc1, bestc2, bestv)
	}
	return bestv
}

func nameInPath(name string, path []string) bool {
	for _, p := range path {
		if p == name {
			return true
		}
	}
	return false
}

func valueAtStep(g *graph.G[string], name string, step int) int {
	w, _ := g.Node(name)
	return w * (*maxSteps - step)
}

func parseInput(path string) (*graph.G[string], error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	g := graph.New[string]()

	s := bufio.NewScanner(f)
	var nlines int
	for s.Scan() {
		nlines++
		m := specRE.FindStringSubmatch(s.Text())
		if m == nil {
			return nil, fmt.Errorf("line %d: invalid input", nlines)
		}
		rate, _ := strconv.Atoi(m[2])
		g.SetNode(m[1], rate)
		for _, out := range strings.Split(m[3], ", ") {
			g.Edge(m[1], out)
		}
	}
	return g, s.Err()
}
