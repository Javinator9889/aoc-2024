package astar

import (
	"container/heap"
	"fmt"
	"strings"

	"github.com/Javinator9889/aoc-2024/2024/day07/ops"
)

// A node represents a possible state in the grid
type Node struct {
	finder PathFinder
	cost   int
	index  int
	rank   int
	prev   *Node
	open   bool
	closed bool
}

// A Grid represents the problem to solve
type Grid struct {
	Goal        int            // The goal to reach by applying operations to the numbers
	Numbers     []int          // The numbers to apply operations to
	ValidOps    []ops.Op       // The valid operations to apply
	IsValidCost func(int) bool // A function to validate the cost of the operation
}

// A Path represents a sequence of operations to reach the goal
type Path []*Node

func (p Path) String() string {
	var sb strings.Builder
	for _, n := range p {
		op := ""
		if n.finder.operation != "" {
			op = " " + n.finder.operation.String() + " "
		}
		sb.WriteString(fmt.Sprintf("%s%v", op, n.finder.value))
	}
	return sb.String()
}

type nodeMap map[PathFinder]*Node

func (nm nodeMap) Get(p PathFinder) *Node {
	res, ok := nm[p]
	if !ok {
		res = &Node{
			finder: p,
			rank:   -1,
		}
		nm[p] = res
	}
	return res
}

// Runs the A* algorithm to find the path to the goal, if possible. The algorithm considers the
// next number in the array as the reference. The next states are built by combining the set of
// valid operations with the next number. The algorithm stops when the goal is reached or there
// are no more nodes to explore.
//
// If `exhaustive` is false, the algorithm will try to find any path that yields the result.
// Otherwise, it will try to find the path that uses all the numbers in the array. The `cost`
// always refers the cumulative result of the operations. Depending on the operations available
// in the grid, it's possible to have multiple paths to the goal. But it's also possible that
// a path is unreachable from a certain state (e.g. the cost is greater than the goal and there
// are no "sub" operations). The `IsValidCost` function is used to discard states - thus optimize
// the search - by checking if the cost is valid. Simply returning `true` will make the algorithm
// to explore all the states, which could take longer.
func (g *Grid) AStar(exhaustive bool) Path {
	nm := nodeMap{}
	nq := &PriorityQueue{}
	heap.Init(nq)
	// We consider the `cost` the cumulative result of the operations
	from := PathFinder{value: g.Numbers[0]}
	fromNode := nm.Get(from)
	fromNode.open = true
	fromNode.cost = from.value
	fromNode.prev = nil
	heap.Push(nq, fromNode)

	for {
		// There are no more nodes to explore
		if nq.Len() == 0 {
			break
		}
		current := heap.Pop(nq).(*Node)
		current.open = false
		current.closed = true

		// We reached the goal. Verify that we used all the numbers
		if (current.cost == g.Goal && !exhaustive) ||
			(current.cost == g.Goal && exhaustive && current.finder.pos == len(g.Numbers)-1) {
			var path Path
			for current != nil {
				path = append(Path{current}, path...)
				current = current.prev
			}
			return path
		}

		// Get the neighbors of the current node
		for _, neighbor := range current.finder.Neighbors(g) {
			cost := neighbor.operation.Cal(current.cost, neighbor.value)
			// Skip if the cost is invalid
			if !g.IsValidCost(cost) {
				continue
			}
			neighborNode := nm.Get(neighbor)
			// If we are closer to the goal...
			if cost > neighborNode.cost {
				if neighborNode.open {
					heap.Remove(nq, neighborNode.index)
				}
				neighborNode.open = false
				neighborNode.closed = false
			}
			if !neighborNode.open && !neighborNode.closed {
				neighborNode.open = true
				neighborNode.rank = neighbor.EstimatedCost(g) - cost
				neighborNode.cost = cost
				neighborNode.prev = current
				heap.Push(nq, neighborNode)
			}
		}
	}
	return nil
}
