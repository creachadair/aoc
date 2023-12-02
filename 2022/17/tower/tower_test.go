package tower_test

import (
	"fmt"
	"testing"

	"github.com/creachadair/aoc/2022/17/tower"
)

func TestTower(t *testing.T) {
	w := tower.New(7, ">>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>")
	w.Drop(tower.HLine())
	fmt.Print(w, "\n\n")

	w.Drop(tower.Plus())
	fmt.Print(w, "\n\n")

	w.Drop(tower.Angle())
	fmt.Print(w, "\n\n")

	w.Drop(tower.VLine())
	fmt.Print(w, "\n\n")

	w.Drop(tower.Box())
	fmt.Print(w, "\n\n")

	w.Drop(tower.HLine())
	fmt.Print(w, "\n\n")
}
