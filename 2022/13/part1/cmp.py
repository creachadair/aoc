#!/usr/bin/env python3

def compare(a, b):
    ta, tb = type(a), type(b)
    if ta == tb == list:
        i = 0
        while i < len(a) and i < len(b):
            cur = compare(a[i], b[i])
            if cur != 0:
                return cur
            i += 1
        return zcmp(len(a), len(b))
    elif ta == tb == int:
        return zcmp(a, b)
    elif ta == int:
        return compare([a], b)
    else:
        return compare(a, [b])
        
def zcmp(a, b):
    d = a - b
    return d if d == 0 else d // abs(d)

__export__ = ('compare', 'zcmp')
