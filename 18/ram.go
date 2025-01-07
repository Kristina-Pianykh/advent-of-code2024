package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const MAX_INT = int((uint(1) << 63) - 1)

type Coordinate struct {
	x int
	y int
}

func (c Coordinate) string() string {
	return fmt.Sprintf("x: %d, y: %d", c.x, c.y)
}

type Grid [][]byte
type Maze struct {
	grid  *Grid
	start Coordinate
	end   Coordinate
}

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

func main() {
	lines := readLinesFromStream(os.Stdin)
	cs := parse(lines)
	byteLimit := 1024

	for i := byteLimit; i < len(cs); i++ {
		maze := initGrid(cs, i)
		cnt := maze.dijkstra()

		if i == byteLimit {
			fmt.Printf("part 1 | steps : %d\n", cnt)
		} else {
			if cnt == -1 {
				fmt.Printf("part 2 | byte: %s\n", cs[i-1].string())
				break
			}
		}
	}
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

	pq := make(PriorityQueue, 0, rows*cols)
	heap.Init(&pq)

	// Start node
	startNode := &Node{
		c:    m.start,
		cost: 0,
	}
	heap.Push(&pq, startNode)
	drawGrid[startNode.c.y][startNode.c.x] = 'O'
	dist[m.start.y][m.start.x] = 0

	for pq.Len() > 0 {
		u := heap.Pop(&pq).(*Node)

		if u.c.x == m.end.x && u.c.y == m.end.y {
			return u.cost
		}

		neighbors := m.getNeighbors(u.c)
		for _, n := range neighbors {
			alt := u.cost + 1

			if alt < dist[n.y][n.x] {
				dist[n.y][n.x] = alt
				heap.Push(&pq, &Node{
					c:    n,
					cost: alt,
				})

				drawGrid[n.y][n.x] = 'O'
			}
		}
	}

	print(&drawGrid)
	return -1
}

func (m *Maze) getNeighbors(c Coordinate) []Coordinate {
	dirs := [][]int{
		{1, 0},
		{-1, 0},
		{0, -1},
		{0, 1},
	}
	neighbors := []Coordinate{}
	for _, dir := range dirs {
		newX := c.x + dir[0]
		newY := c.y + dir[1]
		newC := Coordinate{x: newX, y: newY}

		if m.ok(newC) {
			neighbors = append(neighbors, newC)
		}
	}
	return neighbors
}

func (m *Maze) ok(c Coordinate) bool {
	if c.x < 0 || c.x >= len((*m.grid)[0]) || c.y < 0 || c.y >= len(*m.grid) {
		return false
	}
	if (*m.grid)[c.y][c.x] == '#' {
		return false
	}
	return true
}

func print(g *Grid) {
	for y := range *g {
		fmt.Printf("%s\n", string((*g)[y]))
	}
}

func initGrid(cs []Coordinate, limit int) *Maze {
	size := 71
	var g Grid = make([][]byte, size)
	for i := range g {
		g[i] = make([]byte, size)
		for x := range g[i] {
			g[i][x] = '.'
		}
	}

	for i := 0; i < limit; {
		c := cs[i]
		g[c.y][c.x] = '#'
		i++
	}
	return &Maze{grid: &g, start: Coordinate{0, 0}, end: Coordinate{size - 1, size - 1}}
}

func parse(lines []string) []Coordinate {
	reg := regexp.MustCompile("^([[:digit:]]+),([[:digit:]]+)$")
	cs := []Coordinate{}

	for _, line := range lines {
		if len(line) > 0 {
			x, err := strconv.Atoi(string(reg.FindSubmatch([]byte(line))[1]))
			if err != nil {
				panic(err)
			}
			y, err := strconv.Atoi(string(reg.FindSubmatch([]byte(line))[2]))
			if err != nil {
				panic(err)
			}
			cs = append(cs, Coordinate{x: x, y: y})
		}
	}
	return cs
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
