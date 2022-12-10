#!/usr/bin/env python3
import os, sys
input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

def shift(d, a):
    return 0 if a == 0 else a // d

class knot(object):
    def __init__(self, x=0, y=0):
        self.x, self.y = x, y
    def __repr__(self):
        return f'@({self.x}, {self.y})'

    def U(self): self.y += 1
    def D(self): self.y -= 1
    def L(self): self.x -= 1
    def R(self): self.x += 1

    def follow(self, other):
        dx = other.x - self.x; adx = abs(dx)
        dy = other.y - self.y; ady = abs(dy)
        if adx > 1 or ady > 1:
            self.x += shift(dx, adx)
            self.y += shift(dy, ady)

def grid(h, t, nr=5, nc=6):
    g = bytearray(b'.' * nr * nc)
    g[t.y*nc + t.x] = ord(b'T')
    g[h.y*nc + h.x] = ord(b'H')
    lines = list(g[i:i+nc].decode('ascii') for i in range(0, len(g), nc))
    lines.reverse()
    print('\n'.join(lines))
    print()

with open(input_name, 'r') as fp:
    head, tail = knot(), knot()
    vis = set()    
    for line in fp:
        direction, count = line.strip().split()
        for i in range(int(count)):
            getattr(head, direction)()
            tail.follow(head)
            vis.add((tail.x, tail.y))

print(len(vis))
