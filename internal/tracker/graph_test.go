package tracker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildGraph(t *testing.T) {
	t.Run("acyclic graph: success", func(t *testing.T) {
		const (
			expectedVerticesLen = 5
		)

		edges := [][]string{{"A", "B"}, {"B", "C"}, {"C", "E"}, {"C", "D"}, {"D", "A"}}

		var (
			expectedEdgesForA = map[string]bool{"B": true}
			expectedEdgesForB = map[string]bool{"C": true}
			expectedEdgesForC = map[string]bool{"E": true, "D": true}
			expectedEdgesForD = map[string]bool{"A": true}
			expectedEdgesForE = make(map[string]bool)
		)

		graph, err := BuildGraph(edges)
		assert.Empty(t, err)

		assert.Equal(t, expectedVerticesLen, len(graph.vertices))
		for key, vertex := range graph.vertices {
			var metAllEdges bool
			switch key {
			case "A":
				metAllEdges = checkEdges(vertex, expectedEdgesForA)
			case "B":
				metAllEdges = checkEdges(vertex, expectedEdgesForB)
			case "C":
				metAllEdges = checkEdges(vertex, expectedEdgesForC)
			case "D":
				metAllEdges = checkEdges(vertex, expectedEdgesForD)
			case "E":
				metAllEdges = checkEdges(vertex, expectedEdgesForE)
			}
			assert.True(t, metAllEdges)
		}
	})

	t.Run("invalid edge: same vertex ", func(t *testing.T) {
		// E -> E (departure and arrival airport should be different)
		edges := [][]string{{"B", "C"}, {"A", "B"}, {"E", "E"}, {"C", "E"}}

		_, err := BuildGraph(edges)
		assert.Error(t, ErrInvalidEdge, err)
	})

	t.Run("invalid edge: empty vertex ", func(t *testing.T) {
		// "" -> E (departure and arrival airport can't be empty)
		edges := [][]string{{"B", "C"}, {"", "B"}, {"E", "E"}, {"C", "E"}}

		_, err := BuildGraph(edges)
		assert.Error(t, ErrInvalidEdge, err)
	})

	t.Run("invalid edge: 3 vertex ", func(t *testing.T) {
		// you can have only 1 departure and 1 arrival airport
		edges := [][]string{{"B", "C"}, {"C", "B", "D"}, {"D", "E"}, {"C", "E"}}

		_, err := BuildGraph(edges)
		assert.Error(t, ErrInvalidEdge, err)
	})
}

func checkEdges(vertex *Vertex, expectedEdges map[string]bool) bool {
	for val := range expectedEdges {
		_, ok := vertex.edges[val]
		if !ok {
			return false
		}
	}
	return true
}

func TestFindStartAndEndPoint(t *testing.T) {
	t.Run("acyclic short graph: start and end point found ", func(t *testing.T) {
		const (
			expectedStart = "C"
			expectedEnd   = "D"
		)

		edges := [][]string{{"C", "D"}}

		graph, err := BuildGraph(edges)
		assert.Empty(t, err)

		start, end, err := graph.FindStartAndEndPoint()
		assert.Empty(t, err)
		assert.Equal(t, expectedStart, start)
		assert.Equal(t, expectedEnd, end)
	})

	t.Run("acyclic graph: start and end point found ", func(t *testing.T) {
		const (
			expectedStart = "A"
			expectedEnd   = "D"
		)

		// A -> B -> C -> E -> D
		edges := [][]string{{"E", "D"}, {"B", "C"}, {"A", "B"}, {"C", "E"}}

		graph, err := BuildGraph(edges)
		assert.Empty(t, err)

		start, end, err := graph.FindStartAndEndPoint()
		assert.Empty(t, err)
		assert.Equal(t, expectedStart, start)
		assert.Equal(t, expectedEnd, end)
	})

	t.Run("cyclic graph: start and end point found ", func(t *testing.T) {
		const (
			expectedStart = "A"
			expectedEnd   = "D"
		)

		// A -> B -> C -> B -> D
		edges := [][]string{{"B", "D"}, {"A", "B"}, {"B", "C"}, {"B", "D"}, {"C", "B"}}

		graph, err := BuildGraph(edges)
		assert.Empty(t, err)

		start, end, err := graph.FindStartAndEndPoint()
		assert.Empty(t, err)
		assert.Equal(t, expectedStart, start)
		assert.Equal(t, expectedEnd, end)
	})

	t.Run("cyclic graph: start and end point NOT found ", func(t *testing.T) {

		// A -> B -> C -> E -> A  == E -> A -> B -> C -> E
		edges := [][]string{{"E", "A"}, {"B", "C"}, {"A", "B"}, {"C", "E"}}

		graph, err := BuildGraph(edges)
		assert.Empty(t, err)

		_, _, err = graph.FindStartAndEndPoint()
		assert.Error(t, ErrUnableFindStartAndEndPoint, err)
	})

	t.Run("not connected graph: start and end point NOT found ", func(t *testing.T) {

		// A -> B -> C,  E -> D
		edges := [][]string{{"E", "D"}, {"B", "C"}, {"A", "B"}}

		graph, err := BuildGraph(edges)
		assert.Empty(t, err)

		_, _, err = graph.FindStartAndEndPoint()
		assert.Error(t, ErrGraphIsUnconnected, err)
	})
}
