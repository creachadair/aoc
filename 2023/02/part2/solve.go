package main

import (
	"flag"
	"fmt"
	"log"

	lib "github.com/creachadair/aoc/2023/02"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	games, err := lib.ParseGames(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Invalid input: %v", err)
	}

	var sum int
	for _, game := range games {
		min := minSample(game)
		log.Printf("Game %d: min sample is %s (power %d)", game.ID, min, powerOf(min))
		sum += powerOf(min)
	}
	fmt.Println(sum)
}

func minSample(g lib.Game) lib.Sample {
	var min lib.Sample

	for _, s := range g.Samples {
		if s.Red > min.Red {
			min.Red = s.Red
		}
		if s.Green > min.Green {
			min.Green = s.Green
		}
		if s.Blue > min.Blue {
			min.Blue = s.Blue
		}
	}
	return min
}

func powerOf(s lib.Sample) int { return s.Red * s.Green * s.Blue }
