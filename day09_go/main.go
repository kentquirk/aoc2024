package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

type block struct {
	id         int
	size       int
	offset     int
	considered bool
}

const (
	Free    = -1
	Deleted = -2
)

func (b block) split(size int) (block, block) {
	return block{id: b.id, size: size, offset: b.offset},
		block{id: b.id, size: b.size - size, offset: b.offset + size}
}

type blocklist []block

const debug = false

func (bl blocklist) print(s string) {
	ids := "01234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if !debug {
		return
	}
	sb := bl.sort()
	fmt.Printf("%10s: ", s)
	for _, b := range sb {
		if b.offset > 200 {
			fmt.Printf("(etc)")
			break
		}
		ch := "."
		if ix := b.id % len(ids); ix >= 0 {
			ch = ids[ix : ix+1]
		}
		fmt.Printf("%s", strings.Repeat(ch, b.size))
	}
	fmt.Println()
}

func (bl blocklist) validate() bool {
	offset := 0
	for i, b := range bl {
		if b.offset != offset {
			fmt.Println("Offset mismatch at", i, "expected", offset, "got", b.offset)
			return false
		}
		offset += b.size
	}
	return true
}

func (bl blocklist) totalLength() int {
	l := 0
	for _, b := range bl {
		l += b.size
	}
	return l
}

func (bl blocklist) lastOccupiedBlock() (int, block) {
	for i := len(bl) - 1; i >= 0; i-- {
		if bl[i].id >= 0 {
			return i, bl[i]
		}
	}
	return -1, block{}
}

func (bl blocklist) lastUnconsideredBlock() (int, block) {
	for i := len(bl) - 1; i >= 0; i-- {
		if bl[i].id >= 0 && bl[i].considered == false {
			return i, bl[i]
		}
	}
	return -1, block{}
}

func (bl blocklist) checksum() int {
	ck := 0
	for _, b := range bl {
		if b.id >= 0 {
			for p := 0; p < b.size; p++ {
				ck += b.id * (b.offset + p)
			}
		}
	}
	return ck
}

func (b1 *blocklist) insertAt(b block, ix int) {
	*b1 = append(*b1, block{})
	copy((*b1)[ix+1:], (*b1)[ix:])
	(*b1)[ix] = b
}

func (bl blocklist) sort() blocklist {
	sorted := make([]block, len(bl))
	copy(sorted, bl)
	slices.SortFunc[[]block, block](sorted, func(a, b block) int {
		return a.offset - b.offset
	})
	return sorted
}

func parseData(data string) blocklist {
	offset := 0
	bl := make(blocklist, 0)
	for i := 0; i < len(data); i++ {
		osiz := int(data[i] - '0')
		if osiz != 0 {
			if i%2 == 1 {
				// it's a free space
				bl = append(bl, block{id: Free, size: osiz, offset: offset})
			} else {
				bl = append(bl, block{id: i / 2, size: osiz, offset: offset})
			}
			offset += osiz
		}
	}
	return bl
}

func denselyCompact(bl blocklist) blocklist {
	newlist := make(blocklist, 0)
outer:
	for i := 0; i < len(bl); {
		b := bl[i]
		switch b.id {
		case Deleted:
			i++
			continue
		case Free:
			lastBlockIx, lastBlock := bl.lastOccupiedBlock()
			if lastBlockIx == -1 {
				break outer
			}
			switch {
			case lastBlock.size < b.size:
				// split free block
				b1, b2 := b.split(lastBlock.size)
				b1.id = lastBlock.id
				newlist = append(newlist, b1)
				bl[lastBlockIx].id = Deleted
				bl[i] = b2
				// do not increment i because we replaced b[i] with a smaller one
			case lastBlock.size == b.size:
				// free block is exactly the same size as occupied block
				b.id = lastBlock.id
				newlist = append(newlist, b)
				bl[lastBlockIx].id = Deleted
				i++
			case lastBlock.size > b.size:
				// split occupied block
				b1, b2 := lastBlock.split(b.size)
				b1.offset = b.offset
				newlist = append(newlist, b1)
				bl[lastBlockIx] = b2
				i++
			}
		default:
			// we might have to correct the offset
			if len(newlist) > 0 {
				b.offset = newlist[len(newlist)-1].offset + newlist[len(newlist)-1].size
			}
			newlist = append(newlist, b)
			bl[i].id = Deleted
			i++
		}
		newlist.print("newlist")
	}
	return newlist
}

func firstFit(bl blocklist) blocklist {
	moveCount := 0
	for {
		lastBlockIx, lastBlock := bl.lastUnconsideredBlock()
		if lastBlockIx == -1 {
			break
		}
		bl[lastBlockIx].considered = true
		// find first free block that fits
	inner:
		for i, b := range bl {
			if b.offset >= lastBlock.offset {
				break
			}
			if b.id == Free && b.size >= lastBlock.size {
				switch {
				case b.size > lastBlock.size:
					// split free block
					b1, b2 := b.split(lastBlock.size)
					b1.id = lastBlock.id
					b1.considered = true
					bl[lastBlockIx].id = Deleted
					bl[i] = b2
					moveCount++
					bl.insertAt(b1, i)
					break inner
				case lastBlock.size == b.size:
					// free block is exactly the same size as occupied block
					b.id = lastBlock.id
					b.considered = true
					bl[i] = b
					moveCount++
					bl[lastBlockIx].id = Deleted
					break inner
				}
			}
		}
		bl.print("loop")
		bl.validate()
	}
	fmt.Println("Move count:", moveCount)
	return bl
}

func part1(data string) int {
	bl := parseData(data)
	// fmt.Println(bl)
	bl.print("before")

	bl = denselyCompact(bl)
	bl.print("after")

	return bl.checksum()
}

func part2(data string) int {
	bl := parseData(data)
	// fmt.Println(bl)
	bl.print("before")

	bl = firstFit(bl)
	bl.print("after")
	fmt.Println("Total length:", bl.totalLength())

	return bl.checksum()
}

func readData(filename string) string {
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
	data := readData(filename)
	fmt.Println(part1(data))
	fmt.Println(part2(data))
}
