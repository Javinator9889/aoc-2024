package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Javinator9889/aoc-2024/util"
)

//go:embed input.txt
var input string

func init() {
	// do this in init (not main) so test file has same input
	input = strings.TrimRight(input, "\n")
	if len(input) == 0 {
		panic("empty input.txt file")
	}
}

func main() {
	var part int
	var debug bool
	flag.IntVar(&part, "part", 1, "part 1 or 2")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()
	fmt.Println("Running part", part)
	if debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if part == 1 {
		ans := part1(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	} else {
		ans := part2(input)
		util.CopyToClipboard(fmt.Sprintf("%v", ans))
		fmt.Println("Output:", ans)
	}
}

var UP = Dir{0, -1}
var DOWN = Dir{0, 1}
var LEFT = Dir{-1, 0}
var RIGHT = Dir{1, 0}

var dirTransform = map[rune]Dir{
	'^': UP,
	'v': DOWN,
	'<': LEFT,
	'>': RIGHT,
}

type Dir struct {
	x, y int
}

type Position struct {
	x, y     int
	obstacle bool
	visited  bool
}

type Guard struct {
	x, y int
	dir  Dir
}

type Map [][]Position

func (m Map) outOfBounds(pos Position) bool {
	return pos.x < 0 || pos.y < 0 || pos.x >= len(m[0]) || pos.y >= len(m)
}

func (g *Guard) moveForward() {
	g.x += g.dir.x
	g.y += g.dir.y
}

func part1(input string) (uniqueVisited int) {
	mapp, guard := parseInput(input)
	for {
		slog.Debug("Guard", "guard", guard)
		// First, mark out current position as visited
		if !mapp[guard.y][guard.x].visited {
			slog.Debug("Visited", "guard", guard, "position", mapp[guard.y][guard.x])
			mapp[guard.y][guard.x].visited = true
			uniqueVisited++
		}
		// The guard moves depending on the direction. There are two posibilites:
		// 1. There is an obstacle in front of the guard.
		// 2. There is no obstacle in front of the guard.
		// If there is an obstacle, the guard will turn right. If there is no obstacle,
		// the guard will move forward.
		next := Position{x: guard.x + guard.dir.x, y: guard.y + guard.dir.y}
		if mapp.outOfBounds(next) {
			slog.Debug("Out of bounds", "next", next)
			break // The guard has left the map
		}
		slog.Debug("Next", "next", mapp[next.y][next.x])
		if mapp[next.y][next.x].obstacle {
			// Turn right
			guard.dir = Dir{-guard.dir.y, guard.dir.x}
			slog.Debug("Turn right", "guard", guard)
		} else {
			// Move forward
			guard.moveForward()
			slog.Debug("Move forward", "guard", guard)
		}
	}

	return
}

func part2(input string) int {
	return 0
}

func parseInput(input string) (mapp Map, guard Guard) {
	mapp = make(Map, strings.Count(input, "\n")+1)
	for i, line := range strings.Split(input, "\n") {
		slog.Debug("Line", "i", i, "line", line)
		mapp[i] = make([]Position, len(line))
		for j, char := range line {
			switch char {
			case '#':
				mapp[i][j] = Position{x: j, y: i, obstacle: true}
			case '^', 'v', '<', '>':
				guard = Guard{x: j, y: i, dir: dirTransform[char]}
				fallthrough
			case '.':
				mapp[i][j] = Position{x: j, y: i, obstacle: false}
			default:
				msg := fmt.Sprintf(`invalid character in input: '%c' [%d, %d]`, char, i, j)
				panic(msg)
			}
		}
	}
	return
}
