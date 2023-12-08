package main

import (
	"aoc23/lib"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
)

var (
	doPart1      = flag.Bool("part1", false, "calculate part 1 solution")
	doPart2Naive = flag.Bool("part2", false, "calculate part 2 solution")
	doPart2Fast  = flag.Bool("part2fast", false, "calculate part 2 solution")
	doViz        = flag.Bool("viz", false, "graphvizify the (part 2) paths")
	verbose      = flag.Bool("verbose", false, "print steps")
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

func part1(
	startNodeName string,
	stopFn func(string) bool,
	turns []turnDir,
	nodes map[string]*node,
) int {
	steps := 0
	curNodeName := startNodeName
	for {
		nextTurn := turns[steps%len(turns)]
		nextNodeName := nodes[curNodeName].child[nextTurn]
		if *verbose {
			log.Printf("Stepping %s from %s -> %s", nextTurn, curNodeName, nextNodeName)
		}
		if stopFn(curNodeName) {
			break
		}
		curNodeName = nextNodeName
		steps++
	}
	return steps
}

func main() {
	flag.Parse()

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

	// Solve it: part 1.
	if *doPart1 {
		steps := part1(
			startNodeName,
			func(nn string) bool { return nn == endNodeName },
			turns,
			nodes,
		)
		fmt.Printf("Went AAA -> ZZZ in %d steps\n", steps)
	}

	// Solve it: part 2.
	if *doPart2Naive {
		startTime := time.Now()
		steps := 0
		var curNodeNames []string
		for k := range nodes {
			if strings.HasSuffix(k, "A") {
				curNodeNames = append(curNodeNames, k)
			}
		}
		fmt.Printf("Progress: %10d steps. Current nodes: %v\n", steps, curNodeNames)

		// Walk all nodes
		for steps = 1; ; steps++ {
			nextTurn := turns[(steps-1)%len(turns)]
			if *verbose {
				log.Printf("Step %d", steps)
			}
			for i, n := range curNodeNames {
				nextNodeName := nodes[n].child[nextTurn]
				if *verbose {
					log.Printf("\t%s turn %s => %s", nextTurn, n, nextNodeName)
				}
				curNodeNames[i] = nextNodeName
			}

			if steps%1_000_000 == 0 {
				fmt.Printf("Progress: %10d steps. Current nodes: %v; %5s elapsed\n", steps, curNodeNames, time.Since(startTime).Round(time.Second))
			}

			// check end condition
			if func() bool {
				for _, n := range curNodeNames {
					if !strings.HasSuffix(n, "Z") {
						return false
					}
				}
				// All node names end in Z
				return true
			}() {
				break
			}
		}
		fmt.Printf("Final nodes after %d steps: %v. Elapsed time %s\n", steps, curNodeNames, time.Since(startTime))
	}

	if *doPart2Fast {
		startTime := time.Now()
		var startNodeNames []string
		for k := range nodes {
			if strings.HasSuffix(k, "A") {
				startNodeNames = append(startNodeNames, k)
			}
		}
		endNodeNames := make([]string, len(startNodeNames))
		log.Printf("Starting nodes: %v", startNodeNames)

		steps := make([]int, len(startNodeNames))
		for i := range startNodeNames {
			steps[i] = part1(
				startNodeNames[i],
				func(nn string) bool {
					if strings.HasSuffix(nn, "Z") {
						endNodeNames[i] = nn
						return true
					}
					return false
				},
				turns,
				nodes,
			)
		}

		fmt.Printf("Individual step counts:\n")
		for i := range startNodeNames {
			fmt.Printf("\t%s => %s: %d\n", startNodeNames[i], endNodeNames[i], steps[i])
		}

		stepsAll := lib.LCM(steps...)
		fmt.Printf("Combined steps %d (elapsed %s)\n", stepsAll, time.Since(startTime))
	}

	if *doViz {
		startTime := time.Now()
		var startNodeNames []string
		for k := range nodes {
			if strings.HasSuffix(k, "A") {
				startNodeNames = append(startNodeNames, k)
			}
		}

		fmt.Printf("digraph AOC {\n")
		fmt.Printf("\tlayout=\"neato\";\n")
		fmt.Printf("\n")

		for n := range nodes {
			startstop := ""
			if strings.HasSuffix(n, "A") {
				startstop = ", color=\"green\""
			} else if strings.HasSuffix(n, "Z") {
				startstop = ", color=\"red\""
			} else {
				continue
			}
			fmt.Printf("\t\"%s\" [label=\"%s\"%s];\n", n, n, startstop)
		}
		fmt.Printf("\n")

		for _, startNode := range startNodeNames {
			todo := []string{startNode}

			seen := make(map[string]bool)
			for len(todo) > 0 {
				cur := todo[0]
				todo = todo[1:]

				if seen[cur] {
					continue
				}
				seen[cur] = true

				fmt.Printf("\t\"%s\" -> \"%s\";\n", cur, nodes[cur].child[0])
				fmt.Printf("\t\"%s\" -> \"%s\";\n", cur, nodes[cur].child[1])
				todo = append(todo, nodes[cur].child[0])
				todo = append(todo, nodes[cur].child[1])
			}
		}

		fmt.Printf("}\n")
		log.Printf("Generated graph in %s", time.Since(startTime))
	}
}
