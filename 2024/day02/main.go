package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"math"
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
	flag.IntVar(&part, "part", 1, "part 1 or 2")
	flag.Parse()
	fmt.Println("Running part", part)

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

func part1(input string) int {
	parsed := parseInput(input)
	var safe int
	for i := range parsed {
		var safeRow bool
	checker:
		for j := range len(parsed[i]) - 2 {
			check := parsed[i][j : j+3]
			diff := check[2] - check[1]
			prev := check[1] - check[0]
			switch {
			case
				diff == 0,
				prev == 0,
				math.Abs(float64(diff)) > 3,
				math.Abs(float64(prev)) > 3,
				diff < 0 && prev > 0,
				diff > 0 && prev < 0:
				slog.Debug("Not safe", "row", parsed[i], "diff", diff, "prev", prev)
				safeRow = false
				break checker
			default:
				safeRow = true
			}
		}
		if safeRow {
			slog.Debug("Safe", "row", parsed[i])
			safe++
		}
	}

	return safe
}

func part2(input string) int {
	return 0
}

func parseInput(input string) (ans [][]int) {
	for i, line := range strings.Split(input, "\n") {
		numbers := strings.Fields(line)
		ans = append(ans, make([]int, len(numbers)))
		for j, number := range numbers {
			ans[i][j] = cast.ToInt(number)
		}
	}
	return ans
}
