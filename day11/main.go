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
	scaleFactor = flag.Int("scale", 2, "scale factor (part1=>2, part2=>1_000_000)")
)

func manDist(p, q []int) int {
	if len(p) != len(q) {
		return 0
	}
	sum := 0
	for i := range p {
		sum += lib.Abs(p[i] - q[i])
	}
	return sum
}

func main() {
	flag.Parse()

	lines, err := lib.GetInputAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	var grid [][]int
	galaxyID := 1
	for _, line := range lines {
		line = strings.TrimSpace(line)
		gridLine := make([]int, 0, len(line))
		expand := true
		for _, x := range line {
			if x == '#' {
				gridLine = append(gridLine, galaxyID)
				galaxyID++
				expand = false
			} else {
				gridLine = append(gridLine, 0)
			}
		}
		if expand {
			for i := range gridLine {
				gridLine[i] = -1
			}
			grid = append(grid, gridLine)
		} else {
			grid = append(grid, gridLine)
		}
	}
	for x := 0; x < len(grid[0]); x++ {
		expand := true
		for _, gridLine := range grid {
			if gridLine[x] > 0 {
				expand = false
				break
			}
		}
		if expand {
			for i, gridLine := range grid {
				gridLine[x] = -1
				grid[i] = gridLine
			}
		}
	}

	galaxies := make(map[int][]int)
	y := 0
	for _, row := range grid {
		if row[0] == -1 {
			// Line is empty, expand by scale factor
			y += *scaleFactor
		} else {
			x := 0
			for _, el := range row {
				if el > 0 {
					// Add a galaxy, and move once.
					galaxies[el] = []int{x, y}
					x++
				} else if el < 0 {
					// Expand by moving scaleFactor.
					x += *scaleFactor
				} else {
					// Just move once.
					x++
				}
			}
			y++
		}
	}

	var sb strings.Builder
	for _, row := range grid {
		for _, el := range row {
			if el == 0 {
				fmt.Fprintf(&sb, "  .")
			} else {
				fmt.Fprintf(&sb, "%3d", el)
			}
		}
		sb.WriteString("\n")
	}
	log.Printf("Grid:\n%s\n", sb.String())

	log.Printf("Galaxies:")
	keys := lib.Keys(galaxies)
	for _, k := range keys {
		g := galaxies[k]
		log.Printf("  [%2d] %v", k, g)
	}

	sum := 0
	for i, k1 := range keys {
		for _, k2 := range keys[i+1:] {
			d := manDist(galaxies[k1], galaxies[k2])
			log.Printf("Distance %d -> %d = %d", k1, k2, d)
			sum += d
		}
	}
	fmt.Printf("Total distance: %d\n", sum)
}
