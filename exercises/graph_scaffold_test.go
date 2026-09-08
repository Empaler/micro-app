package exercises

import "testing"

type WeightedGraph map[string][]GraphEdge

type GraphEdge struct {
	To     string
	Weight int
}

func ShortestPathInterview(graph WeightedGraph, start, end string) int {
	if start == "" || end == "" {
		return -1
	}
	if len(graph) == 0 {
		return -1
	}
	if start == end {
		return 0
	}
	_, ok := graph[start]
	if !ok {
		return -1
	}
	_, ok = graph[end]
	if !ok {
		return -1
	}

	bestCost := make(map[string]int)

	bestCost[start] = 0
	toDoList := []string{start}

	for len(toDoList) > 0 {
		nodeLabel := toDoList[0]
		for _, node := range graph[nodeLabel] {
			if node.To != "" {
				_, ok := bestCost[node.To]
				if !ok {
					bestCost[node.To] = bestCost[nodeLabel] + node.Weight
					toDoList = append(toDoList, node.To)
				}
				if bestCost[node.To] > node.Weight+bestCost[nodeLabel] {
					bestCost[node.To] = node.Weight + bestCost[nodeLabel]
					toDoList = append(toDoList, node.To)
				}
			}
		}
		toDoList = toDoList[1:]
	}
	_, ok = bestCost[end]
	if ok {
		return bestCost[end]
	} else {
		return -1
	}
}

func TestGraphInterviewDirectEdge(t *testing.T) {
	graph := WeightedGraph{
		"A": {{To: "B", Weight: 5}},
		"B": {},
	}

	got := ShortestPathInterview(graph, "A", "B")
	if got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestGraphInterviewTwoPathsPickShortest(t *testing.T) {
	graph := WeightedGraph{
		"A": {{To: "B", Weight: 10}, {To: "C", Weight: 3}},
		"B": {{To: "D", Weight: 1}},
		"C": {{To: "D", Weight: 2}},
		"D": {},
	}

	// A->C->D costs 3+2=5, which is shorter than A->B->D (10+1=11)
	got := ShortestPathInterview(graph, "A", "D")
	if got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestGraphInterviewUnreachable(t *testing.T) {
	graph := WeightedGraph{
		"A": {{To: "B", Weight: 5}},
		"B": {},
		"C": {},
	}

	got := ShortestPathInterview(graph, "A", "C")
	if got != -1 {
		t.Fatalf("expected -1 (unreachable), got %d", got)
	}
}

func TestGraphInterviewStartEqualsEnd(t *testing.T) {
	graph := WeightedGraph{
		"A": {},
	}

	got := ShortestPathInterview(graph, "A", "A")
	if got != 0 {
		t.Fatalf("expected 0 (same node), got %d", got)
	}
}
