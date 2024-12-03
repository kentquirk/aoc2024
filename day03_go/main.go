package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

func part1(data string) int {
	pat := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	matches := pat.FindAllStringSubmatch(data, -1)
	total := 0
	for _, match := range matches {
		// fmt.Println(match)
		a, _ := strconv.Atoi(match[1])
		b, _ := strconv.Atoi(match[2])
		total += a * b
	}
	return total
}

func part2(data string) int {
	pat := regexp.MustCompile(`(do(?:n't)?)|(mul\((\d+),(\d+)\))`)
	matches := pat.FindAllStringSubmatch(data, -1)
	total := 0
	enabled := true
	for _, match := range matches {
		// fmt.Printf("%#v\n", match)
		if match[1] == "do" {
			enabled = true
		} else if match[1] == "don't" {
			enabled = false
		} else if enabled {
			a, _ := strconv.Atoi(match[3])
			b, _ := strconv.Atoi(match[4])
			total += a * b
		}
	}
	return total
}

func readall(filename string) string {
	f, err := os.Open(fmt.Sprintf("./data/%s.txt", filename))
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func main() {
	args := os.Args[1:]
	filename := "sample"
	if len(args) > 0 {
		filename = args[0]
	}
	data := readall(filename)
	fmt.Println(part1(data))
	fmt.Println(part2(data))
}
