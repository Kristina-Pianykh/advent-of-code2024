package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
)

var (
	// direction: {y, x}
	advance_coordinates map[string][2]int = map[string][2]int{
		"up":    {-1, 0},
		"down":  {1, 0},
		"left":  {0, -1},
		"right": {0, 1},
	}
	DIRECTIONS []string     = []string{"up", "down", "left", "right"}
	HEIGHTS    []Coordinate = []Coordinate{}
	VISITED    []Coordinate = []Coordinate{}
)

type Coordinate struct {
	y int
	x int
}

func (c Coordinate) String() string {
	return fmt.Sprintf("(%d,%d)", c.y, c.x)
}

func (c Coordinate) next(dir string) Coordinate {
	y_inc := advance_coordinates[dir][0]
	x_inc := advance_coordinates[dir][1]
	return Coordinate{x: c.x + x_inc, y: c.y + y_inc}
}

func (c Coordinate) prev(moved_to string) Coordinate {
	y_decr := advance_coordinates[not(moved_to)][0]
	x_decr := advance_coordinates[not(moved_to)][1]
	return Coordinate{x: c.x + x_decr, y: c.y + y_decr}
}

type Grid [][]byte

func (g Grid) val(c Coordinate) byte {
	return g[c.y][c.x]
}

func not(v string) string {
	var opposite string
	switch v {
	case "left":
		opposite = "right"
	case "right":
		opposite = "left"
	case "up":
		opposite = "down"
	case "down":
		opposite = "up"
	default:
		log.Fatalf("unexpected direction %s\n", v)
	}
	return opposite
}

func is_off_grid(c Coordinate, rows, cols int) bool {
	if c.x > cols-1 || c.y > rows-1 || c.x < 0 || c.y < 0 {
		return true
	}
	return false
}

func solve(grid *Grid, rows, cols int) (int, int) {
	heights := make([]Coordinate, 0, rows*cols)

	cum_rating := 0
	for y := range *grid {
		for x := range (*grid)[y] {
			if (*grid)[y][x] == '0' {
				rating := get_rating(grid, Coordinate{x: x, y: y})
				cum_rating += rating
				for _, v := range HEIGHTS {
					heights = append(heights, v)
				}
				HEIGHTS = HEIGHTS[:0]
			}
		}
	}
	return len(heights), cum_rating
}

func main() {
	rows, cols := 60, 60
	grid, err := read_file("input.txt", rows, cols)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	score, rating := solve(&grid, rows, cols)
	fmt.Printf("part 1 | score: %d\n", score)
	fmt.Printf("part 2 | rating: %d\n", rating)
}

func walk(grid *Grid, c Coordinate, dir string, i int) int {
	if is_off_grid(c, len(*grid), len(*grid)) {
		return 0
	}

	if slices.Contains(VISITED, c) {
		return 0
	}

	c_prev := c.prev(dir)
	if grid.val(c_prev)+1 != grid.val(c) {
		return 0
	}

	if grid.val(c) == '9' {
		if !slices.Contains(HEIGHTS, c) {
			HEIGHTS = append(HEIGHTS, c)
		}
		return 1
	}

	VISITED = append(VISITED, c)
	path_count := 0
	for _, next_dir := range DIRECTIONS {
		if next_dir == not(dir) {
			continue
		}
		paths := walk(grid, c.next(next_dir), next_dir, i+1)
		path_count += paths
	}
	VISITED = VISITED[:len(VISITED)-1]
	return path_count
}

func get_rating(grid *Grid, c Coordinate) int {
	rating := 0
	for i, dir := range DIRECTIONS {
		rating += walk(grid, c.next(dir), dir, i)
	}
	return rating
}

func read_file(file_path string, rows, cols int) (Grid, error) {
	file, err := os.Open(file_path)
	if err != nil {
		return nil, errors.New("Error opening file")
	}
	defer file.Close()

	lines := make([][]byte, rows)
	for i := range lines {
		lines[i] = make([]byte, cols)
	}

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		copy(lines[i], line)
		i++
	}

	return lines, nil
}
