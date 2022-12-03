#!/usr/bin/env python3

import sys

input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

def prio(x):
    return '_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'.find(x)

def common(xs):
    s = set(xs[0])
    for x in xs[1:]:
        s &= set(x)
    return s.pop()

with open(input_name, 'r') as fp:
    lines = fp.read().strip().split('\n')
    total = sum(prio(common(group)) for group in
                (lines[i:i+3] for i in range(0, len(lines), 3)))

print(total)
