package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Coordinate struct {
	x int
	y int
}

func (c Coordinate) string() string {
	return fmt.Sprintf("x: %d, y: %d", c.x, c.y)
}

type Robot struct {
	startPos Coordinate
	velocity Coordinate
}

func (r Robot) string() string {
	return fmt.Sprintf("startPost={%s}, velocity={%s}", r.startPos.string(), r.velocity.string())
}

type Grid [][]int

func (g Grid) string() string {
	var sb strings.Builder

	for _, row := range g {
		for _, cell := range row {
			if cell == 0 {
				sb.WriteByte('.')
			} else {
				sb.WriteString(fmt.Sprintf("%d", cell))
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	start := time.Now()
	rows := 103
	cols := 101
	grid := initSparseGrid(rows, cols)
	robots := readLinesFromStream(os.Stdin)

	for i := 100; i < 10000; i++ {
		for _, robot := range robots {
			finalPos := calcPosition(robot, rows, cols, i)
			grid[finalPos.y][finalPos.x]++
		}
		gridStr := grid.string()
		if containsLinedupRobots(gridStr) {
			// writeGrid(gridStr, fmt.Sprintf("tree.txt"))
			fmt.Printf("part 2 | found Xmas tree after %d sec\n", i)
			break
		}
		if i == 100 {
			fmt.Printf("part 1 | safety factor: %d\n", calcSafetyFactor(&grid, rows, cols))
		}
		resetGrid(&grid)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		robot := parse(line)
		finalPos := calcPosition(robot, rows, cols, 100)
		grid[finalPos.y][finalPos.x]++
	}

	fmt.Printf("took: %v\n", time.Now().Sub(start))
}

func calcSafetyFactor(grid *Grid, rows, cols int) int {
	var acc int
	safetyFactor := 1
	acc = 0
	for _, row := range (*grid)[:rows/2] {
		for _, cell := range row[:cols/2] {
			acc += cell
		}
	}
	safetyFactor *= acc

	acc = 0
	for _, row := range (*grid)[rows/2+1:] {
		for _, cell := range row[:cols/2] {
			acc += cell
		}
	}
	safetyFactor *= acc

	acc = 0
	for _, row := range (*grid)[rows/2+1:] {
		for _, cell := range row[cols/2+1:] {
			acc += cell
		}
	}
	safetyFactor *= acc

	acc = 0
	for _, row := range (*grid)[:rows/2] {
		for _, cell := range row[cols/2+1:] {
			acc += cell
		}
	}

	return safetyFactor * acc
}

func parse(line string) Robot {
	var err error
	var startX, startY, velocityX, velocityY int
	robot := Robot{}

	re := regexp.MustCompile(`p=(-?[[:digit:]]+),(-?[[:digit:]]+) v=(-?[[:digit:]]+),(-?[[:digit:]]+)`)
	sX, sY, vX, vY := re.FindStringSubmatch(line)[1], re.FindStringSubmatch(line)[2], re.FindStringSubmatch(line)[3], re.FindStringSubmatch(line)[4]
	startX, err = strconv.Atoi(sX)
	if err != nil {
		panic(err)
	}
	startY, err = strconv.Atoi(sY)
	if err != nil {
		panic(err)
	}
	robot.startPos.x = startX
	robot.startPos.y = startY
	velocityX, err = strconv.Atoi(vX)
	if err != nil {
		panic(err)
	}
	velocityY, err = strconv.Atoi(vY)
	if err != nil {
		panic(err)
	}

	robot.velocity.x = velocityX
	robot.velocity.y = velocityY
	return robot
}

func containsLinedupRobots(grid string) bool {
	limit := 10

	re := regexp.MustCompile(`[^\.]+`)
	robotSeq := re.FindAll([]byte(grid), -1)
	for _, seq := range robotSeq {
		if len(seq) >= limit {
			return true
		}
	}
	return false
}

func calcPosition(r Robot, rows, cols, sec int) Coordinate {
	c := Coordinate{}
	c.x = mod(r.startPos.x+r.velocity.x*sec, cols)
	c.y = mod(r.startPos.y+r.velocity.y*sec, rows)
	return c
}

func initSparseGrid(rows, cols int) Grid {
	grid := make([][]int, rows)
	for i := range grid {
		grid[i] = make([]int, cols)
	}
	return grid
}

func mod(a, b int) int {
	return (a%b + b) % b
}

func readLinesFromStream(file *os.File) []Robot {
	scanner := bufio.NewScanner(file)
	var robots []Robot = []Robot{}

	for scanner.Scan() {
		line := scanner.Text()
		robot := parse(line)
		robots = append(robots, robot)
	}
	return robots
}

func resetGrid(grid *Grid) {
	for y := range *grid {
		for x := range (*grid)[y] {
			(*grid)[y][x] = 0
		}
	}
}

func writeGrid(grid string, filepath string) {
	file, err := os.Create(filepath)
	check(err)
	defer file.Close()

	_, err = file.WriteString(grid)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
