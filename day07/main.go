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

type handType int

// Types of hands listed in increasing order of strength.
const (
	Invalid handType = iota
	HighCard
	OnePair
	TwoPair
	ThreeOfKind
	FullHouse
	FourOfKind
	FiveOfKind
)

func (ht handType) String() string {
	switch ht {
	case Invalid:
		return "invalid"
	case HighCard:
		return "high-card"
	case OnePair:
		return "one-pair"
	case TwoPair:
		return "two-pair"
	case ThreeOfKind:
		return "three-of-a-kind"
	case FullHouse:
		return "full-house"
	case FourOfKind:
		return "four-of-a-kind"
	case FiveOfKind:
		return "five-of-a-kind"
	}
	return "cheater"
}

type card byte

func (c card) String() string {
	return string(c)
}

func (c card) value() int {
	if c >= '1' && c <= '9' {
		return int(byte(c) - '0')
	}
	switch c {
	case 'T':
		return 10
	case 'J':
		return 11
	case 'Q':
		return 12
	case 'K':
		return 13
	case 'A':
		return 14
	}
	return 0
}

type hand struct {
	cards []card
	bid   int

	_rank      handType
	_useJokers bool
}

func (h hand) String() string {
	return fmt.Sprintf("{%s %4d %s}", h.cards, h.bid, h.Rank())
}

func (h *hand) Rank() handType {
	if h._rank > Invalid {
		return h._rank
	}

	seen := make(map[card]int)
	for i := 0; i < 5; i++ {
		seen[h.cards[i]]++
	}

	cur := HighCard
	for c, n := range seen {
		if h._useJokers && c == 'J' {
			continue
		}
		switch {
		case n == 2 && cur < OnePair:
			cur = OnePair
		case n == 2 && cur == OnePair:
			cur = TwoPair
		case n == 2 && cur == ThreeOfKind,
			n == 3 && cur == OnePair:
			cur = FullHouse
		case n == 3 && cur < ThreeOfKind:
			cur = ThreeOfKind
		case n == 4 && cur < FourOfKind:
			cur = FourOfKind
		case n == 5 && cur < FiveOfKind:
			cur = FiveOfKind
		}
	}

	// Account for jokers
	if h._useJokers {
		jokers := seen['J']
		switch cur {
		case FiveOfKind:
			// no extra cards could be a joker.
		case FourOfKind:
			if jokers == 1 {
				cur = FiveOfKind
			}
		case ThreeOfKind:
			if jokers == 1 {
				cur = FourOfKind
			} else if jokers == 2 {
				cur = FiveOfKind
			}
		case FullHouse:
			// no extra cards could be a joker.
		case TwoPair:
			if jokers == 1 {
				cur = FullHouse
			}
		case OnePair:
			if jokers == 1 {
				cur = ThreeOfKind
			} else if jokers == 2 {
				cur = FourOfKind
			} else if jokers == 3 {
				cur = FiveOfKind
			}
		case HighCard:
			if jokers == 1 {
				cur = OnePair
			} else if jokers == 2 {
				cur = ThreeOfKind
			} else if jokers == 3 {
				cur = FourOfKind
			} else if jokers == 4 {
				cur = FiveOfKind
			} else if jokers == 5 {
				cur = FiveOfKind
			}
		default:
			log.Fatalf("invalid state")
		}
	}

	if cur == Invalid {
		log.Fatalf("about to set to invalid")
	}
	h._rank = cur
	return cur
}

func (h hand) Beats(oh *hand) bool {
	if h.Rank() == oh.Rank() {
		log.Printf("Checking %s vs %s, card-by-card", h, oh)
		for i := 0; i < 5; i++ {
			v1 := h.cards[i].value()
			v2 := oh.cards[i].value()
			if h._useJokers {
				// Special case: jokers are shit.
				if h.cards[i] == 'J' {
					v1 = 1
				}
				if oh.cards[i] == 'J' {
					v2 = 1
				}
			}
			log.Printf("\t%s(%d) vs %s(%d)", h.cards[i], v1, oh.cards[i], v2)

			if v1 > v2 {
				log.Printf("\t%s wins", h)
				return true
			} else if v1 < v2 {
				log.Printf("\t%s wins", oh)
				return false
			} else {
				// cards[i] match, continue checking
			}
		}
		log.Fatalf("tie %s vs %s", h, oh)
	}
	return h.Rank() > oh.Rank()
}

func main() {
	lines, err := lib.GetInputAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	var hands []*hand
	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Fields(line)

		h := &hand{
			cards: make([]card, 0, 5),
		}
		for i := 0; i < len(parts[0]); i++ {
			h.cards = append(h.cards, card(parts[0][i]))
		}
		h.bid = int(lib.Must(strconv.ParseInt(parts[1], 10, 64)))
		hands = append(hands, h)
	}

	log.Printf("As dealt:")
	for i, h := range hands {
		log.Printf("Hand %d: %+v", i, h)
	}
	sort.SliceStable(hands, func(i, j int) bool { return !hands[i].Beats(hands[j]) })

	// Part 1: rank hands and tally score.
	log.Printf("Part 1 ranked:")
	score := 0
	for i, h := range hands {
		pts := h.bid * (i + 1)
		score += pts
		log.Printf("Hand %d: %+v, => %d", i, h, pts)
	}
	fmt.Printf("Part 1 score: %d\n", score)

	// Part 2: Activate Jokers.
	for _, h := range hands {
		h._useJokers = true
		h._rank = Invalid // reset rank to recalculate.
	}
	sort.SliceStable(hands, func(i, j int) bool { return !hands[i].Beats(hands[j]) })
	log.Printf("Part 2 ranked:")
	score = 0
	for i, h := range hands {
		pts := h.bid * (i + 1)
		score += pts
		log.Printf("Hand %d: %+v, => %d", i, h, pts)
	}
	fmt.Printf("Part 2 score: %d\n", score)
}
