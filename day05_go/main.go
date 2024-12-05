package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/hmdsefi/gograph"
	"github.com/hmdsefi/gograph/traverse"
)

func parseNumbersFrom(line string) []int {
	pat := regexp.MustCompile(`\d+`)
	parts := pat.FindAllString(line, -1)
	numbers := make([]int, len(parts))
	for i, part := range parts {
		numbers[i], _ = strconv.Atoi(part)
	}
	return numbers
}

func parseConstraints(lines []string) map[int][]int {
	constraints := make(map[int][]int)
	for _, line := range lines {
		if strings.Contains(line, "|") {
			parts := parseNumbersFrom(line)
			constraints[parts[0]] = append(constraints[parts[0]], parts[1])
		}
	}
	return constraints
}

// walk two slices of ordered vertices and determine if the first is a subset of the second in the same order
func isSubset(whole, part []int) bool {
	if len(part) > len(whole) {
		return false
	}
	if len(part) == 0 {
		return true
	}

	p := 0
	w := 0
	for p < len(part) && w < len(whole) {
		if part[p] != whole[w] {
			w++
			continue
		}
		p++
		w++
	}
	// fmt.Println(whole, part, w, p)
	return p == len(part)
}

func bothParts(lines []string) (int, int) {
	constraints := parseConstraints(lines)

	correctTotal := 0
	incorrectTotal := 0
	for _, line := range lines {
		if strings.Contains(line, ",") {
			pages := parseNumbersFrom(line)
			// now we're going to build a graph from the constraints on the pages in the given line
			g := gograph.New[int](gograph.Directed())
			for _, page := range pages {
				v1 := gograph.NewVertex(page)
				for _, constraint := range constraints[page] {
					v2 := gograph.NewVertex(constraint)
					g.AddEdge(v1, v2)
				}
			}
			// get the ordered version of the graph
			iter, err := traverse.NewTopologicalIterator(g)
			if err != nil {
				log.Fatal(err)
			}
			orderedGraph := make([]int, 0)
			for iter.HasNext() {
				v := iter.Next()
				orderedGraph = append(orderedGraph, v.Label())
			}
			if isSubset(orderedGraph, pages) {
				correctTotal += pages[len(pages)/2]
			} else {
				correctOrder := make([]int, 0)
				for _, page := range orderedGraph {
					if slices.Contains(pages, page) {
						correctOrder = append(correctOrder, page)
					}
				}
				// fmt.Println(orderedGraph, pages, correctOrder)
				incorrectTotal += correctOrder[len(correctOrder)/2]
			}
		}
	}

	return correctTotal, incorrectTotal
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
	fmt.Println(bothParts(lines))
}
