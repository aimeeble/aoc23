package main

import (
	"aoc23/lib"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type seq []int

func (s seq) isZero() bool {
	for _, e := range s {
		if e != 0 {
			return false
		}
	}
	return true
}

func (s seq) reduce() seq {
	var r seq
	for i := 1; i < len(s); i++ {
		r = append(r, s[i]-s[i-1])
	}
	return r
}

func (s seq) render(indent int) string {
	var sb strings.Builder
	for i := 0; i < indent; i++ {
		sb.WriteString("   ")
	}
	for i, e := range s {
		if i > 0 {
			sb.WriteString("   ")
		}
		fmt.Fprintf(&sb, "%3d", e)
	}
	return sb.String()
}

func (s seq) extrapolate(indent int) int {
	log.Printf("%s", s.render(indent))

	if s.isZero() {
		return 0
	}
	r := s.reduce()
	re := r.extrapolate(indent + 1)
	return s[len(s)-1] + re
}

func (s seq) pretrapolate(indent int) int {
	log.Printf("%s", s.render(indent))

	if s.isZero() {
		return 0
	}
	r := s.reduce()
	re := r.pretrapolate(indent + 1)
	return s[0] - re
}

func main() {
	flag.Parse()

	lines, err := lib.GetInputAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	var sequences []seq
	for _, line := range lines {
		parts := strings.Fields(strings.TrimSpace(line))
		s := lib.StrToInt(parts)
		sequences = append(sequences, s)
	}

	sum := 0
	for i, s := range sequences {
		e := s.extrapolate(0)
		log.Printf("Seq[%2d] %v => %d", i, s, e)
		sum += e
	}
	fmt.Printf("Part 1 sum: %d\n", sum)

	sum = 0
	for i, s := range sequences {
		e := s.pretrapolate(0)
		log.Printf("Seq[%2d] %v => %d", i, s, e)
		sum += e
	}
	fmt.Printf("Part 2 sum: %d\n", sum)
}
