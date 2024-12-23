package main

import (
	"fmt"
	"testing"
)

func TestGetKeyByVal(t *testing.T) {
	dirs := map[byte][]int{
		'>': {1, 0},
		'<': {-1, 0},
		'^': {0, -1},
		'v': {0, 1},
	}
	val, err := getKeyByValue(dirs, []int{0, -1})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%c\n", val)
}
