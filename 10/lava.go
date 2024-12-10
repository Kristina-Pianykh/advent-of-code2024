package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

// const ROWS = 7
// const COLS = 7

type Coordinate struct {
	y int
	x int
}

func (c Coordinate) String() string {
	return fmt.Sprintf("(%d,%d)", c.y, c.x)
}

func (c Coordinate) next(dir string) Coordinate {
	y_inc := ADVANCED_COORDINATES[dir][0]
	x_inc := ADVANCED_COORDINATES[dir][1]
	return Coordinate{x: c.x + x_inc, y: c.y + y_inc}
}

func (c Coordinate) prev(moved_to string) Coordinate {
	y_decr := ADVANCED_COORDINATES[not(moved_to)][0]
	x_decr := ADVANCED_COORDINATES[not(moved_to)][1]
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

var (
	// direction: {y, x}
	ADVANCED_COORDINATES map[string][2]int = map[string][2]int{
		"up":    {-1, 0},
		"down":  {1, 0},
		"left":  {0, -1},
		"right": {0, 1},
	}
	DIRECTIONS []string     = []string{"up", "down", "left", "right"}
	HEIGHTS    []Coordinate = []Coordinate{}
	VISITED    []Coordinate = []Coordinate{}
)

func is_off_grid(c Coordinate, rows, cols int) bool {
	if c.x > cols-1 || c.y > rows-1 || c.x < 0 || c.y < 0 {
		return true
	}
	return false
}

func solve1(grid *Grid) int {
	score_sum := 0
	for y := range *grid {
		for x := range (*grid)[y] {
			if (*grid)[y][x] == '0' {
				fmt.Printf("trailhead at y=%d, x=%d\n", y, x)

				score := get_score(grid, Coordinate{x: x, y: y})
				HEIGHTS = HEIGHTS[:0]
				fmt.Printf("score for trailhead at y=%d, x=%d: %d\n", y, x, score)
				score_sum += score
			}
		}
	}
	return score_sum
}

func solve2(grid *Grid) int {
	cum_rating := 0
	for y := range *grid {
		for x := range (*grid)[y] {
			if (*grid)[y][x] == '0' {
				fmt.Printf("trailhead at y=%d, x=%d\n", y, x)

				rating := get_rating(grid, Coordinate{x: x, y: y})
				fmt.Printf("score for trailhead at y=%d, x=%d: %d\n", y, x, rating)
				cum_rating += rating
			}
		}
	}
	return cum_rating
}

func main() {
	rows, cols := 60, 60
	grid, err := read_file("input.txt", rows, cols)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("part 1 | score: %d\n", solve1(&grid))
	HEIGHTS = HEIGHTS[:0]
	fmt.Printf("part 2 | rating: %d\n", solve2(&grid))
}

func walk(grid *Grid, c Coordinate, dir string, i int) int {
	if is_off_grid(c, len(*grid), len(*grid)) {
		fmt.Printf("%s(%d,%d) off grid. Return 0\n", strings.Repeat(" ", i), c.y, c.x)
		fmt.Printf("%sVISITED=%v\n", strings.Repeat(" ", i), VISITED)
		return 0
	}

	if slices.Contains(VISITED, c) {
		fmt.Printf("%salready visited %s. Return 0\n", strings.Repeat(" ", i), c)
		return 0
	}

	c_prev := c.prev(dir)
	fmt.Printf("%s%s: %c -> moved %s -> %s: %c\n", strings.Repeat(" ", i), c_prev, grid.val(c_prev), dir, c, grid.val(c))
	if grid.val(c_prev)+1 != grid.val(c) {
		fmt.Printf("%s%c + 1 != %c. Returning 0\n", strings.Repeat(" ", i), grid.val(c_prev), grid.val(c))
		fmt.Printf("%sVISITED=%v\n", strings.Repeat(" ", i), VISITED)
		return 0
	}

	if grid.val(c) == '9' {
		fmt.Printf("%sVISITED=%v\n", strings.Repeat(" ", i), VISITED)
		fmt.Printf("%sreached end of trail at %s: %c. Return 1\n\n", strings.Repeat(" ", i), c, grid.val(c))
		if !slices.Contains(HEIGHTS, c) {
			HEIGHTS = append(HEIGHTS, c)
			fmt.Printf("%sappended %s: %v\n", strings.Repeat(" ", i), c, HEIGHTS)
		}
		return 1
	}

	VISITED = append(VISITED, c)
	fmt.Printf("%supdated VISITED=%v\n", strings.Repeat(" ", i), VISITED)
	path_count := 0
	for _, next_dir := range DIRECTIONS {
		if next_dir == not(dir) {
			fmt.Printf("%sskipping %s for coming from %s\n\n", strings.Repeat(" ", i), not(dir), dir)
			continue
		}
		fmt.Printf("%s%s: %c -> move %s\n", strings.Repeat(" ", i), c, grid.val(c), next_dir)
		paths := walk(grid, c.next(next_dir), next_dir, i+1)
		path_count += paths
		fmt.Printf("%spath_count+=%d\n", strings.Repeat(" ", i), paths)
	}
	fmt.Printf("%spath count: %d\n\n", strings.Repeat(" ", i), path_count)
	VISITED = VISITED[:len(VISITED)-1]
	return path_count
}

func get_score(grid *Grid, c Coordinate) int {
	score := 0
	for i, dir := range DIRECTIONS {
		c_next := c.next(dir)
		fmt.Printf("TRAILHEAD: move %s to %s\n", dir, c_next)
		if paths := walk(grid, c_next, dir, i); paths > 0 {
			score += paths
		}
	}
	return len(HEIGHTS)
	// return score
}

func get_rating(grid *Grid, c Coordinate) int {
	rating := 0
	for i, dir := range DIRECTIONS {
		c_next := c.next(dir)
		fmt.Printf("TRAILHEAD: move %s to %s\n", dir, c_next)
		if paths := walk(grid, c_next, dir, i); paths > 0 {
			rating += paths
		}
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
