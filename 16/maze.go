package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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
	grid    *Grid
	start   Coordinate
	end     Coordinate
	pos     Coordinate
	prevPos Coordinate
}

func main() {
	lines := readLinesFromStream(os.Stdin)
	maze := parseGrid(lines)
	maze.printGrid()
	fmt.Printf("Start pos: %s\n", maze.start.string())
	fmt.Printf("End pos: %s\n", maze.end.string())

	score := maze.shortestPath()
	fmt.Printf("%d\n", score)
}

func (v Visited) get(c Coordinate) bool {
	return v[c.y][c.x]
}

func (m *Maze) get(c Coordinate) byte {
	return (*m.grid)[c.y][c.x]
}

func (v Visited) set(c Coordinate, val bool) {
	v[c.y][c.x] = val
}

func (c Coordinate) diff(co Coordinate) []int {
	return []int{c.x - co.x, c.y - co.y}
}

func (m *Maze) updatePos(x, y int, v byte) {
	m.pos.x = x
	m.pos.y = y
	m.pos.v = v
}

func (m *Maze) updatePrevPos(x, y int, v byte) {
	m.prevPos.x = x
	m.prevPos.y = y
	m.prevPos.v = v
}

func (m *Maze) shortestPath() int {
	var pathScore int
	shortestPathScore := int((uint(1) << 63) - 1)
	visited := initVisited(len(*m.grid))
	visited[m.start.y][m.start.x] = true
	var rec func(c Coordinate, score int, scoreInc int) int

	dirs := map[byte][]int{
		'>': {1, 0},
		'<': {-1, 0},
		'^': {0, -1},
		'v': {0, 1},
	}

	rec = func(c Coordinate, score int, scoreInc int) int {
		if visited.get(c) {
			return score
		}
		if m.get(c) == '#' {
			return score
		}
		if c == m.end {
			return score + scoreInc
		}
		visited.set(c, true)
		score = score + scoreInc

		fmt.Printf("updated pos: %s\n", m.pos.string())

		// prevPos == pos at this point
		delta := c.diff(m.prevPos)
		if len(delta) != 2 {
			panic(fmt.Sprintf("c: %s, m.prevPos: %s, delta: %v\n", c.string(), m.prevPos.string(), delta))
		}

		// fmt.Printf("%v\n", dirs[m.prevPos.v])
		// fmt.Printf("%v\n", delta)
		if dirs[m.prevPos.v][0] == delta[0] && dirs[m.prevPos.v][1] == delta[1] {
			m.updatePos(c.x, c.y, m.pos.v)
			m.updatePrevPos(c.x, c.y, m.prevPos.v)
			scoreInc = 1
		} else {
			val, err := getKeyByValue(dirs, []int{0, -1})
			if err != nil {
				panic(err)
			}
			m.updatePos(c.x, c.y, val)
			m.updatePrevPos(c.x, c.y, val)
			scoreInc = 1001
		}

		for k := range dirs {
			xInc := dirs[k][0]
			yInc := dirs[k][1]
			fmt.Printf("next step: %d %d\n", m.pos.x+xInc, m.pos.y+yInc)
			pathScore = rec(Coordinate{x: m.pos.x + xInc, y: m.pos.y + yInc}, score, scoreInc)

			if m.pos == m.end {
				return pathScore
			} else {
				//backtracking
				visited.set(Coordinate{x: m.pos.x + xInc, y: m.pos.y + yInc}, false)
			}
		}

		return score
	}

	for k := range dirs {
		xInc := dirs[k][0]
		yInc := dirs[k][1]
		fmt.Printf("next step: %d %d\n", m.pos.x+xInc, m.pos.y+yInc)
		pathScore = rec(Coordinate{x: m.pos.x + xInc, y: m.pos.y + yInc}, 0, 1)
		if pathScore < shortestPathScore {
			shortestPathScore = pathScore
		}
	}

	return shortestPathScore
}

func getKeyByValue(dirs map[byte][]int, val []int) (byte, error) {
	for k := range dirs {
		if dirs[k][0] == val[0] && dirs[k][1] == val[1] {
			return k, nil
		}
	}
	return '0', errors.New(fmt.Sprintf("key for %v not found in dirs\n", val))
}

func initVisited(size int) Visited {
	var visited Visited = make([][]bool, size)
	for y := range visited {
		visited[y] = make([]bool, size)
	}
	return visited
}

func (m *Maze) printGrid() {
	for y := range *m.grid {
		for _, cell := range (*m.grid)[y] {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
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
				maze.start = Coordinate{x: i, y: y, v: byte(ch)}
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
	maze.pos, maze.prevPos = maze.start, maze.start
	maze.pos.v = '>'
	maze.prevPos.v = '>'

	return &maze
}

// func (m *Maze) write(moveIdx string) {
// 	filename := fmt.Sprintf("%s.txt", moveIdx)
// 	f, err := os.Create(filename)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
//
// 	writer := bufio.NewWriter(f)
// 	defer func() {
// 		if err := writer.Flush(); err != nil {
// 			panic(err)
// 		}
// 	}()
//
// 	for i := range *m.grid {
// 		fmt.Fprintf(writer, "%s\n", string((*m.grid)[i]))
// 	}
// 	return
// }

func readLinesFromStream(file *os.File) []string {
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines
}
