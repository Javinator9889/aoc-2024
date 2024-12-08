package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Javinator9889/aoc-2024/2024/day07/astar"
	"github.com/Javinator9889/aoc-2024/2024/day07/ops"
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

type Row struct {
	value   int
	numbers []int
}

func part1(input string) (solvable int) {
	// We can consider this exercise as a graph problem, where each node is a number with
	// a set of valid operations to apply to. The goal is to reach a certain number by using
	// all the numbers in the array. We can use the A* algorithm to find the path to the goal
	// by considering the next number in the array as the reference. The next states are built
	// by combining the set of valid operations with the next number in the array.
	// The cost of the operation is the value of the number, and the heuristic is the position
	// of the number in the array. The algorithm stops when the cost of the operation is greater
	// than the goal.
	// The algorithm is implemented in the astar package.
	parsed := parseInput(input)
	for _, row := range parsed {
		grid := astar.Grid{
			Goal:     row.value,
			Numbers:  row.numbers,
			ValidOps: []ops.Op{ops.ADD, ops.MUL},
			IsValidCost: func(cost int) bool {
				return cost <= row.value
			},
		}
		path := grid.AStar(true /* exhaustive */)
		// We have to use all the numbers
		if path != nil && len(path) == len(row.numbers) {
			slog.Debug("Path for", "n", row.value, "path", path)
			solvable += row.value
		} else if path != nil {
			slog.Warn("Not using all numbers", "n", row.value, "path", path, "nums", row.numbers)
		}
	}
	return
}

func part2(input string) int {
	return 0
}

func parseInput(input string) (ans []Row) {
	ans = make([]Row, 0)
	for _, line := range strings.Split(input, "\n") {
		items := strings.Split(line, ": ")
		k, v := cast.ToInt(items[0]), strings.Split(items[1], " ")
		row := Row{
			value:   k,
			numbers: make([]int, len(v)),
		}
		for i := range v {
			row.numbers[i] = cast.ToInt(v[i])
		}
		ans = append(ans, row)
	}
	return
}
