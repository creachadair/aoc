package grid

import (
	"fmt"
)

type Tiles struct {
	Size int
	Tile []Tile
}

func (t Tiles) String() string {
	return fmt.Sprintf("Tiles(size=%d, n=%d, %v)", t.Size, len(t.Tile), t.Tile)
}

type Tile struct {
	T, L, B, R     int     // bounding box in standard coordinates
	NT, NL, NB, NR int     // neighboring tile index (0 == none)
	Label          byte    // cosmetic label (Dorsal, Ventral, Anterior, Posterior, Left, Right)
	Rot            Heading // rotation relative to "upright" (0 == no change)
}

func (t Tile) String() string {
	return fmt.Sprintf("Tile(%d,%d|%d,%d:N[%d,%d,%d,%d]:=%c+%d)",
		t.T, t.L, t.B, t.R, t.NT, t.NL, t.NB, t.NR, t.Label, t.Rot)
}

func findTiles(m *Map) Tiles {
	key := m.At(1, 1) == ' '

	var ts int
	for (m.At(ts+1, ts+1) == ' ') == key {
		ts++
	}

	var tiles []Tile

	for r := 1; r < m.NR; r += ts {
		for c := 1; c < m.NC; c += ts {
			if m.At(r, c) != ' ' {
				tiles = append(tiles, Tile{T: r, L: c, B: r + ts - 1, R: c + ts - 1})
			}
		}
	}

	// Find neighboring tiles.
	for i, a := range tiles {
		for j, b := range tiles {
			if a.T == b.T && a.R+1 == b.L {
				// b[j] is the right neighbor of a[i]
				tiles[i].NR = j + 1
				tiles[j].NL = i + 1
			}
			if a.L == b.L && a.B+1 == b.T {
				// b[j] is the lower neighbor of a[i]
				tiles[i].NB = j + 1
				tiles[j].NT = i + 1
			}

			// We don't have to check the other direction, we'll pick it up when
			// we hit those indices in the opposite order.
		}
	}

	// Solve rotations.
	// Phase 1: Designate one tile as the bottom (ventral), and mark its
	// neighbors.
	//
	//       +---+
	//       | A |
	//   +---+---+---+
	//   | L | V | R |
	//   +---+---+---+
	//       | P |
	//       +---+
	//
	// Note that not all of these may exist, depending on the input; that's
	// fine.  These neighbors do not require rotation relative to V.

	// We designate the first tile we found since the problem wants us to start
	// out in the "normal" orientation.
	tiles[0].Label = 'V'
	if t := tiles[0].NL; t != 0 {
		tiles[t-1].Label = 'L'
	}
	if t := tiles[0].NT; t != 0 {
		tiles[t-1].Label = 'A'
	}
	if t := tiles[0].NR; t != 0 {
		tiles[t-1].Label = 'R'
	}
	if t := tiles[0].NB; t != 0 {
		tiles[t-1].Label = 'P'
	}

	// Phase 2: Figure out what the other unfilled tiles are. Where a tile
	// depends up depends on how the original cube was cut apart, and the
	// relationships are deterministic.
	for {
		var more bool
		for i, tile := range tiles {
			if tile.Label != 0 {
				continue // skip solved tiles
			}

			var r rule
			var trot Heading
			if t := tile.NR; t != 0 && tiles[t-1].Label != 0 {
				r = adjRule[Left][tiles[t-1].Label]
				trot = tiles[t-1].Rot
			} else if t := tile.NB; t != 0 && tiles[t-1].Label != 0 {
				r = adjRule[Up][tiles[t-1].Label]
				trot = tiles[t-1].Rot
			} else if t := tile.NL; t != 0 && tiles[t-1].Label != 0 {
				r = adjRule[Right][tiles[t-1].Label]
				trot = tiles[t-1].Rot
			} else if t := tile.NT; t != 0 && tiles[t-1].Label != 0 {
				r = adjRule[Down][tiles[t-1].Label]
				trot = tiles[t-1].Rot
			} else {
				more = true
				continue
			}

			tiles[i].Label = r.role
			tiles[i].Rot = Heading((trot + Heading(r.rot)) % 4)
		}

		if !more {
			break
		}
	}

	return Tiles{Size: ts, Tile: tiles}
}

/*
Given an unsolved tile X, if X is adjacent to Z, we can determine the role
and orientation of X based on the resolved tiles adjacent to it.

Note that no unsolved tile can be adjacent to V, since we filled those
already.  In this table, +n (-1) means the tile is n*90 degrees rotated
right (left) of Z.

	(Z)

X position   | A    P    L    R    D
-------------|-------------------------
X left of Z  | L+1  L-1  D+2  -    L+2  < role and orientation relative to Z
X right of Z | R-1  R+1  -    D+2  R+2
X above Z    | D+0  -    A-1  A+1  P+0
X below Z    | -    D+0  P+1  P-1  A+2
*/
type rule struct {
	role byte
	rot  byte
}

// N.B. Use +3 instead of -1 so we don't need signed values.
var adjRule = [...][]rule{
	Left:  {'A': {'L', 1}, 'P': {'L', 3}, 'L': {'A', 3}, 'R': {'-', 0}, 'D': {'L', 2}},
	Right: {'A': {'R', 3}, 'P': {'R', 1}, 'L': {'-', 0}, 'R': {'D', 2}, 'D': {'R', 2}},
	Up:    {'A': {'D', 0}, 'P': {'-', 0}, 'L': {'A', 3}, 'R': {'A', 1}, 'D': {'P', 0}},
	Down:  {'A': {'-', 0}, 'P': {'D', 0}, 'L': {'A', 3}, 'R': {'P', 3}, 'D': {'A', 2}},
}
