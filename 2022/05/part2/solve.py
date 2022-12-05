#!/usr/bin/env python3

import sys

input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

def parse_stacks(spec):
    lines = spec.split('\n')
    count = len(lines.pop().strip().split())
    stacks = [list() for _ in range(count)]
    for line in lines:
        n = 0
        for i in range(0, len(line), 4):
            elt = line[i:i+3]
            if elt.startswith('['):
                stacks[n].append(elt[1:-1])
            n += 1
    return stacks

def parse_moves(spec):
    moves = []
    for line in spec.split('\n'):
        rule = line.split()
        #            move <>        from <>         to <>
        moves.append((int(rule[1]), int(rule[3])-1, int(rule[5])-1))
    return moves

def word(stacks):
    return ''.join(s[0] if len(s) > 0 else '' for s in stacks)

with open(input_name, 'r') as fp:
    input = fp.read().rstrip()
    stack_spec, moves_spec = input.split('\n\n', 1)
    stacks = parse_stacks(stack_spec)

    for n, src, dst in parse_moves(moves_spec):
        top, stacks[src] = stacks[src][:n], stacks[src][n:]
        stacks[dst] = top + stacks[dst]

print(word(stacks))
