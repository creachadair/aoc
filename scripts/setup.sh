#!/usr/bin/env bash
#
# Set up a skeleton for an AoC puzzle.
#
set -euo pipefail

# If it's before 9pm, use today's date; otherwise tomorrow.
if [[ "$(date +%k)" -lt 21 ]] ; then
    today="$(date +%Y/%d)"
else
    today="$(date -v +1d +%Y/%d)"
fi
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
var OK string
EOF
goimports >"${today}/part1/solve.go" <<EOF
package main
import lib "github.com/creachadair/aoc/${today}"

var _ = lib.OK

func main() {
  flag.Parse()
  m, err := aoc.ParseMap(aoc.MustReadLines())
  if err != nil {
    log.Fatalf("Parse map: %v", err)
  }
  fmt.Println(m)
}
EOF
cp "${today}/part1/solve.go" "${today}/part2"
