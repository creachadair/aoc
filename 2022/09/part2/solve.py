#!/usr/bin/env python3
import os, sys
from grid import grid
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

class rope(object):
    def __init__(self, knots):
        self.knots = knots

    def U(self): self.knots[0].U(); self.pull()
    def D(self): self.knots[0].D(); self.pull()
    def L(self): self.knots[0].L(); self.pull()
    def R(self): self.knots[0].R(); self.pull()

    def tail(self): return self.knots[-1]
    def pull(self):
        for i in range(1, len(self.knots)):
            self.knots[i].follow(self.knots[i-1])

with open(input_name, 'r') as fp:
    rope = rope([knot() for _ in range(10)])
    vis = set()
    for line in fp:
        direction, count = line.strip().split()
        for i in range(int(count)):
            getattr(rope, direction)()
            tail = rope.tail()
            vis.add((tail.x, tail.y))

grid(rope.knots, vis)
print(len(vis))
