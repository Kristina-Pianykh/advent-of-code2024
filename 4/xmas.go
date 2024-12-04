package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	content, err := read_file("input.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	lines := make([][]byte, len(content))
	for i, line := range content {
		lines[i] = []byte(line)
	}

	var (
		m1        map[int]byte = map[int]byte{0: 'X', 1: 'M', 2: 'A', 3: 'S'}
		xmas_cnt1 int          = 0
		xmas_cnt2 int          = 0
	)

	// part 1
	for y, line := range lines {
		for x := 0; x < len(line); x++ {

			if line[x] == 'X' {
				if walk_left(line, m1, x) {
					xmas_cnt1++
				}
				if walk_right(line, m1, x) {
					xmas_cnt1++
				}
				if walk_up(lines, m1, x, y) {
					xmas_cnt1++
				}
				if walk_down(lines, m1, x, y) {
					xmas_cnt1++
				}
				if walk_up_left(lines, m1, x, y) {
					xmas_cnt1++
				}
				if walk_up_right(lines, m1, x, y) {
					xmas_cnt1++
				}
				if walk_down_left(lines, m1, x, y) {
					xmas_cnt1++
				}
				if walk_down_right(lines, m1, x, y) {
					xmas_cnt1++
				}
			}
		}
	}
	fmt.Printf("XMAS count: %d\n", xmas_cnt1)

	// part 2
	rotated_matrix := lines
	for i := 0; i < 4; i++ {
		for y := 1; y < len(rotated_matrix)-1; y++ {
			line := rotated_matrix[y]
			for x := 1; x < len(line)-1; x++ {
				if line[x] == 'A' && check_pattern(rotated_matrix, x, y) {
					xmas_cnt2++
				}
			}
		}
		rotated_matrix = rotate(rotated_matrix)
	}
	fmt.Printf("X-MAS count: %d\n", xmas_cnt2)
}

func rotate(matrix [][]byte) [][]byte {
	n := len(matrix)
	m := len(matrix[0])
	new_matrix := make([][]byte, n)
	for i := 0; i < n; i++ {
		new_matrix[i] = make([]byte, m)
	}
	for i := 0; i < n; i++ {
		for j, val := range matrix[i] {
			new_matrix[j][n-i-1] = val
		}
	}
	return new_matrix
}

func check_pattern(lines [][]byte, x int, y int) bool {
	if lines[y-1][x-1] == 'M' && lines[y-1][x+1] == 'M' && lines[y+1][x-1] == 'S' && lines[y+1][x+1] == 'S' {
		return true
	}
	return false
}

func walk_left(line []byte, m map[int]byte, x int) bool {
	if x < len(m)-1 {
		return false
	}
	for i := 1; i < len(m); i++ {
		if line[x-i] != m[i] {
			return false
		}
	}
	return true
}

func walk_right(line []byte, m map[int]byte, x int) bool {
	if x > len(line)-len(m) {
		return false
	}
	for i := 1; i < len(m); i++ {
		if line[x+i] != m[i] {
			return false
		}
	}
	return true
}

func walk_up(lines [][]byte, m map[int]byte, x int, y int) bool {
	if y < len(m)-1 {
		return false
	}
	for i := 1; i < len(m); i++ {
		if lines[y-i][x] != m[i] {
			return false
		}
	}
	return true
}

func walk_down(lines [][]byte, m map[int]byte, x int, y int) bool {
	if y > len(lines)-len(m) {
		return false
	}
	for i := 1; i < len(m); i++ {
		if lines[y+i][x] != m[i] {
			return false
		}
	}
	return true
}

func walk_up_left(lines [][]byte, m map[int]byte, x int, y int) bool {
	if y < len(m)-1 || x < len(m)-1 {
		return false
	}
	for i := 1; i < len(m); i++ {
		if lines[y-i][x-i] != m[i] {
			return false
		}
	}
	return true
}

func walk_up_right(lines [][]byte, m map[int]byte, x int, y int) bool {
	if y < len(m)-1 || x > len(lines[y])-len(m) {
		return false
	}
	for i := 1; i < len(m); i++ {
		if lines[y-i][x+i] != m[i] {
			return false
		}
	}
	return true
}

func walk_down_left(lines [][]byte, m map[int]byte, x int, y int) bool {
	if y > len(lines)-len(m) || x < len(m)-1 {
		return false
	}
	for i := 1; i < len(m); i++ {
		if lines[y+i][x-i] != m[i] {
			return false
		}
	}
	return true
}

func walk_down_right(lines [][]byte, m map[int]byte, x int, y int) bool {
	if y > len(lines)-len(m) || x > len(lines[y])-len(m) {
		return false
	}
	for i := 1; i < len(m); i++ {
		if lines[y+i][x+i] != m[i] {
			return false
		}
	}
	return true
}

func read_file(file_path string) ([]string, error) {
	file, err := os.Open(file_path)
	if err != nil {
		return nil, errors.New("Error opening file")
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, errors.New("Error getting file stats")
	}

	size := stat.Size()
	content := make([]byte, size)
	_, err = file.Read(content)
	if err != nil {
		return nil, errors.New("Error reading file")
	}
	stringified := string(content)
	if strings.HasSuffix(stringified, "\n") {
		stringified = stringified[:len(stringified)-1]
	}
	return strings.Split(strings.ReplaceAll(stringified, " ", ""), "\n"), nil
}
