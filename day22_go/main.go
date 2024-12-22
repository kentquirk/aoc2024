package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type monkey struct {
	prices map[string]int
}

func newMonkey(secret int) *monkey {
	const nsteps = 2001
	m := &monkey{prices: map[string]int{}}
	secrets := make([]int, nsteps+1)
	for i := 0; i < nsteps+1; i++ {
		secrets[i] = secret
		secret = step(secret)
	}
	for i := 4; i < nsteps+1; i++ {
		d1 := secrets[i-3]%10 - secrets[i-4]%10
		d2 := secrets[i-2]%10 - secrets[i-3]%10
		d3 := secrets[i-1]%10 - secrets[i-2]%10
		d4 := secrets[i]%10 - secrets[i-1]%10
		key := fmt.Sprintf("%x", []int{d1, d2, d3, d4})
		if _, ok := m.prices[key]; !ok {
			// only store the first one
			m.prices[key] = secrets[i] % 10
		}
	}
	return m
}

func (m *monkey) print() {
	fmt.Println("Monkey")
	fmt.Println(m.prices)
	for k, v := range m.prices {
		fmt.Printf("%s: %d\n", k, v)
	}
}

func step(x int) int {
	x1 := ((x << 6) ^ x) & 0xFFFFFF
	x2 := ((x1 >> 5) ^ x1) & 0xFFFFFF
	x3 := ((x2 << 11) ^ x2) & 0xFFFFFF
	// fmt.Printf("%x %x %x %x\n", x, x1, x2, x3)
	return x3
}

func nsteps(x int, n int) int {
	for i := 0; i < n; i++ {
		x = step(x)
	}
	return x
}

func part1(lines []string) int {
	total := 0
	for _, line := range lines {
		x, _ := strconv.Atoi(line)
		x1 := nsteps(x, 2000)
		// fmt.Println(x1)
		total += x1
	}
	return total
}

func part2(lines []string) int {
	monkeys := make([]*monkey, len(lines))
	for i, line := range lines {
		x, _ := strconv.Atoi(line)
		m := newMonkey(x)
		monkeys[i] = m
		// m.print()
	}
	allkeys := make(map[string]struct{})
	for _, m := range monkeys {
		for k := range m.prices {
			allkeys[k] = struct{}{}
		}
	}

	totals := make(map[string]int)
	for k := range allkeys {
		total := 0
		for _, m := range monkeys {
			total += m.prices[k]
		}
		totals[k] = total
	}

	maxkey := ""
	maxvalue := 0
	for k, v := range totals {
		if v > maxvalue {
			maxkey = k
			maxvalue = v
		}
	}
	fmt.Println(maxkey, maxvalue)
	return maxvalue
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
	fmt.Println("Part 1: ", part1(lines))
	fmt.Println("Part 2: ", part2(lines))
}
