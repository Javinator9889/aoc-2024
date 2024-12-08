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
	pather Pather
	cost   int
	index  int
	rank   int
	prev   *Node
	open   bool
	closed bool
}

type nodeMap map[Pather]*Node

func (nm nodeMap) Get(p Pather) *Node {
	res, ok := nm[p]
	if !ok {
		res = &Node{
			pather: p,
			rank:   -1,
		}
		nm[p] = res
	}
	return res
}

type Pather struct {
	pos       int
	value     int
	operation Op
	origin    *Pather
}

type Grid struct {
	goal        int
	numbers     []int
	validOps    []Op
	IsValidCost func(int) bool
}
type Path []*Node

func (p Path) String() string {
	var sb strings.Builder
	for _, n := range p {
		op := ""
		if n.pather.operation != "" {
			op = " " + n.pather.operation.String() + " "
		}
		sb.WriteString(fmt.Sprintf("%s%v", op, n.pather.value))
	}
	return sb.String()
}

func (p Pather) Neighbors(g *Grid) []Pather {
	neighbors := make([]Pather, 0)
	// Check if we are at the end of the numbers
	if p.pos == len(g.numbers)-1 {
		return neighbors
	}
	for _, op := range g.validOps {
		neighbor := Pather{
			value:     g.numbers[p.pos+1],
			operation: op,
			pos:       p.pos + 1,
			origin:    &p,
		}
		neighbors = append(neighbors, neighbor)
	}
	return neighbors
}

func (p Pather) EstimatedCost(g *Grid) int {
	return -p.pos
}

func (g *Grid) AStar() Path {
	nm := nodeMap{}
	nq := &priorityQueue{}
	heap.Init(nq)
	// We consider the `cost` the cumulative result of the operations
	from := Pather{value: g.numbers[0]}
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
		if current.cost == g.goal && current.pather.pos == len(g.numbers)-1 {
			var path Path
			for current != nil {
				path = append(Path{current}, path...)
				current = current.prev
			}
			return path
		}

		// Get the neighbors of the current node
		for _, neighbor := range current.pather.Neighbors(g) {
			cost := neighbor.operation.Cal(current.cost, neighbor.value)
			if !g.IsValidCost(cost) {
				// Skip if the rank is invalid
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

type Row struct {
	value   int
	numbers []int
}

func part1(input string) (solvable int) {
	// We can consider this exercise as a graph problem, where each node is a number with
	// a set of valid operations to apply to. The goal is to reach a certain number with
	// the minimum number of operations.
	parsed := parseInput(input)
	for _, row := range parsed {
		grid := Grid{
			goal:     row.value,
			numbers:  row.numbers,
			validOps: []Op{ADD, MUL},
			IsValidCost: func(cost int) bool {
				return cost <= row.value
			},
		}
		path := grid.AStar()
		// We have to use all the numbers
		if path != nil && len(path) == len(row.numbers) {
			slog.Debug("Path for", "n", row.value, "path", path)
			solvable += row.value
		} else if path != nil {
			slog.Warn("Not using all numbers", "n", row.value, "path", path, "nums", row.numbers)
		}
	}
	return
}

func part2(input string) int {
	return 0
}

func parseInput(input string) (ans []Row) {
	ans = make([]Row, 0)
	for _, line := range strings.Split(input, "\n") {
		items := strings.Split(line, ": ")
		k, v := cast.ToInt(items[0]), strings.Split(items[1], " ")
		row := Row{
			value:   k,
			numbers: make([]int, len(v)),
		}
		for i := range v {
			row.numbers[i] = cast.ToInt(v[i])
		}
		ans = append(ans, row)
	}
	return
}
