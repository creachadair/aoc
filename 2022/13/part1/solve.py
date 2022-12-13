#!/usr/bin/env python3
import json, os, sys
from cmp import *
input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

with open(input_name, 'r') as fp:
    pairs = list(chunk.split('\n')
                 for chunk in fp.read().strip().split('\n\n'))
    packets = list((json.loads(a), json.loads(b))
                   for a, b in pairs)

print(sum(i+1 for i, (a, b) in enumerate(packets)
          if compare(a, b) < 0))
    
