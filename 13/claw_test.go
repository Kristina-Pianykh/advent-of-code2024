package main

import (
	"fmt"
	"testing"
)

func TestGcdAndLcm(t *testing.T) {
	var a, b, g int

	a = 17
	b = 86
	g = gcd(a, b)
	fmt.Printf("gcd(%d, %d)=%d\n", a, b, g)
	fmt.Printf("lcm(%d, %d, %d)=%d\n", a, b, g, lcm(a, b, g))

	a = 94
	b = 34
	g = gcd(a, b)
	fmt.Printf("gcd(%d, %d)=%d\n", a, b, g)
	fmt.Printf("lcm(%d, %d, %d)=%d\n", a, b, g, lcm(a, b, g))

	a = 26
	b = 66
	g = gcd(a, b)
	fmt.Printf("gcd(%d, %d)=%d\n", a, b, g)
	fmt.Printf("lcm(%d, %d, %d)=%d\n", a, b, g, lcm(a, b, g))
}

func TestGetButtonCoordinates(t *testing.T) {
	lines := []string{
		"Button A: X+94, Y+34",
		"Button B: X+22, Y+67",
		// "Prize: X=8400, Y=5400",
	}
	x := 0
	y := 0
	fmt.Println(x)
	fmt.Println(y)
	for _, line := range lines {
		x, y = getButtonCoordinates(line)
		fmt.Printf("x=%d, y=%d\n", x, y)
	}
}

func TestGetPrizeCoordinates(t *testing.T) {
	lines := []string{
		"Prize: X=8400, Y=5400",
		"Prize: X=12748, Y=12176",
		"Prize: X=7870, Y=6450",
		"Prize: X=18641, Y=10279",
	}
	x := 0
	y := 0

	for _, line := range lines {
		x, y = getPrizeCoordinates(line)
		fmt.Printf("x=%d, y=%d\n", x, y)
	}
}

func TestSolveEquation(t *testing.T) {
	var aXCoef, aYCoef, bXCoef, bYCoef, prizeX, prizeY, a, b int

	aXCoef, aYCoef = getButtonCoordinates("Button A: X+94, Y+34")
	bXCoef, bYCoef = getButtonCoordinates("Button B: X+22, Y+67")
	prizeX, prizeY = getPrizeCoordinates("Prize: X=8400, Y=5400")
	a, b = solveEquation(aXCoef, aYCoef, bXCoef, bYCoef, prizeX, prizeY)
	fmt.Printf("A: %d, B: %d\n", a, b)

	aXCoef, aYCoef = getButtonCoordinates("Button A: X+94, Y+34")
	bXCoef, bYCoef = getButtonCoordinates("Button B: X+22, Y+67")
	prizeX, prizeY = getPrizeCoordinates("Prize: X=8400, Y=5400")
	a, b = solveEquation(aXCoef, aYCoef, bXCoef, bYCoef, prizeX, prizeY)
	fmt.Printf("A: %d, B: %d\n", a, b)

	aXCoef, aYCoef = getButtonCoordinates("Button A: X+94, Y+34")
	bXCoef, bYCoef = getButtonCoordinates("Button B: X+22, Y+67")
	prizeX, prizeY = getPrizeCoordinates("Prize: X=8400, Y=5400")
	a, b = solveEquation(aXCoef, aYCoef, bXCoef, bYCoef, prizeX, prizeY)
	fmt.Printf("A: %d, B: %d\n", a, b)
}
