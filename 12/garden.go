package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
)

// Optional is type wrapper for any value
type Optional[T any] struct {
	value T
	some  bool
}

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
	// return fmt.Sprintf("plotType=%c, area=%d, perimeter=%d, price=%d, len(plots)=%d, plots=%v\n", r.plotType, r.area, r.perimeter, r.price, len(r.plots), r.plots)
	return fmt.Sprintf("plotType=%c, area=%d, perimeter=%d, sides=%d, len(plots)=%d\n", r.plotType, r.area, r.perimeter, r.sides, len(r.plots))
}

// ?? pass region by pointer? would it work otherwise?
func (r *Region) calcArea() {
	(*r).area = len((*r).plots)
}

// ?? pass region by pointer? would it work otherwise?
func (r *Region) calcPerimeter(grid *Grid) {
	acc := 0
	incLayout := [][]int{
		{0, -1},
		{1, 0},
		{0, 1},
		{-1, 0},
	}
	for _, c := range (*r).plots {
		// fmt.Printf("%s\n", c.string())
		for j := 0; j < 4; j++ {
			x := c.x + incLayout[j][0]
			y := c.y + incLayout[j][1]
			ch, err := (*grid).getVal(Coordinate{x: x, y: y})

			if err == nil {
				if ch != (*r).plotType {
					acc++
					// fmt.Printf("inc perimeter for %c at (%d, %d); acc=%d\n", c.v, x, y, acc)
				}
			} else {
				acc++
				// fmt.Printf("out of bounds for (%d, %d); acc=%d\n", x, y, acc)
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
		// fmt.Printf("region %c has %d plots\n", r.plotType, len(r.plots))
		r.calcArea()
		r.calcPerimeter(grid)
		// fmt.Printf("%s\n", r.string())
		acc += r.area * r.perimeter
	}
	return acc
}

func solve2(grid *Grid, regions *[]Region) int {
	acc := 0
	for _, r := range *regions {
		// fmt.Printf("region %c has %d plots\n", r.plotType, len(r.plots))
		r.calcArea()
		r.calcSides(grid)
		fmt.Printf("%s\n", r.string())
		acc += r.area * r.sides
	}
	return acc
}

func main() {
	rows, cols := 140, 140
	grid, err := readFile("input.txt", rows, cols)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	regions := walkGrid(&grid, rows, cols)
	// fmt.Printf("found %d regions\n", len(regions))

	fmt.Printf("part 1 | price: %d\n", solve1(&grid, &regions))
	fmt.Printf("part 2 | price: %d\n", solve2(&grid, &regions))
	//98311 too low
	//860382 too high
}

func (c Coordinate) getNeighbors(grid *Grid) []Coordinate {
	if (*grid).isOffGrid(c) {
		log.Fatalf("%s is off grid\n", c.string())
	}

	fmt.Printf("search neighbors for %s\n", c.string())
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
			fmt.Printf("%s is off grid\n", potentialNeighbor.string())
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
		fmt.Printf("toVisit: %d\n", toVisit)
		for len(toVisit) > 0 {
			c := toVisit[0]
			toVisit = remove(toVisit, c)

			fmt.Printf("current node: %s\n", c.string())

			if slices.Contains(visited, c) {
				fmt.Printf("has already been visited, goto start\n")
				continue
				// return
			}
			visited = append(visited, c)

			neighbors := c.getNeighbors(grid)
			fmt.Printf("neighbors: %v\n", neighbors)
			for _, neighbor := range neighbors {
				if !slices.Contains(visited, neighbor) {
					if !slices.Contains(toVisit, neighbor) {
						toVisit = append(toVisit, neighbor)
					}
				}
			}
			fmt.Printf("toVisit: %d\n", toVisit)

			localOutsiders := c.getOutsiders(grid)
			fmt.Printf("sides before comparison: %d\n", sides)
			fmt.Printf("len(outsiders) before: %d\n", len(outsiders))
			fmt.Printf("localOutsiders: %v\n", localOutsiders)
			for _, o := range localOutsiders {

				if o.incLayoutIdx == nil {
					log.Fatalf("incLayoutIdx is not set for %s\n", o.string())
				}

				for _, out := range outsiders {

					if o == out && *o.incLayoutIdx != *out.incLayoutIdx {
						fmt.Printf("%s : new inner boder\n", o.string())
						sides++
						// break
						goto nextLocalOutsider
					}

					if o.isAdjacent(out) && *o.incLayoutIdx == *out.incLayoutIdx {
						fmt.Printf("%s is adjacent to %s\n", o.string(), out.string())
						// continue
						goto nextLocalOutsider
					}
				}
				sides++ // if the given local outsider is not adjacent to any found outsiders from the same side
				outsiders = append(outsiders, o)
			nextLocalOutsider:
			}
			fmt.Printf("sides after comparison: %d\n", sides)

			// merge outsiders with localOutsiders but avoid duplicates
			outsiders = addDistinct(outsiders, localOutsiders)
			fmt.Printf("len(outsiders) after: %d\n\n", len(outsiders))
		}
	}

	toVisit = append(toVisit, (*r).plots[0])
	fmt.Printf("first toVisit: %d\n", toVisit)
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

	// might potentially cause trouble
	// might make sense to return a copy?
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
			// fmt.Printf("%s\n", region.string())
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

func readFile(file_path string, rows, cols int) (Grid, error) {
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
