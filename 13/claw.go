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
	var aXCoef, aYCoef, bXCoef, bYCoef, prizeX, prizeY, tokens int

	lineIdx := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Printf("processing: %s\n", line)
		switch {
		case strings.HasPrefix(string(line), "Button A"):
			aXCoef, aYCoef = getButtonCoordinates(line)
		case strings.HasPrefix(string(line), "Button B"):
			bXCoef, bYCoef = getButtonCoordinates(line)
		case strings.HasPrefix(string(line), "Prize"):
			prizeX, prizeY = getPrizeCoordinates(line)
			a, b := solveEquation(aXCoef, aYCoef, bXCoef, bYCoef, prizeX, prizeY)
			// fmt.Printf("a=%d, b=%d\n", a, b)
			if a > 0 && b > 0 {
				fmt.Printf("solved equation for %s\n", line)
			}
			tokens += (a*3 + b)
		default:
			continue
		}
		lineIdx++
	}
	fmt.Printf("Part 1 | tokens to spend: %d\n", tokens)
	fmt.Printf("took: %v\n", time.Now().Sub(start))
	// 5112831749757268814 too high
}

func getButtonCoordinates(line string) (int, int) {
	var err error
	var xInt, yInt int

	re := regexp.MustCompile(`X\+([[:digit:]]+), Y\+([[:digit:]]+)`)
	x, y := re.FindStringSubmatch(line)[1], re.FindStringSubmatch(line)[2]
	// fmt.Printf("x=%s, y=%s\n", x, y)
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
	// fmt.Printf("x=%s, y=%s\n", x, y)
	xInt, err = strconv.Atoi(x)
	if err != nil {
		log.Fatalf("failed to convert to int: %s\n", x)
	}
	yInt, err = strconv.Atoi(y)
	if err != nil {
		log.Fatalf("failed to convert to int: %s\n", y)
	}
	return xInt + 10000000000000, yInt + 10000000000000
}

// returns 0 on non-solvable
func solveEquation(aXCoef, aYCoef, bXCoef, bYCoef, prizeX, prizeY int) (int, int) {
	var a, b int

	// aXCoef*a + bXCoef*b = prizeX
	// aYCoef*a + bYCoef*b = prizeY

	lcmV := lcm(aXCoef, aYCoef, gcd(aXCoef, aYCoef))
	// fmt.Printf("LCM: %d\n", lcmV)
	multiplierX := lcmV / aXCoef
	multiplierY := lcmV / aYCoef
	// fmt.Printf("multiplierX: %d\n", multiplierX)
	// fmt.Printf("multiplierY: %d\n", multiplierY)

	// aYCoef*a*multiplierY + bYCoef*b*multiplierY = prizeY * multiplierY

	if aXCoef*multiplierX-aYCoef*multiplierY != 0 {
		log.Fatal("something went wrong")
	}

	// bXCoef*b*multiplierX - bYCoef*b*multiplierY = prizeX*multiplierX - prizeY*multiplierY
	// b(bXCoef*multiplierX - bYCoef*multiplierY) = prizeX*multiplierX - prizeY*multiplierY
	// b * Z = prizeX*multiplierX - prizeY*multiplierY
	// b = (prizeX*multiplierX - prizeY*multiplierY) / Z
	// b = (prizeX*multiplierX - prizeY*multiplierY) / (bXCoef*multiplierX - bYCoef*multiplierY)

	bF := (float64(prizeX)*float64(multiplierX) - float64(prizeY)*float64(multiplierY)) / (float64(bXCoef)*float64(multiplierX) - float64(bYCoef)*float64(multiplierY))
	// fmt.Printf("b as float: %f\n", bF)
	// fmt.Printf("b as float64(int(bF)): %f\n", float64(int(bF)))
	// fmt.Printf("bF/float64(int(bF)): %f\n", bF/float64(int(bF)))

	if bF/float64(int(bF)) != 1.0 {
		// fmt.Printf("we are here")
		return 0, 0
	}
	b = int(bF)

	aF := (float64(prizeX) - float64(b)*float64(bXCoef)) / float64(aXCoef)

	if aF/float64(int(aF)) != 1.0 {
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

func lcm(a, b, gcd int) int {
	return a * (b / gcd)
}
