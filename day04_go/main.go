package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type pair struct {
	r int
	c int
}

var masOffsets = [][3]pair{
	{{0, 1}, {0, 2}, {0, 3}},       // horizontal right
	{{1, 1}, {2, 2}, {3, 3}},       // diagonal down right
	{{1, 0}, {2, 0}, {3, 0}},       // vertical down
	{{-1, 1}, {-2, 2}, {-3, 3}},    // diagonal down left
	{{0, -1}, {0, -2}, {0, -3}},    // horizontal left
	{{-1, -1}, {-2, -2}, {-3, -3}}, // diagonal up left
	{{-1, 0}, {-2, 0}, {-3, 0}},    // vertical up
	{{1, -1}, {2, -2}, {3, -3}},    // diagonal up right
}

func getLetterAt(lines []string, r, c int) (byte, bool) {
	if r < 0 || r >= len(lines) || c < 0 || c >= len(lines[0]) {
		return 0, false
	}
	return lines[r][c], true
}

func hasLetterAt(lines []string, letter byte, r, c int) bool {
	if r < 0 || r >= len(lines) || c < 0 || c >= len(lines[0]) {
		return false
	}
	return lines[r][c] == letter
}

func countXMASesFrom(lines []string, r int, c int) int {
	letters := "XMAS"
	count := 0
	letter := lines[r][c]
	if letter != 'X' {
		return 0
	}
	for _, offsets := range masOffsets {
		hasLetters := true
		for i, offset := range offsets {
			if !hasLetterAt(lines, letters[i+1], r+offset.r, c+offset.c) {
				hasLetters = false
				break
			}
		}
		if hasLetters {
			count++
		}
	}
	return count
}

func countXMASes(lines []string) int {
	count := 0
	for r := 0; r < len(lines); r++ {
		for c := 0; c < len(lines[0]); c++ {
			count += countXMASesFrom(lines, r, c)
		}
	}
	return count
}

func countMASXesFrom(lines []string, r int, c int) int {
	alt := map[byte]byte{'M': 'S', 'S': 'M'}
	letter := lines[r][c]
	if letter != 'A' {
		return 0
	}

	letter1, ok1 := getLetterAt(lines, r-1, c-1)
	letter2, ok2 := getLetterAt(lines, r+1, c+1)
	if !ok1 || !ok2 {
		return 0
	}
	if letter1 != alt[letter2] || letter2 != alt[letter1] {
		return 0
	}

	letter1, ok1 = getLetterAt(lines, r-1, c+1)
	letter2, ok2 = getLetterAt(lines, r+1, c-1)
	if !ok1 || !ok2 {
		return 0
	}
	if letter1 != alt[letter2] || letter2 != alt[letter1] {
		return 0
	}

	return 1
}

func countMASXes(lines []string) int {
	count := 0
	for r := 0; r < len(lines); r++ {
		for c := 0; c < len(lines[0]); c++ {
			count += countMASXesFrom(lines, r, c)
		}
	}
	return count
}

func part1(lines []string) int {
	return countXMASes(lines)
}

func part2(lines []string) int {
	return countMASXes(lines)
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
