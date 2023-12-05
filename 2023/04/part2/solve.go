package main

import (
	"flag"
	"fmt"
	"log"

	lib "github.com/creachadair/aoc/2023/04"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	cards, err := lib.ParseCards(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Parse cards: %v", err)
	}
	have := make(map[int]int)
	for _, c := range cards {
		have[c.ID] = 1
	}

	var total int
	for _, c := range cards {
		s := c.Matches()
		if s == 0 {
			log.Printf("card %d does nothing", c.ID)
			total += have[c.ID]
			continue
		}
		end := c.ID + s
		if end > len(cards) {
			end = len(cards)
		}
		for end > c.ID {
			have[end] += have[c.ID]
			log.Printf("card %d adds %d of card %d", c.ID, have[c.ID], end)
			end--
		}

		// Once we are done with a card, we cannot add any more of it.
		total += have[c.ID]
	}
	fmt.Println(total)
}
