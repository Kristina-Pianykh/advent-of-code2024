package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"os"
	"slices"
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
	x int
	y int
	v byte
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

func main() {
	lines := readLinesFromStream(os.Stdin)
	maze := parseGrid(lines)
	res := maze.dijkstra()
	fmt.Printf("res: %d\n", res)
	// drawGrid := initDrawGrid(*maze.grid)
}

func (v Visited) get(c Coordinate) bool {
	return v[c.y][c.x]
}

func (v Visited) set(c Coordinate, val bool) {
	v[c.y][c.x] = val
}

func (g *Grid) get(c Coordinate) byte {
	return (*g)[c.y][c.x]
}

func (g *Grid) set(c Coordinate, val byte) {
	(*g)[c.y][c.x] = val
}

func (c Coordinate) diff(co Coordinate) []int {
	return []int{c.x - co.x, c.y - co.y}
}

func (m *Maze) dijkstra() int {
	size := len(*m.grid)
	dist := make([][]int, size)
	for y := range *m.grid {
		dist[y] = make([]int, size)
		for x := range (*m.grid)[y] {
			dist[y][x] = MAX_INT
		}
	}

	pq := make(PriorityQueue, 0, size*size)
	heap.Init(&pq)

	// Start node
	startNode := &Node{
		c:    m.start,
		cost: 0,
	}
	heap.Push(&pq, startNode)
	dist[m.start.y][m.start.x] = 0
	// node := heap.Pop(&pq).(*Node)
	// fmt.Printf("%s\n", node.string())

	for pq.Len() > 0 {

		u := heap.Pop(&pq).(*Node)

		if u.c.x == m.end.x && u.c.y == m.end.y {
			return u.cost
		}

		neighbors := u.c.getNeighbors(m.grid)
		for _, n := range neighbors {
			// if !pq.Contains(n) {
			// 	continue
			// }

			alt := u.cost + getCost(u.c, n)

			fmt.Printf("dist of %s: %d\n", n.string(), dist[n.y][n.x])
			if alt < dist[n.y][n.x] {
				// prev[n.y][n.x] = u.c
				dist[n.y][n.x] = alt
				heap.Push(&pq, &Node{
					c:    n,
					cost: alt,
				})
			}
		}
	}

	return -1
}

func getCost(u, v Coordinate) int {
	if u.v == v.v {
		return 1
	}
	return 1001
	// switch u.v {
	// case '>':
	// 	if v.x == u.x+1 && v.y == u.y {
	// 		return 1
	// 	}
	// case '<':
	// 	if v.x == u.x-1 && v.y == u.y {
	// 		return 1
	// 	}
	// case 'v':
	// 	if v.x == u.x && v.y == u.y+1 {
	// 		return 1
	// 	}
	// case '^':
	// 	if v.x == u.x && v.y == u.y-1 {
	// 		return 1
	// 	}
	// default:
	// 	panic(fmt.Sprintf("undefined direction of current node: %c\n", u.v))
	// }
	// return 1001
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
	// fmt.Printf("%c\n", c.v)
	straight := Coordinate{x: c.x + dirs[c.v][0], y: c.y + dirs[c.v][1]}
	if grid.get(straight) != '#' {
		straight.v = c.v
		neighbors = append(neighbors, straight)
	}

	for _, rot := range rotations[c.v] {
		swirved := Coordinate{x: c.x + dirs[rot][0], y: c.y + dirs[rot][1]}
		if grid.get(swirved) != '#' {
			swirved.v = rot
			neighbors = append(neighbors, swirved)
		}
	}

	for _, n := range neighbors {
		fmt.Printf("%s\n", n.string())
	}
	fmt.Println()
	return neighbors
}

func nodeWithMinDist(distGrid *[][]int) Coordinate {
	minDist := MAX_INT
	c := Coordinate{}
	for y := range *distGrid {
		for x, cell := range (*distGrid)[y] {
			if cell < minDist {
				minDist = cell
				c.y = y
				c.x = x
			}
		}
	}
	return c
}

func remove(arr []Coordinate, toRemove Coordinate) []Coordinate {
	if len(arr) == 0 {
		return arr
	}
	if !slices.Contains(arr, toRemove) {
		return arr
	}

	var new_arr []Coordinate = make([]Coordinate, len(arr)-1)
	idx := 0
	for _, c := range arr {
		if c == toRemove {
			continue
		}
		new_arr[idx] = c
		idx++
	}
	return new_arr
}

func getKeyByValue(dirs map[byte][]int, val []int) (byte, error) {
	for k := range dirs {
		if dirs[k][0] == val[0] && dirs[k][1] == val[1] {
			return k, nil
		}
	}
	return '0', errors.New(fmt.Sprintf("key for %v not found in dirs\n", val))
}

func initDrawGrid(grid Grid) Grid {
	size := len(grid)
	var drawGrid Grid = make([][]byte, size)
	for y := range grid {
		drawGrid[y] = make([]byte, size)
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

// func (grid *Grid) printGrid() {
// 	for y := range *grid {
// 		for _, cell := range (*grid)[y] {
// 			fmt.Printf("%c", cell)
// 		}
// 		fmt.Println()
// 	}
// }

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
