package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type Page struct {
	v     int
	rules *[][]int
}

type byRules []Page

func (p Page) String() string {
	return fmt.Sprintf("%d", p.v)
}

func (pages byRules) Len() int      { return len(pages) }
func (pages byRules) Swap(i, j int) { pages[i], pages[j] = pages[j], pages[i] }

func (pages byRules) Less(i, j int) bool {
	page_i := pages[i]
	page_j := pages[j]
	return slices.Contains((*page_i.rules)[page_i.v], page_j.v)
}

func main() {
	file, err := os.Open("input.txt")

	if err != nil {
		log.Fatal(err)
	}

	var (
		scanner    *bufio.Scanner = bufio.NewScanner(file)
		re_rules   *regexp.Regexp = regexp.MustCompile(`[[:digit:]]{2}\|[[:digit:]]{2}`)
		re_updates *regexp.Regexp = regexp.MustCompile(`([[:digit:]]{2},)+[[:digit:]]{2}`)
		rules      [][]int        = make([][]int, 100)
		acc1       int            = 0
		acc2       int            = 0
	)
	for i := 0; i < len(rules); i++ {
		rules[i] = []int{}
	}

	updates := [][]Page{}

	for scanner.Scan() {
		line := scanner.Bytes()
		if rule := re_rules.Find(line); len(rule) > 0 {
			before, after := parse_rules(rule)
			rules[before] = append(rules[before], after)
		}

		if update_str := re_updates.Find(line); len(update_str) > 0 {
			update := parse_update(update_str, &rules)
			updates = append(updates, update)

			if is_valid_update(rules, update) {
				acc1 += update[len(update)/2].v
			} else {
				sort.Sort(byRules(update))
				acc2 += update[len(update)/2].v
			}
		}
	}

	err = scanner.Err()
	if err != nil {
		log.Fatalf("scanner encountered an err: %s", err)
	}
	fmt.Printf("part 1 | result: %d\n", acc1)
	fmt.Printf("part 2 | result: %d\n", acc2)
}

func is_valid_update(rules [][]int, update []Page) bool {
	for i := 1; i < len(update); i++ {
		prev := update[i-1].v
		curr := update[i].v
		prev_followed_by := rules[prev]
		if !slices.Contains(prev_followed_by, curr) {
			return false
		}
	}

	if len(update)%2 == 0 {
		log.Fatal(fmt.Sprintf("number of updates is an even number: %v | len(update)=%d\n", update, len(update)))
	}
	return true
}

func parse_rules(str []byte) (int, int) {
	vals := [2]int{}
	for i, v := range strings.Split(string(str), "|") {
		int_v, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Fatal("failed to convert ", v, "to int: ", err)
		}
		vals[i] = int(int_v)
	}
	return vals[0], vals[1]
}

func parse_update(str []byte, rules *[][]int) []Page {
	updates := []Page{}
	for _, v := range strings.Split(string(str), ",") {
		int_v, err := strconv.ParseInt(strings.Replace(v, " ", "", -1), 10, 64)
		if err != nil {
			log.Fatal("failed to convert ", v, "to int: ", err)
		}
		updates = append(updates, Page{int(int_v), rules})
	}
	return updates
}
