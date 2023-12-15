package lib_test

import (
	"testing"

	lib "github.com/creachadair/aoc/2023/12"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"???.### 1,1,3", 1},
		{".??..??...?##. 1,1,3", 4},
		{"?#?#?#?#?#?#?#? 1,3,1,6", 1},
		{"????.#...#... 4,1,1", 1},
		{"????.######..#####. 1,6,5", 4},
		{"?###???????? 3,2,1", 10},
	}
	for _, tc := range tests {
		r, err := lib.ParseRecords([]byte(tc.input))
		if err != nil {
			t.Fatalf("Invalid input: %v", err)
		} else if len(r) != 1 {
			t.Fatalf("Found %d records, want 1", len(r))
		}
		got := lib.Solve(r[0])
		if got != tc.want {
			t.Errorf("Record %v: got %d, want %d", r[0], got, tc.want)
		}
	}
}
