package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"testing"
)

func TestReadFile(t *testing.T) {
	rows, cols := 8, 8
	grid, err := read_file("test_sample.txt", rows, cols)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	if len(grid) != rows {
		log.Fatalf("expected rows: %d, got: %d\n", rows, len(grid))
	}
	for i, row := range grid {
		if len(row) != cols {

			log.Fatalf("expected length of row at idx %d: %d, got: %d\n", i, cols, len(grid[i]))
		}
	}
}

func TestIncByte(t *testing.T) {
	var a byte = byte('1')
	fmt.Println(a)
	fmt.Println(a + 1)
}

func TestGridValAt(t *testing.T) {
	rows, cols := 8, 8
	grid, err := read_file("test_sample.txt", rows, cols)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("%c\n", grid.val(Coordinate{2, 3}))
	fmt.Printf("%c\n", grid.val(Coordinate{7, 7}))
	fmt.Printf("%c\n", grid.val(Coordinate{0, 0}))
	fmt.Printf("%c\n", grid.val(Coordinate{5, 1}))
}

func TestGridNext(t *testing.T) {
	c := Coordinate{3, 6}
	fmt.Printf("next coordinate: %s\n", c.next("up"))
	fmt.Printf("next coordinate: %s\n", c.next("down"))
	fmt.Printf("next coordinate: %s\n", c.next("left"))
	fmt.Printf("next coordinate: %s\n", c.next("right"))
}

func TestGridPrev(t *testing.T) {
	c := Coordinate{3, 6}
	fmt.Printf("next coordinate: %s\n", c.prev("up"))
	fmt.Printf("next coordinate: %s\n", c.prev("down"))
	fmt.Printf("next coordinate: %s\n", c.prev("left"))
	fmt.Printf("next coordinate: %s\n", c.prev("right"))
}

func TestSliceContainsCoordinate(t *testing.T) {
	slc := []Coordinate{{1, 2}, {2, 4}}
	fmt.Println(slices.Contains(slc, Coordinate{1, 2}))
	fmt.Println(slices.Contains(slc, Coordinate{2, 2}))
}

func TestSolve(t *testing.T) {
	inputs := map[string][2]int{
		"test_sample13.txt":  {7, 7},
		"test_sample227.txt": {6, 6},
		"test_sample.txt":    {8, 8},
	}
	for filename, size := range inputs {
		grid, err := read_file(filename, size[0], size[1])
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		score, rating := solve(&grid, size[0], size[1])
		fmt.Printf("input %s | score: %d\n", filename, score)
		fmt.Printf("input %s | rating: %d\n\n", filename, rating)
	}
}
