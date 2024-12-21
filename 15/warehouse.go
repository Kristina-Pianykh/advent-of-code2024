package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
	"time"
)

type Coordinate struct {
	x int
	y int
	v byte
}

func (c Coordinate) string() string {
	return fmt.Sprintf("x: %d, y: %d, v: %c", c.x, c.y, c.v)
}

type ByY []Coordinate

func (a ByY) Len() int           { return len(a) }
func (a ByY) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByY) Less(i, j int) bool { return a[i].y < a[j].y }

type Grid [][]byte
type Warehouse struct {
	grid     *Grid
	moves    *[]byte
	robotPos Coordinate
}

func main() {
	lines := readLinesFromStream(os.Stdin)
	grid, robotPos := parseGrid(lines)
	w := Warehouse{grid: grid, robotPos: robotPos, moves: parseMoves(lines)}

	// for y := range *w.grid {
	// 	for _, cell := range (*w.grid)[y] {
	// 		fmt.Printf("%c", cell)
	// 	}
	// 	fmt.Println()
	// }
	for _, move := range *w.moves {
		w.update(move)
		// w.write(fmt.Sprintf("%d", i))
		w.printGrid()
	}
	fmt.Printf("res: %d\n", w.calcGps2())
}

func (w *Warehouse) canPushHorizontally(move byte) (int, bool) {
	switch move {
	case '<':
		if (*w.grid)[w.robotPos.y][w.robotPos.x-1] == ']' {
			i := 1
			for {
				if (*w.grid)[w.robotPos.y][w.robotPos.x-i*2-1] == '#' {
					return 0, false
				}
				if (*w.grid)[w.robotPos.y][w.robotPos.x-i*2-1] == '.' {
					break
				}
				i++
			}
			return i, true
		}
	case '>':
		if (*w.grid)[w.robotPos.y][w.robotPos.x+1] == '[' {
			i := 1
			for {
				if (*w.grid)[w.robotPos.y][w.robotPos.x+i*2+1] == '#' {
					return 0, false
				}
				if (*w.grid)[w.robotPos.y][w.robotPos.x+i*2+1] == '.' {
					break
				}
				i++
			}
			return i, true
		}
	default:
		panic("unrecognized move character\n")
	}
	return 0, false
}

// if has crate ahead, check if can push
func (w *Warehouse) canPushVertically(move byte) ([]Coordinate, bool) {
	var collectNodes func(Coordinate, byte) bool
	nodes := []Coordinate{}

	collectNodes = func(c Coordinate, move byte) bool {
		c.v = (*w.grid)[c.y][c.x]
		switch c.v {
		case '.':
			return true
		case '#':
			return false
		case ']':
			if !slices.Contains(nodes, c) {
				nodes = append(nodes, c)
			}
			closing := Coordinate{x: c.x - 1, y: c.y, v: '['}
			if !slices.Contains(nodes, closing) {
				nodes = append(nodes, closing)
			}
			if move == '^' {
				return collectNodes(Coordinate{x: c.x, y: c.y - 1}, '^') && collectNodes(Coordinate{x: c.x - 1, y: c.y - 1}, '^')
			} else if move == 'v' {
				return collectNodes(Coordinate{x: c.x, y: c.y + 1}, 'v') && collectNodes(Coordinate{x: c.x - 1, y: c.y + 1}, 'v')
			}
		case '[':
			if !slices.Contains(nodes, c) {
				nodes = append(nodes, c)
			}
			closing := Coordinate{x: c.x + 1, y: c.y, v: ']'}
			if !slices.Contains(nodes, closing) {
				nodes = append(nodes, closing)
			}
			if move == '^' {
				return collectNodes(Coordinate{x: c.x, y: c.y - 1}, '^') && collectNodes(Coordinate{x: c.x + 1, y: c.y - 1}, '^')
			} else if move == 'v' {
				// fmt.Printf("continue with going v\n")
				return collectNodes(Coordinate{x: c.x, y: c.y + 1}, 'v') && collectNodes(Coordinate{x: c.x + 1, y: c.y + 1}, 'v')
			}
		default:
			panic(fmt.Sprintf("unrecognized character: %c\n", c.v))
		}
		return false
	}

	switch move {
	case '^':
		if collectNodes(Coordinate{x: w.robotPos.x, y: w.robotPos.y - 1}, '^') {
			sort.Sort(ByY(nodes))
			return nodes, true
		}
	case 'v':
		if collectNodes(Coordinate{x: w.robotPos.x, y: w.robotPos.y + 1}, 'v') {
			sort.Sort(ByY(nodes))
			return nodes, true
		}
	default:
		panic("unrecognized move character\n")
	}
	return nil, false
}

func (w *Warehouse) updateRobotPos(xInc, yInc int) {
	(*w.grid)[w.robotPos.y][w.robotPos.x] = '.'
	w.robotPos.x += xInc
	w.robotPos.y += yInc
	(*w.grid)[w.robotPos.y][w.robotPos.x] = '@'
}

func (w *Warehouse) calcGps2() int {
	acc := 0
	for y := range *w.grid {
		for x, cell := range (*w.grid)[y] {
			if cell == '[' {
				acc += 100*y + x
			}
		}
	}
	return acc
}

func (w *Warehouse) calcGps1() int {
	acc := 0
	for y := range *w.grid {
		for x, cell := range (*w.grid)[y] {
			if cell == 'O' {
				acc += 100*y + x
			}
		}
	}
	return acc
}

func (w *Warehouse) update(move byte) {
	switch move {
	case '<':
		// can't walk if next to obstacle
		if (*w.grid)[w.robotPos.y][w.robotPos.x-1] == '#' {
			break
		}

		if (*w.grid)[w.robotPos.y][w.robotPos.x-1] == '.' {
			w.updateRobotPos(-1, 0)
			break
		}

		if crateN, ok := w.canPushHorizontally(move); ok {
			for j := crateN; j > 0; j-- {
				(*w.grid)[w.robotPos.y][w.robotPos.x-j*2-1] = '['
				(*w.grid)[w.robotPos.y][w.robotPos.x-j*2] = ']'
			}
			w.updateRobotPos(-1, 0)
		}

	case '>':
		if (*w.grid)[w.robotPos.y][w.robotPos.x+1] == '#' {
			break
		}

		if (*w.grid)[w.robotPos.y][w.robotPos.x+1] == '.' {
			w.updateRobotPos(1, 0)
			break
		}

		if crateN, ok := w.canPushHorizontally(move); ok {

			for j := crateN; j > 0; j-- {
				(*w.grid)[w.robotPos.y][w.robotPos.x+j*2+1] = ']'
				(*w.grid)[w.robotPos.y][w.robotPos.x+j*2] = '['
			}
			w.updateRobotPos(1, 0)
		}

	case '^':
		if (*w.grid)[w.robotPos.y-1][w.robotPos.x] == '#' {
			break
		}

		if (*w.grid)[w.robotPos.y-1][w.robotPos.x] == '.' {
			w.updateRobotPos(0, -1)
			break
		}

		if nodes, ok := w.canPushVertically(move); ok {
			for j := 0; j < len(nodes); j++ {
				x := nodes[j].x
				y := nodes[j].y
				(*w.grid)[y-1][x] = nodes[j].v
				(*w.grid)[y][x] = '.'
			}
			w.updateRobotPos(0, -1)
		}
	case 'v':
		if (*w.grid)[w.robotPos.y+1][w.robotPos.x] == '#' {
			break
		}

		if (*w.grid)[w.robotPos.y+1][w.robotPos.x] == '.' {
			w.updateRobotPos(0, 1)
			break
		}

		if nodes, ok := w.canPushVertically(move); ok {
			for j := len(nodes) - 1; j >= 0; j-- {
				x := nodes[j].x
				y := nodes[j].y
				(*w.grid)[y+1][x] = nodes[j].v
				(*w.grid)[y][x] = '.'
			}
			w.updateRobotPos(0, 1)
		}
	default:
		panic("unrecognized move character\n")
	}
}

func parseGrid(lines []string) (*Grid, Coordinate) {
	var robotPos Coordinate
	var grid Grid = [][]byte{}

	for y, line := range lines {
		if len(line) == 0 {
			break
		}
		row := []byte{}
		for i, ch := range lines[y] {
			switch ch {
			case '#':
				row = append(row, []byte("##")...)
			case 'O':
				row = append(row, []byte("[]")...)
			case '@':
				row = append(row, []byte("@.")...)
				robotPos = Coordinate{x: i * 2, y: y}
			case '.':
				row = append(row, []byte("..")...)
			default:
				panic(fmt.Sprintf("unrecognized character: %c\n", ch))
			}
		}
		grid = append(grid, row)
	}

	return &grid, robotPos
}

func parseMoves(lines []string) *[]byte {
	var lineIdx int
	moves := []byte{}

	for y, line := range lines {
		if len(line) == 0 {
			lineIdx = y + 1
		}
	}
	for _, line := range lines[lineIdx:] {
		moves = append(moves, []byte(line)...)
	}
	return &moves
}

func (w *Warehouse) write(moveIdx string) {
	filename := fmt.Sprintf("%s.txt", moveIdx)
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

	for i := range *w.grid {
		fmt.Fprintf(writer, "%s\n", string((*w.grid)[i]))
	}
	return
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

func (w *Warehouse) printGrid() {
	fmt.Print("\033[H\033[2J")
	// padding := 2

	var sb strings.Builder
	for y := range *w.grid {
		for _, ch := range (*w.grid)[y] {

			switch ch {
			case '#':
				sb.WriteString("🧱")
			case '[':
				sb.WriteString(fmt.Sprintf("%s", "【"))
			case ']':
				sb.WriteString(fmt.Sprintf("%s", "】"))
				// continue
			case '.':
				// sb.WriteString(fmt.Sprintf("%s", "⬛"))
				sb.WriteString(fmt.Sprintf("%s", ". "))
			case '@':
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
