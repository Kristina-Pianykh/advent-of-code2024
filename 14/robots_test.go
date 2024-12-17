package main

import (
	"fmt"
	"testing"
)

func TestModulo(t *testing.T) {
	var a, b int
	a = 7
	b = -2
	fmt.Printf("%d mod %d = %d\n", a, b, a%b)

	a = 7
	b = -5
	fmt.Printf("%d mod %d = %d\n", a, b, a%b)

	a = 7
	b = -8
	fmt.Printf("%d mod %d = %d\n", a, b, a%b)

	a = -11
	b = 7
	fmt.Printf("%d mod %d = %d\n", a, b, a%b)
}

func TestInitSparseGrid(t *testing.T) {
	grid := initSparseGrid(5, 5)
	fmt.Printf("%v\n", grid)
}

func TestParse(t *testing.T) {
	var robot Robot

	strs := []string{
		"p=0,4 v=3,-3",
		"p=6,3 v=-1,-3",
		"p=10,3 v=-1,2",
		"p=2,0 v=2,-1",
		"p=0,0 v=1,3",
		"p=3,0 v=-2,-2",
		"p=7,6 v=-1,-3",
		"p=3,0 v=-1,-2",
		"p=9,3 v=2,3",
		"p=7,3 v=-1,2",
		"p=2,4 v=2,-3",
		"p=9,5 v=-3,-3",
	}
	for _, v := range strs {
		robot = parse(v)
		fmt.Printf("%s\n", robot.string())
	}
}

func TestMod(t *testing.T) {
	fmt.Printf("%d\n", mod(-11, 7))
	fmt.Printf("%d\n", mod(-296, 7))
	a := -11
	b := 7
	fmt.Println((a%b + b) % b)
	fmt.Println(a % b)
	// fmt.Printf("%d\n", mod(7, -304))
}

func TestFindRobotSeq(t *testing.T) {
	str := "................1...........................2.1.....1............................................882838235828585...."
	fmt.Printf("%v\n", containsLinedupRobots(str))
}
