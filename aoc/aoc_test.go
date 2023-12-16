package aoc_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/creachadair/aoc/aoc"
	"github.com/google/go-cmp/cmp"
)

type thing string

func (t *thing) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		*t = ""
		return nil
	} else if len(data) < 2 || data[0] != '<' || data[len(data)-1] != '>' {
		return errors.New("bad thing")
	}
	*t = thing(data[1 : len(data)-1])
	return nil
}

func TestScanx(t *testing.T) {
	re := regexp.MustCompile(`Z=(\d+), *ZL=(?P<comma>\d+(?:,\d+)*), *S=([^,]+), *L=(?P<words>\w+(?: *\w+)*), *T=(.+)$`)

	type data struct {
		Z  int
		Zs []int
		S  string
		Ss []string
		T  thing
	}

	var got data
	if err := aoc.Scanx(re, `Z=122,ZL=1,15,19, S=alpha beta gamma, L=foo bar  baz,     T=<wheep>`,
		&got.Z, &got.Zs, &got.S, &got.Ss, &got.T); err != nil {
		t.Fatalf("Scanx: unexpected error: %v", err)
	}
	if diff := cmp.Diff(got, data{
		Z:  122,
		Zs: []int{1, 15, 19},
		S:  "alpha beta gamma",
		Ss: []string{"foo", "bar", "baz"},
		T:  "wheep",
	}); diff != "" {
		t.Errorf("Scanx (-got, +want):\n%s", diff)
	}
}
