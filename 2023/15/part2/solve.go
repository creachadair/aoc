package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	lib "github.com/creachadair/aoc/2023/15"
	"github.com/creachadair/aoc/aoc"
)

func main() {
	flag.Parse()
	input := strings.TrimSpace(string(aoc.MustReadInput()))
	rules := strings.Split(input, ",")

	type lens struct {
		label string
		focal byte
	}
	box := make([][]lens, 256)
	remove := func(label string) {
		h := lib.Hash(label)
		for j := 0; j < len(box[h]); j++ {
			if box[h][j].label == label {
				box[h] = append(box[h][:j], box[h][j+1:]...)
				return
			}
		}
	}
	replace := func(label string, focal int) {
		h := lib.Hash(label)
		for j := 0; j < len(box[h]); j++ {
			if box[h][j].label == label {
				box[h][j] = lens{label: label, focal: byte(focal)}
				return
			}
		}
		box[h] = append(box[h], lens{label: label, focal: byte(focal)})
	}

	for i, rule := range rules {
		label, focal, ok := strings.Cut(rule, "=")
		if !ok {
			remove(strings.TrimSuffix(label, "-"))
		} else if v, err := strconv.Atoi(focal); err != nil {
			log.Fatalf("Rule %d (%q): invalid focal length: %v", i+1, rule, err)
		} else {
			replace(label, v)
		}
	}

	var totalPower int
	for i, b := range box {
		var boxPower int
		for j, e := range b {
			boxPower += (i + 1) * (j + 1) * int(e.focal)
		}
		if boxPower != 0 {
			log.Printf("Box %d: power=%d", i+1, boxPower)
		}
		totalPower += boxPower
	}
	fmt.Println(totalPower)
}
