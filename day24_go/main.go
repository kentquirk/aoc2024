package main

import (
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
)

type signal struct {
	name     string
	value    bool
	valid    bool
	bitIndex int
}

func (s signal) String() string {
	v := "?"
	if s.valid {
		v = fmt.Sprintf("%t", s.value)
	}
	if s.bitIndex != -1 {
		return fmt.Sprintf("%s(%d): %s", s.name, s.bitIndex, v)
	}
	return fmt.Sprintf("%s: %s", s.name, v)
}

type gate struct {
	op     string
	inputs []string
	output string
}

func (g gate) String() string {
	return fmt.Sprintf("%s %s %s -> %s", g.inputs[0], g.op, g.inputs[1], g.output)
}

func (g gate) graph() []string {
	result := []string{}
	result = append(result, fmt.Sprintf("%s -> %s\n", g.inputs[0], g.output))
	result = append(result, fmt.Sprintf("%s -> %s\n", g.inputs[1], g.output))
	shape := ""
	style := "filled"
	fillcolor := "white"
	switch g.op {
	case "AND":
		shape = "house"
		fillcolor = "lightblue"
	case "OR":
		shape = "octagon"
		fillcolor = "lightgreen"
	case "XOR":
		shape = "diamond"
		fillcolor = "pink"
	}
	if strings.HasPrefix(g.output, "z") {
		fillcolor = "yellow"
	}
	line := fmt.Sprintf("%s [shape=%s, style=%s, fillcolor=%s]\n", g.output, shape, style, fillcolor)
	result = append(result, line)
	return result
}

type system struct {
	signals map[string]signal
	gates   []gate
}

func (s *system) String() string {
	var out strings.Builder
	keys := slices.Collect(maps.Keys(s.signals))
	slices.Sort(keys)
	for _, k := range keys {
		out.WriteString(s.signals[k].String())
		out.WriteString("\n")
	}
	for _, g := range s.gates {
		out.WriteString(g.String())
		out.WriteString("\n")
	}
	return out.String()
}

func (s *system) graph() string {
	result := []string{}
	for _, g := range s.gates {
		result = append(result, g.graph()...)
	}

	slices.SortFunc(result, func(a, b string) int {
		keyfunc := func(s string) string {
			prefix := "2"
			if strings.Contains(s, "[") {
				prefix = "9"
			} else if s[0] == 'x' || s[0] == 'y' {
				prefix = "1"
			} else if s[0] == 'z' {
				prefix = "8"
			}
			return fmt.Sprintf("%s%s", prefix, s)
		}
		return strings.Compare(keyfunc(a), keyfunc(b))
	})

	var out strings.Builder
	out.WriteString("digraph G {\n")
	for _, line := range result {
		out.WriteString(line)
	}
	out.WriteString("}\n")
	return out.String()
}

func (s *system) addGate(op string, inputs []string, output string) {
	if inputs[0] > inputs[1] {
		inputs[0], inputs[1] = inputs[1], inputs[0]
	}
	s.gates = append(s.gates, gate{op: op, inputs: inputs, output: output})
}

func (s *system) addSignal(name string) {
	// don't add if already exists
	if _, ok := s.signals[name]; ok {
		return
	}
	index := -1
	if name[0] == 'x' || name[0] == 'y' || name[0] == 'z' {
		index, _ = strconv.Atoi(name[1:])
	}
	s.signals[name] = signal{name: name, value: false, valid: false, bitIndex: index}
}

func (s *system) setSignal(name string, value bool) {
	sig := s.signals[name]
	sig.value = value
	sig.valid = true
	s.signals[name] = sig
}

func (s *system) getValue() int {
	var value int
	for _, sig := range s.signals {
		if sig.name[0] == 'z' && sig.bitIndex != -1 && sig.value {
			value |= 1 << sig.bitIndex
		}
	}
	return value
}

func (s *system) set(prefix string, value int) {
	for _, sig := range s.signals {
		if sig.bitIndex != -1 && strings.HasPrefix(sig.name, prefix) {
			sig.value = (value & (1 << sig.bitIndex)) != 0
			sig.valid = true
			s.signals[sig.name] = sig
		}
	}
}

func (s *system) step() bool {
	changed := false
	for _, g := range s.gates {
		// don't recompute signals
		if s.signals[g.output].valid {
			continue
		}
		// don't compute if inputs are not valid
		if !s.signals[g.inputs[0]].valid || !s.signals[g.inputs[1]].valid {
			continue
		}
		output := false
		if g.op == "AND" {
			output = s.signals[g.inputs[0]].value && s.signals[g.inputs[1]].value
		} else if g.op == "OR" {
			output = s.signals[g.inputs[0]].value || s.signals[g.inputs[1]].value
		} else if g.op == "XOR" {
			output = s.signals[g.inputs[0]].value != s.signals[g.inputs[1]].value
		}
		s.setSignal(g.output, output)
		changed = true
	}
	return changed
}

func parseLines(lines []string) *system {
	s := &system{make(map[string]signal), []gate{}}
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ": ")
			name := parts[0]
			value := parts[1]
			s.addSignal(name)
			if value == "1" {
				s.setSignal(name, true)
			} else {
				s.setSignal(name, false)
			}
		}

		if strings.Contains(line, "->") {
			parts := strings.Fields(line)
			s.addGate(parts[1], []string{parts[0], parts[2]}, parts[4])
			s.addSignal(parts[0])
			s.addSignal(parts[2])
			s.addSignal(parts[4])
		}
	}
	return s
}

func part1(lines []string, x, y int) int {
	s := parseLines(lines)
	s.set("x", x)
	s.set("y", y)
	fmt.Println(s)
	for s.step() {
	}
	fmt.Println(s)
	return s.getValue()
}

func part2(lines []string) int {
	s := parseLines(lines)

	fmt.Println(s.graph())
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

// I did part 2 manually by feeding the data to a dot file and looking for inconsistencies using graphviz
// on GraphvizOnline. The dot engine was useful, but the fdp engine was better for this particular problem.
func main() {
	args := os.Args[1:]
	filename := "inputFixed"
	if len(args) > 0 {
		filename = args[0]
	}
	lines := readlines(filename)
	// x := 0x3
	// y := 0x3
	// if len(args) > 1 {
	// 	x, _ = strconv.Atoi(args[1])
	// 	y, _ = strconv.Atoi(args[2])
	// }
	// answer := part1(lines, x, y)
	// fmt.Printf("%d + %d = %d\n", x, y, answer)
	// fmt.Printf("%x + %x = %x\n", x, y, answer)
	part2(lines)
}
