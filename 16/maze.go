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
	drawGrid := initDrawGrid(*maze.grid)
	res := maze.dijkstra(&drawGrid)
	fmt.Printf("part 1 | score: %d\n", res)
}

func (m *Maze) dijkstra(drawGrid *Grid) int {
	rows := len(*m.grid)
	cols := len((*m.grid)[0])
	dist := make([][]int, rows)
	paths := make([][][]Coordinate, rows)
	mem := make([][]map[byte]int, rows)

	for y := range *m.grid {
		dist[y] = make([]int, rows)
		paths[y] = make([][]Coordinate, cols)
		mem[y] = make([]map[byte]int, cols)
		for x := range (*m.grid)[y] {
			dist[y][x] = MAX_INT
			paths[y][x] = []Coordinate{}
			mem[y][x] = map[byte]int{
				'>': MAX_INT,
				'<': MAX_INT,
				'^': MAX_INT,
				'v': MAX_INT,
				'E': MAX_INT,
			}
		}
	}

	pq := make(PriorityQueue, 0, rows*rows)
	heap.Init(&pq)

	// Start node
	startNode := &Node{
		c:    m.start,
		cost: 0,
	}
	heap.Push(&pq, startNode)
	dist[m.start.y][m.start.x] = 0
	// paths[m.start.y][m.start.x][m.start] = true

	for pq.Len() > 0 {
		// fmt.Printf("pq:\n")
		// for _, v := range pq {
		// 	fmt.Printf("  %s\n", v.string())
		// }

		u := heap.Pop(&pq).(*Node)
		mem[u.c.y][u.c.x][u.c.v] = u.cost
		// fmt.Printf("next u: %s\n", u.string())

		if u.c.x == m.end.x && u.c.y == m.end.y {
			paths[m.end.y][m.end.x] = append(paths[m.end.y][m.end.x], u.c)
			continue
			// return u.cost
		}

		neighbors := u.c.getNeighbors(m.grid, false)
		for _, n := range neighbors {
			alt := u.cost + getCost(u.c, n)
			// fmt.Printf("dist of %s: %d\n", n.string(), dist[n.y][n.x])

			// fmt.Printf("mem[n.y][n.x][n.v] = %d\n", mem[n.y][n.x][n.v])
			// fmt.Printf("alt = %d\n", alt)
			if alt < mem[n.y][n.x][n.v] {
				dist[n.y][n.x] = alt
				heap.Push(&pq, &Node{
					c:    n,
					cost: alt,
				})
				// fmt.Printf("pushing %s\n", n.string())
				paths[n.y][n.x] = append(paths[n.y][n.x], u.c)
			}
		}
	}

	pp := []Coordinate{}
	var rec func(c Coordinate, prevC Coordinate, currScore int) bool
	visited := make([][]bool, rows)
	for y := range *m.grid {
		visited[y] = make([]bool, cols)
	}

	rec = func(c Coordinate, prevC Coordinate, currScore int) bool {
		rotations := map[byte][]byte{
			'>': {'^', 'v'},
			'<': {'v', '^'},
			'^': {'<', '>'},
			'v': {'>', '<'},
		}
		var newScore int
		if visited[c.y][c.x] {
			return false
		}
		if currScore < 0 {
			return false
		}

		if c.v != prevC.v && c.v != rotations[prevC.v][0] && c.v != rotations[prevC.v][1] {
			return false
		}
		// fmt.Printf("c: %s\n", c.string())

		visited[c.y][c.x] = true
		defer func() { visited[c.y][c.x] = false }()
		if c.v == prevC.v {
			newScore = currScore - 1
		} else {
			newScore = currScore - 1001
		}

		if c.x == m.start.x && c.y == m.start.y {
			fmt.Printf("reached start %s with score: %d\n", c.string(), newScore)
			if c.v == '>' && newScore == 0 {
				fmt.Printf("we are here\n")
				return true
			}
			if c.v != '>' && newScore == 1000 {
				fmt.Printf("we are here\n")
				return true
			}
			return false
		}

		all := []bool{}
		for _, n := range c.getNeighbors(m.grid, true) {
			res := rec(n, c, newScore)
			all = append(all, res)
			if res {
				if !contains(pp, n) {
					pp = append(pp, n)
				}
			}
		}
		for _, res := range all {
			if res {
				return true
			}
		}
		return false
	}

	minScore := MAX_INT
	for key := range mem[m.end.y][m.end.x] {
		if mem[m.end.y][m.end.x][key] < minScore {
			minScore = mem[m.end.y][m.end.x][key]
		}
	}

	fmt.Printf("end: %s\n", m.end.string())
	pp = append(pp, m.end)
	for key := range mem[m.end.y][m.end.x] {
		if mem[m.end.y][m.end.x][key] == minScore {
			fmt.Printf("%c\n", key)
			prevC := m.end
			prevC.v = key
			for _, n := range prevC.getNeighbors(m.grid, true) {
				fmt.Printf("high level n: %s, prevC: %s\n", n.string(), prevC.string())
				if rec(n, prevC, minScore) {
					pp = append(pp, n)
				}
			}
		}
	}
	fmt.Printf("micScore: %d\n", minScore)
	fmt.Printf("cells: %d\n", len(pp))

	for y := range *drawGrid {
		for x, cell := range (*drawGrid)[y] {
			if contains(pp, Coordinate{x: x, y: y}) {
				fmt.Printf("O")
				continue
			}
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
	return -1
}

func contains(lst []Coordinate, c Coordinate) bool {
	for i := range lst {
		if lst[i].x == c.x && lst[i].y == c.y {
			return true
		}
	}
	return false
}

func getCost(u, v Coordinate) int {
	if u.v == v.v {
		return 1
	}
	return 1001
}

func (c Coordinate) getNeighbors(grid *Grid, backtrack bool) []Coordinate {
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
	var straight Coordinate
	var swirved Coordinate

	neighbors := []Coordinate{}
	if backtrack {
		straight = Coordinate{x: c.x - dirs[c.v][0], y: c.y - dirs[c.v][1]}
	} else {
		straight = Coordinate{x: c.x + dirs[c.v][0], y: c.y + dirs[c.v][1]}
	}

	if (*grid)[straight.y][straight.x] != '#' {
		straight.v = c.v
		neighbors = append(neighbors, straight)
	}

	for _, rot := range rotations[c.v] {
		if backtrack {
			swirved = Coordinate{x: c.x - dirs[rot][0], y: c.y - dirs[rot][1]}
		} else {
			swirved = Coordinate{x: c.x + dirs[rot][0], y: c.y + dirs[rot][1]}
		}
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
