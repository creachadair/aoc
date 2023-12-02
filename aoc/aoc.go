package aoc

import (
	"flag"
	"io"
	"log"
	"os"
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
