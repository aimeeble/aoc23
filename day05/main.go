package main

import (
	"aoc23/lib"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type modeState int

const (
	lookingMode modeState = iota
	mapTypeMode
	mapContentMode
)

func (ms modeState) String() string {
	switch ms {
	case lookingMode:
		return "lookingMode"
	case mapTypeMode:
		return "mapTypeMode"
	case mapContentMode:
		return "mapContentMode"
	default:
		return "invalidMode"
	}
}

type thingRange struct {
	dstLo int
	srcLo int
	count int
}

type thingMap struct {
	ranges []thingRange
}

func (tm *thingMap) sortRanges() {
	sort.Slice(tm.ranges, func(i, j int) bool {
		return tm.ranges[i].srcLo < tm.ranges[j].srcLo
	})
}

// Get maps srcThing[id] to dstThing[id].
func (tm *thingMap) Get(id int) int {
	//tm.sortRanges()
	for _, r := range tm.ranges {
		if id >= r.srcLo && id <= r.srcLo+r.count {
			return r.dstLo + (id - r.srcLo)
		}
	}
	return id
}

type onEachFun func(origId int, dstType string, dstId int)

func printMaps(
	startType string,
	things map[string]map[string]*thingMap,
	ids []int,
	onEach onEachFun,
) {
	var sb strings.Builder
	next := startType
	fmt.Fprintf(&sb, "%15s", "seed")
	for len(things[next]) > 0 {
		for x := range things[next] {
			next = x
			fmt.Fprintf(&sb, "%15s", x)
			break
		}
	}
	sb.WriteString("\n")

	var printMapsInner func(origId int, srcType string, id int, sb *strings.Builder, onEach onEachFun)
	printMapsInner = func(origId int, srcType string, id int, sb *strings.Builder, onEach onEachFun) {
		fmt.Fprintf(sb, "%15d", id)
		srcMap := things[srcType]
		for dstType, curMap := range srcMap {
			//log.Printf("%s-to-%s %d = %d", srcType, dstType, id, curMap.Get(id))
			nextId := curMap.Get(id)
			printMapsInner(origId, dstType, nextId, sb, onEach)
			onEach(origId, dstType, nextId)
			break
		}
	}

	for _, id := range ids {
		printMapsInner(id, startType, id, &sb, onEach)
		sb.WriteString("\n")
	}
	fmt.Printf("%s\n", sb.String())
}

func main() {
	lines, err := lib.GetInputAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	things := make(map[string]map[string]*thingMap)

	// line 1 is seeds
	parts := strings.Fields(lines[0])
	if parts[0] != "seeds:" {
		log.Fatalf("Unexpected line 1: %q", lines[0])
	}
	var seeds []int
	for _, s := range parts[1:] {
		seeds = append(seeds, int(lib.Must(strconv.ParseInt(s, 10, 64))))
	}

	var (
		mode   modeState
		curMap *thingMap
	)
	for linen, line := range lines[1:] {
		line = strings.TrimSpace(line)

		log.Printf("Line %2d (%14s): %q", linen, mode, line)
		if line == "" {
			mode = mapTypeMode
			continue // skip blank lines
		}

		switch mode {
		case mapTypeMode:
			var ok bool

			mode = mapContentMode
			parts = strings.Fields(line)
			parts = strings.Split(parts[0], "-")
			if len(parts) != 3 {
				log.Fatalf("Unexpected map name %q", line)
			}
			srcThing := parts[0]
			dstThing := parts[2]
			log.Printf("src = %s", srcThing)
			log.Printf("dst = %s", dstThing)

			srcMap, ok := things[srcThing]
			if !ok {
				srcMap = make(map[string]*thingMap)
				things[srcThing] = srcMap
			}
			curMap, ok = srcMap[dstThing]
			if !ok {
				curMap = new(thingMap)
				srcMap[dstThing] = curMap
			}

		case mapContentMode:
			parts := strings.Fields(line)
			dstLo := int(lib.Must(strconv.ParseInt(parts[0], 10, 64)))
			srcLo := int(lib.Must(strconv.ParseInt(parts[1], 10, 64)))
			count := int(lib.Must(strconv.ParseInt(parts[2], 10, 64)))

			curMap.ranges = append(curMap.ranges, thingRange{
				dstLo: dstLo,
				srcLo: srcLo,
				count: count,
			})
		}
	}

	for srcType, srcMap := range things {
		for dstType, thingMap := range srcMap {
			//thingMap.sortRanges()
			log.Printf("%s-to-%s", srcType, dstType)
			for _, r := range thingMap.ranges {
				log.Printf("\t%v", r)
			}
		}
	}

	log.Printf("Seeds: %v", seeds)
	min := lib.MaxUint
	printMaps("seed", things, seeds, func(id int, dstType string, dstId int) {
		if dstType == "location" {
			min = lib.Min(uint(dstId), min)
		}
	})
	fmt.Printf("Minimum location: %d\n", min)
}
