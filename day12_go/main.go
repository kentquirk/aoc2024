package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

const (
	north = 0
	east  = 1
	south = 2
	west  = 3
)

type plot struct {
	plant  rune
	r      int
	c      int
	region int
	fences []int
}

func (p *plot) addFence(fence int) {
	if !slices.Contains(p.fences, fence) {
		p.fences = append(p.fences, fence)
	}
}

func (p *plot) hasFence(fence int) bool {
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
}

func newGarden(lines []string) *garden {
	plots := make([][]*plot, len(lines))
	for r, line := range lines {
		plots[r] = make([]*plot, len(line))
		for c, plant := range line {
			plots[r][c] = &plot{plant: plant, r: r, c: c}
		}
	}
	return &garden{w: len(lines[0]), h: len(lines), plots: plots}
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
}

func (g *garden) addFence(r, c int, fence int) {
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

func (g *garden) priceRegions(debug bool) int {
	regions := make(map[int][]*plot)
	for r := 0; r < g.h; r++ {
		for c := 0; c < g.w; c++ {
			regions[g.plot(r, c).region] = append(regions[g.plot(r, c).region], g.plot(r, c))
		}
	}
	perimeters := make(map[int]int)
	for region, plots := range regions {
		totalP := 0
		for _, p := range plots {
			totalP += len(p.fences)
		}
		perimeters[region] = totalP
	}

	price := 0
	for region, plots := range regions {
		price += len(plots) * perimeters[region]
		if debug {
			fmt.Printf("Region %d (perimeter %d):\n", region, perimeters[region])
			for _, p := range plots {
				fmt.Printf("  %c(%d,%d) [%d]\n", p.plant, p.r, p.c, len(p.fences))
			}
		}
	}
	return price
}

func part1(lines []string) int {
	g := newGarden(lines)
	g.addFences()
	g.regionize()
	total := g.priceRegions(false)
	// g.print()
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
