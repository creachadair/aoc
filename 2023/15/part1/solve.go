package main

import (
	"flag"
	"fmt"
	"strings"

	lib "github.com/creachadair/aoc/2023/15"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()
	input := strings.TrimSpace(string(aoc.MustReadInput()))
	rules := strings.Split(input, ",")

	var sum int
	for _, rule := range rules {
		sum += lib.Hash(rule)
	}
	fmt.Println(sum)
}
