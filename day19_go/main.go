package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

func parse(lines []string) ([]string, []string) {
	towels := strings.Split(lines[0], ", ")
	slices.SortFunc(towels, func(a, b string) int { return len(b) - len(a) })
	return towels, lines[2:]
}

func countFromRight(requirement string, towels []string, level int) int {
	possibleTowels := make([]string, 0)
	for _, t := range towels {
		if strings.HasSuffix(requirement, t) {
			possibleTowels = append(possibleTowels, t)
		}
	}
	slices.SortFunc(possibleTowels, func(a, b string) int { return len(b) - len(a) })
	for _, t := range possibleTowels {
		if len(t) == len(requirement) {
			return 1
		}
		if countFromRight(requirement[:len(requirement)-len(t)], towels, level+1) == 1 {
			return 1
		}
	}
	return 0
}

var cache map[string]int = make(map[string]int)

func combosFromRight(requirement string, towels []string, level int) int {
	if count, ok := cache[requirement]; ok {
		return count
	}

	possibleTowels := make([]string, 0)
	for _, t := range towels {
		if strings.HasSuffix(requirement, t) {
			possibleTowels = append(possibleTowels, t)
		}
	}
	slices.SortFunc(possibleTowels, func(a, b string) int { return len(b) - len(a) })
	combos := 0
	for _, t := range possibleTowels {
		if len(t) == len(requirement) {
			combos++
		} else {
			count := combosFromRight(requirement[:len(requirement)-len(t)], towels, level+1)
			if count != 0 {
				combos += count
			}
		}
	}

	cache[requirement] = combos
	return combos
}

func part1(lines []string) int {
	towels, requirements := parse(lines)
	count := 0
	for _, r := range requirements {
		if countFromRight(r, towels, 1) != 0 {
			fmt.Println("\ryes: ", r)
			count++
		} else {
			fmt.Println("\r no: ", r)
		}
	}
	return count
}

func part2(lines []string) int {
	towels, requirements := parse(lines)
	total := 0
	for _, r := range requirements {
		combos := combosFromRight(r, towels, 1)
		if combos != 0 {
			fmt.Println("\ryes: ", r, combos)
			total += combos
		} else {
			fmt.Println("\r no: ", r, combos)
		}
	}
	return total
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
