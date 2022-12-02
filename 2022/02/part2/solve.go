package main

import (
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
	elfMove   = map[string]int{"A": 1, "B": 2, "C": 3}
	wantScore = map[string]int{"X": 0, "Y": 3, "Z": 6}
)

func parseRule(rule string) (elf int, score int, _ error) {
	move := strings.Fields(rule)
	if len(move) != 2 {
		return 0, 0, fmt.Errorf("wrong number or fields (got %d, want 2)", len(move))
	}
	elf = elfMove[move[0]]
	if elf <= 0 {
		return 0, 0, fmt.Errorf("invalid elf move %q", move[0])
	}
	score, ok := wantScore[move[1]]
	if !ok {
		return 0, 0, fmt.Errorf("invalid target score %q", move[1])
	}
	return
}

func prev(move int) int { return (move+1)%3 + 1 }
func succ(move int) int { return (move+3)%3 + 1 }

func score(elf, score int) int {
	// rock == 1, paper == 2, scissors == 3
	// Each item beats the one before it in circular order.
	switch score {
	case 0:
		return score + prev(elf) // play to a loss
	case 3:
		return score + elf // play to a draw
	case 6:
		return score + succ(elf) // play to a win
	default:
		panic("impossible")
	}
}
