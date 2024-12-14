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

type Coordinate struct {
	x            int
	y            int
	v            byte
	incLayoutIdx *int
}

func (c Coordinate) string() string {
	return fmt.Sprintf("x: %d, y: %d, v: %c", c.x, c.y, c.v)
}

type Region struct {
	plots     []Coordinate
	plotType  byte
	area      int
	perimeter int
	sides     int
	price     int
}

func (r Region) string() string {
	return fmt.Sprintf("plotType=%c, area=%d, perimeter=%d, sides=%d, len(plots)=%d", r.plotType, r.area, r.perimeter, r.sides, len(r.plots))
}

func (r *Region) calcArea() {
	(*r).area = len((*r).plots)
}

func (r *Region) calcPerimeter(grid *Grid) {
	acc := 0
	incLayout := [][]int{
		{0, -1},
		{1, 0},
		{0, 1},
		{-1, 0},
	}
	for _, c := range (*r).plots {
		for j := 0; j < 4; j++ {
			x := c.x + incLayout[j][0]
			y := c.y + incLayout[j][1]
			ch, err := (*grid).getVal(Coordinate{x: x, y: y})

			if err == nil {
				if ch != (*r).plotType {
					acc++
				}
			} else {
				acc++
			}

		}
	}
	(*r).perimeter = acc
}

func (r *Region) calcPrice() {
	(*r).price = (*r).area * (*r).perimeter
}

type Grid [][]byte
type VisitedGrid [][]bool

func (g *Grid) isOffGrid(c Coordinate) bool {
	rows := len(*g)
	cols := len((*g)[0])
	if c.x > cols-1 || c.y > rows-1 || c.x < 0 || c.y < 0 {
		return true
	}
	return false
}

func (g *Grid) getVal(c Coordinate) (byte, error) {
	if (*g).isOffGrid(c) {
		return 0, errors.New(fmt.Sprintf("%v is off grid\n", c))
	}
	return (*g)[c.y][c.x], nil
}

func solve1(grid *Grid, regions *[]Region) int {
	acc := 0
	for _, r := range *regions {
		r.calcArea()
		r.calcPerimeter(grid)
		acc += r.area * r.perimeter
	}
	return acc
}

func solve2(grid *Grid, regions *[]Region) int {
	acc := 0
	for _, r := range *regions {
		r.calcArea()
		r.calcSides(grid)
		fmt.Printf("%s\n", r.string())
		acc += r.area * r.sides
	}
	return acc
}

func main() {
	grid := readLinesFromStream(os.Stdin)
	rows := len(grid)
	cols := len(grid[0])
	regions := walkGrid(&grid, rows, cols)
	fmt.Printf("part 1 | price: %d\n", solve1(&grid, &regions))
	fmt.Printf("part 2 | price: %d\n", solve2(&grid, &regions))
}

func (c Coordinate) getNeighbors(grid *Grid) []Coordinate {
	if (*grid).isOffGrid(c) {
		log.Fatalf("%s is off grid\n", c.string())
	}

	incLayout := [][]int{
		{0, -1},
		{1, 0},
		{0, 1},
		{-1, 0},
	}
	neighbors := []Coordinate{}
	for i := 0; i < 4; i++ {
		potentialNeighbor := Coordinate{x: c.x + incLayout[i][0], y: c.y + incLayout[i][1]}
		ch, err := (*grid).getVal(potentialNeighbor)
		if err != nil {
			continue
		}
		if ch == c.v {
			potentialNeighbor.v = ch
			neighbors = append(neighbors, potentialNeighbor)
		}
	}
	return neighbors
}

func (c Coordinate) getOutsiders(grid *Grid) []Coordinate {
	if (*grid).isOffGrid(c) {
		return nil
	}

	incLayout := [][]int{
		{0, -1},
		{1, 0},
		{0, 1},
		{-1, 0},
	}
	outsiders := []Coordinate{}
	for i := 0; i < 4; i++ {
		potentialOutsider := Coordinate{x: c.x + incLayout[i][0], y: c.y + incLayout[i][1]}
		ch, err := (*grid).getVal(Coordinate{x: c.x + incLayout[i][0], y: c.y + incLayout[i][1]})

		if err != nil || ch != c.v {
			potentialOutsider.v = ch
			potentialOutsider.incLayoutIdx = &i
			outsiders = append(outsiders, potentialOutsider)
		}
	}
	return outsiders
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

func (r *Region) calcSides(grid *Grid) {
	visited := []Coordinate{}
	toVisit := []Coordinate{}
	outsiders := []Coordinate{}
	sides := 0
	var bfs func()

	bfs = func() {
		for len(toVisit) > 0 {
			c := toVisit[0]
			toVisit = remove(toVisit, c)

			if slices.Contains(visited, c) {
				continue
			}
			visited = append(visited, c)

			neighbors := c.getNeighbors(grid)
			for _, neighbor := range neighbors {
				if !slices.Contains(visited, neighbor) {
					if !slices.Contains(toVisit, neighbor) {
						toVisit = append(toVisit, neighbor)
					}
				}
			}

			localOutsiders := c.getOutsiders(grid)

			for _, o := range localOutsiders {
				if o.incLayoutIdx == nil {
					log.Fatalf("incLayoutIdx is not set for %s\n", o.string())
				}
				for _, out := range outsiders {
					if o.isAdjacent(out) && *o.incLayoutIdx == *out.incLayoutIdx {
						goto nextLocalOutsider
					}
				}

				sides++ // if the given local outsider is not adjacent to any found outsiders from the same side
			nextLocalOutsider:
			}

			// merge outsiders with localOutsiders but avoid duplicates
			outsiders = addDistinct(outsiders, localOutsiders)
		}
	}

	toVisit = append(toVisit, (*r).plots[0])
	bfs()
	(*r).sides = sides
}

func addDistinct(dst []Coordinate, src []Coordinate) []Coordinate {
	for _, o := range src {
		if !slices.Contains(dst, o) {
			dst = append(dst, o)
		}
	}
	return dst
}

func remove(arr []Coordinate, toRemove Coordinate) []Coordinate {
	if len(arr) == 0 {
		return arr
	}
	if !slices.Contains(arr, toRemove) {
		return arr
	}

	var new_arr []Coordinate = make([]Coordinate, len(arr)-1)
	idx := 0
	for _, c := range arr {
		if c == toRemove {
			continue
		}
		new_arr[idx] = c
		idx++
	}
	return new_arr
}

func walkGrid(grid *Grid, rows, cols int) []Region {
	visited := initVisitedGrid(rows, cols)
	regions := []Region{}
	var region Region

	var findRegion func(grid *Grid, c Coordinate, region *Region)
	findRegion = func(grid *Grid, c Coordinate, region *Region) {
		cell, err := (*grid).getVal(c)
		if err != nil {
			return
		}

		if (*region).plotType != cell {
			return
		}

		if visited[c.y][c.x] {
			return
		}

		visited[c.y][c.x] = true
		c.v = cell
		region.plots = append(region.plots, c)
		findRegion(grid, Coordinate{x: c.x, y: c.y - 1}, region) // up
		findRegion(grid, Coordinate{x: c.x + 1, y: c.y}, region) // right
		findRegion(grid, Coordinate{x: c.x, y: c.y + 1}, region) // down
		findRegion(grid, Coordinate{x: c.x - 1, y: c.y}, region) // left
	}

	for y := range *grid {
		for x, cell := range (*grid)[y] {
			if !visited[y][x] {
				region = Region{plots: []Coordinate{}, plotType: cell}
				findRegion(grid, Coordinate{x: x, y: y, v: cell}, &region)
				regions = append(regions, region)
			}
			visited[y][x] = true
			fmt.Printf("%s\n", region.string())
		}
	}
	return regions
}

func initVisitedGrid(rows, cols int) VisitedGrid {
	dir_grid := make([][]bool, rows)
	for y := range dir_grid {
		dir_grid[y] = make([]bool, cols)
	}
	return dir_grid
}

func readLinesFromStream(file *os.File) Grid {
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	var byteArr [][]byte
	for y := 0; y < len(lines); y++ {
		lines[y] = strings.TrimSpace(lines[y])
		byteArr = append(byteArr, []byte(lines[y]))
	}
	return byteArr
}
