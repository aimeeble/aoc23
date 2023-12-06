package main

import (
	"aoc23/lib"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	lines, err := lib.GetInputAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	sum := 0
	for _, line := range lines {
		var numRunes []rune
		for _, r := range line {
			if lib.IsNum(r) {
				numRunes = append(numRunes, r)
			}
		}

		log.Printf("Got %d digits", len(numRunes))
		numStr := fmt.Sprintf("%d%d", lib.RuneToDigit(numRunes[0]), lib.RuneToDigit(numRunes[len(numRunes)-1]))
		num, err := strconv.ParseInt(numStr, 10, 32)
		if err != nil {
			log.Fatalf("Error converting numbers: %v", err)
		}
		log.Printf("Got number %d", num)
		sum += int(num)
	}

	log.Printf("Sum: %d", sum)
}
