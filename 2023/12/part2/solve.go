package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

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
		r2 := lib.Record{Pattern: explode(r.Pattern)}
		for i := 0; i < 5; i++ {
			r2.Groups = append(r2.Groups, r.Groups...)
		}
		nr := lib.Solve(r2)
		sum += nr
		fmt.Println(nr)
	}
	fmt.Println(sum)
}

func explode(pat string) string {
	ss := make([]string, 5)
	for i := range ss {
		ss[i] = pat
	}
	return strings.Join(ss, "?")
}
