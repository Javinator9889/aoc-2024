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

const (
	EMPTY     = "."
	WALL      = "#"
	OBJECT    = "O"
	ROBOT     = "@"
	BOX_SIDE1 = "["
	BOX_SIDE2 = "]"
	UP        = "^"
	DOWN      = "v"
	LEFT      = "<"
	RIGHT     = ">"
)

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

type Coordinates struct {
	x, y int
}

func (c Coordinates) Add(other Coordinates) Coordinates {
	return Coordinates{c.x + other.x, c.y + other.y}
}

type Element interface {
	Move(dir Coordinates, grid Grid) bool
	WouldMove(dir Coordinates, grid Grid) bool
	Position() Coordinates
	String() string
}

type Wall struct {
	pos Coordinates
}

func (w *Wall) Move(dir Coordinates, grid Grid) bool {
	return false
}

func (w *Wall) WouldMove(dir Coordinates, grid Grid) bool {
	return false
}

func (w *Wall) Position() Coordinates {
	return w.pos
}

func (w *Wall) String() string {
	return WALL
}

type Object struct {
	pos Coordinates
}

func (o *Object) WouldMove(dir Coordinates, grid Grid) bool {
	dst := Coordinates{o.pos.x + dir.x, o.pos.y + dir.y}
	if dst.x < 0 || dst.y < 0 || dst.x >= len(grid) || dst.y >= len(grid[0]) {
		return false
	}
	if grid[dst.x][dst.y] == nil {
		return true
	}
	return grid[dst.x][dst.y].WouldMove(dir, grid)
}

func (o *Object) Move(dir Coordinates, grid Grid) bool {
	dst := Coordinates{o.pos.x + dir.x, o.pos.y + dir.y}
	if dst.x < 0 || dst.y < 0 || dst.x >= len(grid) || dst.y >= len(grid[0]) {
		return false
	}
	// If there is free space, just move
	if grid[dst.x][dst.y] == nil {
		grid[dst.x][dst.y] = o
		grid[o.pos.x][o.pos.y] = nil
		o.pos = dst
		return true
	}
	// Otherwise, we should try to push any object in the way
	if !grid[dst.x][dst.y].Move(dir, grid) {
		return false
	}
	// If we could move the object, we can move
	grid[dst.x][dst.y] = o
	grid[o.pos.x][o.pos.y] = nil
	o.pos = dst
	return true
}

func (o *Object) Position() Coordinates {
	return o.pos
}

func (o *Object) String() string {
	return OBJECT
}

// Robot is an Object that can move on its own
type Robot Object

func (r *Robot) Move(dir Coordinates, grid Grid) bool {
	moved := (*Object)(r).Move(dir, grid)
	if moved {
		// Cast to Robot to avoid infinite recursion
		grid[r.pos.x][r.pos.y] = (*Robot)(r)
	}
	return moved
}

func (r *Robot) WouldMove(dir Coordinates, grid Grid) bool {
	return (*Object)(r).WouldMove(dir, grid)
}

func (r *Robot) Position() Coordinates {
	return r.pos
}

func (r *Robot) String() string {
	return ROBOT
}

type Box struct {
	o    Object
	side string
	next *Box
}

func (b *Box) moveObject(dir Coordinates, grid Grid) bool {
	ret := b.o.Move(dir, grid)
	if ret {
		grid[b.o.Position().x][b.o.Position().y] = b
	}
	return ret
}

func (b *Box) Move(dir Coordinates, grid Grid) bool {
	// Check if the other side is in the same direction as the one we're moving
	dst := b.o.Position().Add(dir)
	if dst == b.next.o.Position() {
		moved := b.next.moveObject(dir, grid)
		if moved {
			b.moveObject(dir, grid)
			return true
		}
		return false
	}
	// If the other side is not in the same direction, we should verify first if we can move
	if !b.WouldMove(dir, grid) {
		return false
	}
	// If we can move, we should move the other side first
	b.next.moveObject(dir, grid)
	b.moveObject(dir, grid)
	return true
}

func (b *Box) WouldMove(dir Coordinates, grid Grid) bool {
	dst := b.o.Position().Add(dir)
	if dst == b.next.o.Position() {
		return b.next.o.WouldMove(dir, grid)
	}
	return b.o.WouldMove(dir, grid) && b.next.o.WouldMove(dir, grid)
}

func (b *Box) Position() Coordinates {
	return b.o.Position()
}

func (b *Box) String() string {
	return b.side
}

type Grid [][]Element

func (g Grid) String() string {
	var sb strings.Builder
	for _, row := range g {
		for _, elem := range row {
			if elem == nil {
				sb.WriteString(EMPTY)
			} else {
				sb.WriteString(elem.String())
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func part1(input string) (coordinates int) {
	grid, robot, moves := parseInput(input)
	slog.Debug("Parsed Grid", "grid", grid, "robot", robot, "moves", moves)
	fmt.Println(grid)

	for _, move := range moves {
		dst := robot.Add(move)
		rbt := grid[robot.x][robot.y]
		moved := rbt.Move(move, grid)
		slog.Debug(
			"Robot attempted to move",
			"origin", robot,
			"dst", dst,
			"moved", moved,
		)
		if moved {
			robot = dst
		}
	}
	slog.Debug("Final Grid", "grid", grid)
	fmt.Println(grid)
	for _, row := range grid {
		for _, elem := range row {
			if elem == nil || elem.String() != OBJECT {
				continue
			}
			pos := elem.Position()
			coordinates += (pos.x * 100) + pos.y
		}
	}

	return
}

func part2(input string) (coordinates int) {
	grid, robot, moves := parseInput2(input)
	slog.Debug("Parsed Grid", "grid", grid, "robot", robot, "moves", moves)
	fmt.Println(grid)

	for _, move := range moves {
		dst := robot.Add(move)
		rbt := grid[robot.x][robot.y]
		moved := rbt.Move(move, grid)
		slog.Debug(
			"Robot attempted to move",
			"origin", robot,
			"dst", dst,
			"moved", moved,
		)
		if moved {
			fmt.Println(grid)
			robot = dst
		}
	}
	slog.Debug("Final Grid", "grid", grid)
	fmt.Println(grid)
	for _, row := range grid {
		for _, elem := range row {
			if elem == nil || elem.String() != BOX_SIDE1 {
				continue
			}
			pos := elem.Position()
			coordinates += (pos.x * 100) + pos.y
		}
	}

	return
}

func parseInput(input string) (grid Grid, robot Coordinates, moves []Coordinates) {
	for x, line := range strings.Split(input, "\n") {
		row := make([]Element, 0)
		for y, char := range line {
			switch string(char) {
			case EMPTY:
				row = append(row, nil)
			case WALL:
				row = append(row, &Wall{Coordinates{x, y}})
			case OBJECT:
				row = append(row, &Object{Coordinates{x, y}})
			case ROBOT:
				robot = Coordinates{x, y}
				row = append(row, &Robot{robot})
			case UP:
				moves = append(moves, Coordinates{-1, 0})
			case DOWN:
				moves = append(moves, Coordinates{1, 0})
			case LEFT:
				moves = append(moves, Coordinates{0, -1})
			case RIGHT:
				moves = append(moves, Coordinates{0, 1})
			default:
				slog.Warn("Ignoring unknown character", "char", char)
			}
		}
		if len(row) > 0 {
			grid = append(grid, row)
		}
	}
	return
}

func parseInput2(input string) (grid Grid, robot Coordinates, moves []Coordinates) {
	for x, line := range strings.Split(input, "\n") {
		row := make([]Element, 0)
		for y, char := range line {
			// Everything is duplicated, so we need to multiply by 2
			realY := y * 2
			switch string(char) {
			case EMPTY:
				row = append(row, nil, nil)
			case WALL:
				row = append(row, &Wall{Coordinates{x, realY}}, &Wall{Coordinates{x, realY + 1}})
			case OBJECT:
				side1 := &Box{Object{Coordinates{x, realY}}, BOX_SIDE1, nil}
				side2 := &Box{Object{Coordinates{x, realY + 1}}, BOX_SIDE2, side1}
				side1.next = side2
				row = append(row, side1, side2)
			case ROBOT:
				robot = Coordinates{x, realY}
				row = append(row, &Robot{robot}, nil)
			case UP:
				moves = append(moves, Coordinates{-1, 0})
			case DOWN:
				moves = append(moves, Coordinates{1, 0})
			case LEFT:
				moves = append(moves, Coordinates{0, -1})
			case RIGHT:
				moves = append(moves, Coordinates{0, 1})
			default:
				slog.Warn("Ignoring unknown character", "char", char)
			}
		}
		if len(row) > 0 {
			slog.Debug("Row", "row", row)
			grid = append(grid, row)
		}
	}
	return
}
