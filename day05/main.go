package main

import (
	"aoc23/lib"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
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
	srcLo int
	dstLo int
	count int
}

func (tr thingRange) String() string {
	return fmt.Sprintf(
		"{src [%3d-%3d] => dst [%3d-%3d]",
		tr.srcLo, tr.srcLo+tr.count-1,
		tr.dstLo, tr.dstLo+tr.count-1,
	)
}

type thingMap struct {
	ranges []*thingRange
}

func (tm *thingMap) sortRanges() {
	sort.SliceStable(tm.ranges, func(i, j int) bool {
		return tm.ranges[i].srcLo < tm.ranges[j].srcLo
	})
}

// Get maps srcThing[id] to dstThing[id].
func (tm *thingMap) Get(id int, typeName string) int {
	// Find an optimal starting point (first range that id could be in).
	idx := sort.Search(len(tm.ranges), func(i int) bool {
		if id < tm.ranges[i].srcLo+tm.ranges[i].count {
			//log.Printf("Range[%d] %+v good for id %d", i, tm.ranges[i], id)
			return true
		}
		//log.Printf("Range[%d] %+v  bad for id %d", i, tm.ranges[i], id)
		return false
	})
	if idx == -1 {
		// "this shouldn't happen" 😅
		log.Fatalf("Bad index for %d ", id)
	}
	if idx == len(tm.ranges) {
		// There is no range that could hold id, so this is the
		// identity mapping.
		return id
	}

	// Search remaining ranges, starting at our first possibly valid index.
	for _, r := range tm.ranges[idx : idx+1] {
		if id >= r.srcLo && id < r.srcLo+r.count {
			//log.Printf("ID %s-%d matched range %+v", typeName, id, r)
			return r.dstLo + (id - r.srcLo)
		}
	}

	// Slow-fail: some ranges looked promising, but it didn't pan out.
	// Another identity mapping.
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
	//fmt.Fprintf(&sb, "%15s", "seed")
	for len(things[next]) > 0 {
		for x := range things[next] {
			next = x
			//fmt.Fprintf(&sb, "%15s", x)
			break
		}
	}
	//sb.WriteString("\n")

	var printMapsInner func(origId int, srcType string, id int, sb *strings.Builder, onEach onEachFun)
	printMapsInner = func(origId int, srcType string, id int, sb *strings.Builder, onEach onEachFun) {
		//fmt.Fprintf(sb, "%15d", id)
		srcMap := things[srcType]
		for dstType, curMap := range srcMap {
			nextId := curMap.Get(id, srcType)
			printMapsInner(origId, dstType, nextId, sb, onEach)
			onEach(origId, dstType, nextId)
			break
		}
	}

	for _, id := range ids {
		printMapsInner(id, startType, id, &sb, func(origId int, dstType string, dstId int) {
			onEach(origId, dstType, dstId)
		})
		//sb.WriteString("\n")
	}
	//fmt.Printf("%s\n", sb.String())
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
		_ = linen

		//log.Printf("Line %2d (%14s): %q", linen, mode, line)
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

			curMap.ranges = append(curMap.ranges, &thingRange{
				dstLo: dstLo,
				srcLo: srcLo,
				count: count,
			})
		}
	}

	for srcType, srcMap := range things {
		for dstType, thingMap := range srcMap {
			thingMap.sortRanges()
			log.Printf("%s-to-%s", srcType, dstType)
			for _, r := range thingMap.ranges {
				log.Printf("\t%+v", r)
			}
		}
	}

	log.Printf("Seeds: %v", seeds)
	min := lib.MaxUint
	startTime := time.Now()
	printMaps("seed", things, seeds, func(id int, dstType string, dstId int) {
		if dstType == "location" {
			min = lib.Min(uint(dstId), min)
		}
	})
	log.Printf("Took %s", time.Since(startTime))
	fmt.Printf("Minimum location (part 1): %d\n", min)

	// Part2: Expand seed list.
	if len(seeds)%2 != 0 {
		log.Fatalf("Unbalanced seed count: must be even, got %d", len(seeds))
	}
	min = lib.MaxUint
	log.Print("Expanding seeds (batches)")

	var wg sync.WaitGroup
	var mu sync.Mutex
	c := make(chan intPair)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batchRange := range c {
				batchSet := makeBatch(batchRange)
				printMaps("seed", things, batchSet, func(id int, dstType string, dstId int) {
					if dstType == "location" {
						mu.Lock()
						min = lib.Min(uint(dstId), min)
						mu.Unlock()
					}
				})
			}
		}()
	}

	startTime = time.Now()
	for i := 0; i < len(seeds); i += 2 {
		lo := seeds[i]
		run := seeds[i+1]
		log.Printf("\tExpanding [%10d-%10d] (len=%10d)", lo, lo+run-1, run)
		batchSize := 1_000_000
		batches := makeBatchSpec(lo, lo+run, batchSize)
		for _, batchRange := range batches {
			c <- batchRange
			//	batchSet := makeBatch(batchRange)
			//	printMaps("seed", things, batchSet, func(id int, dstType string, dstId int) {
			//		if dstType == "location" {
			//			min = lib.Min(uint(dstId), min)
			//		}
			//	})
			//	log.Printf("Batch %d/%d (%d elements) took %s", batchId, len(batches), len(batchSet), time.Since(startTime))
		}
	}
	close(c)
	wg.Wait()

	log.Printf("Took %s", time.Since(startTime))
	fmt.Printf("Minimum location (part 2): %d\n", min)
}

type intPair struct{ a, b int }

func makeBatchSpec(lo, hi, batchSize int) []intPair {
	var res []intPair
	for i := lo; i < hi; i += batchSize {
		res = append(res, intPair{i, lib.Min(i+batchSize, hi)})
	}
	return res
}

func makeBatch(ip intPair) []int {
	n := ip.b - ip.a
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = ip.a + i
	}
	return res
}
