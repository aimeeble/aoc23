package main

import (
	"aoc23/lib"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines, err := lib.GetInputAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	sum := 0
	// Initialise with one copy of each card.
	copies := make(map[int]int)
	for x := 0; x < len(lines); x++ {
		copies[x] = 1
	}

	for linen, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.SplitN(line, ": ", 2)
		parts = strings.SplitN(parts[1], " | ", 2)

		// Build list of picks and wins.
		winpartsStr := strings.Fields(parts[0])
		pickpartsStr := strings.Fields(parts[1])
		winparts := make([]int, 0, len(winpartsStr))
		pickparts := make([]int, 0, len(pickpartsStr))
		wins := make(map[int]bool)
		for _, x := range winpartsStr {
			n, _ := strconv.ParseInt(x, 10, 32)
			winparts = append(winparts, int(n))
			wins[int(n)] = true
		}
		for _, x := range pickpartsStr {
			n, _ := strconv.ParseInt(x, 10, 32)
			pickparts = append(pickparts, int(n))
		}

		log.Printf("Card %d (x%2d): %v | %v", linen, copies[linen]+1, winparts, pickparts)

		// Look for wins.
		matches := 0
		pts := 0
		for _, x := range pickparts {
			if wins[x] {
				matches++
				if pts == 0 {
					pts = 1
				} else {
					pts *= 2
				}
			}
		}
		log.Printf("\tMatches: %d x %d copies", matches, copies[linen])

		// Duplicate cards based on number of matches.
		for x := linen + 1; x <= linen+matches && x < len(copies); x++ {
			log.Printf("\tcard %d cloned %d times", x+1, copies[linen])
			copies[x] += copies[linen]
		}
		sum += pts
	}
	totalCopies := 0
	for i := 0; i < len(copies); i++ {
		totalCopies += copies[i]
	}

	log.Printf("")
	log.Printf("Total points: %d", sum)
	log.Printf("Card counts:")
	for i := 0; i < len(copies); i++ {
		log.Printf("\tCard %d: %d", i+1, copies[i])
	}
	log.Printf("Total cards: %d", totalCopies)
}
