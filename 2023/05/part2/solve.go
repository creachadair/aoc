package main

import (
	"flag"
	"fmt"
	"log"

	lib "github.com/creachadair/aoc/2023/05"
	"github.com/creachadair/aoc/aoc"
	"github.com/creachadair/taskgroup"
)

func main() {
	flag.Parse()

	a, err := lib.ParseAlmanac(aoc.MustReadInput())
	if err != nil {
		log.Fatalf("Parse almanac: %v", err)
	}

	mins := expandSeeds(a.Seeds, func(seed int, min *int) {
		loc := a.Track(seed)["location"]
		if *min < 0 || loc < *min {
			*min = loc
		}
	})
	for _, v := range mins {
		if v < mins[0] {
			mins[0] = v
		}
	}
	fmt.Println(mins[0])
}

func expandSeeds(spec []int, f func(seed int, min *int)) []int {
	mins := make([]int, len(spec)/2)
	for i := range mins {
		mins[i] = -1
	}
	g := taskgroup.New(nil)
	for i := 0; i+1 < len(spec); i += 2 {
		i, min := i, &mins[i/2]
		g.Run(func() {
			log.Printf("Seed range: (%d..%d)", spec[i], spec[i]+spec[i+1])
			for j := 0; j < spec[i+1]; j++ {
				f(spec[i]+j, min)
			}
			log.Printf(" done range (%d..%d)", spec[i], spec[i]+spec[i+1])
		})
	}
	g.Wait()
	return mins
}
