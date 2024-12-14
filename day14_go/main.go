package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

func parseNumbersFrom(line string) []int {
	pat := regexp.MustCompile(`[0-9-]+`)
	parts := pat.FindAllString(line, -1)
	numbers := make([]int, len(parts))
	for i, part := range parts {
		numbers[i], _ = strconv.Atoi(part)
	}
	return numbers
}

type robot struct {
	pos point
	vel point
}

func (r *robot) positionAtTime(t int, siz point) point {
	pos := point{(r.pos.x + r.vel.x*t) % siz.x, (r.pos.y + r.vel.y*t) % siz.y}
	if pos.x < 0 {
		pos.x += siz.x
	}
	if pos.y < 0 {
		pos.y += siz.y
	}
	return pos
}

type floor struct {
	siz    point
	robots []robot
}

func newFloor(lines []string) floor {
	nums := parseNumbersFrom(lines[0])
	f := floor{point{nums[0], nums[1]}, []robot{}}
	for _, line := range lines[1:] {
		nums := parseNumbersFrom(line)
		f.robots = append(f.robots, robot{point{nums[0], nums[1]}, point{nums[2], nums[3]}})
	}
	return f
}

func (f *floor) quadrant(pos point) int {
	if pos.x == f.siz.x/2 || pos.y == f.siz.y/2 {
		return -1
	}
	if pos.x < f.siz.x/2 {
		if pos.y < f.siz.y/2 {
			return 0
		}
		return 1
	}
	if pos.y < f.siz.y/2 {
		return 2
	}
	return 3
}

func (f *floor) dangerLevel(t int) int {
	positions := make([]point, len(f.robots))
	for i, r := range f.robots {
		positions[i] = r.positionAtTime(t, f.siz)
	}
	if t == 7037 {
		f.print(positions)
	}

	quads := map[int]int{}
	for _, pos := range positions {
		quads[f.quadrant(pos)]++
	}
	return quads[0] * quads[1] * quads[2] * quads[3]
}

func (f *floor) print(positions []point) {
	floor := make([][]int, f.siz.y)
	for y := 0; y < f.siz.y; y++ {
		floor[y] = make([]int, f.siz.x)
	}
	for _, pos := range positions {
		floor[pos.y][pos.x] = 1
	}
	for y := 0; y < f.siz.y; y++ {
		for x := 0; x < f.siz.x; x++ {
			if floor[y][x] == 1 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func part1(lines []string) int {
	f := newFloor(lines)
	return f.dangerLevel(100)
}

func part2(lines []string) int {
	f := newFloor(lines)
	minDanger := 1000000000
	for t := 0; t < 10000; t++ {
		danger := f.dangerLevel(t)
		if danger <= minDanger {
			// fmt.Println(t, danger)
			minDanger = danger
		}
	}
	return minDanger
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
