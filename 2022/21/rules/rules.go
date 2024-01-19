package rules

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/creachadair/mds/stack"
)

// A Rule is the parsed format of a puzzle input. As read from the input,
// either Op == "", in which case value is the constant reported by that node,
// or LHS and RHS give the names of the inputs to the operator.
//
// Valid input operators are "+", "-", "*", and "/".
type Rule struct {
	Name     string
	Op       string
	LHS, RHS string
	Value    int
}

func (r Rule) String() string {
	if r.Op == "" {
		return fmt.Sprintf("%s: %d", r.Name, r.Value)
	}
	return fmt.Sprintf("%s: %s %s %s", r.Name, r.LHS, r.Op, r.RHS)
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
	Values map[string]float64  // :: name → value (if known)
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

	q := stack.New[string]()
	q.Add("root")

	for !q.IsEmpty() {
		next := q.Top()
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

// PreSolve resolves the values of all the nodes that can be statically solved
// without knowing the values of any unknowns. If there are no "?" nodes in the
// graph, this should solve every node.
func (g *Graph) PreSolve() {
	g.Values = make(map[string]float64)

	for _, id := range g.topoSort() {
		r, ok := g.Rules[id]
		if !ok {
			log.Fatalf("No rule matching %q", id)
		}
		if r.Op == "" {
			g.Values[id] = float64(r.Value)
			continue
		} else if r.Op == "?" || r.Op == "=" {
			continue // placeholder for unresolved questions
		}

		// If we don't have both operands for a combining expression, just keep
		// going for now; we'll come back to it during the resolution step.  But
		// we should have at least one, otherwise the input is ill-formed: The
		// puzzle has only one unknown and the input forms a tree.
		_, lok := g.Values[r.LHS]
		_, rok := g.Values[r.RHS]
		if !lok && !rok {
			log.Fatalf("Both operands missing for %q (%s)", id, r.Op)
		} else if !lok || !rok {
			continue
		}

		// N.B. We've already filtered out "=" above.
		g.evaluate(r)
	}
}

// evaluate computes and caches the value for an operator rule r.  It panics if
// either operand of r is not known.
func (g *Graph) evaluate(r *Rule) float64 {
	lhs, ok := g.Values[r.LHS]
	if !ok {
		panic(fmt.Sprintf("missing lhs %q for %q", r.LHS, r.Name))
	}
	rhs, ok := g.Values[r.RHS]
	if !ok {
		panic(fmt.Sprintf("missing rhs %q for %q", r.RHS, r.Name))
	}

	switch r.Op {
	case "+":
		g.Values[r.Name] = lhs + rhs
	case "-":
		g.Values[r.Name] = lhs - rhs
	case "*":
		g.Values[r.Name] = lhs * rhs
	case "/":
		g.Values[r.Name] = lhs / rhs
	case "=":
		if lhs == rhs {
			g.Values[r.Name] = lhs
		} else {
			g.Values[r.Name] = math.NaN()
		}
	default:
		panic(fmt.Sprintf("invalid operator %q for %q", r.Op, r.Name))
	}
	return g.Values[r.Name]
}

// Solve computes the value for the named node that will cause the dependents
// of that node to compute the specified target value.
func (g *Graph) Solve(name string, target float64) (float64, error) {
	// Case 1: We already have the value for this node (e.g., a constant).
	// Report success as long as that value is the target.
	if v, ok := g.Values[name]; ok {
		if v == target {
			return v, nil
		}
		return v, fmt.Errorf("value of %q is %v, not %v", name, v, target)
	}

	rule, ok := g.Rules[name]
	if !ok {
		return 0, fmt.Errorf("no rule for %q", rule)
	}

	// Case 2: The rule is an unresolved variable. We can resolve it directly to
	// the target value. Note that if we already resolved it, we would catch
	// that in Case 1 (and the unification check is there).
	if rule.Op == "?" {
		g.Values[name] = target
		return target, nil
	}

	// Case 3: Otherwise this is an operator node with exactly one unknown
	// argument.  If both arguments were known, we would have computed this
	// value during pre-solution. If neither argument is known, pre-solution
	// would have failed (an optimization for this puzzle, where there is at
	// most one unknown and no shared paths).
	var comb func(bool, float64) float64
	switch rule.Op {
	case "+":
		// target = ? + v or v + ?, so ? = target - v.
		comb = func(_ bool, v float64) float64 { return target - v }

	case "*":
		// target = ? * v or v * ?, so ? = target / v.
		comb = func(_ bool, v float64) float64 { return target / v }

	case "=":
		// Set the target for the unresolved branch to the resolved branch.
		comb = func(_ bool, v float64) float64 { return v }

	case "-":
		// either target = v - ?, and ? = v - target
		// or     target = ? - v, and ? = target + v
		comb = func(left bool, v float64) float64 {
			if left {
				return v - target
			}
			return target + v
		}

	case "/":
		// either target = v / ?, and ? = v / target
		// or     target = ? / v, and ? = target * v
		comb = func(left bool, v float64) float64 {
			if left {
				return v / target
			}
			return target * v
		}

	default:
		panic(fmt.Sprintf("Invalid operator %q for %q", rule.Op, name))
	}

	// Choose the next rule that needs a solution (one of the children of this
	// node), update the target, and recur. If this completes, we can then solve
	// this node.
	next, nextTarget := g.update(rule.LHS, rule.RHS, comb)
	if v, err := g.Solve(next, nextTarget); err != nil {
		return 0, err
	} else {
		g.Values[next] = v // resolve the missing branch
	}
	return g.evaluate(rule), nil
}

// update finds which of lhs and rhs already has a value, and calls comb with
// that value. If neither has a value, update panics. The comb function should
// produce a new target value given the operand.
//
// The flag tells comb whether the known operand value is the left (true) or
// right (false).
func (g *Graph) update(lhs, rhs string, comb func(left bool, v float64) float64) (string, float64) {
	if v, ok := g.Values[lhs]; ok {
		return rhs, comb(true, v) // N.B. return the name of the OTHER branch!
	}
	if v, ok := g.Values[rhs]; ok {
		return lhs, comb(false, v)
	}
	panic("unreached")
}

func (g *Graph) Dot(w io.Writer) {
	fmt.Fprintf(w, "digraph G {\n")
	fmt.Fprintln(w, `  node [shape=record];`)
	defer fmt.Fprintln(w, "}")

	for name, rule := range g.Rules {
		if rule.Op != "" {
			fmt.Fprintf(w, "  %s [label=\"{%s|%s %s %s}\"]\n", name, name, rule.LHS, rule.Op, rule.RHS)
		} else {
			fmt.Fprintf(w, "  %s [label=\"%s|%d\"]\n", name, name, rule.Value)
		}
	}
	for name, deps := range g.Deps {
		fmt.Fprintf(w, "  %s -> {%s}\n", name, strings.Join(deps, " "))
	}
}
