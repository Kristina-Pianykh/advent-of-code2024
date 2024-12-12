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

var STONES int = 0
var LEAVES_BY_DEPTH map[int]map[int]int = make(map[int]map[int]int, 1000000)

func main() {
	start := time.Now()
	stones, err := read_file("input.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	momoize(75)
	blink := 25
	fmt.Printf("part 1 | stones: %d\n", solve(blink, stones))
	STONES = 0
	blink = 75
	fmt.Printf("part 2 | stones: %d\n", solve(blink, stones))
	fmt.Printf("took %v\n", time.Now().Sub(start))
}

func momoize(upper_bound int) {
	for i := 0; i < 100; i++ {
		// initialize LEAVES_BY_DEPTH
		LEAVES_BY_DEPTH[i] = make(map[int]int, upper_bound)

		for depth_target := 1; depth_target <= upper_bound; depth_target++ {
			LEAVES_BY_DEPTH[i][depth_target] = cache(0, i, i, depth_target)
		}
	}
}

func cache(depth, starting_n, current_n, depth_target int) int {
	if depth == depth_target {
		return 1
	}

	if leaves, ok := LEAVES_BY_DEPTH[current_n][depth_target-depth]; ok {
		return leaves
	}

	acc := 0
	for _, res := range calc(current_n) {
		acc += cache(depth+1, starting_n, res, depth_target)
	}

	if _, exists := LEAVES_BY_DEPTH[current_n]; !exists {
		LEAVES_BY_DEPTH[current_n] = make(map[int]int)
	}

	LEAVES_BY_DEPTH[current_n][depth_target-depth] = acc
	return acc
}

func solve(blink_n int, stones [8]int) int {
	for _, v := range stones {
		dfs(0, v, blink_n)
	}
	return STONES
}

func calc(n int) []int {
	switch {
	case n == 0:
		return []int{1}
	case even_digits(n):
		n1, n2 := split_n(n)
		return []int{n1, n2}
	default:
		return []int{n * 2024}
	}
}

func dfs(depth, n, depth_target int) {
	if depth == depth_target {
		STONES++
		return
	}
	if leaves, ok := LEAVES_BY_DEPTH[n][depth_target-depth]; ok {
		STONES += leaves
		return
	}

	for _, res := range calc(n) {
		if leaves, ok := LEAVES_BY_DEPTH[res][depth_target-depth-1]; ok {
			STONES += leaves
			continue
		}
		dfs(depth+1, res, depth_target)
	}
}

func split_n(num int) (int, int) {
	digits := count_digits(num)
	var num1, num2 int

	n := num
	j := 0
	i := 0
	for n > 0 {
		j = (n%10)*pow10(i) + j
		n = n / 10
		i++
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

func read_file(file_path string) ([8]int, error) {
	nums := [8]int{}
	file, err := os.Open(file_path)
	if err != nil {
		return nums, errors.New("Error opening file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	for i, v := range strings.Split(line, " ") {
		int_v, err := strconv.Atoi(strings.Trim(v, " "))
		if err != nil {
			return nums, errors.New(fmt.Sprintf("failed to parse int from %s\n", v))
		}
		nums[i] = int_v
	}
	return nums, nil
}
