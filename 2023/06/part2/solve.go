package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	lib "github.com/creachadair/aoc/2023/06"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	races, err := lib.ParseRaces(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Parse races: %v", err)
	}
	var tbuf, dbuf strings.Builder
	for _, r := range races {
		fmt.Fprintf(&tbuf, "%d", r.Time)
		fmt.Fprintf(&dbuf, "%d", r.Dist)
	}
	ttime, err := strconv.Atoi(tbuf.String())
	if err != nil {
		log.Fatalf("Total time: %v", err)
	}
	tdist, err := strconv.Atoi(dbuf.String())
	if err != nil {
		log.Fatalf("Total dist: %v", err)
	}

	min := lib.MinGreater(ttime, tdist)
	max := ttime - min
	nd := (max + 1) - min
	log.Printf("Big race: %d ≤ i ≤ %d (%d)", min, max, nd)
	fmt.Println(nd)
}
