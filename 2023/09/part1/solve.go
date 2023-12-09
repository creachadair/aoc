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
		next := s.Project()
		fmt.Println(s, next)
		sum += next
	}
	fmt.Println(sum)
}
