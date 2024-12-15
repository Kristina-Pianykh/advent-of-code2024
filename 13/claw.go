package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	var aXCoef, aYCoef, bXCoef, bYCoef, prizeX, prizeY, tokens1, tokens2, a, b int

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(string(line), "Button A"):
			aXCoef, aYCoef = getButtonCoordinates(line)
		case strings.HasPrefix(string(line), "Button B"):
			bXCoef, bYCoef = getButtonCoordinates(line)
		case strings.HasPrefix(string(line), "Prize"):
			prizeX, prizeY = getPrizeCoordinates(line)
			a, b = solveEquation(aXCoef, aYCoef, bXCoef, bYCoef, prizeX, prizeY)
			tokens1 += (a*3 + b)
			a, b = solveEquation(aXCoef, aYCoef, bXCoef, bYCoef, prizeX+10000000000000, prizeY+10000000000000)
			tokens2 += (a*3 + b)
		default:
			continue
		}
	}
	fmt.Printf("Part 1 | tokens to spend: %d\n", tokens1)
	fmt.Printf("Part 2 | tokens to spend: %d\n", tokens2)
	fmt.Printf("took: %v\n", time.Now().Sub(start))
}

func getButtonCoordinates(line string) (int, int) {
	var err error
	var xInt, yInt int

	re := regexp.MustCompile(`X\+([[:digit:]]+), Y\+([[:digit:]]+)`)
	x, y := re.FindStringSubmatch(line)[1], re.FindStringSubmatch(line)[2]
	xInt, err = strconv.Atoi(x)
	if err != nil {
		log.Fatalf("failed to convert to int: %s\n", x)
	}
	yInt, err = strconv.Atoi(y)
	if err != nil {
		log.Fatalf("failed to convert to int: %s\n", y)
	}
	return xInt, yInt
}

func getPrizeCoordinates(line string) (int, int) {
	var err error
	var xInt, yInt int

	re := regexp.MustCompile(`X\=([[:digit:]]+), Y\=([[:digit:]]+)`)
	x, y := re.FindStringSubmatch(line)[1], re.FindStringSubmatch(line)[2]
	xInt, err = strconv.Atoi(x)
	if err != nil {
		log.Fatalf("failed to convert to int: %s\n", x)
	}
	yInt, err = strconv.Atoi(y)
	if err != nil {
		log.Fatalf("failed to convert to int: %s\n", y)
	}
	return xInt, yInt
}

// returns 0, 0 on non-solvable
func solveEquation(aXCoef, aYCoef, bXCoef, bYCoef, prizeX, prizeY int) (int, int) {
	var a, b int

	lcmV := lcm(aXCoef, aYCoef)
	multiplierX := lcmV / aXCoef
	multiplierY := lcmV / aYCoef

	if aXCoef*multiplierX-aYCoef*multiplierY != 0 {
		log.Fatal("something went wrong")
	}

	bF := (float64(prizeX)*float64(multiplierX) - float64(prizeY)*float64(multiplierY)) / (float64(bXCoef)*float64(multiplierX) - float64(bYCoef)*float64(multiplierY))
	if float64(int(bF)) != bF {
		return 0, 0
	}
	b = int(bF)

	aF := (float64(prizeX) - float64(b)*float64(bXCoef)) / float64(aXCoef)

	if float64(int(aF)) != aF {
		return 0, 0
	}
	a = int(aF)

	return a, b
}

func gcd(a, b int) int {
	for a != b {
		if a > b {
			a = a - b
		} else {
			b = b - a
		}
	}
	return a
}

func lcm(a, b int) int {
	return a * (b / gcd(a, b))
}
