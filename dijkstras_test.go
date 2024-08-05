package functions

import (
	"testing"
)

func TestParseGraph(t *testing.T) {
	rawdata := "A B 5\nA C 10\nB A 5\nB C 2\nC A 10\nC B 2\nC D 1\nD C 1"

	vA := &Vertex{Key: "A", Edges: make(map[*Vertex]int)}
	vB := &Vertex{Key: "B", Edges: make(map[*Vertex]int)}
	vC := &Vertex{Key: "C", Edges: make(map[*Vertex]int)}
	vD := &Vertex{Key: "D", Edges: make(map[*Vertex]int)}

	// Create edges
	vA.Edges[vB] = 5
	vA.Edges[vC] = 10
	vB.Edges[vA] = 5
	vB.Edges[vC] = 2
	vC.Edges[vA] = 10
	vC.Edges[vB] = 2
	vC.Edges[vD] = 1
	vD.Edges[vC] = 1

	// Create graph
	expectedGraph := &Graph{
		Vertices: map[string]*Vertex{
			"A": vA,
			"B": vB,
			"C": vC,
			"D": vD,
		},
	}

	graph, elapsedTime := ParseGraph(rawdata)

	if len(graph.Vertices) != len(expectedGraph.Vertices) {
		t.Errorf("Expected graph length %d, got %d", len(expectedGraph.Vertices), len(graph.Vertices))
	}

	for nodeKey, expectedVertex := range expectedGraph.Vertices {
		vertex, exists := graph.Vertices[nodeKey]
		if !exists {
			t.Errorf("Expected vertex %s to exist", nodeKey)
			continue
		}
		if len(vertex.Edges) != len(expectedVertex.Edges) {
			t.Errorf("Expected edges length for node %s: %d, got %d", nodeKey, len(expectedVertex.Edges), len(vertex.Edges))
		}
		for connectedVertex, expectedWeight := range expectedVertex.Edges {
			if weight, exists := vertex.Edges[connectedVertex]; !exists || weight != expectedWeight {
				t.Errorf("Expected weight for edge %s -> %s: %d, got %d", nodeKey, connectedVertex.Key, expectedWeight, weight)
			}
		}
	}

	if elapsedTime < 0 {
		t.Errorf("Expected elapsed time to be greater than or equal to 0, got %f", elapsedTime)
	}
}

func TestDijkstras(t *testing.T) {
	vA := &Vertex{Key: "A", Edges: make(map[*Vertex]int)}
	vB := &Vertex{Key: "B", Edges: make(map[*Vertex]int)}
	vC := &Vertex{Key: "C", Edges: make(map[*Vertex]int)}
	vD := &Vertex{Key: "D", Edges: make(map[*Vertex]int)}

	// Create edges
	vA.Edges[vB] = 5
	vA.Edges[vC] = 10
	vB.Edges[vA] = 5
	vB.Edges[vC] = 2
	vC.Edges[vA] = 10
	vC.Edges[vB] = 2
	vC.Edges[vD] = 1
	vD.Edges[vC] = 1

	// Create graph
	graph := &Graph{
		Vertices: map[string]*Vertex{
			"A": vA,
			"B": vB,
			"C": vC,
			"D": vD,
		},
	}
	startNode := "A"
	expectedDistances := map[string]int{
		"A": 0,
		"B": 5,
		"C": 7,
		"D": 8,
	}

	distances, elapsedTime, err := graph.Dijkstra(startNode)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(distances) != len(expectedDistances) {
		t.Errorf("Expected distances length %d, got %d", len(expectedDistances), len(distances))
	}

	for node, distance := range expectedDistances {
		if distances[node] != distance {
			t.Errorf("Expected distance for node %s: %d, got %d", node, distance, distances[node])
		}
	}

	if elapsedTime < 0 {
		t.Errorf("Expected elapsed time to be greater than or equal to 0, got %f", elapsedTime)
	}
}
