#!/usr/bin/env python3
import os, sys
input_name = sys.argv[1] if len(sys.argv) > 1 else 'input.txt'

def parse(lines):
    cwd = []
    def path(*rest):
        return '/' + '/'.join(cwd+list(rest))

    smap = {}  # path -> size, directs only
    dmap = {}  # path -> [dir], directs only
    for line in lines:
        tag, name, *rest = line.split()
        if tag == '$':
            if name == 'cd':
                if rest[0] == '/':
                    cwd = []
                elif rest[0] == '..':
                    cwd.pop()
                else:
                    cwd.append(rest[0])
                dmap.setdefault(path(), [])
                smap.setdefault(path(), 0)
            elif name != 'ls':
                raise KeyError(name)
        elif tag == 'dir':
            dmap[path()].append(path(name))
        else:
            size = int(tag)
            smap[path()] += size

    class tree(object):
        def sizeof(self, path):
            size = smap[path]
            for sub in dmap[path]:
                size += self.sizeof(os.path.join(path, sub))
            return size
        def dirs(self):
            return smap.keys()

    return tree()


with open(input_name, 'r') as fp:
    tree = parse(fp.read().strip().split('\n'))

size = sum(v for v in (tree.sizeof(path) for path in tree.dirs()) if v <= 100000)
print(size)
