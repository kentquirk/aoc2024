package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

type graph struct {
	nodes map[string]map[string]struct{}
}

func newGraph() *graph {
	return &graph{nodes: make(map[string]map[string]struct{})}
}

func (g *graph) addPair(n, m string) {
	if _, ok := g.nodes[n]; !ok {
		g.nodes[n] = make(map[string]struct{})
	}
	g.nodes[n][m] = struct{}{}
	if _, ok := g.nodes[m]; !ok {
		g.nodes[m] = make(map[string]struct{})
	}
	g.nodes[m][n] = struct{}{}
}

// findTriangles finds all triangles in the graph where a triangle
// is a set of three nodes connected by edges to each other.
func (g *graph) findTriangles() map[string][]string {
	triples := make(map[string][]string)
	for s := range g.nodes {
		for t := range g.nodes[s] {
			for u := range g.nodes[t] {
				if _, ok := g.nodes[u][s]; ok {
					triple := []string{s, t, u}
					sort.Strings(triple)
					key := strings.Join(triple, ",")
					triples[key] = triple
				}
			}
		}
	}
	return triples
}

// extendGroups takes a map of groups (such as those returned by findTriangles or by
// this function) and tries to extend each one by walking the set of nodes in the graph
// and adding any nodes that are connected to all nodes in the group.
func (g *graph) extendGroups(groups map[string][]string) map[string][]string {
	extended := make(map[string][]string)
	for _, group := range groups {
		// find all nodes connected to all nodes in the group
		connected := make(map[string]struct{})
		for _, node := range group {
			for n := range g.nodes[node] {
				connected[n] = struct{}{}
			}
		}
		// remove nodes already in the group
		for _, node := range group {
			delete(connected, node)
		}
		// check if the connected nodes are connected to all nodes in the group
		for n := range connected {
			connectedToAll := true
			for _, node := range group {
				if _, ok := g.nodes[n][node]; !ok {
					connectedToAll = false
					break
				}
			}
			if connectedToAll {
				newGroup := make([]string, len(group)+1)
				copy(newGroup, group)
				newGroup[len(group)] = n
				sort.Strings(newGroup)
				key := strings.Join(newGroup, ",")
				extended[key] = newGroup
			}
		}
	}
	return extended
}

func part1(lines []string) int {
	g := newGraph()
	for _, line := range lines {
		parts := strings.Split(line, "-")
		node := parts[0]
		dest := parts[1]
		g.addPair(node, dest)
	}
	tris := g.findTriangles()
	// fmt.Println(tris)
	countTs := 0
	for k := range tris {
		// fragile but it works
		if k[0] == 't' || k[3] == 't' || k[6] == 't' {
			countTs++
		}
	}
	return countTs
}

func part2(lines []string) int {
	g := newGraph()
	for _, line := range lines {
		parts := strings.Split(line, "-")
		node := parts[0]
		dest := parts[1]
		g.addPair(node, dest)
	}
	groups := g.findTriangles()
	for {
		nextGroups := g.extendGroups(groups)
		if len(nextGroups) == 0 {
			break
		}
		fmt.Println(len(nextGroups))
		groups = nextGroups
	}
	fmt.Println(groups)
	return len(groups)
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
