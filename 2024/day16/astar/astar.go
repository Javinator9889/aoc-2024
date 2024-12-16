package astar

import (
	"container/heap"
	"strings"
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

var (
	NORTH = Location{X: -1, Y: 0}
	SOUTH = Location{X: 1, Y: 0}
	WEST  = Location{X: 0, Y: -1}
	EAST  = Location{X: 0, Y: 1}
)

func (l Location) Add(other Location) Location {
	return Location{
		X: l.X + other.X,
		Y: l.Y + other.Y,
	}
}

type Grid [][]rune

func (g Grid) String() string {
	var sb strings.Builder
	for _, row := range g {
		for _, cell := range row {
			sb.WriteRune(cell)
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

type Reindlympics struct {
	End   Location
	Start Location
	Grid  Grid
}

// A Path represents a sequence of operations to reach the goal
type Path []*Node

func (p Path) String() string {
	var sb strings.Builder
	for _, node := range p {
		sb.WriteString(node.finder.String())
		sb.WriteRune(' ')
	}
	return sb.String()
}

func (p Path) Cost() int {
	return p[len(p)-1].cost
}

func (p Path) Uniq() Path {
	seen := make(map[PathFinder]bool)
	var res Path
	for _, node := range p {
		if _, ok := seen[node.finder]; !ok {
			seen[node.finder] = true
			res = append(res, node)
		}
	}
	return res
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

func (r *Reindlympics) AStar() Path {
	nm := nodeMap{}
	nq := &PriorityQueue{}
	heap.Init(nq)
	// The reindeer starts at the start position, facing east
	from := PathFinder{pos: r.Start, dir: EAST}
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

		if current.finder.pos == r.End {
			var path Path
			for current != nil {
				path = append(Path{current}, path...)
				current = current.prev
			}
			return path
		}

		// Get the neighbors of the current node
		for _, neighbor := range current.finder.Neighbors(r) {
			cost := current.cost + neighbor.moveCost
			neighborNode := nm.Get(neighbor)
			// If we are closer to the goal...
			if cost < neighborNode.cost {
				if neighborNode.open {
					heap.Remove(nq, neighborNode.index)
				}
				neighborNode.open = false
				neighborNode.closed = false
			}
			if !neighborNode.open && !neighborNode.closed {
				neighborNode.open = true
				neighborNode.rank = cost + neighbor.EstimatedCost(r)
				neighborNode.cost = cost
				neighborNode.prev = current
				heap.Push(nq, neighborNode)
			}
		}
	}
	return nil
}

func (r *Reindlympics) AStarRecursive(target int, actual int, initialDir Location) Path {
	nm := nodeMap{}
	nq := &PriorityQueue{}
	heap.Init(nq)
	// The reindeer starts at the start position, facing east
	from := PathFinder{pos: r.Start, dir: initialDir}
	fromNode := nm.Get(from)
	fromNode.open = true
	fromNode.cost = 0
	fromNode.prev = nil
	heap.Push(nq, fromNode)
	var this Path

	for {
		// There are no more nodes to explore
		if nq.Len() == 0 {
			break
		}
		current := heap.Pop(nq).(*Node)
		current.open = false
		current.closed = true

		if current.finder.pos == r.End {
			var tmp Path
			count := 0
			for current != nil {
				tmp = append(Path{current}, tmp...)
				current = current.prev
				count++
			}
			if count < target {
				this = append(this, tmp...)
				return this
			}
		}
		if current.finder.pos == r.Start {
			return nil
		}

		// Get the neighbors of the current node
		for _, neighbor := range current.finder.Neighbors(r) {

			if current.finder.dir != neighbor.dir {
				cTmp := current
				count := 0
				for cTmp != nil {
					count++
					cTmp = cTmp.prev
				}
				rTmp := *r
				rTmp.Start = neighbor.pos

				res := rTmp.AStarRecursive(target, count, neighbor.dir)
				if res != nil {
					this = append(this, res...)
				}
				continue
			}

			cost := current.cost + neighbor.moveCost
			neighborNode := nm.Get(neighbor)
			// If we are closer to the goal...
			if cost < neighborNode.cost {
				if neighborNode.open {
					heap.Remove(nq, neighborNode.index)
				}
				neighborNode.open = false
				neighborNode.closed = false
			}
			if !neighborNode.open && !neighborNode.closed {
				neighborNode.open = true
				neighborNode.rank = cost + neighbor.EstimatedCost(r)
				neighborNode.cost = cost
				neighborNode.prev = current
				heap.Push(nq, neighborNode)
			}
		}
	}
	return nil
}
