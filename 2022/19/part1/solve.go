package main

import (
	"flag"
	"fmt"
	"log"
	"sync/atomic"

	"aoc/2022/19/robot"

	"github.com/creachadair/taskgroup"
)

var (
	inputFile  = flag.String("input", "input.txt", "Input file path")
	maxMinutes = flag.Int("minutes", 24, "Maximum minutes to run")
	maxTasks   = flag.Int("tasks", 8, "Concurrent tasks")
	doVerbose  = flag.Bool("v", false, "Verbose logging")
)

func main() {
	flag.Parse()

	ps := robot.MustParse(*inputFile)
	fmt.Printf("Loaded %d blueprints from %q\n", len(ps), *inputFile)

	g, run := taskgroup.New(nil).Limit(*maxTasks)

	var total int64
	for _, plan := range ps {
		s := robot.NewState(plan, *maxMinutes)
		run(func() error {
			log.Printf("Begin solving plan %d", s.Plan.ID)

			stats, best := robot.Solve(s)
			score := stats.Best * s.Plan.ID
			fmt.Printf("Plan %d yields maximum %d geodes, score=%d\n",
				s.Plan.ID, stats.Best, score)
			fmt.Println(stats)
			atomic.AddInt64(&total, int64(score))

			if *doVerbose {
				for _, cur := range best {
					log.Printf("- %v", cur)
				}
			}
			return nil
		})
	}
	g.Wait()
	fmt.Println(atomic.LoadInt64(&total))
}
