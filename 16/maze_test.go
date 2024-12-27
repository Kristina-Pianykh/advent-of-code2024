package main

import (
	"container/heap"
	"fmt"
	"testing"
)

func TestRemove(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	fmt.Printf("%v\n", remove(arr, 0))
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
