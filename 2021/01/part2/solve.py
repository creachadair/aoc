#!/usr/bin/env python3
import sys
input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

with open(input_name, 'r') as fp:
    samples = list(int(x.strip()) for x in fp)
    prev, inc = 0, -1
    for i in range(len(samples)-2):
        next = sum(samples[i:i+3])
        if next > prev:
            inc += 1
        prev = next

print(inc)
