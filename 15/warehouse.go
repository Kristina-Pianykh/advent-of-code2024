package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"sort"
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
func (a ByY) Swap(i, j int)      { a[i], a[j] = a[j], a[i] } // problematic???
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

	for y := range *w.grid {
		for _, cell := range (*w.grid)[y] {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
	for i, move := range *w.moves {
		w.update(i, move)
		if i == len(*w.moves)-1 {
			w.write(fmt.Sprintf("%d", i), '0')
		} else {
			w.write(fmt.Sprintf("%d", i), (*w.moves)[i+1])
		}
	}
	// 1579161 too low
	// 1582688 is correct
	// 1765303 too high
	fmt.Printf("res: %d\n", w.calcGps2())
	// fmt.Printf("%s\n", string(*moves))
	// fmt.Printf("num of moves: %d\n", len(*moves))
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
			// fmt.Printf("move < crates for cells: %d\n", i)
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
		// fmt.Printf("c: %s\n", c.string())
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
			for _, n := range nodes {
				fmt.Printf("%s\n", n.string())
			}
			fmt.Println()
			sort.Sort(ByY(nodes))
			for _, n := range nodes {
				fmt.Printf("%s\n", n.string())
			}
			return nodes, true
		}
	case 'v':
		// fmt.Printf("going v\n")
		if collectNodes(Coordinate{x: w.robotPos.x, y: w.robotPos.y + 1}, 'v') {
			// for _, n := range nodes {
			// 	fmt.Printf("%s\n", n.string())
			// }
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
				// acc += 100*(len(*w.grid)-y-1) + x
				acc += 100*y + x
				// fmt.Printf("O at x=%d, y=%d, their value: %d\n", x, y, 100*(len(*w.grid)-y-1)+x)
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
				// acc += 100*(len(*w.grid)-y-1) + x
				acc += 100*y + x
			}
		}
	}
	return acc
}

func (w *Warehouse) update(moveIdx int, move byte) {
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
			// fmt.Printf("crateN: %d\n", crateN)
			// ..[]@
			// .[]@.
			for j := crateN; j > 0; j-- {
				(*w.grid)[w.robotPos.y][w.robotPos.x-j*2-1] = '['
				(*w.grid)[w.robotPos.y][w.robotPos.x-j*2] = ']'
			}
			// (*w.grid)[w.robotPos.y][w.robotPos.x-crateN-1] = 'O'
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
			// .@[].
			// ..@[]

			for j := crateN; j > 0; j-- {
				(*w.grid)[w.robotPos.y][w.robotPos.x+j*2+1] = ']'
				(*w.grid)[w.robotPos.y][w.robotPos.x+j*2] = '['
			}
			// (*w.grid)[w.robotPos.y][w.robotPos.x+crateN+1] = 'O'
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
			// fmt.Printf("%d\n", len(nodes))
			for j := 0; j < len(nodes); j++ {
				fmt.Printf("shifting cell: %s\n", nodes[j].string())
				x := nodes[j].x
				y := nodes[j].y
				(*w.grid)[y-1][x] = nodes[j].v
				(*w.grid)[y][x] = '.'
				fmt.Printf("x=%d, y=%d, (*w.grid)[y-1][x]: %c\n", x, y, (*w.grid)[y-1][x])
				fmt.Printf("x=%d, y=%d, (*w.grid)[y][x]: %c\n\n", x, y, (*w.grid)[y][x])
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
			// fmt.Printf("row=%d, col=%d, ch=%c\n", y, i, ch)
			switch ch {
			case '#':
				row = append(row, '#')
				row = append(row, '#')
				fmt.Printf("%d\n", len(row))
			case 'O':
				row = append(row, '[')
				row = append(row, ']')
			case '@':
				row = append(row, '@')
				row = append(row, '.')
				robotPos = Coordinate{x: i * 2, y: y}
			case '.':
				row = append(row, '.')
				row = append(row, '.')
			default:
				panic(fmt.Sprintf("unrecognized character: %c\n", ch))
			}
		}
		// fmt.Printf("%s\n", string(row))
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

func (w *Warehouse) write(moveIdx string, move byte) {
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

	fmt.Fprintf(writer, "next move: %c, idx: %s\n", move, moveIdx)
	fmt.Fprintf(writer, "current pos: %s\n", w.robotPos.string())
	for i := range *w.grid {
		// fmt.Fprintf(writer, "%s\n", string(cleanRow(row)))
		// fmt.Fprintf(writer, "%s\n", string(bytes.Trim((*w.grid)[i], "\x00"))) // Remove null bytes
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
