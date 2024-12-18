package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/beefsack/go-astar"
)

func parseNumbersFrom(line string) []int {
	pat := regexp.MustCompile(`[0-9-]+`)
	parts := pat.FindAllString(line, -1)
	numbers := make([]int, len(parts))
	for i, part := range parts {
		numbers[i], _ = strconv.Atoi(part)
	}
	return numbers
}

type point struct {
	x, y int
}

type memory struct {
	w       int
	h       int
	pairs   []point
	open    map[point]int
	blocked map[point]int
	nodes   map[point]*node
}

func newMemory(w, h int) *memory {
	m := &memory{w: w, h: h, pairs: make([]point, 0), open: make(map[point]int), blocked: make(map[point]int)}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// -1 means open and not part of the path
			// a nonnegative integer indicates which step on the path
			m.open[point{x, y}] = -1
		}
	}
	return m
}

func (m *memory) block(p point, t int) bool {
	if m.open[p] == -1 {
		m.blocked[p] = t
		delete(m.open, p)
		return true
	}
	m.Print(p)
	fmt.Println(m.open)
	log.Fatal("Cannot block point", p, "at time", t)
	return false
}

func (m *memory) addToPath(p point, t int) bool {
	if m.open[p] == -1 {
		m.open[p] = t
		return true
	}
	return false
}

func (m *memory) resetPath() {
	for k, v := range m.open {
		if v >= 0 {
			m.open[k] = -1
		}
	}
}

type node struct {
	p  point
	t  int
	ns []*node
}

func (m *memory) generateNodes() {
	nodes := make(map[point]*node)
	for y := 0; y < m.h; y++ {
		for x := 0; x < m.w; x++ {
			if _, ok := m.blocked[point{x, y}]; ok {
				continue
			}
			n := &node{p: point{x, y}, t: -1, ns: make([]*node, 0)}
			nodes[point{x, y}] = n
		}
	}
	// add neighbors
	for _, n := range nodes {
		x, y := n.p.x, n.p.y
		if _, ok := nodes[point{x, y - 1}]; ok {
			n.ns = append(n.ns, nodes[point{x, y - 1}])
		}
		if _, ok := nodes[point{x, y + 1}]; ok {
			n.ns = append(n.ns, nodes[point{x, y + 1}])
		}
		if _, ok := nodes[point{x - 1, y}]; ok {
			n.ns = append(n.ns, nodes[point{x - 1, y}])
		}
		if _, ok := nodes[point{x + 1, y}]; ok {
			n.ns = append(n.ns, nodes[point{x + 1, y}])
		}
	}

	m.nodes = nodes
}

// ensure that node implements astar.Pather
var _ astar.Pather = &node{}

func (n *node) PathNeighbors() []astar.Pather {
	neighbors := make([]astar.Pather, len(n.ns))
	for i, n := range n.ns {
		neighbors[i] = n
	}
	return neighbors
}

func (n *node) PathNeighborCost(to astar.Pather) float64 {
	return 1
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (n *node) PathEstimatedCost(to astar.Pather) float64 {
	other := to.(*node)
	return float64(abs(n.p.x-other.p.x) + abs(n.p.y-other.p.y))
}

func (m *memory) findPathWithAstar() int {
	m.generateNodes()
	nstart := m.nodes[point{0, 0}]
	nend := m.nodes[point{m.w - 1, m.h - 1}]
	path, distance, found := astar.Path(nstart, nend)
	for i, n := range path {
		m.addToPath(n.(*node).p, i)
	}
	if found {
		return int(distance)
	}
	return -1
}

func (m memory) Print(mark point) {
	for y := 0; y < m.h; y++ {
		for x := 0; x < m.w; x++ {
			pt := point{x, y}
			if pt == mark {
				fmt.Print("!")
			} else if _, ok := m.blocked[pt]; ok {
				fmt.Print("#")
			} else if i, ok := m.open[pt]; ok && i >= 0 {
				fmt.Print("o")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func part1(lines []string, size int, maxTime int) int {
	m := newMemory(size, size)
	for _, line := range lines {
		numbers := parseNumbersFrom(line)
		pt := point{numbers[0], numbers[1]}
		m.pairs = append(m.pairs, pt)
	}
	for i, pt := range m.pairs {
		if i >= maxTime {
			break
		}
		m.block(pt, i)
	}
	d := m.findPathWithAstar()
	// m.Print()
	return d
}

func part2BS(lines []string, size int) string {
	// binary search for the lowest value of maxTime that returns -1 when calling part1
	minValue := 1024
	maxValue := len(lines)
	for minValue < maxValue {
		midValue := (minValue + maxValue) / 2
		if part1(lines, size, midValue) == -1 {
			maxValue = midValue
		} else {
			minValue = midValue + 1
		}
	}
	return lines[minValue]
}

func part2(lines []string, size int) (point, int) {
	m := newMemory(size, size)
	for _, line := range lines {
		numbers := parseNumbersFrom(line)
		pt := point{numbers[0], numbers[1]}
		m.pairs = append(m.pairs, pt)
	}
	for i, pt := range m.pairs {
		m.resetPath()
		m.block(pt, i)
		if d := m.findPathWithAstar(); d == -1 {
			m.Print(pt)
			return pt, i
		}
	}
	m.Print(point{-1, -1})
	fmt.Println("No solution found", len(m.pairs), len(lines), len(m.open), len(m.blocked))
	return point{-1, -1}, -1
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
	maxIndex := 6
	if len(args) > 1 {
		maxIndex, _ = strconv.Atoi(args[1])
	}
	maxTime := 12
	if len(args) > 2 {
		maxTime, _ = strconv.Atoi(args[2])
	}
	size := maxIndex + 1 // add 1 to account for 0-based index
	lines := readlines(filename)
	fmt.Println(part1(lines, size, maxTime))
	fmt.Println(part2(lines, size))
}
