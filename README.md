# Graph

Simple implementation of a directed graph in Go, with just a topological sort
algorithm using Kahn's algorithm. Intended to be used for dependency resolution.

## Usage

```go
g := graph.New[string]()
g.Add("b", []string{"c"})
g.Add("c", []string{"a"})

s, err := g.Sort() // []string{"b", "c", "a"}
```