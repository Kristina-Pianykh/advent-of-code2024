package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	lines := readLinesFromStream(os.Stdin)
	patterns := parsePatterns(lines)
	designs := lines[2:]
	cnt1 := 0
	cnt2 := 0
	for _, d := range designs {
		cnt1 += cntValidDesigns(patterns, d)
		cnt2 += cntCombinations(patterns, d)
	}
	fmt.Printf("part 1 | valid designs: %d\n", cnt1)
	fmt.Printf("part 2 | combinations: %d\n", cnt2)
	fmt.Printf("took: %v\n", time.Now().Sub(start))
}

func cntValidDesigns(patterns map[byte][]string, design string) int {
	var search func(startIdx int) bool

	search = func(startIdx int) bool {
		if startIdx == len(design) {
			return true
		}

		for _, word := range patterns[design[startIdx]] {
			if strings.HasPrefix(design[startIdx:], word) {
				if search(startIdx + len(word)) {
					return true
				}
			}
		}
		return false
	}
	if search(0) {
		return 1
	}

	return 0
}

func cntCombinations(patterns map[byte][]string, design string) int {
	mem := make(map[int]int, len(design))
	var search func(startIdx int) int

	search = func(startIdx int) int {
		if startIdx == len(design) {
			return 1
		}

		if combCnt, ok := mem[startIdx]; ok {
			return combCnt
		}

		combinations := 0
		for _, word := range patterns[design[startIdx]] {
			if strings.HasPrefix(design[startIdx:], word) {
				combinations += search(startIdx + len(word))
			}
		}
		mem[startIdx] = combinations
		return combinations
	}

	return search(0)
}

func add(set []byte, item byte) []byte {
	for _, el := range set {
		if el == item {
			return set
		}
	}
	set = append(set, item)
	return set
}

func parsePatterns(lines []string) map[byte][]string {
	patterns := map[byte][]string{}

	for _, word := range strings.Split(lines[0], ",") {
		word = strings.Replace(word, " ", "", -1)
		if _, ok := patterns[word[0]]; !ok {
			patterns[word[0]] = []string{}
		}
		patterns[word[0]] = append(patterns[word[0]], word)

	}
	return patterns
}

func readLinesFromStream(file *os.File) []string {
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines
}
