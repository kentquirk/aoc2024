package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseNumbersFrom(line string) []int {
	parts := strings.Fields(line)
	numbers := make([]int, len(parts))
	for i, part := range parts {
		numbers[i], _ = strconv.Atoi(part)
	}
	return numbers
}

func makeDeltas(numbers []int) []int {
	deltas := make([]int, len(numbers)-1)
	for i := 0; i < len(numbers)-1; i++ {
		deltas[i] = numbers[i+1] - numbers[i]
	}
	return deltas
}

func testSafe(deltas []int) bool {
	dir := 0
	for _, delta := range deltas {
		if delta > 3 || delta < -3 || delta == 0 {
			return false
		}
		if delta*dir < 0 {
			return false
		}
		dir = delta
	}
	return true
}

func part1(lines []string) int {
	nsafe := 0
	for _, line := range lines {
		data := parseNumbersFrom(line)
		deltas := makeDeltas(data)
		if testSafe(deltas) {
			nsafe++
		}
	}
	return nsafe
}

func part2(lines []string) int {
	nsafe := 0
	for _, line := range lines {
		data := parseNumbersFrom(line)
		deltas := makeDeltas(data)
		if testSafe(deltas) {
			nsafe++
		} else {
			for skip := 0; skip < len(data); skip++ {
				// we need to make a new slice without the skipped element
				var testdata []int
				testdata = append(testdata, data[:skip]...)
				testdata = append(testdata, data[skip+1:]...)
				testdeltas := makeDeltas(testdata)
				if testSafe(testdeltas) {
					nsafe++
					break
				}
			}
		}
	}
	return nsafe
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
