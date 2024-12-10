package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Javinator9889/aoc-2024/cast"
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

type Position struct {
	height  int
	visited map[Coordinate]struct{}
}

type Grid [][]*Position

func (g Grid) OutOfBounds(c Coordinate) bool {
	return c.i < 0 || c.i >= len(g) || c.j < 0 || c.j >= len(g[0])
}

func (g Grid) String() string {
	var sb strings.Builder
	for _, row := range g {
		for _, pos := range row {
			sb.WriteString(fmt.Sprintf("%d", pos.height))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

type Coordinate struct {
	i, j int
}

func (c Coordinate) Add(other Coordinate) Coordinate {
	return Coordinate{c.i + other.i, c.j + other.j}
}

var UP = Coordinate{-1, 0}
var DOWN = Coordinate{1, 0}
var LEFT = Coordinate{0, -1}
var RIGHT = Coordinate{0, 1}
var COORDS = []Coordinate{UP, DOWN, LEFT, RIGHT}

func Trailhead(from Coordinate, start Coordinate, grid Grid, countTotal bool) int {
	current := grid[start.i][start.j]
	if current.height == 9 {
		if countTotal {
			return 1
		}
		if _, ok := current.visited[from]; !ok {
			current.visited[from] = struct{}{}
			return 1
		}
		return 0
	}
	paths := make([]Coordinate, 0)
	for _, c := range COORDS {
		next := start.Add(c)
		if grid.OutOfBounds(next) {
			continue
		}
		if grid[next.i][next.j].height != current.height+1 {
			continue
		}
		paths = append(paths, next)
	}
	if len(paths) == 0 {
		return 0
	}
	ans := 0
	for _, p := range paths {
		ans += Trailhead(from, p, grid, countTotal)
	}
	return ans
}

func part1(input string) (reachable int) {
	parsed := parseInput(input)
	slog.Debug("grid", "grid", parsed)
	for i := range parsed {
		for j := range parsed[i] {
			pos := parsed[i][j]
			if pos.height != 0 {
				continue
			}
			from := Coordinate{i, j}
			reachable += Trailhead(from, from, parsed, false)
		}
	}

	return
}

func part2(input string) (reachable int) {
	parsed := parseInput(input)
        slog.Debug("grid", "grid", parsed)
        for i := range parsed {
                for j := range parsed[i] {
                        pos := parsed[i][j]
                        if pos.height != 0 {
                                continue
                        }
                        from := Coordinate{i, j}
                        reachable += Trailhead(from, from, parsed, true)
                }
        }

        return
}

func parseInput(input string) (ans Grid) {
	for i, line := range strings.Split(input, "\n") {
		ans = append(ans, make([]*Position, 0, len(line)))
		for _, c := range line {
			pos := &Position{height: cast.ToInt(string(c))}
			if pos.height == 9 {
				pos.visited = make(map[Coordinate]struct{}, 0)
			}
			ans[i] = append(ans[i], pos)
		}
	}
	return
}
