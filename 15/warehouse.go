package main

import (
	"bufio"
	"fmt"
	"os"
)

type Coordinate struct {
	x int
	y int
	v byte
}

func (c Coordinate) string() string {
	return fmt.Sprintf("x: %d, y: %d", c.x, c.y)
}

type Grid [][]byte
type Warehouse struct {
	grid     *Grid
	moves    *[]byte
	robotPos Coordinate
}

func main() {
	w := readLinesFromStream(os.Stdin)
	for y := range *w.grid {
		for _, cell := range (*w.grid)[y] {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
	fmt.Printf("robot start pos: %s\n", w.robotPos.string())

	w.write("00", (*w.moves)[0])
	for i, move := range *w.moves {
		w.update(move)
		if i == len(*w.moves)-1 {
			w.write(fmt.Sprintf("%d", i), '0')
		} else {
			w.write(fmt.Sprintf("%d", i), (*w.moves)[i+1])
		}
		fmt.Printf("updated robot position: %s\n", w.robotPos.string())
	}
	// 1579161 too low
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
		fmt.Printf("c: %s\n", c.string())
		switch c.v {
		case '.':
			return true
		case '#':
			return false
		case ']':
			nodes = append(nodes, c)
			closing := Coordinate{x: c.x - 1, y: c.y, v: '['}
			nodes = append(nodes, closing)
			if move == '^' {
				return collectNodes(Coordinate{x: c.x, y: c.y - 1}, '^') && collectNodes(Coordinate{x: c.x - 1, y: c.y - 1}, '^')
			} else if move == 'v' {
				return collectNodes(Coordinate{x: c.x, y: c.y + 1}, 'v') && collectNodes(Coordinate{x: c.x - 1, y: c.y + 1}, 'v')
			}
		case '[':
			nodes = append(nodes, c)
			closing := Coordinate{x: c.x + 1, y: c.y, v: ']'}
			nodes = append(nodes, closing)
			if move == '^' {
				return collectNodes(Coordinate{x: c.x, y: c.y - 1}, '^') && collectNodes(Coordinate{x: c.x + 1, y: c.y - 1}, '^')
			} else if move == 'v' {
				fmt.Printf("continue with going v\n")
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
			return nodes, true
		}
	case 'v':
		fmt.Printf("going v\n")
		if collectNodes(Coordinate{x: w.robotPos.x, y: w.robotPos.y + 1}, 'v') {
			for _, n := range nodes {
				fmt.Printf("%s\n", n.string())
			}
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
				fmt.Printf("O at x=%d, y=%d, their value: %d\n", x, y, 100*(len(*w.grid)-y-1)+x)
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
			fmt.Printf("crateN: %d\n", crateN)
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
			for j := len(nodes) - 1; j >= 0; j-- {
				fmt.Println("we are here")
				fmt.Printf("shifting cell: %s\n", nodes[j].string())
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

func parseForGrid(line string, lineIdx int, grid *Grid) {
	for i, ch := range line {
		// if ch == '\n' {
		// 	continue
		// }
		// fmt.Printf("i=%d, ch=%c\n", i, ch)
		switch ch {
		case '#':
			(*grid)[lineIdx][i*2] = '#'
			(*grid)[lineIdx][i*2+1] = '#'
		case 'O':
			(*grid)[lineIdx][i*2] = '['
			(*grid)[lineIdx][i*2+1] = ']'
		case '@':
			(*grid)[lineIdx][i*2] = '@'
			(*grid)[lineIdx][i*2+1] = '.'
		case '.':
			(*grid)[lineIdx][i*2] = '.'
			(*grid)[lineIdx][i*2+1] = '.'
		default:
			panic(fmt.Sprintf("unrecognized character: %c\n", ch))
		}
	}
}

func parseForMoves(line string, lineIdx int, grid *Grid) {
	for i, ch := range line {
		(*grid)[lineIdx][i] = byte(ch)
	}
}

func (w *Warehouse) write(moveIdx string, move byte) {
	filename := fmt.Sprintf("%s.txt", moveIdx)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	writer := bufio.NewWriter(f)
	fmt.Fprintf(writer, fmt.Sprintf("next move: %c, idx: %s\n", move, moveIdx))
	fmt.Fprintf(writer, fmt.Sprintf("current pos: %s\n", w.robotPos.string()))
	for y, row := range *w.grid {
		fmt.Fprintf(writer, "%s", string(row))
		if y < len(*w.grid) {
			fmt.Fprintln(writer)
		}
	}
	writer.Flush()
}

func readLinesFromStream(file *os.File) *Warehouse {
	scanner := bufio.NewScanner(file)
	// rows, cols := 8, 8
	// rows, cols := 7, 7
	// rows, cols := 10, 10
	rows, cols := 50, 50

	var robotPos Coordinate
	var grid Grid = make([][]byte, rows)
	for i := range grid {
		grid[i] = make([]byte, cols*2)
	}
	//  0   1   2   3
	// 0 1 2 3 4 5 6 7
	readGrid := false
	moves := []byte{}

	lineIdx := 0
	moveCnt := 0
	for scanner.Scan() {
		line := scanner.Text()

		for i, ch := range line {
			if ch == '@' {
				robotPos = Coordinate{x: i * 2, y: lineIdx}
			}
		}
		// fmt.Printf("len(line)=%d\n", len(line))

		if lineIdx == rows {
			readGrid = true
		}

		if readGrid && len(line) > 1 {
			// fmt.Printf("we are here")
			moves = append(moves, []byte(line)...)
			moveCnt += len(line)
		} else if !readGrid && len(line) > 1 && lineIdx < rows {
			fmt.Printf("lineIdx = %d\n", lineIdx)
			parseForGrid(line, lineIdx, &grid)
		}
		lineIdx++
	}
	// fmt.Printf("number of chars: %d\n", moveCnt)
	return &Warehouse{grid: &grid, moves: &moves, robotPos: robotPos}
}
