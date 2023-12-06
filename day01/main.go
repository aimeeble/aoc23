package main

import (
	"aoc23/lib"
	"log"
	"os"
)

func main() {
	lines, err := lib.GetInputAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	sum := 0
	for ln, line := range lines {
		var nums []int
		log.Printf("---- Line (%2d): %q", ln+1, line)
		for i, r := range line {
			if lib.IsNum(r) {
				n := lib.RuneToDigit(r)
				nums = append(nums, n)
				log.Printf("Adding %d (int)", n)
			} else if n, ok := lib.WordToNum(line[0 : i+1]); ok {
				nums = append(nums, n)
				log.Printf("Adding %d (str)", n)
			}
		}

		log.Printf("Got %d digits: %v", len(nums), nums)
		num := nums[0]*10 + nums[len(nums)-1]
		log.Printf("Got number %d", num)
		sum += int(num)
	}

	log.Printf("Sum: %d", sum)
}
