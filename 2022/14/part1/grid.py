class grid(object):
    def __init__(self, input):
        self.parse(input)

    def parse(self, input, sx=500, sy=0):
        self._pts = {}
        self._seed = (sx, sy)
        self._pts[self._seed] = '+'
        for line in input.split('\n'):
            last = None
            for next in line.split(' -> '):
                h, v = [int(x) for x in next.split(',', 1)]
                if last is None:
                    last = (h, v)
                    continue

                hmin, hmax = min(last[0], h),  max(last[0], h)
                for i in range(hmin, hmax+1):
                    vmin, vmax = min(last[1], v), max(last[1], v)
                    for j in range(vmin, vmax+1):
                        self._pts[i,j] = '#'
                last = (h, v)

        self.vmin = min(0, min(v for _, v in self._pts))
        self.vmax = max(v for _, v in self._pts)
        self._floor = -1

    def set_floor(self, depth):
        self._floor = depth

    def is_clear(self, h, v):
        return (h, v) not in self._pts and (self._floor <= 0 or v < self._floor)

    def in_bounds(self, h, v):
        return v >= self.vmin and (
            v <= self.vmax or (self._floor >= 0 and v <= self._floor))

    def count(self, item='o'):
        return sum(1 for x in self._pts.values() if x == item)

    def drop(self, item='o'):
        h, v = self._seed
        while self.in_bounds(h, v):
            if self.is_clear(h, v+1):
                v += 1 # move straight down
            elif self.is_clear(h-1, v+1):
                h -= 1; v += 1 # move down-left
            elif self.is_clear(h+1, v+1):
                h += 1; v += 1 # move down-right
            elif (h, v) == self._seed:
                self._pts[h,v] = item
                return False
            else:
                self._pts[h,v] = item
                return True # blocked, give up
        # If we get here, the item fell off the map
        return False
    
    def __str__(self):
        hmin, hmax = minmax(h for h, _ in self._pts)
        width = hmax + 1 - hmin
        vmin, vmax = minmax(v for _, v in self._pts)
        if self._floor > vmax:
            vmax = self._floor
        height = vmax + 1 - vmin
        buf = bytearray(b'.'*width*height)
        for (h, v), p in self._pts.items():
            pos = (v - vmin)*width + (h - hmin)
            buf[pos] = ord(p)
        return b'\n'.join(buf[i:i+width]
                          for i in range(0, len(buf), width)).decode('ascii')

def minmax(vals):
    min = next(vals)
    max = min
    for elt in vals:
        if elt < min:
            min = elt
        if elt > max:
            max = elt
    return min, max
