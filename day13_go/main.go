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

const (
	costForA = 3
	costForB = 1
)

type point struct {
	x int
	y int
}

type problem struct {
	a     point
	b     point
	prize point
}

// find the integer solutions A and B to the system where
// A*a.x + B*b.x = prize.x
// A*a.y + B*b.y = prize.y
func (p problem) Solve() []point {
	solutions := make([]point, 0)

	// calculate the max value of A
	var maxA int
	if p.a.x < p.b.x {
		maxA = p.prize.x / p.a.x
	} else {
		maxA = p.prize.x / p.b.x
	}
	for a := maxA; a >= 0; a-- {
		r := p.prize.x - a*p.a.x
		if r%p.b.x == 0 {
			b := r / p.b.x
			if a*p.a.y+b*p.b.y == p.prize.y {
				if a > 100 || b > 100 {
					continue
				}
				// fmt.Println(a, b)
				// fmt.Println(a*p.a.x+b*p.b.x, a*p.a.y+b*p.b.y)
				solutions = append(solutions, point{x: a, y: b})
			}
		}
	}
	return solutions
}

func (p problem) Solve2() *point {
	B := (p.prize.y*p.a.x - p.prize.x*p.a.y) / (p.b.y*p.a.x - p.b.x*p.a.y)
	A := (p.prize.x - B*p.b.x) / p.a.x
	if A*p.a.x+B*p.b.x == p.prize.x && A*p.a.y+B*p.b.y == p.prize.y {
		// fmt.Println(p, A, B)
		return &point{x: A, y: B}
	}
	// fmt.Println(p, A, B, "no solution")
	// fmt.Println(A*p.a.x+B*p.b.x, p.prize.x, A*p.a.y+B*p.b.y, p.prize.y)

	return nil
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func parseProblems(lines []string) []problem {
	pat := regexp.MustCompile(`\d+`)
	problems := make([]problem, 0)
	for i := 0; i < len(lines); {
		if lines[i] == "" {
			i++
			continue
		}
		a := pat.FindAllString(lines[i], -1)
		b := pat.FindAllString(lines[i+1], -1)
		prize := pat.FindAllString(lines[i+2], -1)
		problems = append(problems, problem{
			a:     point{x: toInt(a[0]), y: toInt(a[1])},
			b:     point{x: toInt(b[0]), y: toInt(b[1])},
			prize: point{x: toInt(prize[0]), y: toInt(prize[1])},
		})
		i += 3
	}
	return problems
}

func cost(p point) int {
	return costForA*p.x + costForB*p.y
}

func bestSolution(solutions []point) point {
	best := solutions[0]
	for _, s := range solutions {
		if cost(s) < cost(best) {
			best = s
		}
	}
	return best
}

func part1(lines []string) int {
	problems := parseProblems(lines)
	total := 0
	for _, p := range problems {
		solutions := p.Solve()
		if len(solutions) == 0 {
			fmt.Println(p, "no solution")
			continue
		}
		best := bestSolution(solutions)
		total += cost(best)

		// check solve2
		solution := p.Solve2()
		if solution != nil {
			if cost(*solution) != cost(best) {
				fmt.Println("failure", p, best, cost(best), solution, cost(*solution))
			}
		}
		fmt.Println(p, best, cost(best))
	}
	return total
}

func part2(lines []string) int {
	problems := parseProblems(lines)
	total := 0
	for _, p := range problems {
		p.prize.x += 10_000_000_000_000
		p.prize.y += 10_000_000_000_000
		solution := p.Solve2()
		if solution == nil {
			// fmt.Println(p, "no solution")
			continue
		}
		total += cost(*solution)
		// fmt.Println(p, solution, cost(*solution))
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
	fmt.Println("---")
	fmt.Println(part2(lines))
}
