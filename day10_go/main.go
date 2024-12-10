package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type position struct {
	r int
	c int
}

func (p position) add(d position) position {
	return position{p.r + d.r, p.c + d.c}
}

type adjacency map[position][]position

type adjacencies struct {
	adj         []adjacency
	totalRoutes int
}

func getCh(lines []string, r int, c int) byte {
	if r < 0 || r >= len(lines) {
		return 0
	}
	if c < 0 || c >= len(lines[r]) {
		return 0
	}
	return lines[r][c]
}

func parse(lines []string) *adjacencies {
	deltas := []position{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	adjacencies := &adjacencies{adj: make([]adjacency, 9)}
	for i := 0; i < 9; i++ {
		adjacencies.adj[i] = make(adjacency)
		for r, line := range lines {
			for c := range line {
				if line[c] == byte('0')+byte(i) {
					adjacencies.adj[i][position{r, c}] = []position{}
					for _, d := range deltas {
						if getCh(lines, r+d.r, c+d.c) == byte('0')+byte(i+1) {
							adjacencies.adj[i][position{r, c}] = append(adjacencies.adj[i][position{r, c}], position{r + d.r, c + d.c})
						}
					}
				}
			}
		}
	}
	return adjacencies
}

// calculate the score from a single trailhead
func (a *adjacencies) CountRoutesFrom(p position, i int) map[position]struct{} {
	if i == 9 {
		a.totalRoutes++
		return map[position]struct{}{p: struct{}{}}
	}
	destinations := make(map[position]struct{})
	for _, d := range a.adj[i][p] {
		for k := range a.CountRoutesFrom(d, i+1) {
			destinations[k] = struct{}{}
		}
	}
	return destinations
}

func part1(lines []string) int {
	adjacencies := parse(lines)
	totalScore := 0
	for p := range adjacencies.adj[0] {
		endpoints := adjacencies.CountRoutesFrom(p, 0)
		totalScore += len(endpoints)
		// fmt.Println(p, len(endpoints))
	}
	return totalScore
}

func part2(lines []string) int {
	adjacencies := parse(lines)
	for p := range adjacencies.adj[0] {
		adjacencies.CountRoutesFrom(p, 0)
	}
	return adjacencies.totalRoutes
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
