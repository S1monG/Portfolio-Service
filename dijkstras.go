package functions

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Vertex struct {
	Key   string
	Edges map[*Vertex]int
}

type Graph struct {
	Vertices map[string]*Vertex
}

func ParseGraph(rawdata string) (*Graph, float64) {
	startTime := time.Now()

	graph := &Graph{Vertices: make(map[string]*Vertex)}

	lines := strings.Split(rawdata, "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		node, connectedNode, weightStr := parts[0], parts[1], parts[2]
		weight, _ := strconv.Atoi(weightStr)

		if _, exists := graph.Vertices[node]; !exists {
			graph.Vertices[node] = &Vertex{Key: node, Edges: make(map[*Vertex]int)}
		}
		if _, exists := graph.Vertices[connectedNode]; !exists {
			graph.Vertices[connectedNode] = &Vertex{Key: connectedNode, Edges: make(map[*Vertex]int)}
		}

		graph.Vertices[node].Edges[graph.Vertices[connectedNode]] = weight
	}

	elapsedTime := time.Since(startTime)

	return graph, elapsedTime.Seconds()
}

func (g *Graph) Dijkstra(startKey string) (distances map[string]int, runTime float64, err error) {
	startTime := time.Now()

	_, ok := g.Vertices[startKey]
	if !ok {
		return nil, 0, fmt.Errorf("start vertex %v not found", startKey)
	}

	distances = make(map[string]int)
	for key := range g.Vertices {
		distances[key] = math.MaxInt32
	}
	distances[startKey] = 0

	var vertices []*Vertex
	for _, vertex := range g.Vertices {
		vertices = append(vertices, vertex)
	}

	for len(vertices) != 0 {
		sort.SliceStable(vertices, func(i, j int) bool {
			return distances[vertices[i].Key] < distances[vertices[j].Key]
		})

		vertex := vertices[0]
		vertices = vertices[1:]

		for adjacent, cost := range vertex.Edges {
			alt := distances[vertex.Key] + cost
			if alt < distances[adjacent.Key] {
				distances[adjacent.Key] = alt
			}
		}
	}

	runTime = time.Since(startTime).Seconds()

	return distances, runTime, nil
}

/*
func dijkstras(graph map[string]map[string]int, startNode string) (map[string]int, float64) {
	startTime := time.Now()

	// Initialize the distance map
	distances := make(map[string]int)
	nodes := make([]string, 0, len(graph))

	for node := range graph {
		distances[node] = -1
		nodes = append(nodes, node)
	}

	distances[startNode] = 0

	visited := make(map[string]bool)

	for len(nodes) > 0 {
		smallestDistance, smallestNodeIdx := -1, 0
		for index, node := range nodes {
			if smallestDistance == -1 || (distances[node] != -1 && distances[node] < smallestDistance) {
				smallestDistance = distances[node]
				smallestNodeIdx = index
			}
		}

		closestNode := nodes[smallestNodeIdx]
		nodes = append(nodes[:smallestNodeIdx], nodes[smallestNodeIdx+1:]...)

		if distances[closestNode] == -1 {
			break
		}

		visited[closestNode] = true

		for neighbor := range graph[closestNode] {
			if exists := visited[neighbor]; !exists {
				newDistance := distances[closestNode] + graph[closestNode][neighbor]

				if newDistance < distances[closestNode] {
					distances[neighbor] = newDistance
				}
			}
		}
	}

	elapsedTime := time.Since(startTime)

	return distances, elapsedTime.Seconds()
}
*/
