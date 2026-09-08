package exercises

/*
import (
	"fmt"
	"testing"
)

type Edge struct {
	to     string
	weight int
}

type Graph map[string][]Edge


func ShortestPathSimple(graph Graph, startNode, endNode string) int {
	bestCost := -1
    return bestCost
}



func ShortestPathSimple2(graph Graph, startNode, endNode string) int {
	bestCost := -1
	nodesOnCurrentPath := map[string]bool{startNode: true}

	var explorePath func(currentNode string, currentCost int)
	explorePath = func(currentNode string, currentCost int) {
		if bestCost != -1 && currentCost >= bestCost {
			return
		}
		if currentNode == endNode {
			bestCost = currentCost
			return
		}

		for _, edge := range graph[currentNode] {
			nextNode := edge.to
			if nodesOnCurrentPath[nextNode] {
				continue
			}
			nodesOnCurrentPath[nextNode] = true
			explorePath(nextNode, currentCost+edge.weight)
			nodesOnCurrentPath[nextNode] = false
		}
	}

	explorePath(startNode, 0)
	return bestCost
}

func TestGraph(t *testing.T) {
	g := Graph{
		"A": {{"B", 2}, {"C", 4}},
		"B": {{"C", 1}, {"D", 7}},
		"C": {{"E", 3}},
		"D": {{"E", 1}},
		"E": {},
	}
	fmt.Println("Shortest path from A to E:", ShortestPathSimple(g, "A", "E"))
	if ShortestPathSimple(g, "A", "E") != 6 {
		t.Errorf("expected 6, got %d", ShortestPathSimple(g, "A", "E"))
	}
}
*/
