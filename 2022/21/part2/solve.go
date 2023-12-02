package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/creachadair/aoc/2022/21/rules"
)

var (
	inputFile = flag.String("input", "input.txt", "Input file path")
	rootNode  = flag.String("root", "root", "Name of the root")
	varNode   = flag.String("var", "humn", "Name of the variable node")
)

func main() {
	flag.Parse()

	rs := rules.MustParse(*inputFile)
	g := rules.NewGraph(rs)
	g.Rules[*rootNode].Op = "=" // mark the root as an equality test
	g.Rules[*varNode].Op = "?"  // mark the variable node as unresolved
	g.PreSolve()

	v, err := g.Solve(*rootNode, 0)
	if err != nil {
		log.Fatalf("Solve %q failed: %v", *rootNode, err)
	}

	fmt.Println("Root:", int(v))
	fmt.Println("Var:", int(g.Values[*varNode]))
}
