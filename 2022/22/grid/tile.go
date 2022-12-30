package grid

import (
	"fmt"
	"log"
)

type Tiles struct {
	Size int
	Tile []Tile
}

func (t Tiles) getTile(label byte) Tile {
	for _, tile := range t.Tile {
		if tile.Label == label {
			return tile
		}
	}
	panic("tile not found")
}

func (t Tiles) Next(label byte, dir Heading) (Tile, Heading) {
	next := nextTile[label][dir]
	return t.getTile(next.label), next.side
}

func (t Tiles) Normalize(m *Map) {
	for _, tile := range t.Tile {
		for i := tile.Rot; i > 0; i-- {
			tile.rotateLeft(m)
		}
	}
}

func (t Tiles) String() string {
	return fmt.Sprintf("Tiles(size=%d, n=%d, %v)", t.Size, len(t.Tile), t.Tile)
}

type Tile struct {
	T, L, Size     int     // top-left and size
	NT, NL, NB, NR int     // neighboring tile index (0 == none)
	Label          byte    // cosmetic label (Dorsal, Ventral, Anterior, Posterior, Left, Right)
	Rot            Heading // rotation relative to "upright" (0 == no change)
}

func (t Tile) rotateLeft(m *Map) {
	t.transpose(m)
	t.vflip(m)
}

func (t Tile) transpose(m *Map) {
	for d := 0; d < t.Size; d++ {
		r, c := t.T+d, t.L+d
		for i := 0; i < t.Size-d; i++ {
			p1, p2 := m.pos(r, c+i), m.pos(r+i, c)
			m.Data[p1], m.Data[p2] = m.Data[p2], m.Data[p1]
		}
	}
}

func (t Tile) vflip(m *Map) {
	for c := t.L; c <= t.L+t.Size; c++ {
		lo, hi := t.T, t.T+t.Size-1
		for lo < hi {
			p1, p2 := m.pos(lo, c), m.pos(hi, c)
			m.Data[p1], m.Data[p2] = m.Data[p2], m.Data[p1]
			lo++
			hi--
		}
	}
}

func (t Tile) Translate(src Tile, r, c int, oldHead, newHead Heading) (nr, nc int) {
	// One of the dimensions is out of range, the other is in-range and
	// determines which row or column we wind up in.
	dim := r - src.T
	if dim < 0 || dim >= src.T {
		dim = c - src.L
	}

	// If we are making a 90° turn we have to invert the dimension.
	if (oldHead.R() == newHead && (oldHead == Left || oldHead == Right)) ||
		(oldHead.L() == newHead && (oldHead == Up || oldHead == Down)) {
		dim = src.Size - dim + 1
	}

	log.Printf("TRANSLATE (%v,%v) [%v] entering %v facing %v", r, c, dim, t, newHead)
	switch newHead {
	case Down:
		nr, nc = t.T, t.L+dim
	case Up:
		nr, nc = t.T+t.Size-1, t.L+dim
	case Right:
		nr, nc = t.T+dim, t.L
	case Left:
		nr, nc = t.T+dim, t.L+t.Size-1
	}
	return
}

func (t Tile) Contains(r, c int) bool {
	return r >= t.T && r < t.T+t.Size && c >= t.L && c < t.L+t.Size
}

func (t Tile) String() string {
	return fmt.Sprintf("Tile(%d|%d,%d:N[%d,%d,%d,%d]:=%c+%d)",
		t.Size, t.T, t.L, t.NT, t.NL, t.NB, t.NR, t.Label, t.Rot)
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
				tiles = append(tiles, Tile{T: r, L: c, Size: ts})
			}
		}
	}

	// Find neighboring tiles.
	for i, a := range tiles {
		for j, b := range tiles {
			if a.T == b.T && a.L+a.Size == b.L {
				// b[j] is the right neighbor of a[i]
				tiles[i].NR = j + 1
				tiles[j].NL = i + 1
			}
			if a.L == b.L && a.T+a.Size == b.T {
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
var adjRule = [...][]rule{ // :: dir → T → (X, rot)
	Left:  {'A': {'L', 1}, 'P': {'L', 3}, 'L': {'A', 3}, 'R': {'-', 0}, 'D': {'L', 2}},
	Right: {'A': {'R', 3}, 'P': {'R', 1}, 'L': {'-', 0}, 'R': {'D', 2}, 'D': {'R', 2}},
	Up:    {'A': {'D', 0}, 'P': {'-', 0}, 'L': {'A', 3}, 'R': {'A', 1}, 'D': {'P', 0}},
	Down:  {'A': {'-', 0}, 'P': {'D', 0}, 'L': {'A', 3}, 'R': {'P', 3}, 'D': {'A', 2}},
}

type tlink struct {
	label byte
	side  Heading
}

var nextTile = [...][]tlink{ // :: tile → dir → (tile, side)
	'V': {Up: {'A', Up}, Down: {'P', Down}, Left: {'L', Left}, Right: {'R', Right}},
	'D': {Up: {'P', Up}, Down: {'A', Down}, Left: {'L', Right}, Right: {'R', Left}},
	'A': {Up: {'D', Up}, Down: {'V', Down}, Left: {'L', Down}, Right: {'R', Down}},
	'P': {Up: {'V', Up}, Down: {'D', Down}, Left: {'L', Up}, Right: {'R', Up}},
	'L': {Up: {'A', Right}, Down: {'P', Right}, Left: {'D', Right}, Right: {'V', Right}},
	'R': {Up: {'A', Left}, Down: {'P', Left}, Left: {'V', Left}, Right: {'D', Left}},
}
