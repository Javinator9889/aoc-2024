package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
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

func part1(input string) (res int) {
	// We're looking for "mul(x,y)" where x and y are integers with 1-3 digits
	// See: https://regex101.com/r/jjrwRa/1
	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	parsed := parseInput(input)
	for _, line := range parsed {
		matches := re.FindAllStringSubmatch(line, -1)
		slog.Debug("matches", "matches", matches, "line", line)
		for match := range matches {
			slog.Debug("match", "match", match)
			// We're looking for pairs of integers, so we can skip the first match
			// Don't check for errors, we know they're integers (otherwise our regex is wrong)
			x, _ := strconv.Atoi(matches[match][1])
			y, _ := strconv.Atoi(matches[match][2])
			res += x * y
		}
	}

	return res
}

func part2(input string) int {
	return 0
}

func parseInput(input string) []string {
	return strings.Split(input, "\n")
}
