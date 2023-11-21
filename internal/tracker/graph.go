package tracker

// Graph represents passenger's path, where each vertex is airport code and edge is flight
type Graph struct {
	vertices map[string]*Vertex
}

func newGraph() *Graph {
	return &Graph{vertices: make(map[string]*Vertex)}
}

type Vertex struct {
	val   string           // airport code
	edges map[string]*Edge // flights
}

type Edge struct {
	vertex *Vertex
}

func (g *Graph) addVertex(val string) {
	if _, ok := g.vertices[val]; !ok {
		g.vertices[val] = &Vertex{val: val, edges: make(map[string]*Edge)}
	}
}

func (g *Graph) addEdge(srcKey, destKey string) {
	// check if src & dest exist
	if _, ok := g.vertices[srcKey]; !ok {
		return
	}

	if _, ok := g.vertices[destKey]; !ok {
		return
	}

	// add edge src --> dest
	g.vertices[srcKey].edges[destKey] = &Edge{vertex: g.vertices[destKey]}
}

// Since we are dealing with flights, the edges (flights) must meet the conditions:
// 1. airport code can't be empty
// 2. departure and arrival airport can't be same
// 3. you have only 1 departure and 1 arrival airport
func invalidEdge(edge []string) bool {
	if len(edge) != 2 {
		return true
	}

	if edge[0] == "" || edge[1] == "" || edge[0] == edge[1] {
		return true
	}

	return false
}

// BuildGraph builds graph using slice of edges
func BuildGraph(edgesPairs [][]string) (*Graph, error) {
	graph := newGraph()

	for _, edge := range edgesPairs {
		if invalidEdge(edge) {
			return nil, ErrInvalidEdge
		}

		graph.addVertex(edge[0])
		graph.addVertex(edge[1])
		graph.addEdge(edge[0], edge[1])
	}
	return graph, nil
}

// dfs tries to visit every vertex, passing over the same edge only once, and counts incoming and
// outgoing edges for each.
func (g *Graph) dfs(
	currentVertex string,
	visitedVertices map[string]bool,
	visitedEdges map[*Edge]bool,
	outEdges map[string]int,
	inEdges map[string]int,
) {
	visitedVertices[currentVertex] = true
	for key, edge := range g.vertices[currentVertex].edges {
		if visitedEdges[edge] {
			continue
		}

		visitedEdges[edge] = true
		// if there is an edge between vertices A and B,
		// then vertex A (the current vertex) has an outgoing edge,
		// and vertex B has an incoming edge
		edgesNum, ok := outEdges[currentVertex]
		if ok {
			edgesNum++
			outEdges[currentVertex] = edgesNum
		} else {
			outEdges[currentVertex] = 1
		}

		edgesNum, ok = inEdges[edge.vertex.val]
		if ok {
			edgesNum++
			inEdges[edge.vertex.val] = edgesNum
		} else {
			inEdges[edge.vertex.val] = 1
		}

		if _, found := visitedVertices[key]; !found {
			g.dfs(edge.vertex.val, visitedVertices, visitedEdges, outEdges, inEdges)
		}
	}
}

// FindStartAndEndPoint finds start and end point of the path, which is represented by graph.
// Start point is the vertex, which has 1 less incoming edges, than outgoing.
// End point is the vertex, which has 1 less out-coming, than incoming.
//
// Graph should meet conditions:
//
//  1. You can find at least 1 way to visit all vertices, passing over the same edge only once.
//
//  2. Graph doesn't represent round trip by start and end points: since flights are unordered,
//     if you have pairs like [[Paris, Milan] [Milan, Paris]],
//     otherwise you can's say, if direction was Paris-Milan-Paris or Milan-Paris-Milan.
func (g *Graph) FindStartAndEndPoint() (string, string, error) {
	visitedVertices := make(map[string]bool)
	outEdges := make(map[string]int)
	inEdges := make(map[string]int)

	allVerticesAccessible := false

	// check if exists at least 1 path to visit all vertices,
	for key, vertex := range g.vertices {
		allVerticesAccessible = true

		// starting from each vertex,
		// try to reach all vertices in 1 dfs run
		if _, found := visitedVertices[key]; !found {
			visitedVertices[key] = true

			// track visitedVertices vertices in current iteration
			curVisitedVertices := make(map[string]bool)
			curVisitedEdges := make(map[*Edge]bool)

			// clean up edge cnt - if we would find the way to visit all vertices,
			outEdges = make(map[string]int)
			inEdges = make(map[string]int)

			// we will get full info about incoming and outgoing edges
			g.dfs(vertex.val, curVisitedVertices, curVisitedEdges, outEdges, inEdges)

			for _, v := range g.vertices {
				if !curVisitedVertices[v.val] {
					allVerticesAccessible = false
					break
				}
			}

			if allVerticesAccessible {
				// visitedVertices all vertices and got all info about edges,
				// no sense to run dfsModified again
				break
			}
		}
	}

	if !allVerticesAccessible {
		// graph has a gap
		return "", "", ErrGraphIsUnconnected
	}

	var (
		startVertex         string
		endVertex           string
		probablyStartVertex string
		probablyEndVertex   string
	)

	for key, vertex := range g.vertices {
		numOutEdges, outEdgesExist := outEdges[key]
		numInEdges, inEdgesExist := inEdges[key]

		switch {
		case !outEdgesExist:
			// vertex without outgoing edges is end point 100%
			endVertex = vertex.val

		case !inEdgesExist:
			// vertex without incoming edges is start point 100%
			startVertex = vertex.val

		case numInEdges-numOutEdges == 1:
			// if vertex has 1 less incoming edge, than outgoing,
			// it can be an end point in the trip in case if trip contains cycles,
			// but start and end point aren't round trip
			probablyEndVertex = vertex.val

		case numOutEdges-numInEdges == 1:
			// if vertex has 1 less outgoing edge, than incoming,
			// it can be a start point in the trip in case if trip contains cycles,
			// but start and end point aren't round trip
			probablyStartVertex = vertex.val
		}
	}

	if startVertex == "" {
		startVertex = probablyStartVertex
	}

	if endVertex == "" {
		endVertex = probablyEndVertex
	}

	if startVertex == "" || endVertex == "" || startVertex == endVertex {
		return "", "", ErrUnableFindStartAndEndPoint
	}
	return startVertex, endVertex, nil
}
