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
)

func main() {
	stones, err := read_file("input.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", stones)
	blinks := 25
	fmt.Printf("part 1 | stones: %d\n", solve(blinks, &stones))
	// blinks := 75
	// fmt.Printf("part 2 | stones: %d\n", solve(blinks, &stones))
}

func solve(blink_n int, stones *[]int) int {
	var new_stones []int

	for blink := 0; blink < blink_n; blink++ {
		fmt.Printf("blink i: %d; stones: %d\n", blink, len(*stones))
		new_stones = make([]int, 0, len(*stones)*2)
		// fmt.Printf("new_stones=%v\n", new_stones)

		for _, v := range *stones {
			// fmt.Printf("%d\n", v)
			switch {
			case v == 0:
				// fmt.Printf("%d --> 0\n", v)
				new_stones = append(new_stones, 1)
			case even_digits(v):
				n1, n2 := split_n(v)
				// fmt.Printf("split %d: %d %d\n", v, n1, n2)
				new_stones = append(new_stones, n1)
				new_stones = append(new_stones, n2)
			default:
				// fmt.Printf("%d * 2024\n", v)
				new_stones = append(new_stones, v*2024)
			}
		}

		// fmt.Printf("new_stones=%v\n", new_stones)
		if cap(*stones) <= len(new_stones) { // TODO: perhaps < is enough?
			*stones = make([]int, 0, cap(*stones)*2)
		}
		custom_copy(stones, &new_stones)
		// fmt.Printf("stones=%v\n\n", new_stones)
	}

	return len(*stones)
}

func split_n(num int) (int, int) {
	// 28676032
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

// why? because builtin copy copies the minimum of min(len(slice1), len(slice2))
func custom_copy(dst *[]int, src *[]int) {
	*dst = (*src)[:0]
	*dst = append(*dst, *src...)
	if len(*dst) != len(*src) {
		log.Fatalf("failed to copy items from %v to %v\n", src, dst)
	}
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

// func inc_capacity(slc []int) []int {
// 	new_slc := make([]int, len(slc), len(slc)*2)
// 	res := copy(new_slc, slc)
// 	if res != len(slc) {
// 		log.Fatalf("failed to copy all elements from %v\n", new_slc)
// 	}
// 	return new_slc
// }

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
