package aoc

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// MustReadInput reads the contents of the first command-line argument, or if
// none is specified it fully consumes the contents of stdin.  In case of
// error, MustReadInput calls log.Fatal.
func MustReadInput() []byte {
	var data []byte
	var err error
	if flag.NArg() == 0 {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(flag.Arg(0))
	}
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}
	return data
}

// SplitLines splits input into lines.
func SplitLines(input []byte) []string {
	return strings.Split(strings.TrimSpace(string(input)), "\n")
}

// ParseInts parses strings as integers.
func ParseInts(ss []string) ([]int, error) {
	out := make([]int, len(ss))
	for i, s := range ss {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("value %d: %w", i+1, err)
		}
		out[i] = v
	}
	return out, nil
}
