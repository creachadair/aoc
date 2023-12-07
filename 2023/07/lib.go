package lib

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/creachadair/aoc/aoc"
	"golang.org/x/exp/constraints"
)

const cardName = "_23456789TJQKA"

var (
	joker    = byte(strings.IndexByte(cardName, 'J'))
	ace      = byte(strings.IndexByte(cardName, 'A'))
	bestHand = Hand{ace, ace, ace, ace, ace}
)

type Hand [5]byte

// Best returns the strongest hand equivalent to h with all J cards optimally
// replaced as jokers.
func (h Hand) Best() Hand {
	m, _ := mapHand(h)
	nj := m[joker]
	if nj == 0 {
		return h // no jokers to replace
	} else if nj == 5 {
		return bestHand // all jokers; make the best 5-hand
	}

	// Remove the jokers and consider the best most-frequent remaining card in
	// the hand: Replacing the jokers with that card optimizes the result.
	//
	// For example:
	//    XbcdJ    high card X becomes XXbcd one pair
	//    XXbcJ    one pair XX becomes XXXbc three of a kind
	//    XXbJJ    one pair XX becomes XXXXb four of a kind
	//    XXYYJ    two pair XX YY becomes XXXYY full house
	//    XXXyJ    three of a kind XXX becomes XXXXy four of a kind
	//
	// etc.

	delete(m, joker)
	mc := argMax(m)
	cp := h
	for i, c := range cp {
		if c == joker {
			cp[i] = mc
		}
	}
	return cp
}

func mapHand(h Hand) (map[byte]int, int) {
	m := make(map[byte]int)
	var max int
	for _, c := range h {
		m[c]++
		if m[c] > max {
			max = m[c]
		}
	}
	return m, max
}

func (h Hand) Type() HandType {
	m, max := mapHand(h)
	switch len(m) {
	case 1:
		return 6 // five of a kind
	case 2:
		// either four of a kind XXXXa or full house XXXaa
		if max == 4 {
			return 5 // four of a kind
		}
		return 4 // full house
	case 3:
		// either three of a kind XXXab or two pair XXYYc
		if max == 3 {
			return 3 // three of a kind
		}
		return 2 // two pair
	case 4:
		return 1 // one pair (XXabc)
	default:
		return 0 // high card
	}
}

func (h Hand) String() string {
	var buf [5]byte
	for i, c := range h {
		buf[i] = cardName[int(c)]
	}
	return string(buf[:])
}

type HandType int

func (t HandType) String() string {
	if t >= 0 && int(t) < len(handName) {
		return handName[t]
	}
	return "invalid"
}

var handName = []string{"high-card", "one-pair", "two-pair", "3-of-kind", "full-house", "4-of-kind", "5-of-kind"}

func argMax[T constraints.Ordered](m map[T]int) T {
	var max T
	var mc int
	for v, n := range m {
		if n > mc || (n == mc && v > max) {
			max, mc = v, n
		}
	}
	return max
}

func parseHand(s string) (Hand, error) {
	if len(s) != 5 {
		return Hand{}, fmt.Errorf("got %d cards, want 5", len(s))
	}
	h := Hand([]byte(s[:5]))
	for i := 0; i < len(h); i++ {
		p := strings.IndexByte(cardName, h[i])
		if p < 0 {
			return Hand{}, fmt.Errorf("offset %d: invalid card %c", i, h[i])
		}
		h[i] = byte(p)
		/* Weirdly, the puzzle assumes input order is significant.
		for j := i; j > 0; j-- { // insertion sort
			if h[j] <= h[j-1] {
				break
			}
			h[j-1], h[j] = h[j], h[j-1]
		}
		*/
	}
	return h, nil
}

// CompareHands compares a and b treating J cards as Jacks.
func CompareHands(a, b Hand) int {
	ta, tb := a.Type(), b.Type()
	if ta < tb {
		return -1
	} else if ta > tb {
		return 1
	}
	return bytes.Compare(a[:], b[:])
}

// CompareHands compares a and be treating J cards as Jokers.
func CompareHandsWild(a, b Hand) int {
	ta, tb := a.Best().Type(), b.Best().Type()
	if ta < tb {
		return -1
	} else if ta > tb {
		return 1
	}
	// To break ties, we treat jokers as zeroes.
	for i := 0; i < len(a); i++ {
		if a[i] == joker {
			a[i] = 0
		}
		if b[i] == joker {
			b[i] = 0
		}
	}
	return bytes.Compare(a[:], b[:])
}

type Bid struct {
	Hand  Hand
	Value int
}

func ParseBids(input []byte) ([]Bid, error) {
	var out []Bid
	for i, line := range aoc.SplitLines(input) {
		fs := strings.Fields(line)
		if len(fs) != 2 {
			return nil, fmt.Errorf("line %d: got %d fields, want 2", i+1, len(fs))
		}
		h, err := parseHand(fs[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid hand: %w", i+1, err)
		}
		bid, err := strconv.Atoi(fs[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid bid: %w", i+1, err)
		}
		out = append(out, Bid{Hand: h, Value: bid})
	}
	return out, nil
}
