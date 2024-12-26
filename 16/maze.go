package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const MAX_INT = int((uint(1) << 63) - 1)

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
	// 265832 too high
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

func (c Coordinate) diff(co Coordinate) []int {
	return []int{c.x - co.x, c.y - co.y}
}

func (m *Maze) shortestPath(drawGrid *Grid) int {
	shortestPathScore := MAX_INT
	visited := initVisited(len(*m.grid))
	visited[m.start.y][m.start.x] = true
	var pos Coordinate = m.start
	fileIdx := 0

	var rec func(pos, prevPos Coordinate, score int)

	dirs := map[byte][]int{
		'>': {1, 0},
		'<': {-1, 0},
		'^': {0, -1},
		'v': {0, 1},
	}

	rec = func(pos Coordinate, prevPos Coordinate, score int) {
		if score >= 88432 {
			return
		}

		if score >= shortestPathScore {
			return
		}

		if visited.get(pos) {
			return
		}

		if m.grid.get(pos) == '#' {
			return
		}

		if pos.x == m.end.x && pos.y == m.end.y {
			score += 1
			fmt.Printf("reached end with score: %d. shortestPathScore: %d\n", score, shortestPathScore)
			if score < shortestPathScore {
				shortestPathScore = score
				fmt.Printf("new shortestPathScore: %d\n", shortestPathScore)
			}
			return
		}
		visited.set(pos, true)
		defer visited.set(pos, false)

		delta := pos.diff(prevPos)
		// fmt.Printf("delta: %v\n", delta)
		if len(delta) != 2 {
			panic(fmt.Sprintf("pos: %s, prevPos: %s, delta: %v\n", pos.string(), prevPos.string(), delta))
		}

		defer drawGrid.set(pos, m.grid.get(pos))

		if dirs[prevPos.v][0] == delta[0] && dirs[prevPos.v][1] == delta[1] {
			pos.v = prevPos.v
			score += 1
			drawGrid.set(pos, pos.v)
		} else {
			val, err := getKeyByValue(dirs, delta)
			if err != nil {
				panic(err)
			}
			// fmt.Printf("changed direction from %c to %c\n", pos.v, val)
			pos.v = val
			// fmt.Printf("pos updated direction now: %s\n", pos.string())
			score += 1001
			drawGrid.set(pos, val)
		}

		fileIdx++
		// drawGrid.printGrid()
		// drawGrid.write(fileIdx, score)

		for dir := range dirs {
			xInc := dirs[dir][0]
			yInc := dirs[dir][1]
			nextPos := Coordinate{x: pos.x + xInc, y: pos.y + yInc, v: dir}
			// fmt.Printf("next step: %s\n", nextPos.string())
			rec(nextPos, pos, score)
		}
	}

	(*drawGrid).set(m.start, pos.v)
	// drawGrid.printGrid()
	// drawGrid.write(fileIdx, 0)

	for dir := range dirs {
		// fmt.Printf("next dir: %c\n", dir)
		xInc := dirs[dir][0]
		yInc := dirs[dir][1]
		nextPos := Coordinate{x: pos.x + xInc, y: pos.y + yInc, v: dir}
		// fmt.Printf("first dir: %s\n", nextPos.string())
		rec(nextPos, pos, 0)
		visited[m.start.y][m.start.x] = false
	}

	return shortestPathScore
}

func (g *Grid) printGrid() {
	fmt.Print("\033[H\033[2J")
	// padding := 2

	var sb strings.Builder
	for y := range *g {
		for _, ch := range (*g)[y] {

			switch ch {
			case '#':
				sb.WriteString("🧱")
			case '>':
				sb.WriteString(fmt.Sprintf("%s", "> "))
			case '<':
				sb.WriteString(fmt.Sprintf("%s", "< "))
			case '^':
				sb.WriteString(fmt.Sprintf("%s", "^ "))
			case 'v':
				sb.WriteString(fmt.Sprintf("%s", "v "))
			case '.':
				// sb.WriteString(fmt.Sprintf("%s", "⬛"))
				sb.WriteString(fmt.Sprintf("%s", ". "))
			case 'S':
				sb.WriteString("🤖")
			case 'E':
				sb.WriteString("🤖")
			default:
				panic(fmt.Sprintf("unrecorgnized character in the grid: %c\n", ch))
			}
		}
		sb.WriteString("\n")
	}
	fmt.Println(sb.String())
	time.Sleep(80 * time.Millisecond)
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
