package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/creachadair/aoc/2022/17/tower"
)

var (
	inputFile = flag.String("input", "input.txt", "Input file path")
	numDrops  = flag.Int("drop", 2022, "Number of shapes to drop")
	doPrint   = flag.Bool("v", false, "Print resulting tower")
	doTrim    = flag.Bool("trim", false, "Trim output before printing")
)

func main() {
	flag.Parse()
	data, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}
	jets := strings.TrimSpace(string(data))

	t := tower.New(7, jets)
	next := tower.NewSequence()
	for i := 0; i < *numDrops; i++ {
		t.Drop(next())
	}
	if *doTrim {
		t.Trim()
	}
	if *doPrint {
		fmt.Println(t)
	}
	fmt.Println(t.Height())
}
