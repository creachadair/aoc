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
	var sum int
	for _, c := range cards {
		s := c.Score()
		log.Printf("[%s]; score: %d\n", c, s)
		sum += s
	}
	fmt.Println(sum)
}
