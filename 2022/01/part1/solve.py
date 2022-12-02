#!/usr/bin/env python3

import sys

input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

with open(input_name, 'r') as fp:
    max_cal = max(
        sum(int(v) for v in block.strip().split('\n'))
        for block in fp.read().strip().split('\n\n')
    )

print(max_cal)
