package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	lib "github.com/creachadair/aoc/2023/07"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	bids, err := lib.ParseBids(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Parse bids: %v", err)
	}
	sort.Slice(bids, func(i, j int) bool {
		return lib.CompareHandsWild(bids[i].Hand, bids[j].Hand) < 0
	})
	var total int
	for i, bid := range bids {
		val := (i + 1) * bid.Value
		best := bid.Hand.Best()
		log.Printf("Hand %v best %v type %v bid %d val %d", bid.Hand, best, best.Type(), bid.Value, val)
		total += val
	}
	fmt.Println(total)
}
