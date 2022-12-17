package graph

// G represents a graph in adjacency list format.
type G[T comparable] struct {
	nodes map[T]int
	edges map[pair[T]]bool

	dist map[pair[T]]int
	next map[pair[T]]T
}

// New constructs a new empty graph with node type T.
func New[T comparable]() *G[T] {
	return &G[T]{nodes: make(map[T]int), edges: make(map[pair[T]]bool)}
}

// SetNode adds a node with the given label and weight, or updates the weight
// of an existing node.
func (g *G[T]) SetNode(name T, weight int) {
	if _, ok := g.nodes[name]; !ok {
		g.inval()
	}
	g.nodes[name] = weight
}

// Node reports whether a node with the given name exists, and if so returns
// its weight.
func (g *G[T]) Node(name T) (int, bool) { w, ok := g.nodes[name]; return w, ok }

// Nodes returns a slice of the names of all the nodes in g.  The order of the
// results is unspecified and may change.
func (g *G[T]) Nodes() []T {
	names := make([]T, 0, len(g.nodes))
	for name := range g.nodes {
		names = append(names, name)
	}
	return names
}

// Edge adds an edge from a to b in g. If either a or b does not exist in the
// graph, a node with weight zero is added for it.
func (g *G[T]) Edge(a, b T) {
	g.SetNode(a, g.nodes[a])
	g.SetNode(b, g.nodes[b])
	if e := (pair[T]{a, b}); !g.edges[e] {
		g.edges[e] = true
		g.inval()
	}
}

// Distance reports whether there is a path from a to b in g, and if so returns
// the length of the shortest such path.
func (g *G[T]) Distance(a, b T) (int, bool) {
	g.initPaths()
	n, ok := g.dist[pair[T]{a, b}]
	return n, ok
}

// Path returns a shortest path from a to b in g, or nil if no path exist.
func (g *G[T]) Path(a, b T) []T {
	g.initPaths()
	var path []T
	for a != b {
		next, ok := g.next[pair[T]{a, b}]
		if !ok {
			return nil // no path
		}
		path = append(path, a)
		a = next
	}
	return append(path, b)
}

func (g *G[T]) inval() { g.dist = nil; g.next = nil }

type pair[T comparable] struct{ a, b T }

func (g *G[T]) initPaths() {
	if g.dist != nil && g.next != nil {
		return
	}

	// Set up a Floyd-Warshall all-pairs shortest path tables.
	g.dist = make(map[pair[T]]int) // min distance from a to b
	g.next = make(map[pair[T]]T)   // next step on shortest path from a to b
	for uv := range g.edges {
		g.dist[uv] = 1
		g.next[uv] = uv.b
	}
	for a := range g.nodes { // N.B. run after edges, so self-loops don't clobber
		g.dist[pair[T]{a, a}] = 0
		g.next[pair[T]{a, a}] = a
	}
	for r := range g.nodes {
		for u := range g.nodes {
			for v := range g.nodes {
				ur, ok := g.dist[pair[T]{u, r}]
				if !ok {
					continue
				}
				rv, ok := g.dist[pair[T]{r, v}]
				if !ok {
					continue
				}
				uv, ok := g.dist[pair[T]{u, v}]
				if !ok || uv > ur+rv {
					g.dist[pair[T]{u, v}] = ur + rv
					g.next[pair[T]{u, v}] = g.next[pair[T]{u, r}]
				}
			}
		}
	}
}
