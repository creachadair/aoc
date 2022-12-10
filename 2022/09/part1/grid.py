#!/usr/bin/env python3

def grid(knots, points=()):
    minx = min(min(kt.x for kt in knots), min(x for x, _ in points))
    maxx = max(max(kt.x for kt in knots), max(x for x, _ in points))
    miny = min(min(kt.y for kt in knots), min(y for _, y in points))
    maxy = max(max(kt.y for kt in knots), max(y for _, y in points))
    nr = (maxy - miny) + 1
    nc = (maxx - minx) + 1

    g = bytearray(b'.' * nr * nc)
    def setxy(x, y, ch):
        rx, ry = x-minx, y-miny
        g[ry*nc + rx] = ch

    for i, kt in enumerate(knots):
        setxy(kt.x, kt.y, ord(b'H') if i == 0 else ord(b'0')+i)
    for x, y in points:
        setxy(x, y, ord(b'#'))

    lines = list(g[i:i+nc].decode('ascii') for i in range(0, len(g), nc))
    lines.reverse()
    print('\n'.join(lines))

