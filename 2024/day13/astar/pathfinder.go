package astar

import "math"

// A PathFinder represents the possible states we can reach in the grid
type PathFinder struct {
	pos    Location    // The claw's current position
	button *Button     // The button we have pressed
	origin *PathFinder // The previous state
}

// Neighbors returns the possible states we can reach from the current state
func (p PathFinder) Neighbors(a *Arcade) []PathFinder {
	neighbors := make([]PathFinder, 0)
	for _, op := range a.Buttons {
		neighbor := PathFinder{
			button: op,
			pos:    p.pos.Add(op.Increment),
			origin: &p,
		}
		neighbors = append(neighbors, neighbor)
	}
	return neighbors
}

func (l Location) distance(a Location) float64 {
	return math.Sqrt(math.Pow(float64(l.X-a.X), 2) + math.Pow(float64(l.Y-a.Y), 2))
}

// EstimatedCost returns the estimated cost to reach the goal. The closer the number is to
// the end of the array, the lower the cost
func (p PathFinder) EstimatedCost(a *Arcade) int {
	return int(p.pos.distance(a.Prize))
}
