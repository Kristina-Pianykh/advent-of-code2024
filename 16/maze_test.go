package main

import (
	"container/heap"
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

func TestPriorityQueue(t *testing.T) {
	co := []Coordinate{
		{x: 1, y: -1},
		{x: 2, y: -2},
		{x: 3, y: -3},
		{x: 4, y: -4},
	}

	pq := make(PriorityQueue, 4)
	for i := range pq {
		pq[i] = &Node{
			c:     co[i],
			cost:  i,
			index: i,
		}
	}
	heap.Init(&pq)

	pq.update(pq[0], 10)

	for pq.Len() > 0 {
		node := heap.Pop(&pq).(*Node)
		fmt.Printf("c: %s, cost: %d, idx: %d\n", node.c.string(), node.cost, node.index)
	}
}
