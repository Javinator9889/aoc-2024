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

// We're looking for "mul(x,y)" where x and y are integers with 1-3 digits
// See: https://regex101.com/r/jjrwRa/1
var re = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)

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

func eval(line string) (res int) {
	matches := re.FindAllStringSubmatch(line, -1)
	slog.Debug("matches", "matches", matches, "line", line)
	for match := range matches {
		// We're looking for pairs of integers, so we can skip the first match
		// Don't check for errors, we know they're integers (otherwise our regex is wrong)
		x, _ := strconv.Atoi(matches[match][1])
		y, _ := strconv.Atoi(matches[match][2])
		res += x * y
	}
	return res
}

func part1(input string) (res int) {
	parsed := parseInput(input)
	for _, line := range parsed {
		res += eval(line)
	}

	return res
}

func part2(input string) (res int) {
	parsed := parseInput(input)
	// Look for the `do()` and `don't()` instructions. They enable/disable the multiplication
	// See: https://regex101.com/r/do7Zxf/1
	en := regexp.MustCompile(`(do|don't)\(\)`)
	// We start with the multiplication enabled. If a line disables it, keep disabled until
	// the next `do()` instruction
	enabled := true
	for _, line := range parsed {
		index := 0
		slog.Debug("parsing", "line", line)
		// Check if the line enables or disables the multiplication
		endis := en.FindAllStringSubmatchIndex(line, -1)
		for _, match := range endis {
			slog.Debug("match", "match", line[match[2]:match[3]], "enabled", enabled, "index", index)
			// If the match is "do()", we enable the multiplication
			switch line[match[2]:match[3]] {
			case "do":
				// If not enabled, enable the multiplication and update the index
				if !enabled {
					enabled = true
					index = match[3]
				}
			case "don't":
				// If enabled, calculate the current multiplications
				// and disable until the next "do()"
				if enabled {
					enabled = false
					res += eval(line[index:match[2]])
				}
			default:
				slog.Error("unknown", "match", line[match[2]:match[3]])
			}
		}
		// Match:
		// - Last instruction is "do()"
		// - There's no "don't()" instruction
		// - The multiplication is enabled
		if enabled {
			res += eval(line[index:])
		}
	}
	return res
}

func parseInput(input string) []string {
	return strings.Split(input, "\n")
}
