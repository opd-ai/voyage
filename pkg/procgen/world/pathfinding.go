package world

import "container/heap"

// PathResult holds the result of pathfinding.
type PathResult struct {
	Path  []Point
	Cost  int
	Found bool
}

// FindPath finds a path from start to end using A* algorithm.
func (w *WorldMap) FindPath(start, end Point) PathResult {
	if w.GetTile(start.X, start.Y) == nil || w.GetTile(end.X, end.Y) == nil {
		return PathResult{Found: false}
	}

	openSet := &priorityQueue{}
	heap.Init(openSet)

	cameFrom := make(map[Point]Point)
	gScore := make(map[Point]int)
	gScore[start] = 0

	fScore := make(map[Point]int)
	fScore[start] = heuristic(start, end)

	heap.Push(openSet, &pathNode{point: start, priority: fScore[start]})
	inOpen := map[Point]bool{start: true}

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*pathNode).point
		delete(inOpen, current)

		if current == end {
			return PathResult{
				Path:  reconstructPath(cameFrom, current),
				Cost:  gScore[current],
				Found: true,
			}
		}

		tile := w.GetTile(current.X, current.Y)
		if tile == nil {
			continue
		}

		for _, neighbor := range tile.Connections {
			neighborTile := w.GetTile(neighbor.X, neighbor.Y)
			if neighborTile == nil {
				continue
			}

			info := terrainBaseInfo[neighborTile.Terrain]
			tentativeG := gScore[current] + info.MovementCost

			if oldG, exists := gScore[neighbor]; !exists || tentativeG < oldG {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentativeG
				fScore[neighbor] = tentativeG + heuristic(neighbor, end)

				if !inOpen[neighbor] {
					heap.Push(openSet, &pathNode{point: neighbor, priority: fScore[neighbor]})
					inOpen[neighbor] = true
				}
			}
		}
	}

	return PathResult{Found: false}
}

// HasPath checks if a path exists between two points.
func (w *WorldMap) HasPath(start, end Point) bool {
	return w.FindPath(start, end).Found
}

// heuristic estimates distance between two points (Manhattan distance).
func heuristic(a, b Point) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

// reconstructPath builds the path from start to end.
func reconstructPath(cameFrom map[Point]Point, current Point) []Point {
	path := []Point{current}
	for {
		prev, exists := cameFrom[current]
		if !exists {
			break
		}
		path = append([]Point{prev}, path...)
		current = prev
	}
	return path
}

// Priority queue implementation for A*
type pathNode struct {
	point    Point
	priority int
	index    int
}

type priorityQueue []*pathNode

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*pathNode)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*pq = old[0 : n-1]
	return node
}
