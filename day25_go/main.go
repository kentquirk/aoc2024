package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	KEY int = iota
	LOCK
)

const (
	FULL    = 0b11111
	EMPTY   = 0b00000
	COLUMN0 = 0b10000
)

type shape struct {
	id   string
	cols []int
}

func newShape(cols []int) *shape {
	id := ""
	for i := 0; i < len(cols); i++ {
		id += fmt.Sprintf("%d", cols[i])
	}
	return &shape{id: id, cols: cols}
}

func (s *shape) String() string {
	return fmt.Sprintf("%s %v", s.id, s.cols)
}

// shapeTree is a recursive data that inserts shapes into a tree
// based on the values of successive columns in the shape
type shapeTree struct {
	subtrees map[int]*shapeTree
	shapes   []*shape
}

func newShapeTree() *shapeTree {
	return &shapeTree{}
}

func (t *shapeTree) add(cols []int, shape *shape) {
	if len(cols) == 0 {
		t.shapes = append(t.shapes, shape)
		return
	}
	if t.subtrees == nil {
		t.subtrees = make(map[int]*shapeTree)
	}
	if _, ok := t.subtrees[cols[0]]; !ok {
		t.subtrees[cols[0]] = new(shapeTree)
	}
	t.subtrees[cols[0]].add(cols[1:], shape)
}

func (t *shapeTree) len() int {
	n := len(t.shapes)
	for _, subtree := range t.subtrees {
		n += subtree.len()
	}
	return n
}

func (t *shapeTree) countFits(index int, other *shape) int {
	if index == len(other.cols) {
		return len(t.shapes)
	}
	totalFits := 0
	for ix := 0; ix <= 5-other.cols[index]; ix++ {
		if subtree, ok := t.subtrees[ix]; ok {
			totalFits += subtree.countFits(index+1, other)
		}
	}
	return totalFits
}

func (t *shapeTree) all() []*shape {
	shapes := make([]*shape, 0)
	for _, shape := range t.shapes {
		shapes = append(shapes, shape)
	}
	for _, subtree := range t.subtrees {
		shapes = append(shapes, subtree.all()...)
	}
	return shapes
}

func (t *shapeTree) Print(indent string) {
	for _, shape := range t.shapes {
		fmt.Println(indent, shape)
	}
	for k, subtree := range t.subtrees {
		fmt.Println(indent, k)
		subtree.Print(indent + "  ")
	}
}

func parseOne(lines []string) (*shape, int) {
	rows := make([]int, 0)
	for i := 0; i < len(lines); i++ {
		chars := strings.Split(lines[i], "")
		v := 0
		mask := COLUMN0
		for j := 0; j < len(chars); j++ {
			if chars[j] == "#" {
				v |= mask
			}
			mask >>= 1
		}
		rows = append(rows, v)
	}
	cols := make([]int, 5)
	// ignore the first and last rows
	for row := 1; row < len(rows)-1; row++ {
		i := 0
		for colMask := COLUMN0; colMask > 0; colMask >>= 1 {
			if rows[row]&colMask != 0 {
				cols[i]++
			}
			i++
		}
	}
	if rows[0] == FULL {
		// top row is full so it's a lock
		return newShape(cols), LOCK
	} else {
		// top row is empty so it's a key
		return newShape(cols), KEY
	}
}

func part1(lines []string) int {
	keys := newShapeTree()
	locks := newShapeTree()
	for i := 0; i < len(lines); i += 8 {
		shape, kind := parseOne(lines[i : i+7])
		if kind == KEY {
			keys.add(shape.cols, shape)
		} else {
			locks.add(shape.cols, shape)
		}
	}
	fmt.Println("keys", keys.len())
	// keys.Print("")
	fmt.Println("locks", locks.len())
	// fmt.Println(locks.all())
	// locks.Print("")

	totalFits := 0
	for _, key := range keys.all() {
		totalFits += locks.countFits(0, key)
		// fmt.Println(key.id, locks.countFits(0, key))
	}
	return totalFits
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
