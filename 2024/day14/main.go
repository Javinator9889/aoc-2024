package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/Javinator9889/aoc-2024/cast"
	"github.com/Javinator9889/aoc-2024/util"
)

//go:embed input.txt
var input string

var re = regexp.MustCompile(`p=(\d+),(\d+) v=(-?\d+),(-?\d+)`)
var gridSize = Grid{101, 103}

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

type Location struct {
	x, y int
}

type Velocity Location
type Grid Location
type Quadrant struct {
	topLeft, bottomRight Location
}

func (g Grid) WrapAround(l Location) Location {
	inGrid := Location{(l.x + g.x) % g.x, (l.y + g.y) % g.y}
	if inGrid.x < 0 {
		inGrid.x += g.x
	}
	if inGrid.y < 0 {
		inGrid.y += g.y
	}
	return inGrid
}

func (g Grid) Quadrants() []Quadrant {
	return []Quadrant{
		{
			topLeft:     Location{0, 0},
			bottomRight: Location{g.x / 2, g.y / 2},
		},
		{
			topLeft:     Location{g.x / 2, 0},
			bottomRight: Location{g.x, g.y / 2},
		},
		{
			topLeft:     Location{0, g.y / 2},
			bottomRight: Location{g.x / 2, g.y},
		},
		{
			topLeft:     Location{g.x / 2, g.y / 2},
			bottomRight: Location{g.x, g.y},
		},
	}
}

type Robot struct {
	p Location
	v Velocity
}

func (r Robot) Simulate(t time.Duration) Location {
	return Location{
		x: r.p.x + r.v.x*int(t.Seconds()),
		y: r.p.y + r.v.y*int(t.Seconds()),
	}
}

func part1(input string) int {
	robots := parseInput(input)
	locations := make([]Location, 0)
	for _, r := range robots {
		location := r.Simulate(100 * time.Second)
		locations = append(locations, gridSize.WrapAround(location))
	}

	safetyFactor := 1
	for _, quadrant := range gridSize.Quadrants() {
		count := 0
		for _, location := range locations {
			// Skip the location if it's in the middle of the grid
			if location.x == gridSize.x/2 || location.y == gridSize.y/2 {
				continue
			}
			if location.x >= quadrant.topLeft.x && location.x <= quadrant.bottomRight.x &&
				location.y >= quadrant.topLeft.y && location.y <= quadrant.bottomRight.y {
				count++
			}
		}
		safetyFactor *= count
	}

	return safetyFactor
}

func (g Grid) Display(robots []*Robot) {
	// cmd := exec.Command("clear") // Clear the screen
	// cmd.Stdout = os.Stdout
	// cmd.Run()
	for x := 0; x < g.x; x++ {
		line := make([]int, g.y)
		for _, r := range robots {
			if r.p.x == x {
				line[r.p.y]++
			}
		}
		for _, c := range line {
			if c == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print(c)
			}
		}
		fmt.Println()
	}
}

func part2(input string) int {
	robots := parseInput(input)
	origin := make([]Robot, len(robots))
	for i, r := range robots {
		origin[i] = *r
	}
	i := 0
	for {
		for _, r := range robots {
			r.p = gridSize.WrapAround(r.Simulate(time.Second))
		}
		gridSize.Display(robots)
		i++
		fmt.Printf("Iteration %d\n\n\n", i)
		// Check if all robots are in the same position
		for j, r := range robots {
			if r.p != origin[j].p {
				break
			}
			if j == len(robots)-1 {
				slog.Info("Robots have looped back to their original positions")
				return i
			}
		}
	}
	return 0
}

func parseInput(input string) (robots []*Robot) {
	for _, line := range strings.Split(input, "\n") {
		matches := re.FindStringSubmatch(line)
		if matches == nil {
			panic("invalid input")
		}
		robots = append(robots, &Robot{
			p: Location{
				x: cast.ToInt(matches[1]),
				y: cast.ToInt(matches[2]),
			},
			v: Velocity{
				x: cast.ToInt(matches[3]),
				y: cast.ToInt(matches[4]),
			},
		})
	}
	return
}
