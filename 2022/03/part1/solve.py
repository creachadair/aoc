#!/usr/bin/env python3

import sys

input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

def prio(x):
    return '_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'.find(x)

def common(x, y):
    return (set(x) & set(y)).pop()

with open(input_name, 'r') as fp:
    rucks = [(p[:len(p)//2], p[len(p)//2:])
             for p in fp.read().strip().split('\n')]
    total = sum(prio(common(lhs, rhs)) for lhs, rhs in rucks)

print(total)
