package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

type walldir int

const (
	north walldir = iota
	east
	south
	west
)

func (w walldir) String() string {
	return [...]string{"north", "east", "south", "west"}[w]
}

type plot struct {
	plant  rune
	r      int
	c      int
	region int
	fences []walldir
}

func (p *plot) addFence(fence walldir) {
	if !slices.Contains(p.fences, fence) {
		p.fences = append(p.fences, fence)
	}
}

func (p *plot) hasFence(fence walldir) bool {
	for _, f := range p.fences {
		if f == fence {
			return true
		}
	}
	return false
}

type garden struct {
	w       int
	h       int
	plots   [][]*plot
	regions map[int][]*plot
	sides   map[int]sides
}

func newGarden(lines []string) *garden {
	plots := make([][]*plot, len(lines))
	for r, line := range lines {
		plots[r] = make([]*plot, len(line))
		for c, plant := range line {
			plots[r][c] = &plot{plant: plant, r: r, c: c}
		}
	}
	return &garden{
		w:       len(lines[0]),
		h:       len(lines),
		plots:   plots,
		regions: make(map[int][]*plot),
		sides:   make(map[int]sides),
	}
}

func (g *garden) plot(r, c int) *plot {
	return g.plots[r][c]
}

func (g *garden) addFences() {
	// add fences around the garden
	for r := 0; r < g.h; r++ {
		g.addFence(r, 0, west)
		g.addFence(r, g.w-1, east)
		if r < g.h-1 && g.plot(r, g.w-1).plant != g.plot(r+1, g.w-1).plant {
			g.addFence(r, g.w-1, south)
			g.addFence(r+1, g.w-1, north)
		}
	}
	for c := 0; c < g.w; c++ {
		g.addFence(0, c, north)
		g.addFence(g.h-1, c, south)
		if c < g.w-1 && g.plot(g.h-1, c).plant != g.plot(g.h-1, c+1).plant {
			g.addFence(g.h-1, c, east)
			g.addFence(g.h-1, c+1, west)
		}
	}
	for r := 0; r < g.h-1; r++ {
		for c := 0; c < g.w-1; c++ {
			if g.plot(r, c).plant != g.plot(r, c+1).plant {
				g.addFence(r, c, east)
				g.addFence(r, c+1, west)
			}
			if g.plot(r, c).plant != g.plot(r+1, c).plant {
				g.addFence(r, c, south)
				g.addFence(r+1, c, north)
			}
		}
	}
}

func (g *garden) print() {
	for r := 0; r < g.h; r++ {
		for c := 0; c < g.w; c++ {
			if g.plot(r, c).hasFence(north) {
				fmt.Print(" -")
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Println(" ")
		for c := 0; c < g.w; c++ {
			if g.plot(r, c).hasFence(west) {
				fmt.Print("|")
			} else {
				fmt.Print(" ")
			}
			fmt.Printf("%c", g.plot(r, c).plant)
		}
		fmt.Println("|")
	}
	for c := 0; c < g.w; c++ {
		fmt.Print(" -")
	}
	fmt.Println(" ")
}

func (g *garden) oneRegion(r, c int, region int) {
	if g.plot(r, c).region != 0 {
		return
	}
	g.plot(r, c).region = region
	if c < g.w-1 && g.plot(r, c).plant == g.plot(r, c+1).plant {
		g.oneRegion(r, c+1, region)
	}
	if r < g.h-1 && g.plot(r, c).plant == g.plot(r+1, c).plant {
		g.oneRegion(r+1, c, region)
	}
	if c > 0 && g.plot(r, c).plant == g.plot(r, c-1).plant {
		g.oneRegion(r, c-1, region)
	}
	if r > 0 && g.plot(r, c).plant == g.plot(r-1, c).plant {
		g.oneRegion(r-1, c, region)
	}
}

func (g *garden) regionize() {
	region := 0
	for r := 0; r < g.h; r++ {
		for c := 0; c < g.w; c++ {
			if g.plot(r, c).region == 0 {
				region++
				g.oneRegion(r, c, region)
			}
		}
	}

	for r := 0; r < g.h; r++ {
		for c := 0; c < g.w; c++ {
			g.regions[g.plot(r, c).region] = append(g.regions[g.plot(r, c).region], g.plot(r, c))
		}
	}
}

func (g *garden) addFence(r, c int, fence walldir) {
	g.plot(r, c).addFence(fence)
	switch fence {
	case north:
		if r > 0 {
			g.plot(r-1, c).addFence(south)
		}
	case east:
		if c < g.w-1 {
			g.plot(r, c+1).addFence(west)
		}
	case south:
		if r < g.h-1 {
			g.plot(r+1, c).addFence(north)
		}
	case west:
		if c > 0 {
			g.plot(r, c-1).addFence(east)
		}
	}
}

func (g *garden) pricePerimeters(debug bool) int {
	perimeters := make(map[int]int)
	for region, plots := range g.regions {
		totalP := 0
		for _, p := range plots {
			totalP += len(p.fences)
		}
		perimeters[region] = totalP
	}

	price := 0
	for region, plots := range g.regions {
		price += len(plots) * perimeters[region]
		if debug {
			fmt.Printf("Region %d (perimeter %d):\n", region, perimeters[region])
			for _, p := range plots {
				fmt.Printf("  %c(%d,%d) [%d] [%v]\n", p.plant, p.r, p.c, len(p.fences), p.fences)
			}
		}
	}
	return price
}

type side struct {
	wall     walldir
	startRow int
	endRow   int
	startCol int
	endCol   int
}

func (s side) String() string {
	return fmt.Sprintf("wall %5v, start (%d,%d), end (%d,%d)", s.wall, s.startRow, s.startCol, s.endRow, s.endCol)
}

type sides []*side

func (s *sides) add(r, c int, wall walldir) {
	for _, t := range *s {
		if t.wall != wall {
			continue
		}
		switch wall {
		case north, south: // horizontal so start and end row are the same
			if t.startRow == r {
				if t.startCol == c+1 {
					t.startCol = c
					return
				}
				if t.endCol == c-1 {
					t.endCol = c
					return
				}
			}
		case east, west: // vertical so start and end col are the same
			if t.startCol == c {
				if t.startRow == r+1 {
					t.startRow = r
					return
				}
				if t.endRow == r-1 {
					t.endRow = r
					return
				}
			}
		}
	}
	// we didn't find a matching side so add a new one
	side := &side{wall: wall, startRow: r, endRow: r, startCol: c, endCol: c}
	*s = append(*s, side)
}

func (g *garden) priceSides(debug bool) int {
	for rix, region := range g.regions {
		sides := sides(make([]*side, 0))
		for _, p := range region {
			for _, f := range p.fences {
				sides.add(p.r, p.c, f)
			}
		}
		g.sides[rix] = sides
	}

	price := 0
	for region, plots := range g.regions {
		price += len(plots) * len(g.sides[region])
		if debug {
			fmt.Printf("Region %d (area %d, sides %d):\n", region, len(g.regions[region]), len(g.sides[region]))
			for _, s := range g.sides[region] {
				fmt.Printf("  %v\n", s)
			}
		}
	}
	return price
}

func part1(lines []string) int {
	g := newGarden(lines)
	g.addFences()
	g.regionize()
	total := g.pricePerimeters(false)
	// g.print()
	return total
}

func part2(lines []string) int {
	g := newGarden(lines)
	g.addFences()
	g.regionize()
	total := g.priceSides(false)
	// g.print()
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
