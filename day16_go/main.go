package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

type direction byte

const (
	north direction = '^'
	south direction = 'v'
	east  direction = '<'
	west  direction = '>'
)

func (d direction) String() string {
	return string(d)
}

func (d direction) left() direction {
	switch d {
	case north:
		return west
	case south:
		return east
	case east:
		return north
	case west:
		return south
	}
	return d
}

func (d direction) right() direction {
	switch d {
	case north:
		return east
	case south:
		return west
	case east:
		return south
	case west:
		return north
	}
	return d
}

func (d direction) opposite() direction {
	switch d {
	case north:
		return south
	case south:
		return north
	case east:
		return west
	case west:
		return east
	}
	return d
}

var directionDeltas = map[direction]point{
	north: {r: -1, c: 0},
	south: {r: 1, c: 0},
	east:  {r: 0, c: 1},
	west:  {r: 0, c: -1},
}

type point struct {
	r, c int
}

func (p point) add(q point) point {
	return point{p.r + q.r, p.c + q.c}
}

func (p point) next(d direction) point {
	return p.add(directionDeltas[d])
}

func (p point) String() string {
	return fmt.Sprintf("(%d, %d)", p.r, p.c)
}

type node struct {
	point
	neighbors map[direction]point
}

type maze struct {
	w        int
	h        int
	walls    map[point]struct{}
	deadends map[point]struct{}
	m        map[point]*node
	s        point
	e        point
	bestcost int
}

func (m maze) Print(path []point) {
	for r := 0; r < m.h; r++ {
		for c := 0; c < m.w; c++ {
			p := point{r, c}
			if _, ok := m.walls[p]; ok {
				fmt.Print("#")
			} else if p == m.s {
				fmt.Print("S")
			} else if p == m.e {
				fmt.Print("E")
			} else if path != nil && slices.Contains(path, p) {
				fmt.Print("o")
			} else if _, ok := m.deadends[p]; ok {
				fmt.Print("!")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func parseMap(lines []string) maze {
	maze := maze{
		w:        len(lines[0]),
		h:        len(lines),
		walls:    make(map[point]struct{}),
		deadends: make(map[point]struct{}),
		m:        make(map[point]*node),
		bestcost: 100000,
	}
	for r, line := range lines {
		for c, char := range line {
			if char == '#' {
				maze.walls[point{r, c}] = struct{}{}
				continue
			}
			p := point{r, c}
			n := &node{point: p, neighbors: make(map[direction]point)}
			maze.m[p] = n
			if char == 'S' {
				maze.s = p
			}
			if char == 'E' {
				maze.e = p
			}
		}
	}
	for p, n := range maze.m {
		if _, ok := maze.m[p.next(north)]; ok {
			n.neighbors[north] = p.next(north)
		}
		if _, ok := maze.m[p.next(south)]; ok {
			n.neighbors[south] = p.next(south)
		}
		if _, ok := maze.m[p.next(east)]; ok {
			n.neighbors[east] = p.next(east)
		}
		if _, ok := maze.m[p.next(west)]; ok {
			n.neighbors[west] = p.next(west)
		}
	}
	return maze
}

func cp(path []point, p point) []point {
	q := make([]point, len(path)+1)
	copy(q, path)
	q[len(path)] = p
	return q
}

var stepCounter int

// recursive function to find the lowest-cost path from the start to the end
func (m *maze) step(path []point, cost int, position point, dir direction) ([]point, int) {
	// if we've already found a path that's better than the current cost, return
	if m.bestcost > 0 && cost > m.bestcost {
		return nil, 0
	}
	stepCounter++
	// if stepCounter%100000 == 0 {
	// 	m.Print(path)
	// 	fmt.Println(stepCounter, cost)
	// }
	if path[len(path)-1] == m.e {
		// we found a route, print it and return it
		m.Print(path)
		fmt.Println(stepCounter, cost)
		if m.bestcost == 0 || cost < m.bestcost {
			m.bestcost = cost
		}
		return path, cost
	}
	node := m.m[position]
	if node == nil {
		log.Fatalf("no node at %v!", position)
	}
	// these will contain our successful paths and their costs
	paths := make([][]point, 0)
	costs := make([]int, 0)
	// try going straight, left, and right
	toTry := []direction{dir, dir.left(), dir.right()}
	// the cost of going straight, left, and right
	dircosts := []int{1, 1001, 1001}
	for i, d := range toTry {
		// can we move in the direction we're trying?
		next, ok := node.neighbors[d]
		if !ok {
			continue
		}

		// does next have any neighbors besides us?
		nextNeighbors := m.m[next].neighbors
		if len(nextNeighbors) == 1 {
			// it's a useless neighbor, remove it from the neighbors list
			delete(node.neighbors, d)
			m.deadends[next] = struct{}{}
			continue
		}

		// have we already been to this neighbor? If so, skip it
		if slices.Contains(path, next) {
			continue
		}
		// if we can go there, try it
		p, c := m.step(cp(path, next), cost+dircosts[i], next, d)
		if p != nil {
			paths = append(paths, p)
			costs = append(costs, c)
		}
	}
	if len(paths) == 0 {
		return nil, 0
	}
	// find the lowest-cost path
	minix := 0
	for i := 1; i < len(costs); i++ {
		if costs[i] < costs[minix] {
			minix = i
		}
	}
	return paths[minix], costs[minix]
}

func (m *maze) markDeadends() {
	passes := 0
	for found := true; found; {
		passes++
		fmt.Println("pass", passes)
		found = false
		for p, n := range m.m {
			if _, ok := m.deadends[p]; ok {
				continue
			}
			if len(n.neighbors) == 1 && p != m.s && p != m.e {
				m.deadends[p] = struct{}{}
				// there's only one of these
				for dir, neighbor := range n.neighbors {
					toMe := dir.opposite()
					delete(m.m[neighbor].neighbors, toMe)
				}
				found = true
			}
		}
	}
}

func (m *maze) findPath() ([]point, int) {
	return m.step([]point{m.s}, 0, m.s, east)
}

func part1(lines []string) int {
	m := parseMap(lines)
	m.Print(nil)
	m.markDeadends()
	m.Print(nil)
	p, cost := m.findPath()
	m.Print(p)
	fmt.Println("steps:", stepCounter)
	return cost
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
	filename := "test1"
	if len(args) > 0 {
		filename = args[0]
	}
	lines := readlines(filename)
	fmt.Println(part1(lines))
}
