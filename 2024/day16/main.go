package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Javinator9889/aoc-2024/2024/day16/astar"
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

func part1(input string) (cost int) {
	reindlympics := parseInput(input)
	slog.Debug("parsed", "reindlympics", reindlympics)
	path := reindlympics.AStar()
	if path == nil {
		slog.Debug("no path found")
		return -1
	}
	slog.Debug("found path", "path", path)

	return path.Cost()
}

func part2(input string) int {
	reindlympics := parseInput(input)
	slog.Debug("parsed", "reindlympics", reindlympics)
	optimalPath := reindlympics.AStar()
	if optimalPath == nil {
		slog.Debug("no path found")
		return -1
	}
	path := reindlympics.AStarRecursive(len(optimalPath), 0, astar.EAST)
	if path == nil {
		slog.Debug("no paths found")
		return -1
	}
	slog.Debug("found path", "paths", path)
	uniq := path.Uniq()
	slog.Debug("found unique paths", "paths", uniq)

	return len(uniq)
}

func parseInput(input string) *astar.Reindlympics {
	res := &astar.Reindlympics{Grid: make(astar.Grid, 0)}
	for i, line := range strings.Split(input, "\n") {
		res.Grid = append(res.Grid, make([]rune, 0))
		for j, char := range line {
			res.Grid[i] = append(res.Grid[i], char)
			switch char {
			case 'S':
				res.Start = astar.Location{X: i, Y: j}
			case 'E':
				res.End = astar.Location{X: i, Y: j}
			case '#':
			case '.':
			default:
				panic("invalid character")
			}
		}
	}
	return res
}
