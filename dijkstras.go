package functions

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ParseGraph(rawdata string) (map[string]map[string]int, float64) {
	startTime := time.Now()

	graph := make(map[string]map[string]int)

	lines := strings.Split(rawdata, "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		node, connectedNode, weightStr := parts[0], parts[1], parts[2]
		weight, _ := strconv.Atoi(weightStr)

		if _, exists := graph[node]; !exists {
			graph[node] = make(map[string]int)
		}

		graph[node][connectedNode] = weight
	}

	elapsedTime := time.Since(startTime)

	return graph, elapsedTime.Seconds()
}

func Dijkstra(graph map[string]map[string]int, startKey string) (distances map[string]int, runTime float64, err error) {
	startTime := time.Now()

	_, ok := graph[startKey]
	if !ok {
		return nil, 0, fmt.Errorf("start vertex %v not found", startKey)
	}

	distances = make(map[string]int)
	for key := range graph {
		distances[key] = math.MaxInt32
	}
	distances[startKey] = 0

	var vertices []string // visited vertices
	for vertex := range graph {
		vertices = append(vertices, vertex)
	}

	for len(vertices) != 0 {
		sort.SliceStable(vertices, func(i, j int) bool {
			return distances[vertices[i]] < distances[vertices[j]]
		})

		vertex := vertices[0]
		vertices = vertices[1:]

		for adjacent, cost := range graph[vertex] {
			alt := distances[vertex] + cost
			if alt < distances[adjacent] {
				distances[adjacent] = alt
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
