// © 2025 Rolf van de Krol <rolf@vandekrol.xyz>

// Package graph provides a directed graph data structure and topological sort
// algorithm (O(n) for n = [number of nodes] + [number of edges], using Kahn's
// algorithm).
//
//	g := graph.New[string]()
//	g.Add("b", []string{"c"})
//	g.Add("c", []string{"a"})
//
//	s, err := g.Sort() // []string{"b", "c", "a"}
package graph

import "fmt"

// Graph represents a directed graph.
type Graph[Key comparable] struct {
	nodes map[Key]Edges[Key]
}

// Edges represents the edges of a node in a directed graph.
type Edges[Key comparable] map[Key]bool

// New returns a new graph.
func New[Key comparable]() *Graph[Key] {
	return &Graph[Key]{
		nodes: make(map[Key]Edges[Key]),
	}
}

// Node returns the edges for a node. It creates the node if it does not exist.
func (g *Graph[Key]) Node(key Key) Edges[Key] {
	n, ok := g.nodes[key]
	if !ok {
		n = make(Edges[Key])
		g.nodes[key] = n
	}
	return n
}

// Edge adds an edge to the graph. It creates the nodes if they do not exist.
func (g *Graph[Key]) Edge(from Key, to Key) {
	f := g.Node(from)
	g.Node(to)
	f.add(to)
}

// Add adds a node and its outgoing edges to the graph.
func (g *Graph[Key]) Add(node Key, edges []Key) {
	n := g.Node(node)
	for _, e := range edges {
		g.Node(e)
		n.add(e)
	}
}

// Reverse returns a new graph with all edges reversed.
func (g *Graph[Key]) Reverse() *Graph[Key] {
	r := New[Key]()

	for from, e := range g.nodes {
		r.Node(from)
		for to := range e {
			r.Edge(to, from)
		}
	}

	return r
}

// Copy returns a new graph with the same nodes and edges.
func (g *Graph[Key]) Copy() *Graph[Key] {
	c := New[Key]()

	for from, e := range g.nodes {
		c.Node(from)
		for to := range e {
			c.Edge(from, to)
		}
	}

	return c
}

// Sort returns a topological sorted list of the graph nodes. It returns an
// error if the graph has a cycle. It is an implementation of Kahn's algorithm.
// Sort's time complexity is O(n) for n = [number of nodes] + [number of edges].
func (g *Graph[Key]) Sort() ([]Key, error) {
	// https://en.wikipedia.org/wiki/Topological_sorting#Kahn's_algorithm

	// We need to make a copy of the graph, so we can modify it.
	gg := g.Copy()

	// The original graph's edges are actually outgoing edges. We need to
	// reverse the graph to detect nodes with no incoming edges in an efficient
	// way. Without reversing the graph, we would need to iterate over all
	// nodes and their edges to find nodes with no incoming edges.
	r := g.Reverse()

	// The sorted list of keys, which we will return
	var sorted []Key

	// The list of keys with no incoming edges. We need this to start the
	// algorithm. We construct it using the reversed graph and finding nodes
	// with no outgoing edges.
	var next []Key
	for k, e := range r.nodes {
		if len(e) == 0 {
			next = append(next, k)
		}
	}

	// We iterate over the list of nodes with no incoming edges. This list will
	// be empty when the graph is empty or when the graph has a cycle.
	for len(next) > 0 {
		n := next[0]
		next = next[1:]

		// We add the node n to the sorted list.
		sorted = append(sorted, n)

		// We iterate over the nodes that are connected to the current node n.
		// We only consider outgoing edges, because the node we are visiting
		// has no incoming edges.
		for m := range gg.nodes[n] {
			// We remove the edge from n to m from the graph.
			delete(gg.nodes[n], m)
			delete(r.nodes[m], n)

			// If the node m has no incoming edges left after we removed the
			// edge from n to m, we add it to the list of nodes with no
			// incoming edges, so we can consider it in the next iteration.
			if len(r.nodes[m]) == 0 {
				next = append(next, m)
			}
		}

		// We remove the node n from the graph. This is necessary to detect
		// cycles.
		delete(gg.nodes, n)
		delete(r.nodes, n)
	}

	// If the graph is not empty, it means that there is a cycle in the graph.
	if len(gg.nodes) > 0 {
		return nil, fmt.Errorf("cycle detected")
	}

	return sorted, nil
}

// add adds a key to the edges.
func (n Edges[Key]) add(key Key) {
	n[key] = true
}
