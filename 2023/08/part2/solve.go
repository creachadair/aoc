package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	lib "github.com/creachadair/aoc/2023/08"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()

	pgm, err := lib.ParseProgram(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Parse program: %v", err)
	}

	var nsteps []int
	for _, c := range pgm.Insn {
		if strings.HasSuffix(c.Label, "A") {
			ns, goal := pgm.Steps(c.Label, "Z")
			nsteps = append(nsteps, ns)
			log.Printf("Starting at %q, %d steps to goal %q", c.Label, ns, goal)
		}
	}
	fmt.Println(nsteps)
	fmt.Println(lcm(nsteps))
}

func lcm(vs []int) int {
	prod := vs[0]
	for i := 1; i < len(vs); i++ {
		g := gcd2(prod, vs[i])
		prod = (prod * vs[i]) / g
	}
	return prod
}

func gcd2(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
