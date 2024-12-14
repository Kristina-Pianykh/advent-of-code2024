package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestInitVisitedGrid(t *testing.T) {
	rows, cols := 10, 10
	visitedGrid := initVisitedGrid(rows, cols)
	for _, row := range visitedGrid {
		fmt.Printf("%v\n", row)
	}
}

func TestOffGrid(t *testing.T) {
	rows, cols := 10, 10
	grid, err := readFile("test.txt", rows, cols)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", (&grid).isOffGrid(Coordinate{x: 1, y: 10}))
	fmt.Printf("%v\n", (&grid).isOffGrid(Coordinate{x: 10, y: 0}))
	fmt.Printf("%v\n", (&grid).isOffGrid(Coordinate{x: 0, y: 3}))
	fmt.Printf("%v\n", (&grid).isOffGrid(Coordinate{x: 3, y: 1}))
	fmt.Printf("%v\n", (&grid).isOffGrid(Coordinate{x: -1, y: 8}))
	fmt.Printf("%v\n", (&grid).isOffGrid(Coordinate{x: 1, y: -1}))
}

func TestSolve1(t *testing.T) {
	rows, cols := 5, 5
	grid, err := readFile("test.txt", rows, cols)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	regions := walkGrid(&grid, rows, cols)
	res := solve1(&grid, &regions)
	if res != 772 {
		log.Fatalf("Test failed with res %d\n", res)
	}
	fmt.Printf("TEST part 1 | price: %d\n", solve1(&grid, &regions))
}

// func TestCountSides(t *testing.T) {
// 	rows, cols := 5, 5
// 	grid, err := readFile("test.txt", rows, cols)
// 	if err != nil {
// 		log.Fatal(err)
// 		os.Exit(1)
// 	}
// 	// VISITED = initVisitedGrid(rows, cols)
// 	regions := walkGrid(&grid, rows, cols)
// 	res := solve1(&grid, &regions)
// 	if res != 772 {
// 		log.Fatalf("Test failed with res %d\n", res)
// 	}
// 	fmt.Printf("TEST part 1 | price: %d\n", solve1(&grid, &regions))
// }

func TestGetNeighbors(t *testing.T) {
	rows, cols := 5, 5
	grid, err := readFile("test.txt", rows, cols)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	points := []Coordinate{
		{x: 0, y: 0},
		{x: 2, y: 5},
		{x: 1, y: 0},
		{x: -1, y: 0},
		{x: 8, y: 3},
		{x: 3, y: -1},
	}
	fmt.Println(len(grid))
	fmt.Println(len(grid[0]))

	var neighbors []Coordinate
	for _, p := range points {
		fmt.Printf("neighbors for point %s\n", p.string())
		if grid.isOffGrid(p) {
			continue
		}
		neighbors = p.getNeighbors(&grid)
		for _, n := range neighbors {
			fmt.Printf("  %s | ", n.string())
		}
		fmt.Printf("\n\n")
	}
}

func TestRemove(t *testing.T) {
	points := []Coordinate{
		{x: 0, y: 0},
		{x: 2, y: 5},
		{x: 1, y: 0},
		{x: -1, y: 0},
		{x: 8, y: 3},
		{x: 3, y: -1},
	}
	fmt.Printf("original array: %v\n", points)
	for i, p := range points {
		fmt.Printf("new arr for %s: %v\n", p.string(), remove(points, points[i]))
	}
	points2 := []Coordinate{
		{x: 0, y: 0},
	}
	fmt.Printf("new arr for %s: %v\n", points[1].string(), remove(points2, points[1]))
	fmt.Printf("new arr for %s: %v\n", points[0].string(), remove(points2, points[0]))

	// edge case
	points1 := []Coordinate{}
	fmt.Printf("new arr for %s: %v\n", points[0].string(), remove(points1, points[0]))
}

func TestGetOutsiders(t *testing.T) {
	rows, cols := 5, 5
	grid, err := readFile("test.txt", rows, cols)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	c := Coordinate{x: 0, y: 0, v: 'O'}
	fmt.Printf("outsiders for %s: %v\n", c.string(), c.getOutsiders(&grid))
	c = Coordinate{x: 1, y: 1, v: 'X'}
	fmt.Printf("outsiders for %s: %v\n", c.string(), c.getOutsiders(&grid))
	c = Coordinate{x: 1, y: 2, v: 'O'}
	fmt.Printf("outsiders for %s: %v\n", c.string(), c.getOutsiders(&grid))
	c = Coordinate{x: 4, y: 3, v: 'O'}
	fmt.Printf("outsiders for %s: %v\n", c.string(), c.getOutsiders(&grid))
}

func TestIsAdjacent(t *testing.T) {
	c := Coordinate{x: 0, y: 0, v: 'O'}
	co := Coordinate{x: 1, y: 0, v: 'O'}
	fmt.Printf("%s is adjacent to %s: %v\n", c.string(), co.string(), co.isAdjacent(c))
	c = Coordinate{x: 0, y: 0, v: 'O'}
	co = Coordinate{x: -1, y: 0, v: 'O'}
	fmt.Printf("%s is adjacent to %s: %v\n", c.string(), co.string(), co.isAdjacent(c))
	c = Coordinate{x: 0, y: 0, v: 'O'}
	co = Coordinate{x: 2, y: 0, v: 'O'}
	fmt.Printf("%s is adjacent to %s: %v\n", c.string(), co.string(), co.isAdjacent(c))
	c = Coordinate{x: 0, y: 0, v: 'O'}
	co = c
	fmt.Printf("%s is adjacent to %s: %v\n", c.string(), co.string(), co.isAdjacent(c))
}

func TestSolve2(t *testing.T) {
	var (
		rows, cols int
		grid       Grid
		err        error
		regions    []Region
		res        int
		expected   int
	)

	rows, cols = 5, 5
	grid, err = readFile("test.txt", rows, cols)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	regions = walkGrid(&grid, rows, cols)
	res = solve2(&grid, &regions)
	expected = 436
	fmt.Printf("TEST part 2 | price: %d\n", res)
	if res != expected {
		log.Fatalf("Test failed with res %d; expected: %d\n", res, expected)
	}

	rows, cols = 4, 4
	grid, err = readFile("test1.txt", rows, cols)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	regions = walkGrid(&grid, rows, cols)
	res = solve2(&grid, &regions)
	expected = 80
	fmt.Printf("TEST part 2 | price: %d\n", res)
	if res != expected {
		log.Fatalf("Test failed with res %d; expected: %d\n", res, expected)
	}

	rows, cols = 5, 5
	grid, err = readFile("test2.txt", rows, cols)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	regions = walkGrid(&grid, rows, cols)
	res = solve2(&grid, &regions)
	expected = 236
	fmt.Printf("TEST part 2 | price: %d\n", res)
	if res != expected {
		log.Fatalf("Test failed with res %d; expected: %d\n", res, expected)
	}

	rows, cols = 6, 6
	grid, err = readFile("test3.txt", rows, cols)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	regions = walkGrid(&grid, rows, cols)
	res = solve2(&grid, &regions)
	expected = 368
	fmt.Printf("TEST part 2 | price: %d\n", res)
	if res != expected {
		log.Fatalf("Test failed with res %d; expected: %d\n", res, expected)
	}

	rows, cols = 10, 10
	grid, err = readFile("test4.txt", rows, cols)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	regions = walkGrid(&grid, rows, cols)
	res = solve2(&grid, &regions)
	expected = 1206
	fmt.Printf("TEST part 2 | price: %d\n", res)
	if res != expected {
		log.Fatalf("Test failed with res %d; expected: %d\n", res, expected)
	}
}
