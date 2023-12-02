package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	ndig := regexp.MustCompile(`\D+`)
	input := string(aoc.MustReadInput())
	var sum int
	for _, line := range strings.Split(input, "\n") {
		clean := ndig.ReplaceAllString(line, "")
		if clean == "" {
			continue
		}
		num := int(clean[0]-'0')*10 + int(clean[len(clean)-1]-'0')
		sum += num
	}
	fmt.Println(sum)
}
