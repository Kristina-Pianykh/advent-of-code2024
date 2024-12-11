package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	stones, err := read_file("input.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", stones)

	MEM[0] = []int{1}
	// fmt.Printf("%v\n", MEM)

	blinks := 75
	fmt.Printf("part 1 | stones: %d\n", solve(blinks, &stones))
	// STONES = 0
	// blinks = 75
	// fmt.Printf("part 2 | stones: %d\n", solve(blinks, &stones))
	fmt.Printf("took %v\n", time.Now().Sub(start))
}

func solve(blink_n int, stones *[]int) int {
	for _, v := range *stones {
		fmt.Printf("checking %d\n", v)
		bfs(0, v, blink_n)
	}
	return STONES
}

var STONES int = 0
var MEM map[int][]int = make(map[int][]int, 100000000)

func calc(n int) []int {
	switch {
	case even_digits(n):
		n1, n2 := split_n(n)
		return []int{n1, n2}
	default:
		return []int{n * 2024}
	}
}

func bfs(depth, n, depth_target int) {
	if depth == depth_target {
		STONES++
		// if STONES%1000 == 0 {
		// 	fmt.Printf("%d\n", STONES)
		// }
		return
	}
	val, ok := MEM[n]
	if !ok {
		val = calc(n)
		MEM[n] = val
	}

	for _, v := range val {
		bfs(depth+1, v, depth_target)
	}
	return
}

func split_n(num int) (int, int) {
	digits := count_digits(num)
	// fmt.Printf("number of digits: %d\n", digits)
	var num1, num2 int

	n := num
	j := 0
	i := 0
	for n > 0 {
		j = (n%10)*pow10(i) + j
		n = n / 10
		i++
		// fmt.Printf("num=%d, n=%d, j=%d, i=%d\n", num, n, j, i)
		if i == digits/2 {
			num1 = n
			num2 = j
			break
		}
	}
	return num1, num2
}

func pow10(power int) int {
	return int(math.Pow(10, float64(power)))
}

func count_digits(n int) int {
	return 1 + int(math.Log10(float64(n)))
}

func even_digits(num int) bool {
	digits := count_digits(num)
	if digits%2 == 0 {
		return true
	}
	return false
}

func read_file(file_path string) ([]int, error) {
	file, err := os.Open(file_path)
	if err != nil {
		return nil, errors.New("Error opening file")
	}
	defer file.Close()

	nums := make([]int, 0, 10000000)

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	for _, v := range strings.Split(line, " ") {
		int_v, err := strconv.ParseInt(strings.Trim(v, " "), 10, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to parse int from %s\n", v))
		}
		nums = append(nums, int(int_v))
	}
	return nums, nil
}
