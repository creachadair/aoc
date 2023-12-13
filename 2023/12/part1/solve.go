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
	fmt.Println(recs)

	var sum int
	for i, r := range recs {
		var nh int
		fmt.Printf("Record %d:\n * %s %v\n", i+1, r.Pattern, r.Groups)
		r.Solve(func(soln string) {
			fmt.Printf(" - %s\n", soln)
			nh++
		})
		fmt.Printf("TOTAL: %d\n", nh)
		sum += nh
	}
	fmt.Println(sum)
}
