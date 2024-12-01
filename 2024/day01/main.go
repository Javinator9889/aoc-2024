package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math"
	"slices"
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
	diff := 0
	first, second := parseInput(input)

	// Sort the slices
	slices.Sort(first)
	slices.Sort(second)

	// Iterate over the slices. They should have the same length
	for i := range first {
		diff += int(math.Abs(float64(first[i] - second[i])))
	}

	return diff
}

func part2(input string) int {
	diff := 0
	occurrences := make(map[int]int) // Counts the number of times a number appears
	first, second := parseInput(input)

	// Find the occurrences on the second slice
	for _, num := range second {
		occurrences[num]++
	}

	// For each number in the first slice, check how many times it appears in the second slice
	// and add the number * occurrences to the diff variable
	for _, num := range first {
		diff += num * occurrences[num]
	}

	return diff
}

func parseInput(input string) (first []int, second []int) {
	for _, line := range strings.Split(input, "\n") {
		items := strings.Fields(line)
		first = append(first, cast.ToInt(items[0]))
		second = append(second, cast.ToInt(items[1]))
	}
	return first, second
}
