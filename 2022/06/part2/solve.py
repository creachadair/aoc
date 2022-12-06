#!/usr/bin/env python3

import sys

input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

def find_start(line, n=4):
    for i in range(len(line)):
        if len(set(line[i:i+n])) == n:
            return i+n  # +n for split at end of range
    raise IndexError("not found")

with open(input_name, 'r') as fp:
    for line in fp:
        print(find_start(line.strip(), 14))
