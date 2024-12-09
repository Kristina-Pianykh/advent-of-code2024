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

var results []int = []int{}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal("Error opening file")
	}
	defer file.Close()

	var (
		scanner *bufio.Scanner = bufio.NewScanner(file)
		acc     int            = 0
	)

	for scanner.Scan() {
		line := scanner.Text()
		test_v, nums := parse_line(line)
		// nums = []int{1, 2, 3, 4}
		fmt.Printf("%d: %v\n", test_v, nums)
		is_correct_eq(nums)

		if slices.Contains(results, test_v) {
			acc += test_v
		}
		// reset results
		results = results[:0]
	}
	fmt.Printf("part 1 | result: %d\n", acc)
}

func compute(i int, nums []int, op string, acc int) {
	// fmt.Printf("||| nums: %v\ni: %d\nop: %s\nacc: %d\n", nums, i, op, acc)
	if op == "+" {
		// fmt.Printf("update acc: %d + %d = %d\n\n", acc, nums[i+1], acc+nums[i+1])
		acc += nums[i+1]
	} else {
		// fmt.Printf("update acc: %d * %d = %d\n\n", acc, nums[i+1], acc*nums[i+1])
		acc *= nums[i+1]
	}

	if i == len(nums)-2 {
		results = append(results, acc)
		// fmt.Printf("appending new result %d: %v\n", acc, results)

		if op == "+" {
			acc -= nums[i+1]
		} else {
			acc = acc / nums[i+1]
		}

		i--
		return
	}

	compute(i+1, nums, "+", acc)
	compute(i+1, nums, "*", acc)
}

func is_correct_eq(nums []int) bool {
	compute(0, nums, "+", nums[0])
	compute(0, nums, "*", nums[0])
	fmt.Printf("all results: %v\n", results)
	fmt.Printf("len(results)=%d\n", len(results))
	return true
}

func parse_line(line string) (int, []int) {
	str_nums := strings.Split(line, ":")
	test_v, err := strconv.ParseInt(str_nums[0], 10, 64)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to parse %s with error: %s\n", str_nums[0], err.Error()))
	}

	nums := []int{}

	for _, v := range strings.Split(strings.Trim(str_nums[1], " "), " ") {
		if len(strings.Replace(string(v), " ", "", -1)) < 1 {
			continue
		}
		num, err := strconv.ParseInt(strings.Replace(string(v), " ", "", -1), 10, 64)
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to parse %v with error: %s\n", v, err.Error()))
		}
		nums = append(nums, int(num))
	}
	return int(test_v), nums
}
