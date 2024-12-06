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
	grid      *[][]byte
	rows      int
	cols      int
	pos       Coordinate
	Solve     func()
	MoveOk    func() bool
	IsOffGrid func(x, y int) bool
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
	m.rows = len(*m.grid)
	m.cols = len((*m.grid)[0])

	m.IsOffGrid = func(x, y int) bool {
		if x > m.cols-1 || y > m.rows-1 || x < 0 || y < 0 {
			return true
		}
		return false
	}

	m.MoveOk = func() bool {
		switch m.pos.v {
		case '>':
			if m.IsOffGrid(m.pos.x+1, m.pos.y) {
				return false
			}
			if (*m.grid)[m.pos.y][m.pos.x+1] == '.' || (*m.grid)[m.pos.y][m.pos.x+1] == 'X' {
				m.pos.x++
			} else {
				m.pos.v = 'v'
			}
		case '<':
			if m.IsOffGrid(m.pos.x-1, m.pos.y) {
				return false
			}
			if (*m.grid)[m.pos.y][m.pos.x-1] == '.' || (*m.grid)[m.pos.y][m.pos.x-1] == 'X' {
				m.pos.x--
			} else {
				m.pos.v = '^'
			}
		case '^':
			if m.IsOffGrid(m.pos.x, m.pos.y-1) {
				return false
			}
			if (*m.grid)[m.pos.y-1][m.pos.x] == '.' || (*m.grid)[m.pos.y-1][m.pos.x] == 'X' {
				m.pos.y--
			} else {
				m.pos.v = '>'
			}
		case 'v':
			if m.IsOffGrid(m.pos.x, m.pos.y+1) {
				return false
			}
			if (*m.grid)[m.pos.y+1][m.pos.x] == '.' || (*m.grid)[m.pos.y+1][m.pos.x] == 'X' {
				m.pos.y++
			} else {
				m.pos.v = '<'
			}
		default:
			log.Fatal(fmt.Sprintf("invalid guardian position symbol: %c\n", m.pos.v))
		}
		(*m.grid)[m.pos.y][m.pos.x] = 'X'
		return true
	}

	m.Solve = func() {
		for {
			if !m.MoveOk() {
				break
			}
		}
		return
	}

	m.Solve()
	for _, v := range *m.grid {
		fmt.Printf("%s\n", v)
	}
	count_x(m.grid)
	fmt.Printf("part 1 | move count: %d\n", count_x(m.grid))
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

	for i, _ := range lines {
		bytes[i] = []byte(lines[i])
	}
	return bytes, nil
}
