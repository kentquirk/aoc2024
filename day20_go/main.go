package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/beefsack/go-astar"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type direction byte

const (
	north direction = iota
	east
	south
	west
)

func (d direction) String() string {
	return string(d)
}

func (d direction) left() direction {
	return (d + 3) % 4
}

func (d direction) right() direction {
	return (d + 1) % 4
}

func (d direction) opposite() direction {
	return (d + 2) % 4
}

var directionDeltas = map[direction]point{
	north: {r: -1, c: 0},
	east:  {r: 0, c: 1},
	south: {r: 1, c: 0},
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

func (p point) tunnel(d direction) point {
	return p.add(directionDeltas[d]).add(directionDeltas[d])
}

func (p point) String() string {
	return fmt.Sprintf("(%d, %d)", p.r, p.c)
}

type node struct {
	p         point
	t         int
	neighbors []*node
	tunnels   []*node
}

// ensure that node implements astar.Pather
var _ astar.Pather = &node{}

func (n *node) PathNeighbors() []astar.Pather {
	neighbors := make([]astar.Pather, len(n.neighbors))
	for i, n := range n.neighbors {
		neighbors[i] = n
	}
	return neighbors
}

func (n *node) PathNeighborCost(to astar.Pather) float64 {
	return 1
}

func (n *node) PathEstimatedCost(to astar.Pather) float64 {
	other := to.(*node)
	return float64(abs(n.p.r-other.p.r) + abs(n.p.c-other.p.c))
}

func (n *node) String() string {
	return fmt.Sprintf("%v t=%d nn=%d nt=%d", n.p, n.t, len(n.neighbors), len(n.tunnels))
}

type cpu struct {
	w     int
	h     int
	walls map[point]struct{}
	track map[point]*node
	// nodes map[point]*node
	path []point
	s    point
	e    point
}

func (c cpu) Print(path []point) {
	for row := 0; row < c.h; row++ {
		for col := 0; col < c.w; col++ {
			p := point{row, col}
			if _, ok := c.walls[p]; ok {
				fmt.Print("#")
			} else if p == c.s {
				fmt.Print("S")
			} else if p == c.e {
				fmt.Print("E")
			} else if path != nil && slices.Contains(path, p) {
				// n := c.track[p]
				// fmt.Printf("%c", n.t%10+'0')
				fmt.Print("o")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func (c *cpu) findPathWithAstar() int {
	nstart := c.track[c.s]
	nend := c.track[c.e]
	path, distance, found := astar.Path(nstart, nend)
	if !found {
		return -1
	}
	for i := 0; i < len(path); i++ {
		// path was generated in reverse order
		n := path[len(path)-i-1].(*node)
		n.t = i
		c.path = append(c.path, n.p)
	}
	return int(distance)
}

func (c *cpu) allTunnelsAt(p point, md int) map[point]struct{} {
	// find all tunnels up to a manhattan distance of md from point p
	tunnels := make(map[point]struct{})
	// find all tunnels at this distance
	for dr := 0; dr <= md; dr += 1 {
		dc := md - dr
		pt := point{p.r + dr, p.c + dc}
		if _, ok := c.track[pt]; ok {
			tunnels[pt] = struct{}{}
		}
		pt = point{p.r + dr, p.c - dc}
		if _, ok := c.track[pt]; ok {
			tunnels[pt] = struct{}{}
		}
		pt = point{p.r - dr, p.c + dc}
		if _, ok := c.track[pt]; ok {
			tunnels[pt] = struct{}{}
		}
		pt = point{p.r - dr, p.c - dc}
		if _, ok := c.track[pt]; ok {
			tunnels[pt] = struct{}{}
		}
	}
	return tunnels
}

func parseCPU(lines []string) *cpu {
	cpu := &cpu{
		w:     len(lines[0]),
		h:     len(lines),
		walls: make(map[point]struct{}),
		track: make(map[point]*node),
	}
	for r, line := range lines {
		for c, char := range line {
			if char == '#' {
				cpu.walls[point{r, c}] = struct{}{}
				continue
			}
			p := point{r, c}
			n := &node{p: p, neighbors: make([]*node, 0)}
			cpu.track[p] = n
			if char == 'S' {
				fmt.Println("Start", p)
				cpu.s = p
			}
			if char == 'E' {
				fmt.Println("End", p)
				cpu.e = p
			}
		}
	}
	for p, n := range cpu.track {
		for _, d := range []direction{north, south, east, west} {
			v, ok := cpu.track[p.next(d)]
			if ok {
				n.neighbors = append(n.neighbors, v)
			}
			v, ok = cpu.track[p.tunnel(d)]
			if ok {
				n.tunnels = append(n.tunnels, v)
			}
		}
	}
	return cpu
}

func part1(lines []string) int {
	c := parseCPU(lines)
	// c.Print(nil)
	_ = c.findPathWithAstar()
	c.Print(c.path)
	savings := map[int]int{}
	for _, p := range c.path {
		// look for cheats from this point
		n := c.track[p]
		for _, t := range n.tunnels {
			timesaved := t.t - n.t - 2
			if timesaved > 0 {
				savings[timesaved]++
			}
		}
	}
	fmt.Println("Savings", savings)
	total := 0
	for s, n := range savings {
		if s >= 100 {
			total += n
		}
	}
	return total
}

func part2(lines []string) int {
	c := parseCPU(lines)
	// c.Print(nil)
	_ = c.findPathWithAstar()
	c.Print(c.path)
	savings := map[int]int{}
	for _, p := range c.path {
		// look for cheats from this point
		// no tunnels at distance 1
		for md := 2; md <= 20; md += 1 {
			possibleTunnels := c.allTunnelsAt(p, md)
			// fmt.Println("Possible tunnels from", p, "to", possibleTunnels)
			n := c.track[p]
			for tun := range possibleTunnels {
				t := c.track[tun]
				timesaved := t.t - n.t - md
				if timesaved > 0 {
					savings[timesaved]++
				}
			}
		}
	}
	fmt.Println("Savings", savings)
	total := 0
	for s, n := range savings {
		if s >= 100 {
			total += n
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
