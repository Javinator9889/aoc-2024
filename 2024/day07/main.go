package main

import (
	"container/heap"
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Javinator9889/aoc-2024/cast"
	"github.com/Javinator9889/aoc-2024/util"
)

//go:embed input.txt
var input string

const (
	ADD = "add"
	SUB = "sub"
	MUL = "mul"
	DIV = "div"
)

type Op string

func (o Op) Cal(a, b int) int {
	switch o {
	case ADD:
		return a + b
	case SUB:
		return a - b
	case MUL:
		return a * b
	case DIV:
		return a / b
	}
	panic("invalid operation")
}

func (o Op) String() string {
	switch o {
	case ADD:
		return "+"
	case SUB:
		return "-"
	case MUL:
		return "*"
	case DIV:
		return "/"
	}
	panic("invalid operation")
}

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

// A priorityQueue implements heap.Interface and holds Nodes.  The
// priorityQueue is used to track open nodes by rank.
type priorityQueue []*Node

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].rank < pq[j].rank
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	no := x.(*Node)
	no.index = n
	*pq = append(*pq, no)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	no := old[n-1]
	no.index = -1
	*pq = old[0 : n-1]
	return no
}

type Node struct {
	pos       int
	value     int
	cost      int
	operation Op
	index     int
	rank      int
	prev      *Node
	open      bool
	closed    bool
	neighbors []*Node
}

type Grid struct {
	goal        int
	numbers     []int
	validOps    []Op
	IsValidRank func(int) bool
}
type Path []*Node

func (p Path) String() string {
	var sb strings.Builder
	for _, n := range p {
		op := ""
		if n.operation != "" {
			op = " " + n.operation.String() + " "
		}
		sb.WriteString(fmt.Sprintf("%s%v", op, n.value))
	}
	return sb.String()
}

func (n *Node) Neighbors(g *Grid) []*Node {
	if n.neighbors != nil {
		return n.neighbors
	}
	n.neighbors = make([]*Node, 0)
	// Check if we are at the end of the numbers
	if n.pos == len(g.numbers)-1 {
		return n.neighbors
	}
	for _, op := range g.validOps {
		cost := op.Cal(n.cost, g.numbers[n.pos+1])
		neighbor := &Node{
			value:     g.numbers[n.pos+1],
			cost:      cost,
			operation: op,
			prev:      n,
			open:      false,
			closed:    false,
			rank:      g.goal - cost,
			pos:       n.pos + 1,
		}
		if g.IsValidRank(neighbor.rank) {
			n.neighbors = append(n.neighbors, neighbor)
		}
	}
	return n.neighbors
}

func (g *Grid) AStar() Path {
	// var closedSet Path
	nq := &priorityQueue{}
	heap.Init(nq)
	// We consider the `cost` the cumulative result of the operations
	from := &Node{
		value:     g.numbers[0],
		cost:      g.numbers[0],
		rank:      0,
		prev:      nil,
		open:      true,
		pos:       0,
		neighbors: nil,
	}
	heap.Push(nq, from)

	for {
		// There are no more nodes to explore
		if nq.Len() == 0 {
			break
		}
		current := heap.Pop(nq).(*Node)
		current.open = false
		current.closed = true

		// We reached the goal. Verify that we used all the numbers
		if current.rank == 0 && current.pos == len(g.numbers)-1 {
			var path Path
			for current != nil {
				path = append(Path{current}, path...)
				current = current.prev
			}
			return path
		}

		// Get the neighbors of the current node
		for _, neighbor := range current.Neighbors(g) {
			if !neighbor.open && !neighbor.closed {
				neighbor.open = true
				heap.Push(nq, neighbor)
			} else if neighbor.open {
				// If the neighbor is already open, we check if the rank is worse
				if current.rank < neighbor.rank {
					if neighbor.open {
						heap.Remove(nq, neighbor.index)
					}
					neighbor.open = false
					neighbor.closed = false
				}
			}
		}
	}
	return nil
}

func part1(input string) (solvable int) {
	// We can consider this exercise as a graph problem, where each node is a number with
	// a set of valid operations to apply to. The goal is to reach a certain number with
	// the minimum number of operations.
	parsed := parseInput(input)
	for n, nums := range parsed {
		grid := Grid{
			goal:     n,
			numbers:  nums,
			validOps: []Op{ADD, MUL},
			IsValidRank: func(rank int) bool {
				return rank >= 0
			},
		}
		path := grid.AStar()
		// We have to use all the numbers
		if path != nil && len(path) == len(nums) {
			slog.Debug("Path for", "n", n, "path", path)
			solvable += n
		} else if path != nil {
			slog.Warn("Not using all numbers", "n", n, "path", path, "nums", nums)
		}
	}
	return
}

func part2(input string) int {
	return 0
}

func parseInput(input string) (ans map[int][]int) {
	ans = make(map[int][]int)
	for _, line := range strings.Split(input, "\n") {
		items := strings.Split(line, ": ")
		k, v := cast.ToInt(items[0]), strings.Split(items[1], " ")
		ans[k] = make([]int, len(v))
		for i := range v {
			ans[k][i] = cast.ToInt(v[i])
		}
	}
	return
}
