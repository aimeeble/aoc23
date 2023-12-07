package main

import (
	"aoc23/lib"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type raceInfo struct {
	raceDuration   int
	recordDistance int
}

func DistForHold(holdTime, totalDuration int) int {
	speed := holdTime
	movingDuration := totalDuration - holdTime
	return movingDuration * speed
}

func (ri raceInfo) MinWinHold() int {
	for i := 0; i <= ri.raceDuration; i++ {
		myDistance := DistForHold(i, ri.raceDuration)
		if myDistance > ri.recordDistance {
			return i
		}
	}
	return -1
}

func (ri raceInfo) MaxWinHold() int {
	for i := ri.raceDuration; i >= 0; i-- {
		myDistance := DistForHold(i, ri.raceDuration)
		if myDistance > ri.recordDistance {
			return i
		}
	}
	return -1
}

func (ri raceInfo) NumWaysToWin() int {
	return ri.MaxWinHold() - ri.MinWinHold() + 1
}

func main() {
	flag.Parse()

	lines, err := lib.GetInputAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	log.Printf("line 1: %q", lines[0])
	log.Printf("line 2: %q", lines[1])

	// Part 1: races are separate
	times := lib.StrToInt(strings.Fields(lines[0])[1:])
	distances := lib.StrToInt(strings.Fields(lines[1])[1:])

	races := make([]raceInfo, len(times))
	for i := 0; i < len(times); i++ {
		races[i] = raceInfo{
			raceDuration:   times[i],
			recordDistance: distances[i],
		}
	}

	var waysToWin []int
	for i, r := range races {
		n := r.NumWaysToWin()
		log.Printf("Ways to win race %d: %d", i, n)
		if n > 0 {
			waysToWin = append(waysToWin, n)
		}
	}
	if len(waysToWin) == 0 {
		fmt.Printf("Cannot win\n")
	} else {
		margin := 1
		for _, w := range waysToWin {
			margin *= w
		}
		fmt.Printf("Margin to win: %d\n", margin)
	}

	// Part 2: bad kerning single race.
	time := int(lib.Must(strconv.ParseInt(strings.Join(strings.Fields(lines[0])[1:], ""), 10, 64)))
	distance := int(lib.Must(strconv.ParseInt(strings.Join(strings.Fields(lines[1])[1:], ""), 10, 64)))
	log.Printf("Time: %d", time)
	log.Printf("Dist: %d", distance)
	race := raceInfo{
		raceDuration:   time,
		recordDistance: distance,
	}
	n := race.NumWaysToWin()
	log.Printf("Ways to win big race: %d", n)
}
