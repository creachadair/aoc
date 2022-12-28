package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"aoc/2022/20/cipher"
)

var (
	inputPath = flag.String("input", "input.txt", "Input file path")
	positions = flag.String("pos", "1000,2000,3000", "Positions after initial value")
)

func main() {
	flag.Parse()

	input := cipher.MustParse(*inputPath)

	var posn []int
	for _, p := range strings.Split(*positions, ",") {
		v, err := strconv.Atoi(p)
		if err != nil {
			log.Fatalf("Invalid position %q: %v", p, err)
		}
		posn = append(posn, v)
	}

	input.Mix()
	values := input.Values()
	if values[0] != 0 {
		log.Fatal("Malformed result, does not start with zero")
	}
	for i, v := range values[1:] {
		if v == 0 {
			log.Fatalf("Found a spurious zero at position %d", i+2)
		}
	}

	var sum int
	for _, v := range posn {
		pv := values[v%len(values)]
		fmt.Printf("Offset %v: value=%v\n", v, pv)
		sum += pv
	}
	fmt.Println("Sum:", sum)
}
