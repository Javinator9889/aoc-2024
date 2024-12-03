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

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
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

func Safe(diff int) bool {
	return math.Abs(float64(diff)) <= 3 && diff != 0
}

func SameDirection(diff, prev int) bool {
	return diff < 0 && prev < 0 || diff > 0 && prev > 0
}

func ValidSlice(slice []int) bool {
	return Safe(slice[2]-slice[1]) &&
		Safe(slice[1]-slice[0]) &&
		SameDirection(slice[2]-slice[1], slice[1]-slice[0])
}

func ProcessRow(row []int, removals int, maxRemovals int) bool {
	var safeRow bool
check:
	for j := range len(row) - 2 {
		switch {
		case !ValidSlice(row[j : j+3]):
			slog.Debug("Checking", "row", row, "slice", row[j:j+3])
			if removals < maxRemovals {
				// Try removing all the three elements
				for k := 0; k < 3; k++ {
					niter := append([]int{}, row[max(0, j-1):j+k]...)
					niter = append(niter, row[j+k+1:]...)
					slog.Debug("next iter", "row", niter)
					if ProcessRow(niter, removals+1, maxRemovals) {
						return true
					}
				}
			}
			safeRow = false
			break check
		default:
			safeRow = true
		}
	}
	return safeRow
}

func part1(input string) int {
	parsed := parseInput(input)
	var safe int
	for i := range parsed {
		if ProcessRow(parsed[i], 0, 0) { // no removals in part 1
			slog.Debug("Safe", "row", parsed[i])
			safe++
		} else {
			slog.Debug("Unsafe", "row", parsed[i])
		}
	}

	return safe
}

func part2(input string) int {
	parsed := parseInput(input)
	var safe int
	for i := range parsed {
		if ProcessRow(parsed[i], 0, 1) { // one removal in part 2
			slog.Debug("Safe", "row", parsed[i])
			safe++
		} else {
			slog.Debug("Unsafe", "row", parsed[i])
		}
	}

	return safe
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
