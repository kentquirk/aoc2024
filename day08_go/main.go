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

func (p position) vectorTo(p2 position) position {
	return position{p2.r - p.r, p2.c - p.c}
}

func (p position) add(p2 position) position {
	return position{p.r + p2.r, p.c + p2.c}
}

type antennaMap struct {
	w int
	h int
	m map[byte][]position
}

func (a *antennaMap) isOnMap(p position) bool {
	return p.r >= 0 && p.r < a.h && p.c >= 0 && p.c < a.w
}

// record a slice of positions for each frequency
func parseAntennaMap(lines []string) *antennaMap {
	m := &antennaMap{w: len(lines[0]), h: len(lines), m: make(map[byte][]position)}
	for r, line := range lines {
		for c, char := range line {
			switch char {
			case '.':
				// do nothing
			default:
				m.m[byte(char)] = append(m.m[byte(char)], position{r, c})
			}
		}
	}
	return m
}

// returns a slice of positions that are part1Antinodes for the two given positions
// but only if they're on the map (so it may return 0, 1, or 2 positions)
func (m *antennaMap) part1Antinodes(p1, p2 position) []position {
	antinodes := []position{}
	a1 := p2.add(p1.vectorTo(p2))
	if m.isOnMap(a1) {
		antinodes = append(antinodes, a1)
	}
	a2 := p1.add(p2.vectorTo(p1))
	if m.isOnMap(a2) {
		antinodes = append(antinodes, a2)
	}
	return antinodes
}

// returns a slice of positions that are part2Antinodes for the two given positions
// but only if they're on the map (so assuming that the two positions are on the map,
// it will return at least 2 or possibly many more positions)
func (m *antennaMap) part2Antinodes(p1, p2 position) []position {
	antinodes := []position{p1, p2}
	d1 := p1.vectorTo(p2)
	for a1 := p2.add(d1); m.isOnMap(a1); a1 = a1.add(d1) {
		antinodes = append(antinodes, a1)
	}
	d2 := p2.vectorTo(p1)
	for a2 := p1.add(d2); m.isOnMap(a2); a2 = a2.add(d2) {
		antinodes = append(antinodes, a2)
	}
	return antinodes
}

func bothParts(am *antennaMap, f func(p1, p2 position) []position) int {
	allNodes := map[position]struct{}{}
	count := 0
	for freq, positions := range am.m {
		fcount := 0
		for i, p1 := range positions {
			for _, p2 := range positions[i+1:] {
				antinodes := f(p1, p2)
				count += len(antinodes)
				fcount += len(antinodes)
				for _, a := range antinodes {
					allNodes[a] = struct{}{}
				}
			}
		}
		fmt.Printf("Frequency %c: %d\n", freq, fcount)
	}
	fmt.Printf("Total: %d\n", count)
	fmt.Printf("Unique: %d\n", len(allNodes))
	return len(allNodes)
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
	am := parseAntennaMap(lines)
	fmt.Println(bothParts(am, am.part1Antinodes))
	fmt.Println(bothParts(am, am.part2Antinodes))
}
