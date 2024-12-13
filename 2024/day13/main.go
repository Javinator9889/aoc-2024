package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"math"
	"regexp"
	"strings"

	"github.com/Javinator9889/aoc-2024/2024/day13/astar"
	"github.com/Javinator9889/aoc-2024/cast"
	"github.com/Javinator9889/aoc-2024/util"
)

// See: https://regex101.com/r/fuhDlN/1
var buttonRe = regexp.MustCompile(`Button ([A-Z]): X\+(\d+), Y\+(\d+)`)
var prizeRe = regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

const MAX_PRESSES = 100
const INF = math.MaxInt

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

var ORIGIN = astar.Location{X: 0, Y: 0}

func Cramer(a *astar.Arcade) (int, int, error) {
	// If we think of the game as a linalg problem to solve, Cramer rule is the way to go. Consider
	// the following system of equations:
	// 		A1 * x + A2 * y = Prize1
	// 		B1 * x + B2 * y = Prize2
	// We can solve this system by calculating the following determinants:
	// 		D = A1 * B2 - A2 * B1
	// 		Dx = Prize1 * B2 - Prize2 * A2
	// 		Dy = A1 * Prize2 - A2 * Prize1
	// And the solution is:
	// 		x = Dx / D
	// 		y = Dy / D
	// Where Dx and Dy are the determinants of the system with the Prize values
	a1 := float64(a.Buttons["A"].Increment.X)
	a2 := float64(a.Buttons["A"].Increment.Y)
	b1 := float64(a.Buttons["B"].Increment.X)
	b2 := float64(a.Buttons["B"].Increment.Y)
	px := float64(a.Prize.X)
	py := float64(a.Prize.Y)
	det := a1*b2 - a2*b1
	if det == 0 {
		// We cannot solve the system. Use the A* algorithm to find the solution
		path := a.AStar()
		if path == nil {
			return 0, 0, fmt.Errorf("no path found")
		}
		buttonPresses := path.Count()
		return buttonPresses["A"], buttonPresses["B"], nil
	}
	x := (px*b2 - py*b1) / det
	y := (a1*py - a2*px) / det
	return int(math.Trunc(x)), int(math.Trunc(y)), nil
}

func verify(a *astar.Arcade, x, y, max int) bool {
	if x < 0 || y < 0 || x > max || y > max {
		return false
	}
	return a.Buttons["A"].Increment.X*x+a.Buttons["B"].Increment.X*y == a.Prize.X &&
		a.Buttons["A"].Increment.Y*x+a.Buttons["B"].Increment.Y*y == a.Prize.Y
}

func part1(input string) (cost int) {
	arcades := parseInput(input)
	for _, arcade := range arcades {
		arcade.MaxPresses = MAX_PRESSES
		a, b, err := Cramer(arcade)
		if err != nil {
			slog.Warn("error solving system", "arcade", arcade, "error", err)
			continue
		}
		if !verify(arcade, a, b, MAX_PRESSES) {
			slog.Warn("invalid solution", "arcade", arcade, "a", a, "b", b)
			continue
		}
		slog.Debug("solution", "arcade", arcade, "a", a, "b", b)
		cost += a*arcade.Buttons["A"].Tokens + b*arcade.Buttons["B"].Tokens
	}

	return
}

func part2(input string) (cost int) {
	arcades := parseInput(input)
	for _, arcade := range arcades {
		arcade.MaxPresses = INF
		arcade.Prize.X += 10_000_000_000_000
		arcade.Prize.Y += 10_000_000_000_000
		a, b, err := Cramer(arcade)
		if err != nil {
			slog.Warn("error solving system", "arcade", arcade, "error", err)
			continue
		}
		if !verify(arcade, a, b, INF) {
			slog.Warn("invalid solution", "arcade", arcade, "a", a, "b", b)
			continue
		}
		slog.Debug("solution", "arcade", arcade, "a", a, "b", b)
		cost += a*arcade.Buttons["A"].Tokens + b*arcade.Buttons["B"].Tokens
	}

	return
}

func parseInput(input string) (arcades []*astar.Arcade) {
	current := &astar.Arcade{}
	for _, line := range strings.Split(input, "\n") {
		matches := buttonRe.FindStringSubmatch(line)
		if matches != nil {
			var tokens int
			if matches[1] == "A" {
				tokens = 3
			} else if matches[1] == "B" {
				tokens = 1
			} else {
				panic("invalid button")
			}
			if current.Buttons == nil {
				current.Buttons = make(map[string]*astar.Button)
			}
			current.Buttons[matches[1]] = &astar.Button{
				ID: matches[1],
				Increment: astar.Location{
					X: cast.ToInt(matches[2]),
					Y: cast.ToInt(matches[3]),
				},
				Tokens: tokens,
			}
			continue
		}
		matches = prizeRe.FindStringSubmatch(line)
		if matches != nil {
			current.Prize = astar.Location{
				X: cast.ToInt(matches[1]),
				Y: cast.ToInt(matches[2]),
			}
			// No data to parse, append to the list
			arcades = append(arcades, current)
			current = &astar.Arcade{}
		}
	}
	return
}
