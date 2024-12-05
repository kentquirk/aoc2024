package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
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

func parseGraph(lines []string) gograph.Graph[int] {
	g := gograph.New[int](gograph.Directed())
	for _, line := range lines {
		if strings.Contains(line, "|") {
			parts := parseNumbersFrom(line)
			v1 := gograph.NewVertex(parts[0])
			v2 := gograph.NewVertex(parts[1])
			g.AddEdge(v1, v2)
		}
	}
	return g
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
		if part[p] == whole[w] {
			p++
			w++
			continue
		}
		// iterate through the whole until we find the match to part
		for w < len(whole) && part[p] != whole[w] {
			w++
		}
		if w == len(whole) {
			return false
		}
		// now we can increment p
		p++
		if p == len(part) {
			return true
		}
	}
	return false
}

func part1(lines []string) int {
	g := parseGraph(lines)

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
	fmt.Println(orderedGraph)

	total := 0
	for _, line := range lines {
		if strings.Contains(line, ",") {
			parts := parseNumbersFrom(line)
			for _, part := range parts {
				if g.GetVertexByID(part) == nil {
					fmt.Println("Vertex not found: ", part)
					continue
				}
			}
			if isSubset(orderedGraph, parts) {
				total += parts[len(parts)/2]
			}
		}
	}

	return total
}

func part2(lines []string) int {
	return 0
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
}
