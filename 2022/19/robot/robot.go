package robot

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/creachadair/mds/stack"
)

var printRE = regexp.MustCompile(`(?m)Blueprint (\d+):` +
	`\s+Each ore robot costs (\d+) ore.` +
	`\s+Each clay robot costs (\d+) ore.` +
	`\s+Each obsidian robot costs (\d+) ore and (\d+) clay.` +
	`\s+Each geode robot costs (\d+) ore and (\d+) obsidian.`)

func MustParse(path string) []*Blueprint {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}
	mustParse := func(s string) int {
		v, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("Invalid input %q: %v", s, err)
		}
		return v
	}

	var prints []*Blueprint
	for _, m := range printRE.FindAllStringSubmatch(string(data), -1) {
		prints = append(prints, &Blueprint{
			ID:           mustParse(m[1]),
			OreBotOre:    mustParse(m[2]),
			ClayBotOre:   mustParse(m[3]),
			GlassBotOre:  mustParse(m[4]),
			GlassBotClay: mustParse(m[5]),
			RockBotOre:   mustParse(m[6]),
			RockBotGlass: mustParse(m[7]),
		})
	}
	return prints
}

type Blueprint struct {
	ID           int
	OreBotOre    int
	ClayBotOre   int
	GlassBotOre  int
	GlassBotClay int
	RockBotOre   int
	RockBotGlass int
}

func (b *Blueprint) String() string {
	return fmt.Sprintf("Blueprint(id=%d, ob=%d, cb=%d, gb=%d/%d, rb=%d/%d)",
		b.ID, b.OreBotOre, b.ClayBotOre, b.GlassBotOre, b.GlassBotClay, b.RockBotOre, b.RockBotGlass)
}

func NewState(plan *Blueprint, maxTime int) State {
	return State{Plan: plan, OreBots: 1, MaxTime: maxTime}
}

type State struct {
	Plan *Blueprint

	// Robots active
	OreBots   int
	ClayBots  int
	GlassBots int
	RockBots  int

	// Resources available
	Ore   int
	Clay  int
	Glass int
	Rock  int

	// Time tracking
	Step    string
	Time    int
	MaxTime int
}

func (s State) Options() []State {
	// If we're out of time, there's nothing we can do from here.
	if s.Time >= s.MaxTime {
		return nil
	}

	base := s
	base.Ore += s.OreBots
	base.Clay += s.ClayBots
	base.Glass += s.GlassBots
	base.Rock += s.RockBots
	base.Time++

	// Just waiting is always an option.
	base.Step = "wait"
	opts := []State{base}

	// Don't bother building any robots at the last time step, since it cannot
	// possibly affect the outcome.
	if s.Time+1 < s.MaxTime {
		// Build a geode-cracking robot.
		if s.Ore >= s.Plan.RockBotOre && s.Glass >= s.Plan.RockBotGlass {
			cp := base
			cp.Step = "geode"
			cp.Ore -= s.Plan.RockBotOre
			cp.Glass -= s.Plan.RockBotGlass
			cp.RockBots++
			opts = append(opts, cp)
		}

		// Build an obsidian-collecting robot.
		if s.Ore >= s.Plan.GlassBotOre && s.Clay >= s.Plan.GlassBotClay {
			cp := base
			cp.Step = "glass"
			cp.Ore -= s.Plan.GlassBotOre
			cp.Clay -= s.Plan.GlassBotClay
			cp.GlassBots++
			opts = append(opts, cp)
		}

		// Build a clay robot.
		if s.Ore >= s.Plan.ClayBotOre {
			cp := base
			cp.Step = "clay"
			cp.Ore -= s.Plan.ClayBotOre
			cp.ClayBots++
			opts = append(opts, cp)
		}

		// Build an ore robot.
		if s.Ore >= s.Plan.OreBotOre {
			cp := base
			cp.Step = "ore"
			cp.Ore -= s.Plan.OreBotOre
			cp.OreBots++
			opts = append(opts, cp)
		}
	}

	return opts
}

func (s State) String() string {
	return fmt.Sprintf("State(plan=%d %q ob=%d/%d, cb=%d/%d, gb=%d/%d, rb=%d/%d, t=%d/%d)",
		s.Plan.ID, s.Step, s.OreBots, s.Ore, s.ClayBots, s.Clay, s.GlassBots, s.Glass, s.RockBots, s.Rock,
		s.Time, s.MaxTime)
}

func (s State) isBetter(t State) bool { return s.Rock > t.Rock }

type Entry struct {
	State
	Link *Entry
}

type Stats struct {
	Best    int
	Visited int
	Unique  int
	Seen    int
	Improve int
	Elapsed time.Duration
}

func (s Stats) String() string {
	return fmt.Sprintf("Stats: best=%d, visited=%d, unique=%d, seen=%d, improved=%d (%v elapsed)",
		s.Best, s.Visited, s.Unique, s.Seen, s.Improve, s.Elapsed)
}

func Solve(s State) (stats Stats, _ []State) {
	stats, all := solveAllSeen(s, make(map[State]bool))
	return stats, all[0]
}

func unravel(e *Entry) []State {
	if e == nil {
		return nil
	}
	return append(unravel(e.Link), e.State)
}

func SolveStats(s State) (stats Stats) {
	start := time.Now()
	defer func() { stats.Elapsed = time.Since(start) }()

	addImp := func(s State) {
		if s.Rock > stats.Best {
			stats.Best = s.Rock
			stats.Improve++
		}
	}

	seen := make(map[State]bool)
	mark := func(s State) bool {
		if seen[s] {
			stats.Seen++
			return true
		}
		seen[s] = true
		stats.Unique = len(seen)
		return false
	}

	q := stack.New[State]()
	q.Add(s)
	for !q.IsEmpty() {
		next, _ := q.Pop()
		stats.Visited++

		if next.Time == next.MaxTime {
			addImp(next)
		}
		for _, opt := range next.Options() {
			if !mark(opt) {
				q.Add(opt)
			}
		}
	}

	return stats
}

func solveAllSeen(s State, seen map[State]bool) (stats Stats, _ [][]State) {
	start := time.Now()
	defer func() { stats.Elapsed = time.Since(start) }()

	var imps []Entry
	addImp := func(e Entry) bool {
		if n := len(imps) - 1; n >= 0 {
			if e.isBetter(imps[n].State) {
				imps = imps[:0]
			} else if imps[n].isBetter(e.State) {
				return false
			}
		}
		stats.Best = e.Rock
		imps = append(imps, e)
		return true
	}

	mark := func(s State) bool {
		if seen[s] {
			stats.Seen++
			return true
		}
		seen[s] = true
		stats.Unique = len(seen)
		return false
	}

	q := stack.New[Entry]()
	q.Add(Entry{State: s})
	for !q.IsEmpty() {
		e, _ := q.Pop()
		stats.Visited++

		if e.Time == e.MaxTime && addImp(e) {
			stats.Improve++
		}
		for _, opt := range e.Options() {
			if !mark(opt) {
				q.Add(Entry{
					State: opt,
					Link:  &e,
				})
			}
		}
	}

	var out [][]State
	for _, imp := range imps {
		out = append(out, unravel(&imp))
	}

	return stats, out
}
