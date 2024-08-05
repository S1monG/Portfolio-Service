package functions

import (
	"testing"
)

func TestParseGraph(t *testing.T) {
	rawdata := "A B 5\nA C 10\nB A 5\nB C 2\nC A 10\nC B 2\nC D 1\nD C 1"

	expectedGraph := map[string]map[string]int{
		"A": {"B": 5, "C": 10},
		"B": {"A": 5, "C": 2},
		"C": {"A": 10, "B": 2, "D": 1},
		"D": {"C": 1},
	}

	graph, elapsedTime := ParseGraph(rawdata)

	if len(graph) != len(expectedGraph) {
		t.Errorf("Expected graph length %d, got %d", len(expectedGraph), len(graph))
	}

	for node, edges := range expectedGraph {
		if len(graph[node]) != len(edges) {
			t.Errorf("Expected edges length for node %s: %d, got %d", node, len(edges), len(graph[node]))
		}
		for connectedNode, weight := range edges {
			if graph[node][connectedNode] != weight {
				t.Errorf("Expected weight for edge %s -> %s: %d, got %d", node, connectedNode, weight, graph[node][connectedNode])
			}
		}
	}

	if elapsedTime < 0 {
		t.Errorf("Expected elapsed time to be greater than or equal to 0, got %f", elapsedTime)
	}
}

func TestDijkstras(t *testing.T) {
	graph := map[string]map[string]int{
		"A": {"B": 5, "C": 10},
		"B": {"A": 5, "C": 2},
		"C": {"A": 10, "B": 2, "D": 1},
		"D": {"C": 1},
	}
	startNode := "A"
	expectedDistances := map[string]int{
		"A": 0,
		"B": 5,
		"C": 7,
		"D": 8,
	}

	distances, elapsedTime, err := Dijkstra(graph, startNode)

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
