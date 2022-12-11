#!/usr/bin/env python3
import os, sys
input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

def fcount(samples, pos):
    count = sum(int(s[pos] == '1') for s in samples)
    half = len(samples) / 2
    return int(count >= half)

with open(input_name, 'r') as fp:
    samples = fp.read().strip().split('\n')

most_common = [fcount(samples, i) for i in range(len(samples[0]))]
gamma = int(''.join(str(x) for x in most_common), 2)
epsilon = int(''.join(str(1 - x) for x in most_common), 2)

print(f'gamma: {gamma}, epsilon: {epsilon}, power: {gamma*epsilon}')
