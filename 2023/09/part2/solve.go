package main

import (
	"flag"
	"fmt"
	"log"

	lib "github.com/creachadair/aoc/2023/09"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	seqs, err := lib.ParseSeq(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Parse sequences: %v", err)
	}
	var sum int
	for _, s := range seqs {
		next := s.Reject()
		fmt.Println(next, s)
		sum += next
	}
	fmt.Println(sum)
}
