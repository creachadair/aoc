package main

import (
	"flag"
	"fmt"
	"log"

	lib "github.com/creachadair/aoc/2023/02"
	"github.com/creachadair/aoc/aoc"
)

var (
	numRed   = flag.Int("r", 12, "Number of red cubes")
	numGreen = flag.Int("g", 13, "Number of green cubes")
	numBlue  = flag.Int("b", 14, "Number of blue cubes")
)

func main() {
	flag.Parse()

	games, err := lib.ParseGames(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Invalid input: %v", err)
	}

	var sum int
nextGame:
	for _, game := range games {
		for _, samp := range game.Samples {
			if samp.Red > *numRed || samp.Green > *numGreen || samp.Blue > *numBlue {
				log.Printf("Game %d is impossible (sample: %v)", game.ID, samp)
				continue nextGame
			}
		}
		log.Printf("Game %d is possible", game.ID)
		sum += game.ID
	}
	fmt.Println(sum)
}
