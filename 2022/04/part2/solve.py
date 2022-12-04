#!/usr/bin/env python3

import sys

input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

def parse(line):
    def parse_int(v):
        a, b = v.split('-', 1)
        return (int(a), int(b))

    lhs, rhs = line.split(',', 1)
    return parse_int(lhs), parse_int(rhs)

def laps(a, b):
    return not ((a[0] <= b[0] and a[1] < b[0]) or b[1] < a[0])

with open(input_name, 'r') as fp:
    pairs = (parse(line) for line in fp.read().strip().split('\n'))
    nlaps = sum(1 for a, b in pairs if laps(a, b))

print(nlaps)
