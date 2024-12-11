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
	// stones, err := read_file("blink")
	// if err != nil {
	// 	log.Fatal(err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("%v\n", stones)
	// blinks := 25
	// fmt.Printf("part 1 | stones: %d\n", solve(blinks, &stones))
	blinks := 75
	solve(blinks)
	// fmt.Printf("part 2 | stones: %d\n", solve(blinks, &stones))
}

func solve(blink_n int) {
	var (
		filename_read  string
		filename_write string
		reader         *bufio.Scanner
		writer         *bufio.Writer
		line_idx       int
		val            int
	)

	for blink := 0; blink < blink_n; blink++ {
		line_idx = 0
		filename_read = fmt.Sprintf("blink_%d.txt", blink)
		filename_write = fmt.Sprintf("blink_%d.txt", blink+1)
		file_read, err1 := os.Open(filename_read)
		file_write, err2 := os.Create(filename_write)

		if err1 != nil {
			log.Fatal(err1)
		}

		if err2 != nil {
			log.Fatal(err2)
		}
		reader = bufio.NewScanner(file_read)
		writer = bufio.NewWriter(file_write)

		for reader.Scan() {
			str := reader.Text()
			v, err := strconv.ParseInt(strings.Trim(str, " "), 10, 64)
			if err != nil {
				log.Fatalf("%s | failed to parse int from %s\n", filename_read, str)
			}
			val = int(v)

			switch {
			case val == 0:
				// fmt.Printf("%d --> 0\n", v)
				fmt.Fprintf(writer, "%d\n", 1)
			case even_digits(val):
				n1, n2 := split_n(val)
				// fmt.Printf("split %d: %d %d\n", v, n1, n2)
				fmt.Fprintf(writer, "%d\n", n1)
				fmt.Fprintf(writer, "%d\n", n2)
			default:
				// fmt.Printf("%d * 2024\n", v)
				fmt.Fprintf(writer, "%d\n", val*2024)
			}

			if line_idx%10000000 == 0 {
				fmt.Printf("processed %d0 mln lines\n", line_idx/10000000)
				writer.Flush() // Don't forget to flush!
			}
			line_idx++
		}
		fmt.Printf("%s | processed %d stones\n", filename_read, line_idx)

		if err := reader.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		file_read.Close()
		writer.Flush()
		file_write.Close()
	}
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
