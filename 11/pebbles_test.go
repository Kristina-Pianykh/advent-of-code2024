package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestCopySlice(t *testing.T) {
	a := make([]int, 0, 5)
	fmt.Printf("%v\n", a)
	b := []int{1, 2, 3}
	fmt.Printf("len(b): %d, cap(b): %d\n", len(b), cap(b))
	n := copy(a, b)
	fmt.Printf("copied items: %d\n", n)
	fmt.Printf("%v\n", a)
}

func TestReduceLength(t *testing.T) {
	a := []int{1, 2, 3}
	fmt.Printf("len(b): %d, cap(b): %d\n", len(a), cap(a))
	a = a[:0]
	fmt.Printf("len(b): %d, cap(b): %d\n", len(a), cap(a))

	a = make([]int, 0, 5)
	fmt.Printf("len(b): %d, cap(b): %d\n", len(a), cap(a))
	a = a[:0]
	fmt.Printf("len(b): %d, cap(b): %d\n", len(a), cap(a))
}

func TestSplitDigit(t *testing.T) {
	var a int
	a = 12
	b, c := split_n(a)
	fmt.Printf("%d %d\n", b, c)
	a = 1252
	b, c = split_n(a)
	fmt.Printf("%d %d\n", b, c)
	a = 9152
	b, c = split_n(a)
	fmt.Printf("%d %d\n", b, c)
	a = 9052
	b, c = split_n(a)
	fmt.Printf("%d %d\n", b, c)
	a = 254905
	b, c = split_n(a)
	fmt.Printf("%d %d\n", b, c)
	a = 254000
	b, c = split_n(a)
	fmt.Printf("%d %d\n", b, c)
	a = 254004
	b, c = split_n(a)
	fmt.Printf("%d %d\n", b, c)
	a = 254014
	b, c = split_n(a)
	fmt.Printf("%d %d\n", b, c)
	a = 28676032
	b, c = split_n(a)
	fmt.Printf("%d %d\n", b, c)
}

func TestSolve(t *testing.T) {
	stones, err := read_file("test.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", stones)
	blinks := 6
	fmt.Printf("after %d blinks: %d\n", blinks, solve(blinks, &stones))
}

func TestMapInitializer(t *testing.T) {
	m := make(map[int]map[int]int, 10)
	for i := 0; i < 10; i++ {
		m[i] = make(map[int]int, 5)
		// fmt.Printf("%d, %v\n", i, m[i])
	}
	for k, v := range m {
		fmt.Printf("%d, %v\n", k, v)
	}
}
