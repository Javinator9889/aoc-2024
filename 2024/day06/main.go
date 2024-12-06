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
	dir      Dir
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

// Makes the guard go through the map and returns the number of unique positions visited
// before the guard leaves the map. If the guard enters a loop, it will return an error.
func (g *Guard) Go(mapp Map) (int, error) {
	uniqueVisited := 0
	for {
		// First, mark out current position as visited
		if !mapp[g.y][g.x].visited {
			mapp[g.y][g.x].visited = true
			mapp[g.y][g.x].dir = g.dir
			uniqueVisited++
		} else {
			// The guard has visited this position before. Check if the direction is the same
			if mapp[g.y][g.x].dir == g.dir {
				return uniqueVisited, fmt.Errorf("loop detected at position (%d, %d)", g.x, g.y)
			}
		}
		// The guard moves depending on the direction. There are two posibilites:
		// 1. There is an obstacle in front of the guard.
		// 2. There is no obstacle in front of the guard.
		// If there is an obstacle, the guard will turn right. If there is no obstacle,
		// the guard will move forward.
		next := Position{x: g.x + g.dir.x, y: g.y + g.dir.y}
		if mapp.outOfBounds(next) {
			return uniqueVisited, nil
		}
		if mapp[next.y][next.x].obstacle {
			// Turn right
			g.dir = Dir{-g.dir.y, g.dir.x}
		} else {
			// Move forward
			g.moveForward()
		}
	}
}

func (m Map) Clone() (clone Map) {
	clone = make(Map, len(m))
	for i := range m {
		clone[i] = make([]Position, len(m[i]))
		copy(clone[i], m[i])
	}
	return
}

func part1(input string) int {
	mapp, guard := parseInput(input)
	uniqueVisited, err := guard.Go(mapp)
	if err != nil {
		panic(err)
	}
	return uniqueVisited
}

func part2(input string) (loops int) {
	mapp, guard := parseInput(input)
	// We're going bruteforce here! We have to place an item in a location that causes the
	// guard to go in a loop. We will check if the guard has visited the same position twice
	// with the same direction, and if it has, we will break the loop.
	for i := range mapp {
		for j := range mapp[i] {
			// If the current position matches the guard or it's already an obstacle, skip
			if (guard.x == j && guard.y == i) || mapp[i][j].obstacle {
				continue
			}
			tmp := mapp.Clone()
			tmp[i][j].obstacle = true
			tmpGuard := guard
			_, err := tmpGuard.Go(tmp)
			if err != nil {
				loops++
			}
		}
	}
	return
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
