#!/usr/bin/env python3
import os, sys
from grid import grid
input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

with open(input_name, 'r') as fp:
    input = grid(fp.read().strip())

input.set_floor(input.vmax + 2)

while input.drop():
    pass

print(input)
print(input.count())
