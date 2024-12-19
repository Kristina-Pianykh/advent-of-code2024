package main

import (
	"bufio"
	"fmt"
	"os"
)

type Coordinate struct {
	x int
	y int
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

	w.write(00)
	for _, move := range *w.moves {
		w.update(move)
		// w.write(i)
	}
	fmt.Printf("res: %d\n", w.calcGps())
	// fmt.Printf("%s\n", string(*moves))
	// fmt.Printf("num of moves: %d\n", len(*moves))
}

// if has crate ahead, check if can push
func (w *Warehouse) canPush(move byte) (int, bool) {
	switch move {
	case '<':
		if (*w.grid)[w.robotPos.y][w.robotPos.x-1] == ']' {
			i := 1
			for {
				if (*w.grid)[w.robotPos.y][(w.robotPos.x-i)*2] == '#' {
					return 0, false
				}
				if (*w.grid)[w.robotPos.y][(w.robotPos.x-i)*2] == '.' {
					i--
					break
				}
				i++
			}
			return i, true
		}
	case '>':
		if (*w.grid)[w.robotPos.y][(w.robotPos.x+1)*2] == '[' {
			i := 1
			for {
				if (*w.grid)[w.robotPos.y][(w.robotPos.x+i)*2] == '#' {
					return 0, false
				}
				if (*w.grid)[w.robotPos.y][(w.robotPos.x+i)*2] == '.' {
					i--
					break
				}
				i++
			}
			return i, true
		}
	case '^':
		if (*w.grid)[w.robotPos.y-1][w.robotPos.x] == ']' || (*w.grid)[w.robotPos.y-1][w.robotPos.x] == '[' {
			i := 1
			for {
				if (*w.grid)[w.robotPos.y-i][w.robotPos.x] == '#' {
					return 0, false
				}
				if (*w.grid)[w.robotPos.y-i][w.robotPos.x] == '.' {
					i--
					break
				}
				i++
			}
			return i, true
		}
	case 'v':
		if (*w.grid)[w.robotPos.y+1][w.robotPos.x] == ']' || (*w.grid)[w.robotPos.y+1][w.robotPos.x] == '[' {
			i := 1
			for {
				if (*w.grid)[w.robotPos.y+i][w.robotPos.x] == '#' {
					return 0, false
				}
				if (*w.grid)[w.robotPos.y+i][w.robotPos.x] == '.' {
					i--
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

func (w *Warehouse) updateRobotPos(xInc, yInc int) {
	(*w.grid)[w.robotPos.y][w.robotPos.x] = '.'
	w.robotPos.x += xInc
	w.robotPos.y += yInc
	(*w.grid)[w.robotPos.y][w.robotPos.x] = '@'
}

func (w *Warehouse) calcGps() int {
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

		if crateN, ok := w.canPush(move); ok {
			// ..[]@
			// .[]@.
			(*w.grid)[w.robotPos.y][w.robotPos.x-crateN*2-1] = '['
			(*w.grid)[w.robotPos.y][w.robotPos.x-crateN*2] = ']'
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

		if crateN, ok := w.canPush(move); ok {
			// .@[].
			// ..@[]
			(*w.grid)[w.robotPos.y][w.robotPos.x+crateN*2] = '['
			(*w.grid)[w.robotPos.y][w.robotPos.x+crateN*2+1] = ']'
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

		if crateN, ok := w.canPush(move); ok {
			(*w.grid)[w.robotPos.y-crateN-1][w.robotPos.x] = 'O'
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

		if crateN, ok := w.canPush(move); ok {
			(*w.grid)[w.robotPos.y+crateN+1][w.robotPos.x] = 'O'
			w.updateRobotPos(0, 1)
		}
	default:
		panic("unrecognized move character\n")
	}
}

func (c Coordinate) isAdjacent(co Coordinate) bool {
	incLayout := [][]int{
		{0, -1},
		{1, 0},
		{0, 1},
		{-1, 0},
	}
	for i := 0; i < len(incLayout); i++ {
		if c.x == co.x+incLayout[i][0] && c.y == co.y+incLayout[i][1] {
			return true
		}
	}
	return false
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

func (w *Warehouse) write(moveIdx int) {
	filename := fmt.Sprintf("%d.txt", moveIdx)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	writer := bufio.NewWriter(f)
	// fmt.Fprintf(writer, fmt.Sprintf("next move: %s, idx: %d\n", string((*w.moves)[i+1]), i+1))
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
	rows, cols := 7, 7
	// rows, cols := 10, 10
	// rows, cols := 50, 50

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
				robotPos = Coordinate{x: i, y: lineIdx}
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
