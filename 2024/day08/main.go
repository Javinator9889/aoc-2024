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

const (
	OPEN = '.' // Open space
)

type Loc struct {
	frequency rune
	x, y      int
	antinode  bool
}

type Frequency rune

type Location struct {
	x, y int
}

type Vector Location

func (l Location) Direction(o Location) Vector {
	return Vector{o.x - l.x, o.y - l.y}
}

// The reflection of a point simply takes the direction vector of the line between the two points
// and adds it to the second point.
func (l Location) Reflection(o Location) Location {
	vd := l.Direction(o)
	px, py := o.x+vd.x, o.y+vd.y
	return Location{px, py}
}

type Antennas map[Frequency][]Location

type Grid [][]*Loc

func (g Grid) InBounds(l Location) bool {
	return l.x >= 0 && l.x < len(g) && l.y >= 0 && l.y < len(g[0])
}

func part1(input string) (antinodes int) {
	grid, antennas := parseInput(input)
	// Calculate the antinode of each frequency. The antinode is defined as the point in the
	// grid that is aligned with the antennas of the same frequency. The antennas must be aligned
	// and be twice as far away as the other.
	for freq, locs := range antennas {
		if len(locs) < 2 {
			continue
		}
		// Calculate the reflection of the antennas
		for a := range locs {
			for b := range locs {
				if a == b {
					continue
				}
				reflection := locs[a].Reflection(locs[b])
				if grid.InBounds(reflection) {
					slog.Debug(
						"Reflection",
						"freq", string(freq),
						"a", locs[a],
						"b", locs[b],
						"reflection", reflection,
					)
					if !grid[reflection.x][reflection.y].antinode {
						grid[reflection.x][reflection.y].antinode = true
						antinodes++
					}
				} else {
					slog.Debug(
						"Reflection out of bounds",
						"freq", string(freq),
						"a", locs[a],
						"b", locs[b],
						"reflection", reflection,
					)
				}
			}
		}
	}

	return
}

func part2(input string) (antinodes int) {
	grid, antennas := parseInput(input)
	// For the second part, we have to take into account the resonant harmonics. This simply means
	// the antennas emit in a straight line to any grid position, needed at least two antennas of
	// the same frequency aligned. The antinodes are now located all along the line within the grid
	for freq, locs := range antennas {
		if len(locs) < 2 {
			continue
		}
		// Calculate the reflection of the antennas
		for a := range locs {
			for b := range locs {
				if a == b {
					loc := locs[a]
					if !grid[loc.x][loc.y].antinode {
						grid[loc.x][loc.y].antinode = true
						antinodes++
					}
					continue
				}
				prev := locs[a]
				point := locs[b]
				for {
					reflection := prev.Reflection(point)
					if !grid.InBounds(reflection) {
						slog.Debug(
							"Reflection out of bounds",
							"freq", string(freq),
							"a", prev,
							"b", point,
							"reflection", reflection,
						)
						break
					}
					slog.Debug(
						"Reflection",
						"freq", string(freq),
						"a", prev,
						"b", point,
						"reflection", reflection,
					)
					if !grid[reflection.x][reflection.y].antinode {
						grid[reflection.x][reflection.y].antinode = true
						antinodes++
					}
					prev = point
					point = reflection
				}
			}
		}
	}
	return
}

func parseInput(input string) (Grid, Antennas) {
	res := make(Grid, 0)
	antennas := make(Antennas)
	for x, line := range strings.Split(input, "\n") {
		row := make([]*Loc, len(line))
		for y, r := range line {
			row[y] = &Loc{frequency: r, x: x, y: y}
			if r != OPEN {
				antennas[Frequency(r)] = append(antennas[Frequency(r)], Location{x, y})
			}
		}
		res = append(res, row)
	}
	return res, antennas
}
