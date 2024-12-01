package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	part1_res := part_one("input.txt")
	fmt.Printf("part one result: %d\n", part1_res)

	part2_res := part_two("input.txt")
	fmt.Printf("part two result: %d\n", part2_res)
}

func read_file(file_path string) []string {
	file, err := os.Open(file_path)

	if err != nil {
		log.Fatal(err)
	}

	lines := make([]string, 1000)
	line_idx := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines[line_idx] = line
		// fmt.Printf("line: %d | Contents: %q \n", line_idx, line)
		line_idx++
	}

	err = scanner.Err()
	if err != nil {
		log.Fatalf("scanner encountered an err: %s", err)
	}

	return lines
}

func parse_cols(lines []string) ([]int, []int) {
	left_slice := make([]int, 1000)
	right_slice := make([]int, 1000)

	for i := 0; i < len(lines); i++ {
		parts := strings.Split(lines[i], "   ")

		if s, err := strconv.ParseInt(parts[0], 10, 32); err == nil {
			// fmt.Printf("%T, %v\n", s, s)
			left_slice[i] = int(s)
		}

		if s, err := strconv.ParseInt(parts[1], 10, 32); err == nil {
			// fmt.Printf("%T, %v\n", s, s)
			right_slice[i] = int(s)
		}
	}
	return left_slice, right_slice
}

func part_one(file_path string) int {
	lines := read_file(file_path)
	left_slice, right_slice := parse_cols(lines)

	sort.Ints(left_slice)
	sort.Ints(right_slice)

	acc := 0
	for i := 0; i < len(left_slice); i++ {
		if left_slice[i] > right_slice[i] {
			acc += left_slice[i] - right_slice[i]
		} else {
			acc += right_slice[i] - left_slice[i]
		}
	}

	return acc
}

func part_two(file_path string) int {
	lines := read_file(file_path)
	left_slice, right_slice := parse_cols(lines)

	count := make(map[int]int)
	for i := 0; i < len(right_slice); i++ {
		val := right_slice[i]
		if _, ok := count[val]; !ok {
			count[val] = 0
		}
		count[val]++
	}

	acc := 0
	for i := 0; i < len(left_slice); i++ {
		if val, ok := count[left_slice[i]]; ok == true {
			acc += left_slice[i] * val
		}
	}
	return acc
}
