package astar

import (
	"container/heap"
	"strings"

	"github.com/Javinator9889/aoc-2024/cast"
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

type Location struct {
	X, Y int
}

func (l Location) Add(other Location) Location {
	return Location{
		X: l.X + other.X,
		Y: l.Y + other.Y,
	}
}

var ORIGIN = Location{X: 0, Y: 0}

type Button struct {
	ID        string
	Increment Location
	Tokens    int

	pressed int // The number of times the button has been pressed
}

func (b *Button) String() string {
	return b.ID + ":" + cast.ToString(b.Increment.X) + "+X," + cast.ToString(b.Increment.Y) + "+Y"
}

type Arcade struct {
	Prize      Location
	Buttons    map[string]*Button
	MaxPresses int
}

// A Path represents a sequence of operations to reach the goal
type Path []*Node

func (p Path) String() string {
	var sb strings.Builder
	button := make(map[string]int)
	for _, n := range p {
		if n.finder.button != nil {
			button[n.finder.button.ID]++
		}
	}
	for k, v := range button {
		sb.WriteString(k)
		sb.WriteString(":")
		sb.WriteString(cast.ToString(v))
		sb.WriteString(" ")
	}
	return sb.String()
}

func (p Path) Cost() (cost int) {
	for _, n := range p {
		if n.finder.button != nil {
			cost += n.finder.button.Tokens
		}
	}
	return
}

func (p Path) Count() map[string]int {
	button := make(map[string]int)
	for _, n := range p {
		if n.finder.button != nil {
			button[n.finder.button.ID]++
		}
	}
	return button
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
func (a *Arcade) AStar() Path {
	nm := nodeMap{}
	nq := &PriorityQueue{}
	heap.Init(nq)
	// We consider the `cost` the cumulative result of the operations
	from := PathFinder{pos: ORIGIN}
	fromNode := nm.Get(from)
	fromNode.open = true
	fromNode.cost = 0
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
		if current.finder.button != nil {
			current.finder.button.pressed++
			if current.finder.button.pressed > a.MaxPresses {
				break
			}
		}

		if current.finder.pos == a.Prize {
			var path Path
			for current != nil {
				path = append(Path{current}, path...)
				current = current.prev
			}
			return path
		}

		// Get the neighbors of the current node
		for _, neighbor := range current.finder.Neighbors(a) {
			cost := current.cost + neighbor.button.Tokens
			neighborNode := nm.Get(neighbor)
			// If we are closer to the goal...
			if cost < neighborNode.cost {
				if neighborNode.open {
					heap.Remove(nq, neighborNode.index)
					neighborNode.finder.button.pressed--
				}
				neighborNode.open = false
				neighborNode.closed = false
			}
			if !neighborNode.open && !neighborNode.closed {
				neighborNode.open = true
				neighborNode.rank = cost + neighbor.EstimatedCost(a)
				neighborNode.cost = cost
				neighborNode.prev = current
				heap.Push(nq, neighborNode)
			}
		}
	}
	return nil
}
