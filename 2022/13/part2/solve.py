#!/usr/bin/env python3
import json, os, sys
from cmp import compare
from functools import cmp_to_key
input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

with open(input_name, 'r') as fp:
    lines = fp.read().strip().replace('\n\n', '\n').split('\n')

div1 = [[2]]
div2 = [[6]]

packets = list(json.loads(x) for x in lines)
packets.extend([div1, div2])
packets.sort(key=cmp_to_key(compare))
pos1 = packets.index(div1) + 1
pos2 = packets.index(div2) + 1
print(pos1*pos2)

