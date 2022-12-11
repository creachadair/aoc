package monkey

import (
	"fmt"
	"strconv"
	"strings"
)

type Monkey struct {
	ID          int    // original index (from input)
	Items       []int  // current items held (order matters)
	Op          string // update operator
	Arg         int    // 0 == old; otherwise argument value
	Div         int    // divisibility factor
	True, False int    // targets per divisibility
	Inspected   int    // number of items inspected by this monkey
	Mod         int    // compute worry scores modulo this value
	WDiv        int    // worry divisor
}

// More reports whether monkey m has any further items to consider.
func (m *Monkey) More() bool { return len(m.Items) > 0 }

// Catch receives an item thrown by another monkey.
func (m *Monkey) Catch(item int) { m.Items = append(m.Items, item) }

// Next inspects and throws the next item, giving the weight and target.
// Precondition: m.More() == true.
func (m *Monkey) Next() (item, target int) {
	m.Inspected++
	item = m.Items[0]
	copy(m.Items, m.Items[1:])
	m.Items = m.Items[:len(m.Items)-1]

	arg := m.Arg
	if arg == 0 {
		arg = item
	}
	switch m.Op {
	case "*":
		item = (item * arg) % m.Mod
	case "+":
		item = (item + arg) % m.Mod
	default:
		panic("unknown operator: " + m.Op)
	}

	item /= m.WDiv // reduced worry after inspection
	if item%m.Div == 0 {
		target = m.True
	} else {
		target = m.False
	}
	return
}

func (m *Monkey) String() string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "Monkey %d:\n", m.ID)
	fmt.Fprintf(&buf, "  Starting items: %s\n", joinInts(m.Items, ", "))
	fmt.Fprintf(&buf, "  Operation: new = old %s ", m.Op)
	if m.Arg == 0 {
		fmt.Fprintln(&buf, "old")
	} else {
		fmt.Fprintf(&buf, "%d\n", m.Arg)
	}
	fmt.Fprintf(&buf, "  Test: divisible by %d\n", m.Div)
	fmt.Fprintf(&buf, "    If true: throw to monkey %d\n", m.True)
	fmt.Fprintf(&buf, "    If false: throw to monkey %d\n", m.False)
	if m.Inspected > 0 {
		fmt.Fprintf(&buf, "  Inspected: %d\n", m.Inspected)
	}
	return buf.String()
}

func Parse(input string) (_ []*Monkey, err error) {
	var lnum int
	defer func() {
		if x := recover(); x != nil && err == nil {
			err = fmt.Errorf("line %d: %w", lnum, x.(error))
		}
	}()

	var out []*Monkey

	// Track the product of all the divisors. The divisors in the input are all
	// coprime (in fact, all are prime), which means we can compute in a residue
	// space modulo their product. This works because all our operations are
	// compatible (addition, multiplication), and the fact that the dividends
	// are coprime ensures residues are unique (by CRT).
	//
	// More importantly, working in the residue space keeps values from becoming
	// really big and blowing the precision of machine integers. The
	// computations are too big for arbitrary precision to be feasible.
	mod := 1

	for _, hunk := range strings.Split(strings.TrimSpace(input), "\n\n") {
		cur := &Monkey{WDiv: 3}
		for _, line := range strings.Split(hunk, "\n") {
			lnum++
			tag, rest, ok := strings.Cut(strings.TrimSpace(line), ":")
			if !ok {
				return nil, fmt.Errorf("line %d: invalid rule %q", lnum, line)
			}
			fs := strings.Fields(rest)

			switch strings.TrimSpace(tag) {
			case "Starting items":
				cur.Items = mustSplitInts(strings.TrimSpace(rest))

			case "Operation":
				cur.Op = fs[3]
				if fs[4] != "old" {
					cur.Arg = mustParseInt(fs[4])
				}

			case "Test":
				cur.Div = mustParseInt(fs[2])

			case "If true":
				cur.True = mustParseInt(fs[3])

			case "If false":
				cur.False = mustParseInt(fs[3])

			default:
				id := strings.TrimPrefix(tag, "Monkey ")
				if id == tag {
					return nil, fmt.Errorf("line %d: invalid input line %q", lnum, line)
				}
				cur.ID = mustParseInt(strings.TrimSuffix(id, ":"))
			}
		}
		out = append(out, cur)
		mod *= cur.Div
	}
	for _, m := range out {
		m.Mod = mod // fill in the modulus (now that we've seen all the factors)
	}
	return out, nil
}

func mustParseInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return v
}

func mustSplitInts(s string) []int {
	var out []int
	for _, v := range strings.Split(s, ", ") {
		out = append(out, mustParseInt(v))
	}
	return out
}

func joinInts(vs []int, s string) string {
	buf := make([]string, len(vs))
	for i, v := range vs {
		buf[i] = strconv.Itoa(v)
	}
	return strings.Join(buf, s)
}
