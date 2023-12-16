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

	num := regexp.MustCompile(`\d+`)
	var sum int
	for r, row := range sch.RowStrings() {
		for _, hit := range num.FindAllStringIndex(row, -1) {
			val, err := strconv.Atoi(row[hit[0]:hit[1]])
			if err != nil {
				log.Fatalf("Row %d: invalid number at %d (unexpected): %v", r, hit[0], err)
			}
			ok := sch.IsLabel(r, hit[0], hit[1])
			log.Printf("Row %d (%d:%d) val=%d hit=%v", r, hit[0], hit[1], val, ok)
			if ok {
				sum += val
			}
		}
	}
	fmt.Println(sum)
}
