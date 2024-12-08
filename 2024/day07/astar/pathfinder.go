package astar

import "github.com/Javinator9889/aoc-2024/2024/day07/ops"

// A PathFinder represents the possible states we can reach in the grid
type PathFinder struct {
	pos       int         // The position in the numbers array
	value     int         // The value of the number
	operation ops.Op      // The operation to apply with the previous number
	origin    *PathFinder // The previous state
}

// Neighbors returns the possible states we can reach from the current state
func (p PathFinder) Neighbors(g *Grid) []PathFinder {
	neighbors := make([]PathFinder, 0)
	// Check if we are at the end of the numbers
	if p.pos == len(g.Numbers)-1 {
		return neighbors
	}
	for _, op := range g.ValidOps {
		neighbor := PathFinder{
			value:     g.Numbers[p.pos+1],
			operation: op,
			pos:       p.pos + 1,
			origin:    &p,
		}
		neighbors = append(neighbors, neighbor)
	}
	return neighbors
}

// EstimatedCost returns the estimated cost to reach the goal. The closer the number is to
// the end of the array, the lower the cost
func (p PathFinder) EstimatedCost(g *Grid) int {
	return -p.pos
}
