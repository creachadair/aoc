package rules

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/creachadair/mds/mlink"
)

type Rule struct {
	Name     string
	Op       string
	LHS, RHS string
	Value    int
}

var ruleRE = regexp.MustCompile(`(\w+): (?:(\d+)|(\w+) ([-+*/]) (\w+))`)

func (r *Rule) decode(text string) error {
	m := ruleRE.FindStringSubmatch(text)
	if m == nil {
		return fmt.Errorf("invalid rule format %q", text)
	}
	r.Name = m[1]
	if m[2] != "" {
		v, err := strconv.Atoi(m[2])
		if err != nil {
			return fmt.Errorf("invalid value %q: %w", m[2], err)
		}
		r.Value = v
	} else {
		r.LHS = m[3]
		r.Op = m[4]
		r.RHS = m[5]
	}
	return nil
}

func MustParse(path string) []Rule {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	rules := make([]Rule, len(lines))
	for i, line := range lines {
		if err := rules[i].decode(line); err != nil {
			log.Fatalf("Line %d: %v", i+1, err)
		}
	}
	return rules
}

type Graph struct {
	Rules  map[string]*Rule    // :: name → rule
	Deps   map[string][]string // :: name → direct dependencies
	Values map[string]int      // :: name → value (if known)
}

func NewGraph(rules []Rule) *Graph {
	g := &Graph{Rules: make(map[string]*Rule), Deps: make(map[string][]string)}
	for _, rule := range rules {
		if prev, ok := g.Rules[rule.Name]; ok {
			log.Fatalf("Duplicate rules for %q: %v, %v", rule.Name, prev, rule)
		}
		cp := rule
		g.Rules[rule.Name] = &cp
		if rule.Op != "" {
			g.Deps[rule.Name] = []string{rule.LHS, rule.RHS}
		}
	}
	return g
}

func (g *Graph) topoSort() []string {
	var finished []string
	seen := make(map[string]bool)

	q := mlink.NewStack[string]()
	q.Add("root")

	for !q.IsEmpty() {
		next, _ := q.Top()
		if seen[next] {
			finished = append(finished, next)
			q.Pop()
			continue
		}
		seen[next] = true
		for _, succ := range g.Deps[next] {
			q.Add(succ)
		}
	}

	return finished
}

func (g *Graph) Solve() {
	g.Values = make(map[string]int)

	for _, id := range g.topoSort() {
		r, ok := g.Rules[id]
		if !ok {
			log.Fatalf("No rule matching %q", id)
		}
		if r.Op == "" {
			g.Values[id] = r.Value
			continue
		}

		lhs, ok := g.Values[r.LHS]
		if !ok {
			log.Fatalf("Missing LHS %q for %q", r.LHS, id)
		}
		rhs, ok := g.Values[r.RHS]
		if !ok {
			log.Fatalf("Missing RHS %q for %q", r.RHS, id)
		}
		switch r.Op {
		case "+":
			g.Values[id] = lhs + rhs
		case "-":
			g.Values[id] = lhs - rhs
		case "*":
			g.Values[id] = lhs * rhs
		case "/":
			g.Values[id] = lhs / rhs
		case "=":
			if lhs == rhs {
				g.Values[id] = lhs
			} else {
				g.Values[id] = -1 // sentinel for "unequal"
			}
		default:
			log.Fatalf("Invalid operator %q for %q", r.Op, id)
		}
	}
}
