package main

import (
	"flag"
	"fmt"
	"log"

	lib "github.com/creachadair/aoc/2023/16"
	"github.com/creachadair/aoc/aoc"
)

var _ = lib.OK

func main() {
	flag.Parse()
	m, err := aoc.ParseMap(aoc.MustReadLines())
	if err != nil {
		log.Fatalf("Parse map: %v", err)
	}
	fmt.Println(m)
}
