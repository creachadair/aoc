package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var inputFile = flag.String("input", "", "Input file name")

func main() {
	flag.Parse()
	if *inputFile == "" {
		log.Fatal("You must provide an -input file to read")
	}
	data, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Reading input: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var totalScore int
	for i, rule := range lines {
		elf, me, err := parseRule(rule)
		if err != nil {
			log.Fatalf("Line %d: invalid rule: %v", i+1, err)
		}

		totalScore += score(elf, me)
	}
	fmt.Println(totalScore)
}

var (
	elfMove  = map[string]int{"A": 1, "B": 2, "C": 3}
	selfMove = map[string]int{"X": 1, "Y": 2, "Z": 3}
)

func parseRule(rule string) (elf int, self int, _ error) {
	move := strings.Fields(rule)
	if len(move) != 2 {
		return 0, 0, fmt.Errorf("wrong number or fields (got %d, want 2)", len(move))
	}
	elf, self = elfMove[move[0]], selfMove[move[1]]
	if elf <= 0 {
		return 0, 0, errors.New("invalid elf move")
	}
	if self <= 0 {
		return 0, 0, errors.New("invalid self move")
	}
	return
}

func prev(move int) int { return (move+1)%3 + 1 }

func score(elf, self int) int {
	// rock == 1, paper == 2, scissors == 3
	// Each item beats the one before it in circular order.
	if prev(self) == elf {
		return 6 + self // win
	} else if elf == self {
		return 3 + self // draw
	}
	return 0 + self // loss
}
