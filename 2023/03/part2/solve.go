package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"

	lib "github.com/creachadair/aoc/2023/03"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	sch, err := lib.ParseSchematic(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Parse schematic: %v", err)
	}
	log.Printf("Schematic is %d rows, %d columns", sch.Rows(), sch.Cols())

	adj := make(map[[2]int][]int)
	num := regexp.MustCompile(`\d+`)
	for r, row := range sch {
		for _, hit := range num.FindAllStringIndex(row, -1) {
			val, err := strconv.Atoi(row[hit[0]:hit[1]])
			if err != nil {
				log.Fatalf("Row %d: invalid number at %d (unexpected): %v", r, hit[0], err)
			}
			a, b := sch.LabelOf(r, hit[0], hit[1])
			ok := a >= 0 && b >= 0
			log.Printf("Row %d (%d:%d) val=%d hit=%v", r, hit[0], hit[1], val, ok)
			if ok && sch.At(a, b) == '*' {
				key := [2]int{a, b}
				adj[key] = append(adj[key], val)
			}
		}
	}
	var sum int
	for key, vals := range adj {
		if len(vals) == 2 {
			gr := vals[0] * vals[1]
			log.Printf("Found gear at %v: %v (ratio=%d)", key, vals, gr)
			sum += gr
		}
	}
	fmt.Println(sum)
}
