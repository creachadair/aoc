package lib

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
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

func (c Card) Score() int {
	wins := mapset.New(c.Win...)
	cur := 0
	for _, num := range c.Num {
		if wins.Has(num) {
			if cur == 0 {
				cur = 1
			} else {
				cur *= 2
			}
		}
	}
	return cur
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

func parseCard(data string) (Card, error) {
	m := cardRE.FindStringSubmatch(data)
	if m == nil {
		return Card{}, errors.New("invalid card format")
	}
	c := Card{ID: mustParseInt(m[1])}
	for _, win := range strings.Fields(m[2]) {
		c.Win = append(c.Win, mustParseInt(win))
	}
	for _, num := range strings.Fields(m[3]) {
		c.Num = append(c.Num, mustParseInt(num))
	}
	return c, nil
}

func mustParseInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return v
}

func ParseCards(input []byte) (Cards, error) {
	var cards Cards
	for i, line := range aoc.SplitLines(input) {
		c, err := parseCard(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid card: %w", i+1, err)
		}
		cards = append(cards, c)
	}
	return cards, nil
}
