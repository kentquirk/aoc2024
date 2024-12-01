package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func parseNumbersFrom(line string) []int {
	parts := strings.Fields(line)
	numbers := make([]int, len(parts))
	for i, part := range parts {
		numbers[i], _ = strconv.Atoi(part)
	}
	return numbers
}

func part1(lines []string) int {
	left := make([]int, len(lines))
	right := make([]int, len(lines))
	for i, line := range lines {
		parts := parseNumbersFrom(line)
		left[i] = parts[0]
		right[i] = parts[1]
	}
	sort.Ints(left)
	sort.Ints(right)
	dist := 0
	for i := 0; i < len(left); i++ {
		d := right[i] - left[i]
		if d < 0 {
			d = -d
		}
		dist += d
	}
	return dist
}

func part2(lines []string) int {
	left := make([]int, len(lines))
	right := make(map[int]int)
	for i, line := range lines {
		parts := parseNumbersFrom(line)
		left[i] = parts[0]
		right[parts[1]] += 1
	}
	score := 0
	for i := 0; i < len(left); i++ {
		score += left[i] * right[left[i]]
	}
	return score
}

func readlines(filename string) []string {
	f, err := os.Open(fmt.Sprintf("./data/%s.txt", filename))
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(b), "\n")
}

func main() {
	args := os.Args[1:]
	filename := "sample"
	if len(args) > 0 {
		filename = args[0]
	}
	lines := readlines(filename)
	fmt.Println(part1(lines))
	fmt.Println(part2(lines))
}
