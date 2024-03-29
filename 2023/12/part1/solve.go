package main

import (
	"flag"
	"fmt"
	"log"

	lib "github.com/creachadair/aoc/2023/12"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	recs, err := lib.ParseRecords(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Parse records: %v", err)
	}

	var sum int
	for _, r := range recs {
		nr := lib.Solve(r)
		sum += nr
		fmt.Println(nr)
	}
	fmt.Println(sum)
}
