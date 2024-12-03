package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")

	if err != nil {
		log.Fatal(err)
	}

	var (
		scanner  *bufio.Scanner = bufio.NewScanner(file)
		sum_acc1 int            = 0
		sum_acc2 int            = 0
	)

	for scanner.Scan() {
		line := scanner.Bytes()
		sum_acc1 += sum_muls(line)
		sum_acc2 += strict_sum_muls(line)
	}

	err = scanner.Err()
	if err != nil {
		log.Fatalf("scanner encountered an err: %s", err)
	}
	fmt.Printf("part 1 | result: %d\n", sum_acc1)
	fmt.Printf("part 2 | result: %d\n", sum_acc2)
}

func extract_operands(mul_op []byte) []int {
	var (
		re_ops       *regexp.Regexp = regexp.MustCompile(`[[:digit:]]{1,3}`)
		operands     [][]byte       = re_ops.FindAll(mul_op, -1)
		operands_int []int          = make([]int, 2)
	)

	if len(operands) != 2 {
		log.Fatal("failed to extract two integers from ", string(mul_op))
	}

	for i, v := range operands {
		int_v, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			log.Fatal("failed to convert ", string(v), "to int: ", err)
		}
		operands_int[i] = int(int_v)
	}

	return operands_int
}

func sum_muls(line []byte) int {
	re := regexp.MustCompile(`mul\([[:digit:]]{1,3},[[:digit:]]{1,3}\)`)
	acc := 0

	for _, v := range re.FindAll(line, -1) {
		operands := extract_operands(v)
		acc += operands[0] * operands[1]
	}
	return acc
}

var ENABLED bool = true

func strict_sum_muls(line []byte) int {
	re := regexp.MustCompile(`mul\([[:digit:]]{1,3},[[:digit:]]{1,3}\)|do\(\)|don't\(\)`)
	acc := 0

	for _, v := range re.FindAll(line, -1) {
		if strings.HasPrefix(string(v), "do()") {
			ENABLED = true
			continue
		} else if strings.HasPrefix(string(v), "don't()") {
			ENABLED = false
			continue
		} else {
			if ENABLED {
				operands := extract_operands(v)
				acc += operands[0] * operands[1]
			}
		}
	}
	return acc
}
