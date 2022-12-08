#!/usr/bin/env python3
import sys
input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

with open(input_name, 'r') as fp:
    prev, inc = 0, -1
    for line in fp:
        next = int(line.strip())
        if next > prev:
            inc += 1
        prev = next

print(inc)
