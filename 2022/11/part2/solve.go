package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/creachadair/aoc/2022/11/monkey"
)

var (
	inputFile = flag.String("input", "input.txt", "Input file path")
	numRounds = flag.Int("rounds", 10000, "Number of rounds")
	worryDiv  = flag.Int("worry", 1, "Worry reduction factor")
)

func main() {
	flag.Parse()
	input, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}

	monkeys, err := monkey.Parse(string(input))
	if err != nil {
		log.Fatalf("Parse monkeys: %v", err)
	}
	for _, m := range monkeys {
		m.WDiv = *worryDiv
	}

	for r := 0; r < *numRounds; r++ {
		for _, m := range monkeys {
			for m.More() {
				item, target := m.Next()
				monkeys[target].Catch(item)
			}
		}
	}

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].Inspected > monkeys[j].Inspected
	})
	fmt.Println(monkeys[0])
	fmt.Println(monkeys[1])
	fmt.Println(monkeys[0].Inspected * monkeys[1].Inspected)
}
