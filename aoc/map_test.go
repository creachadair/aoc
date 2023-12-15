package aoc_test

import (
	"testing"

	"github.com/creachadair/aoc/aoc"
)

func TestMapTranspose(t *testing.T) {
	tests := []struct {
		r, c int
		text string
		want string
	}{
		{1, 1, "a", "a"},
		{1, 3, "abc", "abc"},
		{3, 1, "abc", "abc"},
		{2, 2, "abcd", "acbd"},
		{3, 5, "abcdefghijklmno", "afkbglchmdinejo"},
		{4, 3, "abcdefghijkl", "aeibfjcgkdhl"},
	}
	for _, tc := range tests {
		got := aoc.NewMap(tc.r, tc.c, []byte(tc.text)).Transpose()
		want := aoc.NewMap(tc.c, tc.r, []byte(tc.want))
		if !got.Equal(want) {
			t.Errorf("Transpose: got (%d, %d)\n%s; want (%d, %d)\n%s",
				got.Rows(), got.Cols(), got.String(),
				want.Rows(), want.Cols(), want.String())
		}
	}
}
