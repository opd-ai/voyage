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
	if !w.isValidPathEndpoints(start, end) {
		return PathResult{Found: false}
	}

	state := w.initPathfindingState(start, end)

	for state.openSet.Len() > 0 {
		current := heap.Pop(state.openSet).(*pathNode).point
		delete(state.inOpen, current)

		if current == end {
			return w.buildPathResult(state, current)
		}

		w.processNeighbors(state, current, end)
	}

	return PathResult{Found: false}
}

// pathfindingState holds A* algorithm state.
type pathfindingState struct {
	openSet  *priorityQueue
	cameFrom map[Point]Point
	gScore   map[Point]int
	fScore   map[Point]int
	inOpen   map[Point]bool
}

// isValidPathEndpoints checks if start and end points are valid tiles.
func (w *WorldMap) isValidPathEndpoints(start, end Point) bool {
	return w.GetTile(start.X, start.Y) != nil && w.GetTile(end.X, end.Y) != nil
}

// initPathfindingState initializes A* algorithm data structures.
func (w *WorldMap) initPathfindingState(start, end Point) *pathfindingState {
	openSet := &priorityQueue{}
	heap.Init(openSet)

	state := &pathfindingState{
		openSet:  openSet,
		cameFrom: make(map[Point]Point),
		gScore:   map[Point]int{start: 0},
		fScore:   map[Point]int{start: heuristic(start, end)},
		inOpen:   map[Point]bool{start: true},
	}

	heap.Push(openSet, &pathNode{point: start, priority: state.fScore[start]})
	return state
}

// processNeighbors evaluates all neighbors of the current node.
func (w *WorldMap) processNeighbors(state *pathfindingState, current, end Point) {
	tile := w.GetTile(current.X, current.Y)
	if tile == nil {
		return
	}

	for _, neighbor := range tile.Connections {
		w.evaluateNeighbor(state, current, neighbor, end)
	}
}

// evaluateNeighbor processes a single neighbor in A* search.
func (w *WorldMap) evaluateNeighbor(state *pathfindingState, current, neighbor, end Point) {
	neighborTile := w.GetTile(neighbor.X, neighbor.Y)
	if neighborTile == nil {
		return
	}

	info := terrainBaseInfo[neighborTile.Terrain]
	tentativeG := state.gScore[current] + info.MovementCost

	if oldG, exists := state.gScore[neighbor]; !exists || tentativeG < oldG {
		state.cameFrom[neighbor] = current
		state.gScore[neighbor] = tentativeG
		state.fScore[neighbor] = tentativeG + heuristic(neighbor, end)

		if !state.inOpen[neighbor] {
			heap.Push(state.openSet, &pathNode{point: neighbor, priority: state.fScore[neighbor]})
			state.inOpen[neighbor] = true
		}
	}
}

// buildPathResult constructs the final path result.
func (w *WorldMap) buildPathResult(state *pathfindingState, current Point) PathResult {
	return PathResult{
		Path:  reconstructPath(state.cameFrom, current),
		Cost:  state.gScore[current],
		Found: true,
	}
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
// Includes a maximum iteration guard to prevent infinite loops (M-013).
func reconstructPath(cameFrom map[Point]Point, current Point) []Point {
	path := []Point{current}
	// Safety guard: max iterations = map size (M-013)
	maxIter := len(cameFrom) + 1
	for i := 0; i < maxIter; i++ {
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
