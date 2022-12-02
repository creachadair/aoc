#!/usr/bin/env python3

import sys

input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

top3 = [0] * 3  # nonincreasing

def insert(v, into=top3):
    for i, old in enumerate(into):
        if v > old:
            into[i] = v
            v = old

with open(input_name, 'r') as fp:
    for block in fp.read().strip().split('\n\n'):
        insert(sum(int(v) for v in block.strip().split('\n')))

print(sum(top3))
