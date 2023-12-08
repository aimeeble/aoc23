package main

import (
	"aoc23/lib"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

const (
	startNodeName = "AAA"
	endNodeName   = "ZZZ"
)

type turnDir int

const (
	turnLeft  turnDir = 0
	turnRight turnDir = 1
)

func (td turnDir) String() string {
	if td == turnLeft {
		return "L"
	}
	return "R"
}

type node struct {
	name  string
	child []string
}

func main() {
	lines, err := lib.GetInputAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	// Load turns
	turns := make([]turnDir, 0, len(lines[0]))
	for i, x := range strings.TrimSpace(lines[0]) {
		if x == 'L' {
			turns = append(turns, turnLeft)
		} else if x == 'R' {
			turns = append(turns, turnRight)
		} else {
			log.Fatalf("Unexpected direction %c at pos %d in %q", x, i, lines[0])
		}
	}

	// Load nodes
	nodes := make(map[string]*node)
	for linen, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.FieldsFunc(line, func(r rune) bool {
			return unicode.IsSpace(r) || r == '(' || r == ')' || r == ','
		})
		if parts[1] != "=" {
			log.Fatalf("Unexepcted format on line %d %q", linen, line)
		}
		n := &node{
			name:  parts[0],
			child: []string{parts[2], parts[3]},
		}
		nodes[parts[0]] = n
	}

	for k, n := range nodes {
		log.Printf("%s = %v", k, n)
	}
	log.Printf("Turns: %v", turns)

	// Solve it
	steps := 0
	for curNodeName := startNodeName; curNodeName != endNodeName; steps++ {
		curNode := nodes[curNodeName]
		nextTurn := turns[steps%len(turns)]
		nextNodeName := curNode.child[nextTurn]
		log.Printf("Stepping %s from %s -> %s", nextTurn, curNodeName, nextNodeName)
		curNodeName = nextNodeName
	}
	fmt.Printf("Went AAA -> ZZZ in %d steps\n", steps)
}
