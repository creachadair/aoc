package lib

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/creachadair/aoc/aoc"
)

const cardName = "23456789TJQKA"

type Hand [5]byte // in nonincreasing order

func (h Hand) Type() HandType {
	m := make(map[byte]int)
	var max int
	for _, c := range h {
		m[c]++
		if m[c] > max {
			max = m[c]
		}
	}
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

func argMax[T comparable](m map[T]int) (T, int) {
	var max T
	var mc int
	for v, n := range m {
		if n > mc {
			max, mc = v, n
		}
	}
	return max, mc
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

func CompareHands(a, b Hand) int {
	ta, tb := a.Type(), b.Type()
	if ta < tb {
		return -1
	} else if ta > tb {
		return 1
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
