package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
)

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

	// maze.grid.printGrid()
	// fmt.Printf("Start pos: %s\n", maze.start.string())
	// fmt.Printf("End pos: %s\n", maze.end.string())

	score := maze.shortestPath(&drawGrid)
	// 452371 too high
	fmt.Printf("%d\n", score)
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

func get(g *[][]int, c Coordinate) int {
	return (*g)[c.y][c.x]
}

func set(g *[][]int, c Coordinate, val int) {
	(*g)[c.y][c.x] = val
}

func (c Coordinate) diff(co Coordinate) []int {
	return []int{c.x - co.x, c.y - co.y}
}

func (m *Maze) shortestPath(drawGrid *Grid) int {
	shortestPathScore := int((uint(1) << 63) - 1)
	visited := initVisited(len(*m.grid))
	scoreGrid := initScoreGrid(len(*m.grid))
	neighbors := make([]Coordinate, 0, 140*140)
	visited[m.start.y][m.start.x] = true
	var pos Coordinate = m.start
	fileIdx := 0

	var rec func(pos Coordinate)

	dirs := map[byte][]int{
		'>': {1, 0},
		'<': {-1, 0},
		'^': {0, -1},
		'v': {0, 1},
	}

	rec = func(pos Coordinate) {
		if visited.get(pos) {
			return
		}

		if m.grid.get(pos) == '#' {
			return
		}

		var score int
		if pos.x == m.end.x && pos.y == m.end.y {
			prevPos := getAdjacent(visited, drawGrid, pos)
			fmt.Printf("prevPos: %s\n", prevPos.string())
			scoreGrid[pos.y][pos.x] = scoreGrid[prevPos.y][prevPos.x] + 1

			if scoreGrid[pos.y][pos.x] < shortestPathScore {
				shortestPathScore = scoreGrid[pos.y][pos.x]
			}

			fmt.Printf("reached end with score: %d\n", scoreGrid[pos.y][pos.x])
			visited.set(m.end, true)
			return
		}
		visited.set(pos, true)

		prevPos := getAdjacent(visited, drawGrid, pos)
		fmt.Printf("prevPos: %s\n", prevPos.string())
		delta := pos.diff(prevPos)
		// fmt.Printf("delta: %v\n", delta)
		if len(delta) != 2 {
			panic(fmt.Sprintf("pos: %s, prevPos: %s, delta: %v\n", pos.string(), prevPos.string(), delta))
		}

		if prevPos.y == m.end.y && prevPos.x == m.end.x {
			// scoreGrid[pos.y][pos.x] = scoreGrid[prevPos.y][prevPos.x] + 1
			// if scoreGrid[pos.y][pos.x] < shortestPathScore {
			// 	shortestPathScore = scoreGrid[pos.y][pos.x]
			// }
		}

		if dirs[prevPos.v][0] == delta[0] && dirs[prevPos.v][1] == delta[1] {
			pos.v = prevPos.v
			score = 1
			drawGrid.set(pos, pos.v)
		} else {
			val, err := getKeyByValue(dirs, delta)
			if err != nil {
				panic(err)
			}
			// fmt.Printf("changed direction from %c to %c\n", pos.v, val)
			pos.v = val
			// fmt.Printf("pos updated direction now: %s\n", pos.string())
			score = 1001
			drawGrid.set(pos, val)
		}
		scoreGrid[pos.y][pos.x] = score + scoreGrid[prevPos.y][prevPos.x]
		fmt.Printf("score for %s: %d\n", pos.string(), scoreGrid[pos.y][pos.x])

		fileIdx++
		drawGrid.write(fileIdx, scoreGrid[pos.y][pos.x], neighbors)

		for _, n := range pos.getNeighbors() {
			if !slices.Contains(neighbors, n) && !visited.get(n) {
				neighbors = append(neighbors, n)
			}
		}

		// for dir := range dirs {
		// 	xInc := dirs[dir][0]
		// 	yInc := dirs[dir][1]
		// 	nextPos := Coordinate{x: pos.x + xInc, y: pos.y + yInc, v: dir}
		// 	// fmt.Printf("next step: %s\n", nextPos.string())
		// 	rec(nextPos, pos, score)
		// }
		// drawGrid.set(pos, m.grid.get(pos)) // reset the cell to original
		// visited.set(pos, false)
	}

	(*drawGrid).set(m.start, pos.v)
	// drawGrid.write(fileIdx, 0)

	for _, n := range pos.getNeighbors() {
		neighbors = append(neighbors, n)
	}
	visited.set(m.start, true)

	for {
		fmt.Printf("len(neighbors): %d\n", len(neighbors))
		if len(neighbors) == 0 {
			break
			// return scoreGrid[m.end.y][m.end.x]
			// break
		}

		nextPos := neighbors[0]
		neighbors = neighbors[1:]

		for _, n := range pos.getNeighbors() {
			if !slices.Contains(neighbors, n) && !visited.get(n) {
				neighbors = append(neighbors, n)
			}
		}
		rec(nextPos)
		drawGrid.printGrid()
		if visited.get(m.end) {
			visited.set(m.end, false)
		}
	}

	return shortestPathScore
	// return scoreGrid[m.end.y][m.end.x]
}

func (c Coordinate) getNeighbors() []Coordinate {
	neighbors := []Coordinate{}
	dirs := map[byte][]int{
		'>': {1, 0},
		'<': {-1, 0},
		'^': {0, -1},
		'v': {0, 1},
	}
	for dir := range dirs {
		xInc := dirs[dir][0]
		yInc := dirs[dir][1]
		neighbor := Coordinate{x: c.x + xInc, y: c.y + yInc, v: dir}
		neighbors = append(neighbors, neighbor)
	}
	return neighbors
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

func initScoreGrid(size int) [][]int {
	var grid [][]int = make([][]int, size)
	for y := range grid {
		grid[y] = make([]int, size)
	}
	return grid
}

func (grid *Grid) printGrid() {
	for y := range *grid {
		for _, cell := range (*grid)[y] {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
}

func getAdjacent(v Visited, drawGrid *Grid, c Coordinate) Coordinate {
	dirs := map[byte][]int{
		'>': {1, 0},
		'<': {-1, 0},
		'^': {0, -1},
		'v': {0, 1},
	}
	var pos Coordinate
	for dir := range dirs {
		xInc := dirs[dir][0]
		yInc := dirs[dir][1]

		pos = Coordinate{x: c.x + xInc, y: c.y + yInc}
		fmt.Printf("pos: %s, drawGrid.get(pos): %c\n", pos.string(), drawGrid.get(pos))
		if v.get(pos) && drawGrid.get(pos) != '#' {
			pos.v = drawGrid.get(pos)
			fmt.Printf("pos.v: %c\n", pos.v)
			return pos
		}
	}
	return pos
}

// func remove(arr []Coordinate, toRemove Coordinate) []Coordinate {
// 	if len(arr) == 0 {
// 		return arr
// 	}
// 	if !slices.Contains(arr, toRemove) {
// 		return arr
// 	}
//
// 	var new_arr []Coordinate = make([]Coordinate, len(arr)-1)
// 	idx := 0
// 	for _, c := range arr {
// 		if c == toRemove {
// 			continue
// 		}
// 		new_arr[idx] = c
// 		idx++
// 	}
// 	return new_arr
// }

func (g *Grid) write(name int, score int, neighbors []Coordinate) {
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
	fmt.Fprintf(writer, "neighbors:\n")
	for _, n := range neighbors {
		fmt.Fprintf(writer, "   %s\n", n.string())
	}
	fmt.Fprintln(writer)
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
