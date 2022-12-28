package cipher_test

import (
	"reflect"
	"testing"

	"aoc/2022/20/cipher"
)

func TestCipher(t *testing.T) {
	s := cipher.NewState([]int{1, 2, -3, 3, -2, 0, 4})
	check := func(pos, delta int, want ...int) {
		t.Helper()
		s.Move(pos, delta)
		got := s.Values()
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Move(%d, %d):\ngot  %v,\nwant %v", pos, delta, got, want)
		}
	}

	// Make sure we got the initial condition.
	check(0, 0 /**/, 0, 4, 1, 2, -3, 3, -2)

	check(0, 1 /**/, 0, 4, 2, 1, -3, 3, -2)
	check(1, 2 /**/, 0, 4, 1, -3, 2, 3, -2)
	check(2, -3 /**/, 0, 4, 1, 2, 3, -2, -3)
	check(3, 3 /**/, 0, 3, 4, 1, 2, -2, -3)
	check(4, -2 /**/, 0, 3, 4, -2, 1, 2, -3)
	check(5, 0 /**/, 0, 3, 4, -2, 1, 2, -3)
	check(6, 4 /**/, 0, 3, -2, 1, 2, -3, 4)
}

func TestExtra(t *testing.T) {
	s := cipher.NewState([]int{8, 2, 32, -41, 6, 29, -4, 6, -8, 8, -3, -8, 3, -5, 0, -1, 2, 1, 10, -9})
	t.Logf("Input: %v", s.Values())
	want := []int{0, -1, -8, -41, -8, 2, -4, 1, -5, -3, 2, 8, 6, 6, 29, 32, 10, 3, -9, 8}

	s.Mix()
	got := s.Values()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Different output:\n got: %v\nwant: %v", got, want)
	}
}
