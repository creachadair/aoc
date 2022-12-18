package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"aoc/2022/16/graph"
)

var (
	inputFile = flag.String("input", "input.txt", "Input file path")
	maxSteps  = flag.Int("steps", 30, "Time steps to simulate")
	startNode = flag.String("start", "AA", "Name of the start node")
	doVerbose = flag.Bool("v", false, "Enable logging")

	specRE = regexp.MustCompile(
		`Valve (\w+) has flow rate=(\d+); tunnels? leads? to valves? ((\w+, )*\w+)\s*$`)
)

func main() {
	flag.Parse()
	g, err := parseInput(*inputFile)
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}

	path, v := best(g, []string{*startNode}, 0, 0)
	fmt.Println(expandPath(g, path))
	fmt.Println(v)
	log.Printf("optimizer calls=%d mhits=%d", ncalls, mhit)
}

func logPrintf(msg string, args ...any) {
	if *doVerbose {
		log.Printf(msg, args...)
	}
}

type mem struct {
	node      string
	step, val int
}

type sval struct {
	suffix []string
	value  int
}

var (
	ncalls, mhit int
	memo         = make(map[mem]sval)
)

func best(g *graph.G[string], path []string, step, value int) ([]string, int) {
	ncalls++
	cur := path[len(path)-1]
	if v, ok := memo[mem{cur, step, value}]; ok {
		mhit++
		return v.suffix, v.value
	}

	bestc, bestv := "", value
	var suffix []string
	logPrintf("[%d] start %v init %v", step, path, value)

	for _, next := range g.Nodes() {
		// Skip cycles. We're choosing a candidate to activate, so anything that
		// makes it into the path is already active and should not be chosen.
		if next == cur || nameInPath(next, path) {
			continue
		}

		d, ok := g.Distance(cur, next)
		if !ok {
			// There should always be a path, but don't get stuck if not.
			continue
		}

		// To get from cur to next costs d steps, plus 1 step to turn on the
		// value at that location. If that takes us past the step limit, we can't
		// include next as a candidate.
		astep := step + d + 1
		if astep > *maxSteps {
			continue
		}

		nexts, nextv := best(g, append(path, next), astep, value+valueAtStep(g, next, astep))
		if nextv > bestv {
			bestv, bestc = nextv, next
			suffix = append([]string{cur}, nexts...)
		}
	}

	memo[mem{cur, step, value}] = sval{suffix, bestv}

	// It's possible no path we try from this point can improve our score.  That
	// basically means we could stop here for the rest of the time period
	// without affecting the outcome.
	if bestc != "" {
		logPrintf("[%d] %v best %v total %v", step, path, bestc, bestv)
	}
	return suffix, bestv
}

func expandPath(g *graph.G[string], path []string) []string {
	if len(path) == 0 {
		return nil
	}
	exp := []string{path[0]}
	last := path[0]
	for _, next := range path[1:] {
		stem := g.Path(last, next)
		exp = append(exp, stem[1:]...)
		exp = append(exp, "*")
		last = next
	}
	return exp
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
