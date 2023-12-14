package main

import (
	"aoc23/lib"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	doSVG   = flag.Bool("svg", false, "render a SVG of the main path")
	doPart1 = flag.Bool("part1", false, "solve part 1")
	doPart2 = flag.Bool("part2", false, "solve part 2")
)

type pipeMap struct {
	grid     [][]string
	dist     [][]int
	startPos []int

	// list of 2-tuples [(x0,y0), (x1,y1), ...], with corners[0] == startPos.
	corners [][]int
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

func (pm *pipeMap) renderSVG() string {
	maxX := 0
	maxY := 0
	for _, pt := range pm.corners {
		if pt[0] > maxX {
			maxX = pt[0]
		}
		if pt[1] > maxY {
			maxY = pt[1]
		}
	}
	maxX *= 10
	maxY *= 10

	var sb strings.Builder
	sb.WriteString("<html><body>\n")
	fmt.Fprintf(&sb, "<svg height=\"%d\" width=\"%d\">\n", maxX+50, maxY+50)

	// Draw with a polygon.
	fmt.Fprintf(&sb, "  <polygon style=\"fill:rgb(0,128,64);stroke:rgb(0,0,0);stroke-width:2\" points=\"")
	for _, pt := range pm.corners {
		fmt.Fprintf(&sb, "%d,%d ", pt[0]*10, pt[1]*10)
	}
	fmt.Fprintf(&sb, "\"></polygon>\n")

	// circle starting position
	fmt.Fprintf(&sb, "  <circle cx=\"%d\" cy=\"%d\" r=\"%d\" style=\"stroke:rgb(0,0,0);stroke-width:2;fill:rgba(0,255,64,0.8)\"></circle>\n",
		pm.corners[0][0]*10, pm.corners[0][1]*10, 5,
	)

	sb.WriteString("</svg>\n")
	sb.WriteString("</body></html>\n")

	return sb.String()
}

func (pm *pipeMap) findCorners() {
	var innerDFS func(x, y int, seen [][]bool)
	innerDFS = func(x, y int, seen [][]bool) {
		if seen[y][x] {
			return
		}
		seen[y][x] = true

		curPipe := pm.grid[y][x]
		//log.Printf("Processing [%d %d] %s", x, y, curPipe)
		switch curPipe {
		case "F":
			pm.corners = append(pm.corners, []int{x, y})

			// valid: right
			if pm.isValid(x, y, +1, 0) {
				innerDFS(x+1, y, seen)
			}
			// valid: down
			if pm.isValid(x, y, 0, +1) {
				innerDFS(x, y+1, seen)
			}

		case "L":
			pm.corners = append(pm.corners, []int{x, y})

			// valid: right
			if pm.isValid(x, y, +1, 0) {
				innerDFS(x+1, y, seen)
			}
			// valid: up
			if pm.isValid(x, y, 0, -1) {
				innerDFS(x, y-1, seen)
			}

		case "7":
			pm.corners = append(pm.corners, []int{x, y})

			// valid: left
			if pm.isValid(x, y, -1, 0) {
				innerDFS(x-1, y, seen)
			}
			// valid: down
			if pm.isValid(x, y, 0, +1) {
				innerDFS(x, y+1, seen)
			}

		case "J":
			pm.corners = append(pm.corners, []int{x, y})

			// valid: left
			if pm.isValid(x, y, -1, 0) {
				innerDFS(x-1, y, seen)
			}
			// valid: up
			if pm.isValid(x, y, 0, -1) {
				innerDFS(x, y-1, seen)
			}

		case "-":
			// valid: left
			if pm.isValid(x, y, -1, 0) {
				innerDFS(x-1, y, seen)
			}
			// valid: right
			if pm.isValid(x, y, +1, 0) {
				innerDFS(x+1, y, seen)
			}

		case "|":
			// valid: up
			if pm.isValid(x, y, 0, -1) {
				innerDFS(x, y-1, seen)
			}
			// valid: down
			if pm.isValid(x, y, 0, +1) {
				innerDFS(x, y+1, seen)
			}

		default:
		}
	}

	seen := make([][]bool, len(pm.grid))
	for i := 0; i < len(pm.grid); i++ {
		seen[i] = make([]bool, len(pm.grid[0]))
	}
	innerDFS(pm.startPos[0], pm.startPos[1], seen)
}

func (pm *pipeMap) dfsDistance(sx, sy int) {
	log.Printf("Starting DFS at [%d %d] -> %s", sx, sy, pm.grid[sy][sx])

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
		//log.Printf("Processing [%d %d] %s", x, y, curPipe)
		switch curPipe {
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

func (pm *pipeMap) replaceStartPos() {
	x, y := pm.startPos[0], pm.startPos[1]
	okayUp := false
	okayDn := false
	okayRi := false
	okayLe := false

	// valid: unknown; must inspect adjacent elements
	if pm.isValid(x, y, 0, -1) {
		switch pm.grid[y-1][x] {
		case "|", "F", "7":
			// Flow down into S, so we can recurse up.
			okayUp = true
		}
	}
	if pm.isValid(x, y, 0, +1) {
		switch pm.grid[y+1][x] {
		case "|", "L", "J":
			// Flow up into S, so we can recurse down.
			okayDn = true
		}
	}
	if pm.isValid(x, y, -1, 0) {
		switch pm.grid[y][x-1] {
		case "-", "F", "L":
			// Flow right into S, so we can recurse left.
			okayLe = true
		}
	}
	if pm.isValid(x, y, +1, 0) {
		switch pm.grid[y][x+1] {
		case "-", "J", "7":
			// Flow left into S, so we can recurse right.
			okayRi = true
		}
	}

	newStart := ""
	switch {
	case okayLe && okayRi:
		newStart = "-"
	case okayUp && okayDn:
		newStart = "|"
	case okayUp && okayRi:
		newStart = "L"
	case okayUp && okayLe:
		newStart = "J"
	case okayDn && okayRi:
		newStart = "F"
	case okayDn && okayLe:
		newStart = "7"
	default:
		log.Fatalf("Cannot infer what S should be replaced with")
	}
	pm.grid[y][x] = newStart
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
	pipes.replaceStartPos()
	pipes.findCorners()

	log.Printf("Corners:")
	for i, corn := range pipes.corners {
		log.Printf("  [%2d] %v", i, corn)
	}
	log.Printf("Grid:\n%s", pipes.renderGrid())

	if *doSVG {
		fmt.Printf("%s\n", pipes.renderSVG())
	}

	if *doPart1 {
		x, y := pipes.startPos[0], pipes.startPos[1]
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
}
