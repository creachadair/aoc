#!/usr/bin/env python3
import os, sys
input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

def fcount(samples, pos):
    count = sum(int(s[pos] == '1') for s in samples)
    half = len(samples) / 2
    return int(count >= half)

with open(input_name, 'r') as fp:
    samples = fp.read().strip().split('\n')

def unique(samples, most_common):
    i = 0
    while len(samples) > 1:
        mc = fcount(samples, i)
        want = str(mc if most_common else 1 - mc)
        samples = [s for s in samples if s[i] == want]
        i += 1
    return samples[0]

oxy = int(unique(samples, most_common=True), 2)
co2 = int(unique(samples, most_common=False), 2)
print(f'oxy: {oxy}, co2: {co2}, life: {oxy*co2}')
