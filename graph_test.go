// © 2025 Rolf van de Krol <rolf@vandekrol.xyz>

package graph

import (
	"reflect"
	"testing"
)

func TestSimpleSort(t *testing.T) {
	g := New[string]()

	// We construct a graph with the following structure: a -> b -> c.
	// The expected topological sort is [a b c]. We explicitly add the nodes in
	// another order to ensure that the sort is not just returning the input.
	g.Node("c")
	g.Node("b")
	g.Node("a")

	g.Edge("a", "b")
	g.Edge("b", "c")

	keys, err := g.Sort()

	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(keys, []string{"a", "b", "c"}) {
		t.Errorf("expected [a b c], got %v", keys)
	}
}

func TestCyclicSort(t *testing.T) {
	g := New[string]()

	// We construct a graph with the following structure: a -> b -> c -> a.
	// The graph has a cycle, so the sort should return an error.
	g.Node("c")
	g.Node("b")
	g.Node("a")

	g.Edge("a", "b")
	g.Edge("b", "c")
	g.Edge("c", "a")

	_, err := g.Sort()

	if err == nil {
		t.Error("expected error")
		return
	}
}

func TestComplexSort(t *testing.T) {
	g := New[string]()

	// We construct a graph with the following structure: a -> b, c -> b.
	// The expected topological sort is [c a b]. We explicitly add the nodes in
	// another order to ensure that the sort is not just returning the input.
	g.Node("c")
	g.Node("b")
	g.Node("a")

	g.Edge("a", "b")
	g.Edge("c", "b")

	keys, err := g.Sort()

	if err != nil {
		t.Error(err)
		return
	}

	if len(keys) != 3 {
		t.Errorf("expected 3, got %v", len(keys))
		return
	}

	if !(reflect.DeepEqual(keys, []string{"a", "c", "b"}) || reflect.DeepEqual(keys, []string{"c", "a", "b"})) {
		t.Errorf("expected [a c b] or [c a b], got %v", keys)
	}
}

// buildGraph creates a graph with a specific structure. It is used to benchmark
// the topological sort algorithm.
// It generates a directed graph with n int nodes, and 2n edges. The graph is
// constructed in such a way that it has a topological sort.
func buildGraph(n int) *Graph[int] {
	g := New[int]()
	for i := 0; i < n; i++ {
		g.Node(i)

		x := i % 5
		if x == 0 {
			continue
		}

		y := i / x
		z := i % y

		for j := 0; j < x; j++ {
			g.Edge(i, j*y+z)
		}
	}
	return g
}

// TestBuildGraph tests the buildGraph function. It ensures that the generated
// graph has a topological sort.
func TestBuildGraph(t *testing.T) {
	g := buildGraph(1000)

	_, err := g.Sort()
	if err != nil {
		t.Error(err)
		return
	}
}

func BenchmarkSort5(b *testing.B) {
	g := buildGraph(5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.Sort()
	}
}

func BenchmarkSort10(b *testing.B) {
	g := buildGraph(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.Sort()
	}
}

func BenchmarkSort100(b *testing.B) {
	g := buildGraph(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.Sort()
	}
}

func BenchmarkSort1000(b *testing.B) {
	g := buildGraph(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.Sort()
	}
}

func BenchmarkSort10000(b *testing.B) {
	g := buildGraph(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.Sort()
	}
}

func BenchmarkSort100000(b *testing.B) {
	g := buildGraph(100000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.Sort()
	}
}
