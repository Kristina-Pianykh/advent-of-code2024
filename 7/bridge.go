package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal("Error opening file")
	}
	defer file.Close()

	var (
		scanner *bufio.Scanner = bufio.NewScanner(file)
		acc1    int            = 0
		acc2    int            = 0
	)

	for scanner.Scan() {
		line := scanner.Text()
		target, nums := parse_line(line)

		if is_correct_eq(target, nums, []string{"+", "*"}) {
			acc1 += target
		}

		if is_correct_eq(target, nums, []string{"+", "*", "||"}) {
			acc2 += target
		}
	}
	fmt.Printf("part 1 | result: %d\n", acc1)
	fmt.Printf("part 2 | result: %d\n", acc2)

	t := time.Now()
	fmt.Printf("Took %v\n", t.Sub(start))
}

func compute(target, i int, nums []int, op string, acc int, ops []string) bool {
	if acc > target {
		return false
	}
	if i == len(nums)-1 {
		return acc == target
	}

	switch op {
	case "+":
		acc += nums[i+1]
	case "*":
		acc *= nums[i+1]
	case "||":
		acc = acc*pow(10, count_digits(nums[i+1])) + nums[i+1]
	default:
		log.Fatal(fmt.Sprintf("unrecognized operator %s\n", op))
	}

	for _, op := range ops {
		if compute(target, i+1, nums, op, acc, ops) {
			return true
		}
	}

	return false
}

func count_digits(n int) int {
	cnt := 0

	for {
		n = n / 10
		cnt++
		if n == 0 {
			break
		}
	}
	return cnt
}

func pow(base, power int) int {
	acc := base
	if power == 0 {
		return 1
	}
	for {
		if power == 1 {
			break
		}
		acc *= base
		power--
	}
	return acc
}

func is_correct_eq(target int, nums []int, ops []string) bool {
	for _, op := range ops {
		if compute(target, 0, nums, op, nums[0], ops) {
			return true
		}
	}
	return false
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
