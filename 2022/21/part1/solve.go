package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"aoc/2022/21/rules"
)

var (
	inputFile  = flag.String("input", "input.txt", "Input file path")
	writeGraph = flag.String("graph", "", "Write graph to this file")
)

func main() {
	flag.Parse()

	rs := rules.MustParse(*inputFile)
	g := rules.NewGraph(rs)
	g.PreSolve()
	v, ok := g.Values["root"]
	if !ok {
		log.Fatal("Missing root value")
	}
	fmt.Printf("%d\n", int(v))

	if *writeGraph != "" {
		f, err := os.Create(*writeGraph)
		if err != nil {
			log.Fatalf("Creating output file: %v", err)
		}
		g.Dot(f)
		if err := f.Close(); err != nil {
			log.Fatalf("Writing graph: %v", err)
		}
	}
}
