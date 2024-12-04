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

const search1 = "XMAS"
const search2 = "MAS"
const search2Rev = "SAM" // To simplify, just reverse the search string to do a full match

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
			if char != search1[0] {
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
				for k := 1; k < len(search1); k++ {
					x, y := i+dir.x*k, j+dir.y*k
					if x < 0 || x >= len(parsed) || y < 0 || y >= len(parsed[0]) {
						break
					}
					if parsed[x][y] != search1[k] {
						break
					}
					if k == len(search1)-1 {
						count++
					}
				}
			}
		}
	}

	return count
}

func part2(input string) (count int) {
	parsed := parseInput(input)
	// We are supposed to find "X-MAS" in the input, forming an 'X' shape, like:
	// 	 M.S
	// 	 .A.
	// 	 M.S
	// Where the dots are any character. It has to be read "MAS" in any direction
	// as long as it forms an 'X' shape. We change the problem a little bit to look for
	// 'A's - which is going to always be in the center - and then check if the surrounding
	// characters form an 'X' shape.
	//
	// For the sake of simplicity, we skip the first and last lines as well as the first and
	// last columns. This way we can check the surrounding characters without worrying about
	// going out of bounds.
	for i := 1; i < len(parsed)-1; i++ {
		for j := 1; j < len(parsed[0])-1; j++ {
			if parsed[i][j] != search2[1] {
				continue
			}
			// Check if the surrounding characters form an 'X' shape
			cross1 := string([]byte{parsed[i-1][j-1], parsed[i][j], parsed[i+1][j+1]})
			cross2 := string([]byte{parsed[i-1][j+1], parsed[i][j], parsed[i+1][j-1]})
			if (cross1 != search2 && cross1 != search2Rev) ||
				(cross2 != search2 && cross2 != search2Rev) {
				continue
			}
			count++
		}
	}

	return count
}

func parseInput(input string) []string {
	return strings.Split(input, "\n")
}
