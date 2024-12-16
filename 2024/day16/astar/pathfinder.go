package astar

import (
	"fmt"
	"math"
)

// A PathFinder represents the possible states we can reach in the grid
type PathFinder struct {
	pos      Location // The reindeer's current position
	dir      Location // The reindeer's current direction
	moveCost int      // The cost to move to this state
}

func (p PathFinder) String() string {
	dirToStr := map[Location]string{
		NORTH: "NORTH",
		SOUTH: "SOUTH",
		WEST:  "WEST",
		EAST:  "EAST",
	}
	return fmt.Sprintf("Pos: %v, Dir: %v", p.pos, dirToStr[p.dir])
}

// Neighbors returns the possible states we can reach from the current state
func (p PathFinder) Neighbors(r *Reindlympics) []PathFinder {
	neighbors := make([]PathFinder, 0)
	for _, move := range []Location{NORTH, EAST, SOUTH, WEST} {
		// Ignore moves >90 degrees
		if move.X*p.dir.X+move.Y*p.dir.Y == -1 {
			continue
		}
		// If we're changing directions, we can just rotate
		if move.X != p.dir.X || move.Y != p.dir.Y {
			neighbors = append(neighbors, PathFinder{
				pos:      p.pos,
				dir:      move,
				moveCost: 1000,
			})
			continue
		}
		// Otherwise, we should move forward
		dst := p.pos.Add(move)
		if dst.X < 0 || dst.X >= len(r.Grid) || dst.Y < 0 || dst.Y >= len(r.Grid[0]) {
			continue
		}
		// We found a wall
		if r.Grid[dst.X][dst.Y] == '#' {
			continue
		}
		neighbors = append(neighbors, PathFinder{
			pos:      dst,
			dir:      p.dir,
			moveCost: 1,
		})
	}

	return neighbors
}

func (l Location) distance(a Location) float64 {
	return math.Sqrt(math.Pow(float64(l.X-a.X), 2) + math.Pow(float64(l.Y-a.Y), 2))
}

// EstimatedCost returns the estimated cost to reach the goal
func (p PathFinder) EstimatedCost(r *Reindlympics) int {
	return int(p.pos.distance(r.End))
}
