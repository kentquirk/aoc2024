package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type opcode byte

const (
	adv opcode = 0
	bxl opcode = 1
	bst opcode = 2
	jnz opcode = 3
	bxc opcode = 4
	out opcode = 5
	bdv opcode = 6
	cdv opcode = 7
)

func (op opcode) String() string {
	return []string{"adv", "bxl", "bst", "jnz", "bxc", "out", "bdv", "cdv"}[op]
}

type operand byte

func (op operand) String() string {
	switch op {
	case 0, 1, 2, 3:
		return fmt.Sprintf("%d", op)
	case 4:
		return "A"
	case 5:
		return "B"
	case 6:
		return "C"
	case 7:
		return "?"
	}
	return "??"
}

type vm struct {
	pc        int
	registers map[string]int
	code      []byte
	output    []int
}

func (v *vm) GetOperandValueAt(ix int) int {
	op := v.code[ix]
	switch op {
	case 0, 1, 2, 3:
		return int(op)
	case 4:
		return v.registers["A"]
	case 5:
		return v.registers["B"]
	case 6:
		return v.registers["C"]
	case 7:
		return 0
	}
	return 0
}

func (v *vm) GetLiteralValueAt(ix int) int {
	return int(v.code[ix])
}

func (v *vm) Step() bool {
	switch opcode(v.code[v.pc]) {
	case adv:
		numer := v.registers["A"]
		denom := 1 << v.GetOperandValueAt(v.pc+1)
		v.registers["A"] = numer / denom
		v.pc += 2
	case bxl:
		result := v.registers["B"] ^ v.GetLiteralValueAt(v.pc+1)
		v.registers["B"] = result
		v.pc += 2
	case bst:
		v.registers["B"] = v.GetOperandValueAt(v.pc+1) & 0x7
		v.pc += 2
	case jnz:
		if v.registers["A"] == 0 {
			v.pc += 2
		} else {
			v.pc = v.GetLiteralValueAt(v.pc + 1)
		}
	case bxc:
		result := v.registers["B"] ^ v.registers["C"]
		v.registers["B"] = result
		v.pc += 2
	case out:
		v.output = append(v.output, v.GetOperandValueAt(v.pc+1)&0x7)
		v.pc += 2
	case bdv:
		numer := v.registers["A"]
		denom := 1 << v.GetOperandValueAt(v.pc+1)
		v.registers["B"] = numer / denom
		v.pc += 2
	case cdv:
		numer := v.registers["A"]
		denom := 1 << v.GetOperandValueAt(v.pc+1)
		v.registers["C"] = numer / denom
		v.pc += 2
	}
	return v.pc < len(v.code)-1
}

func (v *vm) Run() {
	for v.Step() {
	}
}

func (v *vm) Print() {
	for name, value := range v.registers {
		fmt.Printf("%s: %d\n", name, value)
	}
	for i := 0; i < len(v.code); i += 2 {
		caret := " "
		if i == v.pc {
			caret = ">"
		}
		fmt.Printf(" %s  %3s %1s    [%d %d] \n", caret, opcode(v.code[i]), operand(v.code[i+1]), v.code[i], v.code[i+1])
	}
	fmt.Println(v.output)
	fmt.Println()
}

func (v *vm) Reset(registers map[string]int) {
	v.pc = 0
	v.output = make([]int, 0)
	for name, value := range registers {
		v.registers[name] = value
	}
}

func (v *vm) Quine() bool {
	if len(v.output) != len(v.code) {
		return false
	}
	for i := 0; i < len(v.output); i++ {
		if v.output[i] != int(v.code[i]) {
			return false
		}
	}
	return true
}

func (v *vm) RunWith(regs map[string]int, regA int) {
	regs["A"] = regA
	v.Reset(regs)
	v.Run()
}

func (v *vm) Quineiness() int {
	if len(v.output) != len(v.code) {
		return 0
	}
	identical := 0
	for i := 0; i < len(v.output); i++ {
		if v.output[i] != int(v.code[i]) {
			break
		}
		identical++
	}
	return identical
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

func loadProgram(lines []string) *vm {
	vm := &vm{
		registers: make(map[string]int),
		code:      make([]byte, 0),
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "Register") {
			pat := regexp.MustCompile(`Register (\w+): (\d+)`)
			matches := pat.FindStringSubmatch(line)
			vm.registers[matches[1]], _ = strconv.Atoi(matches[2])
			continue
		}
		if strings.HasPrefix(line, "Program") {
			code := parseNumbersFrom(line)
			for i := 0; i < len(code); i++ {
				vm.code = append(vm.code, byte(code[i]))
			}
		}
	}
	return vm
}

func part1(lines []string) int {
	vm := loadProgram(lines)
	vm.Print()
	for vm.Step() {
		// vm.Print()
	}
	vm.Print()
	var outputs []string
	for _, value := range vm.output {
		outputs = append(outputs, fmt.Sprintf("%d", value))
	}
	fmt.Println(strings.Join(outputs, ","))
	return 0
}

// This is ugly. The VM is fine, worked, runs fast, but clearly the problem
// wasn't going to produce a quine quickly, because after watching every
// millionth run, it became clear that we needed very large numbers to get
// enough digits in the output. So I searched for the points where the number of
// digits changed, and then looked at them in hex, and realized after a bit that
// it would be better to look in octal. Given a 3 bit machine, octal made sense.
// Watching the outputs in octal, I realized that we were seeing the values get
// stable from the right. I tried reversing a counter, but it didn't do what I
// wanted, so I wrote a "quininess" function that told me how many consecutive
// digits were correct and randomly sampled until I found a couple that had at
// least half the digits right. I basically kept searching in a range,
// tightening the range, until I managed to find a quine. That was too big, but
// now I had a pattern I could match and test all the possibilities in the
// remaining bits. I don't really feel like coding that search so it's
// repeatable, though.
func part2(lines []string) int {
	vm := loadProgram(lines)
	registers := make(map[string]int)
	registers["A"] = 0
	registers["B"] = vm.registers["B"]
	registers["C"] = vm.registers["C"]

	bot := 0o0777777777777777
	top := 0o7777777777777777
	// bot := 0o3477556042247155
	// top := 0o3477556062247277

	// best := 0
	// for i := bot; i < top; i++ {
	// 	vm.RunWith(registers, i)
	// 	quineiness := vm.Quineiness()
	// 	if quineiness > best {
	// 		best = quineiness
	// 		fmt.Printf("%o, %d\n", i, best)
	// 	}
	// 	if vm.Quine() {
	// 		fmt.Printf("%o\n", i)
	// 		break
	// 	}
	// }
	// os.Exit(0)

	ntests := 500000
	endings := []int{
		0o56052247155,
		0o56052247277,
	}

	// this was a valid quine but not the lowest one
	// so we keep the lower bits and search for the upper bits
	// the correct value was 0o6117156052247277
	theQuine := 0o6517156052247277
	maskBits := 33
	for i := 0; i < 1<<16; i++ {
		test := i<<maskBits | (theQuine & ((1 << maskBits) - 1))
		vm.RunWith(registers, test)
		if vm.Quine() {
			fmt.Printf("found it 0o%o, %d\n", test, test)
			vm.Print()
			os.Exit(0)
		}
		quineiness := vm.Quineiness()
		if quineiness > 14 {
			fmt.Printf("%d\n", quineiness)
		}
	}
	os.Exit(0)

	for {
		var best1, best2 int
		var q1, q2 int
		spread := top - bot
		for i := 0; i < ntests; i++ {
			test := bot + rand.Intn(spread)
			test &= ^0o3777777777
			test |= endings[rand.Intn(len(endings))]
			vm.RunWith(registers, test)
			quineiness := vm.Quineiness()
			if quineiness > best1 {
				best2 = best1
				q2 = q1
				best1 = quineiness
				q1 = test
			} else if quineiness > best2 {
				best2 = quineiness
				q2 = test
			}
			if quineiness >= 11 {
				fmt.Printf("%d %o\n", quineiness, test)
			}
			// 127523764129389
			// 127523764129471
			// if quineiness > 3 {
			// 	fmt.Println(strings.Repeat("=", quineiness))
			// }
		}
		// fmt.Printf("%o (%d), %o(%d)\n", q1, best1, q2, best2)
		if q2 < q1 {
			q1, q2 = q2, q1
		}
		bot = q1
		top = q2
		fmt.Println("now ", bot, top)
		if best1 == 16 {
			break
		}
	}
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
	fmt.Println(part2(lines))
}
