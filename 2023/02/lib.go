package lib

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/creachadair/aoc/aoc"
	"github.com/creachadair/mds/slice"
)

type Game struct {
	ID      int
	Samples []Sample
}

func (g Game) String() string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "Game %d: ", g.ID)
	var ss []string
	for _, s := range g.Samples {
		ss = append(ss, s.String())
	}
	buf.WriteString(strings.Join(ss, "; "))
	return buf.String()
}

type Sample struct {
	Red   int
	Green int
	Blue  int
}

func parseSample(s string) (Sample, error) {
	var samp Sample
	for _, part := range strings.Split(s, ",") {
		fs := strings.Fields(part)
		if len(fs) != 2 {
			return samp, errors.New("invalid sample format")
		}
		n, err := strconv.Atoi(fs[0])
		if err != nil {
			return samp, fmt.Errorf("invalid %q count: %w", fs[1], err)
		}
		switch fs[1] {
		case "red":
			samp.Red = n
		case "green":
			samp.Green = n
		case "blue":
			samp.Blue = n
		default:
			return samp, fmt.Errorf("invalid color %q", fs[1])
		}
	}
	return samp, nil
}

func (s Sample) String() string {
	var ss []string
	if s.Red > 0 {
		ss = append(ss, fmt.Sprintf("%d red", s.Red))
	}
	if s.Green > 0 {
		ss = append(ss, fmt.Sprintf("%d green", s.Green))
	}
	if s.Blue > 0 {
		ss = append(ss, fmt.Sprintf("%d blue", s.Blue))
	}
	return strings.Join(ss, ", ")
}

func ParseGames(input []byte) ([]Game, error) {
	var games []Game
	for i, line := range aoc.SplitLines(input) {
		gid, rest, ok := strings.Cut(line, ":")
		if !ok {
			return nil, fmt.Errorf("line %d: missing game ID label", i+1)
		}
		id, err := strconv.Atoi(slice.At(strings.Fields(gid), -1))
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid game ID: %w", i+1, err)
		}
		game := Game{ID: id}
		for _, s := range strings.Split(strings.TrimSpace(rest), ";") {
			samp, err := parseSample(s)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid sample %q: %w", i+1, s, err)
			}
			game.Samples = append(game.Samples, samp)
		}
		games = append(games, game)
	}
	return games, nil
}
