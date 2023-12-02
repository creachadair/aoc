package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/creachadair/aoc/aoc"
)

// The obvious way to solve this is with a regexp on the spellings, but that
// doesn't work because the puzzle implicitly wants you to separate overlapping
// spellings, e.g., "oneight" should be treated as "18" not "1ight" or "on8".
//
// You will probably only discover this by getting the wrong answer, although
// the instructions carefully permit this interpretation.
var val = map[string]string{
	"one": "1", "1": "1",
	"two": "2", "2": "2",
	"three": "3", "3": "3",
	"four": "4", "4": "4",
	"five": "5", "5": "5",
	"six": "6", "6": "6",
	"seven": "7", "7": "7",
	"eight": "8", "8": "8",
	"nine": "9", "9": "9",
}

func main() {
	flag.Parse()

	input := string(aoc.MustReadInput())
	var sum int
	for _, line := range strings.Split(input, "\n") {
		buf := ""
		for line != "" {
			for p, d := range val {
				if strings.HasPrefix(line, p) {
					buf += d
					break
				}
			}
			line = line[1:]
		}
		if buf == "" {
			continue
		}
		num := int(buf[0]-'0')*10 + int(buf[len(buf)-1]-'0')
		sum += num
	}
	fmt.Println(sum)
}
