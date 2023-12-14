package main

import (
	"aoc23/lib"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type pipeMap struct {
	grid     [][]string
	dist     [][]int
	startPos []int
}

func (pm *pipeMap) renderGrid() string {
	var sb strings.Builder
	for y, gridrow := range pm.grid {
		for x, p := range gridrow {
			if x == pm.startPos[0] && y == pm.startPos[1] {
				fmt.Fprintf(&sb, "\x1b[32m%2s\x1b[0m", p)
			} else {
				fmt.Fprintf(&sb, "%2s", p)
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (pm *pipeMap) renderDist() string {
	var sb strings.Builder
	for y, distrow := range pm.dist {
		for x, dist := range distrow {
			if pm.grid[y][x] == "." {
				sb.WriteString(" .")
			} else {
				if dist > 0 {
					fmt.Fprintf(&sb, "\x1b[32m%2s\x1b[0m", lib.Base10ToBase62(dist))
				} else {
					fmt.Fprintf(&sb, "%2s", lib.Base10ToBase62(dist))
				}
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (pm *pipeMap) isValid(sx, sy, dx, dy int) bool {
	if len(pm.grid) == 0 || len(pm.grid[0]) == 0 {
		return false
	}
	if sx+dx < 0 || sx+dx >= len(pm.grid[0]) {
		return false
	}
	if sy+dy < 0 || sy+dy >= len(pm.grid) {
		return false
	}
	return true
}

func (pm *pipeMap) dfsDistance(sx, sy int) {
	var innerDFS func(x, y, curDist int, seen [][]bool)
	innerDFS = func(x, y, curDist int, seen [][]bool) {
		if seen[y][x] {
			if curDist >= pm.dist[y][x] {
				return
			}
			// We've seen it, but on a longer path. Re-process this
			// location with the shorter path.
		}
		seen[y][x] = true
		pm.dist[y][x] = curDist

		curPipe := pm.grid[y][x]
		log.Printf("Processing [%d %d] %s", x, y, curPipe)
		switch curPipe {
		case "S":
			// valid: unknown; must inspect adjacent elements
			if pm.isValid(x, y, 0, -1) {
				switch pm.grid[y-1][x] {
				case "|", "F", "7":
					// Flow down into S, so we can recurse up.
					innerDFS(x, y-1, curDist+1, seen)
				}
			}
			if pm.isValid(x, y, 0, +1) {
				switch pm.grid[y+1][x] {
				case "|", "L", "J":
					// Flow up into S, so we can recurse down.
					innerDFS(x, y+1, curDist+1, seen)
				}
			}
			if pm.isValid(x, y, -1, 0) {
				switch pm.grid[y][x-1] {
				case "-", "F", "L":
					// Flow right into S, so we can recurse left.
					innerDFS(x-1, y, curDist+1, seen)
				}
			}
			if pm.isValid(x, y, +1, 0) {
				switch pm.grid[y][x+1] {
				case "-", "J", "7":
					// Flow left into S, so we can recurse right.
					innerDFS(x+1, y, curDist+1, seen)
				}
			}
		case "|":
			// valid: up
			if pm.isValid(x, y, 0, -1) {
				innerDFS(x, y-1, curDist+1, seen)
			}
			// valid: down
			if pm.isValid(x, y, 0, +1) {
				innerDFS(x, y+1, curDist+1, seen)
			}
		case "-":
			// valid: left
			if pm.isValid(x, y, -1, 0) {
				innerDFS(x-1, y, curDist+1, seen)
			}
			// valid: right
			if pm.isValid(x, y, +1, 0) {
				innerDFS(x+1, y, curDist+1, seen)
			}
		case "F":
			// valid: right
			if pm.isValid(x, y, +1, 0) {
				innerDFS(x+1, y, curDist+1, seen)
			}
			// valid: down
			if pm.isValid(x, y, 0, +1) {
				innerDFS(x, y+1, curDist+1, seen)
			}
		case "7":
			// valid: left
			if pm.isValid(x, y, -1, 0) {
				innerDFS(x-1, y, curDist+1, seen)
			}
			// valid: down
			if pm.isValid(x, y, 0, +1) {
				innerDFS(x, y+1, curDist+1, seen)
			}
		case "L":
			// valid: right
			if pm.isValid(x, y, +1, 0) {
				innerDFS(x+1, y, curDist+1, seen)
			}
			// valid: up
			if pm.isValid(x, y, 0, -1) {
				innerDFS(x, y-1, curDist+1, seen)
			}
		case "J":
			// valid: left
			if pm.isValid(x, y, -1, 0) {
				innerDFS(x-1, y, curDist+1, seen)
			}
			// valid: up
			if pm.isValid(x, y, 0, -1) {
				innerDFS(x, y-1, curDist+1, seen)
			}
		case ".":
			// do nothing; no pipe.
		default:
			// unknown; error.
			log.Fatalf("Unknown grid element at [%d %d] = %s", x, y, curPipe)
		}
	}

	seen := make([][]bool, len(pm.grid))
	for i := 0; i < len(pm.grid); i++ {
		seen[i] = make([]bool, len(pm.grid[0]))
	}
	innerDFS(sx, sy, 0, seen)
}

func main() {
	flag.Parse()

	lines, err := lib.GetInputAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	pipes := &pipeMap{
		grid: make([][]string, len(lines)),
		dist: make([][]int, len(lines)),
	}

	for y, line := range lines {
		line = strings.TrimSpace(line)
		for x := range line {
			p := line[x : x+1]
			if p == "S" {
				pipes.startPos = []int{x, y}
			}
			pipes.grid[y] = append(pipes.grid[y], p)
			pipes.dist[y] = append(pipes.dist[y], 0)
		}
	}
	if pipes.startPos == nil {
		log.Fatal("Invalid starting position / starting position not specified")
	}

	log.Printf("Grid:\n%s", pipes.renderGrid())

	x, y := pipes.startPos[0], pipes.startPos[1]
	log.Printf("Starting at %v -> %s", pipes.startPos, pipes.grid[y][x])

	pipes.dfsDistance(x, y)
	fmt.Printf("Dist:\n%s\n", pipes.renderDist())

	// Part 1: find longest
	max := 0
	var maxLoc []int
	for y, row := range pipes.dist {
		for x, d := range row {
			if d > max {
				max = d
				maxLoc = []int{x, y}
			}
		}
	}

	fmt.Printf("Max distance is at %v = %d\n", maxLoc, max)
}
