package main

import (
	"aoc23/lib"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type gameInfo struct {
	gameId int
	rounds []gameRound
}

func (gi gameInfo) NumNeeded() (r int, g int, b int) {
	for _, round := range gi.rounds {
		tmpr, tmpg, tmpb := round.NumNeeded()
		r = lib.Max(r, tmpr)
		g = lib.Max(g, tmpg)
		b = lib.Max(b, tmpb)
	}
	return
}

func (gi gameInfo) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Game %3d (%d rounds) = [", gi.gameId, len(gi.rounds))
	for i, r := range gi.rounds {
		something := false
		fmt.Fprintf(&sb, "(Round %d: ", i)
		if r.nRed > 0 {
			if something {
				sb.WriteString(" ")
			}
			fmt.Fprintf(&sb, "%d red", r.nRed)
			something = true
		}
		if r.nBlue > 0 {
			if something {
				sb.WriteString(" ")
			}
			fmt.Fprintf(&sb, "%d blue", r.nBlue)
			something = true
		}
		if r.nGreen > 0 {
			if something {
				sb.WriteString(" ")
			}
			fmt.Fprintf(&sb, "%d green", r.nGreen)
			something = true
		}
		sb.WriteString(")")
		if i < len(gi.rounds)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

type gameRound struct {
	nBlue  int
	nRed   int
	nGreen int
}

func (gr gameRound) NumNeeded() (r, g, b int) {
	return gr.nRed, gr.nGreen, gr.nBlue
}

func parseGame(name string, rounds []string) (*gameInfo, error) {
	if !strings.HasPrefix(name, "Game ") {
		return nil, fmt.Errorf("game does not have prefix 'Game '")
	}
	gamen, err := strconv.ParseInt(name[5:], 10, 32)
	if err != nil {
		return nil, err
	}

	info := &gameInfo{
		gameId: int(gamen),
		rounds: make([]gameRound, len(rounds)),
	}

	for roundn, round := range rounds {
		gr := &info.rounds[roundn]
		parts := strings.Split(round, ",")
		for _, p := range parts {
			nc := strings.SplitN(strings.TrimSpace(p), " ", 2)
			n, err := strconv.ParseInt(nc[0], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("failed to parse result %q: %v", p, err)
			}
			switch strings.ToLower(nc[1]) {
			case "green":
				gr.nGreen = int(n)
			case "red":
				gr.nRed = int(n)
			case "blue":
				gr.nBlue = int(n)
			}
		}
	}

	return info, nil
}

func main() {
	lines, err := lib.GetInputAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	var games []*gameInfo
	for ln, line := range lines {
		log.Printf("----- Line(%2d) %q", ln+1, line)
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			log.Fatalf("malformed line %d: %q", ln+1, line)
		}
		game := parts[0]
		parts = strings.Split(parts[1], ";")
		info, err := parseGame(game, parts)
		if err != nil {
			log.Fatalf("Game on line %d failed to parse: %v", ln+1, err)
		}
		log.Printf("%+v", info)
		games = append(games, info)
	}

	maxRed := 12
	maxGreen := 13
	maxBlue := 14
	minRed := 0
	minGreen := 0
	minBlue := 0
	sum := 0
	powerSum := 0
	var possibleGames []*gameInfo

	log.Printf("Finding possible games with max (r,g,b)=(%d,%d,%d)", maxRed, maxGreen, maxBlue)
	for _, g := range games {
		needr, needg, needb := g.NumNeeded()
		minRed = lib.Max(needr, minRed)
		minGreen = lib.Max(needg, minGreen)
		minBlue = lib.Max(needb, minBlue)
		power := needr * needg * needb
		powerSum += power

		if needr <= maxRed && needg <= maxGreen && needb <= maxBlue {
			possibleGames = append(possibleGames, g)
			sum += g.gameId
			log.Printf("\tGame %3d:     POSSIBLE. Needs (r,g,b)=(%2d,%2d,%2d). Power = %d", g.gameId, needr, needg, needb, power)
		} else {
			log.Printf("\tGame %3d: NOT POSSIBLE. Needs (r,g,b)=(%2d,%2d,%2d). Power = %d", g.gameId, needr, needg, needb, power)
		}
	}

	log.Printf("%d possible games", len(possibleGames))
	log.Printf("Needs for all to be possible: r=%d, g=%d, b=%d", minRed, minGreen, minBlue)
	log.Printf("")
	log.Printf("Sum of possible games: %d", sum)
	log.Printf("Sum of all powers: %d", powerSum)
}
