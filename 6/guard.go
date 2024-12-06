package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type Coordinate struct {
	x int
	y int
	v byte
}

func (c Coordinate) String() string {
	return fmt.Sprintf("x: %d, y: %d, v: %c\n", c.x, c.y, c.v)
}

type LabMap struct {
	grid *[][]byte
	pos  Coordinate
}

var (
	// direction: {y, x}
	advance_coordinates map[byte][2]int = map[byte][2]int{
		'>': {0, 1},
		'<': {0, -1},
		'^': {-1, 0},
		'v': {1, 0},
	}
	dir_change map[byte]byte = map[byte]byte{
		'>': 'v',
		'<': '^',
		'^': '>',
		'v': '<',
	}
)

func (m *LabMap) no_obstacle(dir byte) bool {
	y_inc := advance_coordinates[dir][0]
	x_inc := advance_coordinates[dir][1]
	return (*m.grid)[m.pos.y+y_inc][m.pos.x+x_inc] == '.' || (*m.grid)[m.pos.y+y_inc][m.pos.x+x_inc] == 'X'
}

func (m *LabMap) move_or_turn(curr_dir byte, dir_grid *[][][]byte) {
	y_inc := advance_coordinates[curr_dir][0]
	x_inc := advance_coordinates[curr_dir][1]
	// fmt.Println("hello")
	if m.no_obstacle(curr_dir) {
		m.pos.y += y_inc
		m.pos.x += x_inc
	} else {
		m.pos.v = dir_change[curr_dir]
		if dir_grid != nil {
			(*dir_grid)[m.pos.y][m.pos.x] = append((*dir_grid)[m.pos.y][m.pos.x], m.pos.v)
		}
	}
}

func (m *LabMap) is_off_grid(x, y int) bool {
	if x > COLS-1 || y > ROWS-1 || x < 0 || y < 0 {
		return true
	}
	return false
}

func (m *LabMap) advance(dir_grid *[][][]byte) {
	switch m.pos.v {
	case '>':
		m.move_or_turn('>', dir_grid)
	case '<':
		m.move_or_turn('<', dir_grid)
	case '^':
		m.move_or_turn('^', dir_grid)
	case 'v':
		m.move_or_turn('v', dir_grid)
	default:
		log.Fatal(fmt.Sprintf("invalid guardian position symbol: %c\n", m.pos.v))
	}
	(*m.grid)[m.pos.y][m.pos.x] = 'X'
}

func (m *LabMap) solve() {
	for {
		x_inc := advance_coordinates[m.pos.v][1]
		y_inc := advance_coordinates[m.pos.v][0]
		if m.is_off_grid(m.pos.x+x_inc, m.pos.y+y_inc) {
			break
		}
		m.advance(nil)
	}
}

func loop_detected(arr *[]byte) bool {
	cnt := map[byte]int{}
	for _, v := range *arr {
		cnt[v]++
		if cnt[v] >= 2 {
			return true
		}
	}
	return false
}

func (m *LabMap) solve2(dir_grid *[][][]byte, start_pos *Coordinate) int {
	loop_cnt := 0

	for y := range *m.grid {
		for x := range (*m.grid)[y] {

			if x == start_pos.x && y == start_pos.y {
				continue
			}

			// set new obstacle and remember the replaced char
			prev_ch := (*m.grid)[y][x]
			(*m.grid)[y][x] = 'O'

			// solve until guard lefts the grid
			// or is stuck in a loop
			for {
				y_inc := advance_coordinates[m.pos.v][0]
				x_inc := advance_coordinates[m.pos.v][1]

				if m.is_off_grid(m.pos.x+x_inc, m.pos.y+y_inc) {
					break
				}

				m.advance(dir_grid)

				if loop_detected(&(*dir_grid)[m.pos.y][m.pos.x]) {
					loop_cnt++
					break
				}
			}

			// restore the replaced char before
			// moving on on the grid, the direction grid
			// and the guard's starting position
			(*m.grid)[y][x] = prev_ch
			reset_dir_grid(dir_grid)
			m.pos = *start_pos
		}
	}

	return loop_cnt
}

func main() {
	content, err := read_file("input.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	guard_pos := find_guard(&content)
	if guard_pos == nil {
		log.Fatal("failed to locate the inital position of guard")
		os.Exit(1)
	}

	m := LabMap{grid: &content, pos: *guard_pos}
	(*m.grid)[m.pos.y][m.pos.x] = 'X'

	//part 1
	m.solve()
	count_x(m.grid)
	fmt.Printf("part 1 | move count: %d\n", count_x(m.grid))

	//part 2
	dir_grid := init_dir_grid(130, 130)
	loop_cnt := m.solve2(dir_grid, guard_pos)
	fmt.Printf("part 2 | loop count: %d\n", loop_cnt)
}

func init_dir_grid(rows, cols int) *[][][]byte {
	dir_grid := make([][][]byte, rows)
	for y := range dir_grid {
		dir_grid[y] = make([][]byte, cols)
		for x := range dir_grid[y] {
			dir_grid[y][x] = make([]byte, 0, 4)
		}
	}
	return &dir_grid
}

func find_guard(m *[][]byte) *Coordinate {
	for y, line := range *m {
		for x, ch := range line {
			if ch == '>' || ch == '<' || ch == '^' || ch == 'v' {
				return &Coordinate{x, y, ch}
			}
		}
	}
	return nil
}

func count_x(grid *[][]byte) int {
	acc := 0
	for _, line := range *grid {
		for _, ch := range line {
			if ch == 'X' {
				acc++
			}
		}
	}
	return acc
}

func read_file(file_path string) ([][]byte, error) {
	file, err := os.Open(file_path)
	if err != nil {
		return nil, errors.New("Error opening file")
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, errors.New("Error getting file stats")
	}

	size := stat.Size()

	content := make([]byte, size)
	_, err = file.Read(content)
	if err != nil {
		return nil, errors.New("Error reading file")
	}

	// TODO: don't convert to string, split at byte \n
	stringified := string(content)
	if strings.HasSuffix(stringified, "\n") {
		stringified = stringified[:len(stringified)-1]
	}
	lines := strings.Split(strings.ReplaceAll(stringified, " ", ""), "\n")
	bytes := make([][]byte, len(lines))

	for i := range lines {
		bytes[i] = []byte(lines[i])
	}
	return bytes, nil
}

func reset_dir_grid(grid *[][][]byte) {
	for y := range *grid {
		for x := range (*grid)[y] {
			(*grid)[y][x] = (*grid)[y][x][:0]
		}
	}
}
