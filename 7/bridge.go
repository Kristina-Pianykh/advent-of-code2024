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
		is_correct_eq(test_v, nums)

		if slices.Contains(results, test_v) {
			acc += test_v
		}
		// reset results
		results = results[:0]
	}
	fmt.Printf("result: %d\n", acc)
}

func compute(test_v, i int, nums []int, op string, acc int) {
	if len(results) == 1 {
		return
	}
	switch op {
	case "+":
		acc += nums[i+1]
	case "*":
		acc *= nums[i+1]
	case "||":
		concat := fmt.Sprintf("%d%d", acc, nums[i+1])
		if v, err := strconv.ParseInt(concat, 10, 64); err == nil {
			acc = int(v)
		}
	default:
		log.Fatal(fmt.Sprintf("unrecognized operator %s\n", op))
	}

	if i == len(nums)-2 {
		if acc == test_v {
			results = append(results, acc)
		}

		switch op {
		case "+":
			acc -= nums[i+1]
		case "*":
			acc = acc / nums[i+1]
		case "||":
			acc_str := fmt.Sprintf("%d", acc)
			prev_str := fmt.Sprintf("%d", nums[i+1])
			if v, err := strconv.ParseInt(acc_str[:len(acc_str)-len(prev_str)], 10, 64); err == nil {
				acc = int(v)
			}
		}

		i--
		return
	}

	compute(test_v, i+1, nums, "+", acc)
	compute(test_v, i+1, nums, "*", acc)
	compute(test_v, i+1, nums, "||", acc)
}

func is_correct_eq(test_v int, nums []int) {
	compute(test_v, 0, nums, "+", nums[0])
	compute(test_v, 0, nums, "*", nums[0])
	compute(test_v, 0, nums, "||", nums[0])
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
