package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type position struct {
	r int
	c int
}

type direction byte

const (
	up    direction = '^'
	down  direction = 'v'
	left  direction = '<'
	right direction = '>'
)

type guardState int

const (
	moving guardState = iota
	offmap
	looping
)

func (d direction) String() string {
	return string(d)
}

func (d direction) turnRight() direction {
	switch d {
	case up:
		return right
	case right:
		return down
	case down:
		return left
	case left:
		return up
	}
	return d
}

type directions map[direction]struct{}

func (d directions) add(dir direction) bool {
	if _, ok := d[dir]; ok {
		return false
	}
	d[dir] = struct{}{}
	return true
}

func (d directions) ch() direction {
	lr := false
	ud := false
	for dir := range d {
		switch dir {
		case up, down:
			ud = true
		case left, right:
			lr = true
		}
	}
	if lr && ud {
		return '+'
	}
	if lr {
		return '-'
	}
	return '|'
}

type guard struct {
	pos position
	dir direction
}

func (g guard) String() string {
	return fmt.Sprintf("Guard at %v facing %v", g.pos, g.dir)
}

type lab struct {
	w              int
	h              int
	m              map[position]string
	g              guard
	gOriginal      guard
	guardPositions map[position]directions
}

func NewLab(w, h int) *lab {
	return &lab{
		w:              w,
		h:              h,
		m:              make(map[position]string),
		g:              guard{},
		gOriginal:      guard{},
		guardPositions: make(map[position]directions),
	}
}

func (l *lab) Reset() {
	l.g = l.gOriginal
	l.guardPositions = make(map[position]directions)
}

func (l *lab) Print() {
	for r := 0; r < l.h; r++ {
		for c := 0; c < l.w; c++ {
			pos := position{r, c}
			if l.g.pos == pos {
				fmt.Printf("%c", l.g.dir)
			} else if _, ok := l.m[pos]; ok {
				fmt.Printf("%c", '#')
			} else if _, ok := l.guardPositions[pos]; ok {
				fmt.Printf("%c", l.guardPositions[pos].ch())
			} else {
				fmt.Printf("%c", '.')
			}
		}
		fmt.Println()
	}
}

func (l *lab) numPositions() int {
	return len(l.guardPositions)
}

func (l *lab) guardInBounds() bool {
	return l.g.pos.r >= 0 && l.g.pos.r < l.h && l.g.pos.c >= 0 && l.g.pos.c < l.w
}

// false if we couldn't add a new position/direction combination -- guard is in a loop
func (l *lab) recordGuard() bool {
	if _, ok := l.guardPositions[l.g.pos]; !ok {
		l.guardPositions[l.g.pos] = make(directions)
	}
	return l.guardPositions[l.g.pos].add(l.g.dir)
}

func (l *lab) move() guardState {
	newPos := l.g.pos
	switch l.g.dir {
	case up:
		newPos.r--
	case down:
		newPos.r++
	case left:
		newPos.c--
	case right:
		newPos.c++
	}
	if _, ok := l.m[newPos]; ok {
		l.g.dir = l.g.dir.turnRight()
		l.recordGuard()
		return l.move()
	}
	l.g.pos = newPos
	if !l.guardInBounds() {
		return offmap
	}
	if !l.recordGuard() {
		// fmt.Println("guard is in a loop!")
		return looping
	}
	return moving
}

func parseLab(lines []string) *lab {
	l := NewLab(len(lines[0]), len(lines))
	for r, line := range lines {
		for c, char := range line {
			switch char {
			case '#':
				l.m[position{r, c}] = "#"
			case '^':
				l.g.pos = position{r, c}
				l.g.dir = up
				l.recordGuard()
			default:
				// do nothing
			}
		}
	}
	l.gOriginal = l.g
	return l
}

func part1(lines []string) int {
	l := parseLab(lines)
	state := l.move()
	for ; state == moving; state = l.move() {
		// fmt.Println()
	}
	if state == looping {
		fmt.Println("guard is in a loop!")
	}
	l.Print()
	return l.numPositions()
}

func part2(lines []string) int {
	l := parseLab(lines)
	// do one pass to find the guard's possible positions
	state := l.move()
	for ; state == moving; state = l.move() {
		// fmt.Println()
	}
	if state == looping {
		log.Fatal("it's already broken!")
	}
	// l.Print()
	// now try all the possible positions for a blocker and count the ones that cause a loop
	loopCount := 0
	for pos := range l.guardPositions {
		l.Reset()
		l.m[pos] = "#"
		state := l.move()
		for ; state == moving; state = l.move() {
		}
		if state == looping {
			loopCount++
			// l.Print()
			// fmt.Println()
		}
		delete(l.m, pos)
	}

	return loopCount
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
