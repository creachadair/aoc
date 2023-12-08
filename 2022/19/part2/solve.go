package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/creachadair/aoc/2022/19/robot"
	"github.com/creachadair/mds/mlink"
	"github.com/creachadair/taskgroup"
)

var (
	inputFile  = flag.String("input", "input.txt", "Input file path")
	maxMinutes = flag.Int("minutes", 32, "Maximum minutes to run")
	maxTasks   = flag.Int("tasks", 8, "Concurrent tasks")
	loPlan     = flag.Int("lo", 1, "First plan to simulate")
	hiPlan     = flag.Int("hi", 3, "Last plan to simulate")
)

func main() {
	flag.Parse()

	ps := robot.MustParse(*inputFile)
	fmt.Printf("Loaded %d blueprints from %q\n", len(ps), *inputFile)

	g, run := taskgroup.New(nil).Limit(*maxTasks)

	opt := make([]robot.State, len(ps))
	for i, plan := range ps {
		if (*loPlan > 0 && plan.ID < *loPlan) || (*hiPlan > 0 && plan.ID > *hiPlan) {
			continue
		}

		i, s := i, robot.NewState(plan, *maxMinutes)
		log.Printf("Begin solving plan %v", s.Plan)
		run(func() error {
			opt[i] = search(s, objective, estimate)
			log.Printf("Done solving plan %d: %v", s.Plan.ID, opt[i])
			return nil
		})
	}
	g.Wait()
	prod, score := 1, 0
	for i, best := range opt {
		if (*loPlan == 0 || i+1 >= *loPlan) && (*hiPlan == 0 || i+1 <= *hiPlan) {
			fmt.Printf("Plan %d: optimal solution is %d: %v\n", best.Plan.ID, best.Rock, best)
			if best.Rock > 0 {
				prod *= best.Rock
			}
			score += (best.Rock * best.Plan.ID)
		}
	}
	fmt.Printf("Product: %d, Score: %d\n", prod, score)
}

func objective(s robot.State) int { return -s.Rock }

func estimate(s robot.State) int {
	timeLeft := s.MaxTime - s.Time
	if timeLeft == 0 {
		return objective(s)
	}

	// Bound how many rock bots we could possibly create in the time left by
	// assuming we can create them as fast as their less-expensive reagent in
	// the current plan. Add 1 so the estimate is always nonzero.
	r1, r2 := s.Plan.RockBotOre, s.Plan.RockBotGlass
	if r1 > r2 {
		r1 = r2
	}
	nb := (timeLeft+(r1-1))/r1 + s.RockBots + 1
	return -(s.Rock + timeLeft*nb)
}

func search(s robot.State, obj func(robot.State) int, estimate func(robot.State) int) robot.State {
	bound := 0
	best := s

	q := mlink.NewStack[robot.State]()
	q.Add(s)

	pruned, total, skipped := 0, 0, 0
	seen := make(map[robot.State]bool)
	for !q.IsEmpty() {
		next, _ := q.Pop()
		total++
		if seen[next] {
			skipped++
			continue
		}
		seen[next] = true

		if next.Time >= next.MaxTime {
			if v := obj(next); v < bound {
				log.Printf("- [%d] update bound %d < %d: %v", s.Plan.ID, v, bound, next)
				best = next
				bound = v
			}
			continue
		}

		for _, opt := range next.Options() {
			if v := estimate(opt); v > bound {
				pruned++
			} else {
				q.Add(opt)
			}
		}
	}
	log.Printf("- [%d] done: total %d, pruned %d, skipped %d", s.Plan.ID, total, pruned, skipped)
	return best
}
