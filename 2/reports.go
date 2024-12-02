package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")

	if err != nil {
		log.Fatal(err)
	}

	var (
		scanner *bufio.Scanner = bufio.NewScanner(file)
		acc1    int            = 0
		acc2    int            = 0
		parts   []string
		report  []int
	)

	for scanner.Scan() {
		line := scanner.Text()
		parts = strings.Split(line, " ")
		report = make([]int, len(parts))

		for i := 0; i < len(parts); i++ {
			if s, err := strconv.ParseInt(parts[i], 10, 32); err == nil {
				report[i] = int(s)
			}
		}

		if is_increasing(report) || is_decreasing(report) {
			acc1++
		}
		if is_soft_increasing(report) || is_soft_decreasing(report) {
			acc2++
		}
	}

	err = scanner.Err()
	if err != nil {
		log.Fatalf("scanner encountered an err: %s", err)
	}
	fmt.Printf("Part 1 | result: %d\n", acc1)
	fmt.Printf("Part 2 | result: %d\n", acc2)
}

func is_increasing(report []int) bool {
	for i := 1; i < len(report); i++ {
		if report[i-1] >= report[i] {
			return false
		}
		if report[i]-report[i-1] > 3 {
			return false
		}
	}
	return true
}

func is_soft_increasing(report []int) bool {
	for i := 0; i < len(report); i++ {
		if is_increasing(slices.Concat(report[0:i], report[i+1:])) {
			return true
		}
	}
	return false
}

func is_decreasing(report []int) bool {
	for i := 1; i < len(report); i++ {
		if report[i-1] <= report[i] {
			return false
		}
		if report[i-1]-report[i] > 3 {
			return false
		}
	}
	return true
}

func is_soft_decreasing(report []int) bool {
	for i := 0; i < len(report); i++ {
		if is_decreasing(slices.Concat(report[0:i], report[i+1:])) {
			return true
		}
	}
	return false
}
