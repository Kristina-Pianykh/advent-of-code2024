package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const MAX_INT = int((uint(1) << 63) - 1)

// An Item is something we manage in a priority queue.
type Node struct {
	cost int // The priority of the item in the queue.
	c    Coordinate
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

func (n *Node) string() string {
	return fmt.Sprintf("c: %s, index: %d, cost: %d", n.c.string(), n.index, n.cost)
}

// A PriorityQueue implements heap.Interface and holds Nodes.
type PriorityQueue []*Node

func (pq PriorityQueue) Contains(c Coordinate) bool {
	for i := range pq {
		if pq[i].c.x == c.x && pq[i].c.y == c.y {
			return true
		}
	}
	return false
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest, not highest, priority so we use less than here.
	return pq[i].cost < pq[j].cost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	node := x.(*Node) // type assertion: remove cause we pass *Node anyway?
	node.index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	node.index = -1 // for safety
	*pq = old[0 : n-1]
	return node
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(c Coordinate, cost int) {
	for i := range *pq {
		node := (*pq)[i]
		if node.c.x == c.x && node.c.y == c.y {
			node.cost = cost
			heap.Fix(pq, node.index)
		}
	}
}

type Coordinate struct {
	x    int
	y    int
	v    byte
	cost int
}

func (c Coordinate) string() string {
	return fmt.Sprintf("x: %d, y: %d, v: %c", c.x, c.y, c.v)
}

type Grid [][]byte
type Visited [][]bool

type Maze struct {
	grid  *Grid
	start Coordinate
	end   Coordinate
}

// func remove(arr []any, idxToRemove int) []any {
// 	if len(arr) == 0 {
// 		return arr
// 	}
// 	if !slices.Contains(arr, arr[idxToRemove]) {
// 		return arr
// 	}
//
// 	var new_arr []Coordinate = make([]Coordinate, len(arr)-1)
// 	idx := 0
// 	for i, c := range arr {
// 		if i == idxToRemove {
// 			continue
// 		}
// 		new_arr[idx] = c
// 		idx++
// 	}
// 	return new_arr
// }

func main() {
	lines := readLinesFromStream(os.Stdin)
	maze := parseGrid(lines)
	res := maze.dijkstra()
	fmt.Printf("part 1 | score: %d\n", res)
	// res2 := maze.bfs(res)
	// fmt.Printf("part 2 | res: %d\n", res2)
	// drawGrid := initDrawGrid(*maze.grid)
}

// func (m *Maze) bfs(targetScore int) int {
// 	size := len(*m.grid)
// 	dist := make([][]int, size)
// 	for y := range *m.grid {
// 		dist[y] = make([]int, size)
// 		for x := range (*m.grid)[y] {
// 			dist[y][x] = MAX_INT
// 		}
// 	}
//
// 	visited := make([][]map[bool]int, size)
// 	for y := range visited {
// 		visited[y] = make([]map[bool]int, size)
// 		for x := range visited[y] {
// 			visited[y][x] = map[bool]int{false: 0}
// 		}
// 	}
//
// 	toVisit := []Coordinate{}
// 	toVisit = append(toVisit, Coordinate{
// 		x: m.start.x, y: m.start.y, v: m.start.v, cost: 0},
// 	)
//
// 	dist[m.start.y][m.start.x] = 0
//
// 	for len(toVisit) > 0 {
// 		fmt.Printf("len(toVisit): %d\n", len(toVisit))
//
// 		u := toVisit[0]
// 		// toVisit = toVisit[1:]
// 		toVisit = remove(toVisit, 0)
//
// 		if u.x == m.end.x && u.y == m.end.y {
//
// 			if u.cost == targetScore {
// 				visited[u.y][u.x] = map[bool]int{true: u.cost}
// 			}
// 			continue
// 		}
//
// 		neighbors := u.getNeighbors(m.grid)
// 		for _, n := range neighbors {
// 			alt := u.cost + getCost(u, n)
// 			// fmt.Printf("dist of %s: %d\n", n.string(), dist[n.y][n.x])
//
// 			n.cost = alt
// 			toVisit = append(toVisit, n)
// 			if !contains(toVisit, n) && !containsVisited(visited, n) {
// 				toVisit = append(toVisit, n)
// 			}
// 			visited[n.y][n.x] = map[bool]int{true: n.cost}
// 			// if alt < dist[n.y][n.x] {
// 			// 	dist[n.y][n.x] = alt
// 			// 	n.cost = alt
// 			// 	toVisit = append(toVisit, n)
// 			// }
// 		}
// 	}
//
// 	cnt := 0
// 	for y := range visited {
// 		for _, cell := range visited[y] {
// 			for k := range cell {
// 				if k {
// 					cnt++
// 				}
// 			}
// 		}
// 	}
//
// 	return cnt
// 	// return -1
// }

func contains(toVisit []Coordinate, c Coordinate) bool {
	for _, co := range toVisit {
		if co.x == c.x && co.y == c.y {
			return true
		}
	}
	return false
}

func containsVisited(visited [][]map[bool]int, c Coordinate) bool {
	for k := range visited[c.y][c.x] {
		if k && visited[c.y][c.x][k] == c.cost {
			return true
		}
	}
	return false
}

func (m *Maze) dijkstra() int {
	rows := len(*m.grid)
	cols := len((*m.grid)[0])
	dist := make([][]int, rows)
	for y := range *m.grid {
		dist[y] = make([]int, cols)
		for x := range (*m.grid)[y] {
			dist[y][x] = MAX_INT
		}
	}
	drawGrid := initDrawGrid(*m.grid)

	cnt := 0
	pq := make(PriorityQueue, 0, rows*cols)
	heap.Init(&pq)

	// Start node
	startNode := &Node{
		c:    m.start,
		cost: 0,
	}
	heap.Push(&pq, startNode)
	drawGrid[startNode.c.y][startNode.c.x] = 'O'
	cnt++
	dist[m.start.y][m.start.x] = 0

	for pq.Len() > 0 {
		u := heap.Pop(&pq).(*Node)

		if u.c.x == m.end.x && u.c.y == m.end.y {
			fmt.Printf("hit the end\n")
			// return u.cost
		}

		neighbors := u.c.getNeighbors(m.grid)
		for _, n := range neighbors {
			alt := u.cost + getCost(u.c, n)
			// fmt.Printf("dist of %s: %d\n", n.string(), dist[n.y][n.x])

			if alt < dist[n.y][n.x] {
				dist[n.y][n.x] = alt
				heap.Push(&pq, &Node{
					c:    n,
					cost: alt,
				})

				drawGrid[n.y][n.x] = 'O'
				// drawGrid.printGrid()
				cnt++
			}
		}
	}

	fmt.Printf("cnt: %d\n", cnt)
	// drawGrid.printGrid()
	return -1
}

func (g Grid) printGrid() {
	for y := range g {
		fmt.Printf("%s\n", string(g[y]))
	}
}

func getCost(u, v Coordinate) int {
	if u.v == v.v {
		return 1
	}
	return 1001
}

func (c Coordinate) getNeighbors(grid *Grid) []Coordinate {
	dirs := map[byte][]int{
		'>': {1, 0},
		'<': {-1, 0},
		'^': {0, -1},
		'v': {0, 1},
	}
	var rotations = map[byte][]byte{
		'>': {'v', '^'},
		'<': {'^', 'v'},
		'^': {'>', '<'},
		'v': {'<', '>'},
	}
	neighbors := []Coordinate{}
	straight := Coordinate{x: c.x + dirs[c.v][0], y: c.y + dirs[c.v][1]}
	if (*grid)[straight.y][straight.x] != '#' {
		straight.v = c.v
		neighbors = append(neighbors, straight)
	}

	for _, rot := range rotations[c.v] {
		swirved := Coordinate{x: c.x + dirs[rot][0], y: c.y + dirs[rot][1]}
		if (*grid)[swirved.y][swirved.x] != '#' {
			swirved.v = rot
			neighbors = append(neighbors, swirved)
		}
	}

	// for _, n := range neighbors {
	// 	fmt.Printf("%s\n", n.string())
	// }
	// fmt.Println()
	return neighbors
}

func initDrawGrid(grid Grid) Grid {
	rows := len(grid)
	cols := len(grid[0])
	var drawGrid Grid = make([][]byte, rows)
	for y := range grid {
		drawGrid[y] = make([]byte, cols)
		copy(drawGrid[y], grid[y])
	}
	return drawGrid
}

func initVisited(size int) Visited {
	var visited Visited = make([][]bool, size)
	for y := range visited {
		visited[y] = make([]bool, size)
	}
	return visited
}

func (g *Grid) write(name int, score int) {
	filename := fmt.Sprintf("%d.txt", name)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	defer func() {
		if err := writer.Flush(); err != nil {
			panic(err)
		}
	}()

	fmt.Fprintf(writer, "score: %d\n", score)
	for i := range *g {
		fmt.Fprintf(writer, "%s\n", string((*g)[i]))
	}
	return
}

func parseGrid(lines []string) *Maze {
	var g Grid = [][]byte{}
	maze := Maze{grid: &g}

	for y, line := range lines {
		if len(line) == 0 {
			break
		}
		row := []byte{}
		for i, ch := range lines[y] {
			if ch == 'S' {
				maze.start = Coordinate{x: i, y: y, v: '>'}
				row = append(row, '.')
				continue
			}
			if ch == 'E' {
				maze.end = Coordinate{x: i, y: y, v: byte(ch)}

			}
			row = append(row, byte(ch))
		}
		*maze.grid = append(*maze.grid, row)
	}
	return &maze
}

func readLinesFromStream(file *os.File) []string {
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines
}
