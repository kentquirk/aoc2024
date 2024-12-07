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

type stack []int

func (s *stack) push(vs ...int) {
	*s = append(*s, vs...)
}

func (s *stack) pop() int {
	if len(*s) == 0 {
		return 0
	}
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *stack) len() int {
	return len(*s)
}

func (s *stack) copy() stack {
	return append(stack{}, *s...)
}

func parseNumbersFrom(line string) []int {
	pat := regexp.MustCompile(`\d+`)
	parts := pat.FindAllString(line, -1)
	numbers := make([]int, len(parts))
	for i, part := range parts {
		numbers[i], _ = strconv.Atoi(part)
	}
	return numbers
}

func reverse(s []int) []int {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func combinations(ops string, n int) []string {
	if n == 0 {
		return []string{""}
	}
	if len(ops) == 0 {
		return []string{}
	}
	combs := []string{}
	for i := range ops {
		for _, rest := range combinations(ops, n-1) {
			combs = append(combs, ops[i:i+1]+rest)
		}
	}
	return combs
}

func calcTest(operators string, values stack, target int) bool {
	combos := combinations(operators, values.len()-1)
	for _, ops := range combos {
		vs := values.copy()
		for _, op := range ops {
			v1 := vs.pop()
			v2 := vs.pop()
			switch op {
			case '+':
				vs.push(v1 + v2)
			case '*':
				vs.push(v1 * v2)
			case '|':
				v, _ := strconv.Atoi(fmt.Sprintf("%d%d", v1, v2))
				vs.push(v)
			}
			ops = ops[1:]
		}
		// fmt.Println(i, values, vs, target)
		if vs.pop() == target {
			return true
		}
	}
	return false
}

func doIt(operators string, lines []string) int {
	total := 0
	for _, line := range lines {
		numbers := parseNumbersFrom(line)
		result := numbers[0]
		values := stack{}
		// we need to evaluate l-r so we reverse the numbers
		r := reverse(numbers[1:])
		values.push(r...)
		if calcTest(operators, values, result) {
			total += result
		}
	}
	return total
}

func part1(lines []string) int {
	return doIt("*+", lines)
}

func part2(lines []string) int {
	return doIt("*+|", lines)
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
