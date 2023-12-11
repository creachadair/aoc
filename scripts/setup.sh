#!/usr/bin/env bash
#
# Set up a skeleton for an AoC puzzle.
#
set -euo pipefail

# Use tomorrow's date, since puzzles ship at midnight ET.
today="$(date -v +1d +%Y/%d)"
if [[ -d "$today" ]] ; then
    echo "Directory $today is already set up" 1>&2
    exit 1
fi
mkdir -p "${today}/part1" "${today}/part2"
touch "${today}/puzzle.md" \
      "${today}/input.txt" \
      "${today}/example.txt" \
      "${today}/part2/puzzle.md"
goimports >"${today}/lib.go" <<EOF
package lib
func OK(any) {}
EOF
goimports >"${today}/part1/solve.go" <<EOF
package main
import lib "github.com/creachadair/aoc/${today}"
func main() {
  flag.Parse()
  lib.OK(aoc.MustReadInput())
}
EOF
cp "${today}/part1/solve.go" "${today}/part2"

