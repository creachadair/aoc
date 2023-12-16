package lib

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/creachadair/aoc/aoc"
	"github.com/creachadair/mds/mapset"
)

type Card struct {
	ID  int
	Win []int
	Num []int
}

func (c Card) String() string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "Card %d:", c.ID)
	for _, win := range c.Win {
		fmt.Fprintf(&buf, " %2d", win)
	}
	buf.WriteString(" |")
	for _, num := range c.Num {
		fmt.Fprintf(&buf, " %2d", num)
	}
	return buf.String()
}

func (c Card) Matches() int {
	wins := mapset.New(c.Win...)
	var matches int
	for _, num := range c.Num {
		if wins.Has(num) {
			matches++
		}
	}
	return matches
}

func (c Card) Score() int {
	if m := c.Matches(); m > 0 {
		return 1 << (m - 1)
	}
	return 0
}

type Cards []Card

func (c Cards) String() string {
	cs := make([]string, len(c))
	for i, cd := range c {
		cs[i] = cd.String()
	}
	return strings.Join(cs, "\n")
}

var cardRE = regexp.MustCompile(`^Card +(\d+):((?: +\d+)+) \|((?: +\d+)+)$`)

func ParseCards(input []byte) (Cards, error) {
	var cards Cards
	for i, line := range aoc.SplitLines(input) {
		var c Card
		if err := aoc.Scanx(cardRE, line, &c.ID, &c.Win, &c.Num); err != nil {
			return nil, fmt.Errorf("line %d: invalid card: %w", i+1, err)
		}
		cards = append(cards, c)
	}
	return cards, nil
}
