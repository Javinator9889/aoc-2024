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

const search = "XMAS"

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

func part1(input string) (count int) {
	parsed := parseInput(input)
	for i, line := range parsed {
		for j, char := range []byte(line) {
			// Look for the first character of the search string
			if char != search[0] {
				continue
			}
			// Check if the rest of the search string is in any direction
			for _, dir := range []struct{ x, y int }{
				{1, 0},   // right
				{0, 1},   // down
				{1, 1},   // down-right
				{1, -1},  // up-right
				{-1, 1},  // down-left
				{-1, -1}, // up-left
				{-1, 0},  // left
				{0, -1},  // up
			} {
				for k := 1; k < len(search); k++ {
					x, y := i+dir.x*k, j+dir.y*k
					if x < 0 || x >= len(parsed) || y < 0 || y >= len(parsed[0]) {
						break
					}
					if parsed[x][y] != search[k] {
						break
					}
					if k == len(search)-1 {
						count++
					}
				}
			}
		}
	}

	return count
}

func part2(input string) int {
	return 0
}

func parseInput(input string) []string {
	return strings.Split(input, "\n")
}
