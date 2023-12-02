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
	numDrops  = flag.Int("drop", 1000000000000, "Number of shapes to drop")
	doPrint   = flag.Bool("v", false, "Print resulting tower")
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

	// As we add items to the tower, we will eventually form a cycle within the
	// jet sequence. Once we do that, continuing to add items from that point on
	// will increase the height by the same amount.
	//
	// That is, once we reach time 3, the difference h3 - h2 is the height we
	// will gain by adding items until that cycle repeats again.
	//
	// height:          h1              h2              h3
	// time:    0-------1---------------2---------------3-----...
	// drops:   0       d1              d2              d3
	//
	// Eventually we will reach a point where the remaining drops we want to do
	// are not sufficient to complete a cycle, at which point we'll need to sim
	// the remainder to find out the height change:
	//
	// height:          hr         hn
	// time:    ...-----r----------$
	// drops:           dr         n ← (n - dr) < (d2 - d1)

	// Stage 1: Run until we find the first jet phase cycle.
	nd := 0
	seen := make(map[int]bool)
	for {
		t.Drop(next())
		t.Trim()
		nd++

		if seen[t.Phase()] {
			break
		}
		seen[t.Phase()] = true
	}
	keyPhase := t.Phase()
	initHeight, initLength := t.Height(), nd
	log.Printf("Stage 1: init length=%d height=%d phase=%d", initLength, initHeight, keyPhase)

	// Stage 2: Run the cycle to find its height delta.
	cdStart, chStart := nd, t.Height()
	for {
		t.Drop(next())
		t.Trim()
		nd++

		if t.Phase() == keyPhase {
			break
		}
	}
	cycleHeight, cycleLength := t.Height()-chStart, nd-cdStart
	log.Printf("Stage 2: cycle length=%d height=%d", cycleLength, cycleHeight)

	// Stage 3: Compute the remainder and simulate the trailing drops.
	remLength := *numDrops - initLength // drops after the leader and the first cycle
	cycles, remainder := remLength/cycleLength, remLength%cycleLength
	rhStart := t.Height()
	for i := 0; i < remainder; i++ {
		t.Drop(next())
		t.Trim()
	}
	remHeight := t.Height() - rhStart
	log.Printf("Stage 3: remainder length=%d cycles=%d height=%d", remainder, cycles, remHeight)

	// Stage 4: Patch it all together.
	totalHeight := initHeight + cycles*cycleHeight + remHeight
	log.Printf("Stage 3: total=%d", totalHeight)

	if *doPrint {
		t.Trim()
		fmt.Println(t)
	}
}
