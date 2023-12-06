package main

import (
	"aoc23/lib"
	"log"
	"os"
)

type part struct {
	name    string
	x, y    int
	numbers []*partNumber
}

type partNumber struct {
	x, y   int
	used   bool
	digits int
	value  int
}

func main() {
	lines, err := lib.GetInputAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	// Iteration 1: Extract the numbers.
	var numbers []*partNumber
	for ln, line := range lines {
		exNums := lib.ExtractNumbers(line)
		for _, en := range exNums {
			numbers = append(numbers, &partNumber{
				x:      en.Offset,
				y:      ln,
				used:   false,
				digits: en.Digits,
				value:  en.Value,
			})
		}
	}
	log.Printf("Loaded %d part numbers", len(numbers))

	// Iteration 2: look for part symbols.
	var parts []*part
	for ln, line := range lines {
		for i, x := range line {
			if lib.IsNum(x) || x == '.' || x == '\n' {
				continue
			}
			log.Printf("Part %c at (%d,%d)", x, i, ln)
			p := &part{
				name:    string(x),
				x:       i,
				y:       ln,
				numbers: nil,
			}
			parts = append(parts, p)
		}
	}
	log.Printf("Loaded %d parts", len(parts))

	// Find numbers near parts.
	for _, p := range parts {
		log.Printf("Part %s at (x,y)=(%d,%d)", p.name, p.x, p.y)
		for _, pnum := range numbers {
			if lib.Abs(pnum.y-p.y) > 1 {
				continue // not on same or adjact row
			}
			numX0 := pnum.x
			numX1 := pnum.x + pnum.digits
			if p.x < numX0-1 || p.x > numX1 {
				continue // too far side-to-side
			}
			log.Printf("\tUses part number %d at (x,y)=(%d,%d)", pnum.value, pnum.x, pnum.y)
			p.numbers = append(p.numbers, pnum)
			pnum.used = true
		}
	}

	// Collect orphaned numbers.
	var orphanedNumbers []*partNumber
	sum := 0
	for _, pnum := range numbers {
		if !pnum.used {
			orphanedNumbers = append(orphanedNumbers, pnum)
		} else {
			sum += pnum.value
		}
	}
	log.Printf("Orphaned %d part numbers", len(orphanedNumbers))

	log.Printf("")
	log.Printf("Sum of non-orphaned part numbers: %d", sum)
}
