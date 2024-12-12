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

type Garden [][]*Flower

func (g Garden) outOfBounds(x, y int) bool {
	return x < 0 || x >= len(g) || y < 0 || y >= len((g)[0])
}

var (
	UP    = []int{0, -1}
	DOWN  = []int{0, 1}
	LEFT  = []int{-1, 0}
	RIGHT = []int{1, 0}
)

type Flower struct {
	x, y      int
	i         string
	clustered bool
}

func (f *Flower) getCluster(garden Garden, cluster *Cluster) {
	// A cluster is a set of flowers whose value is the same and are connected. The area
	// is the number of flowers in the cluster, and the perimeter is a virtual fence around
	// the cluster.
	f.clustered = true
	cluster.i = f.i
	cluster.flowers = append(cluster.flowers, f)
	cluster.area++
	prev := f
	for {
		// check if the cluster is connected to another flower
		// if it is, add that flower to the cluster and continue
		// if it is not, break
		// if all flowers are in the cluster, break
		end := true
		for _, dir := range [][]int{UP, DOWN, LEFT, RIGHT} {
			if garden.outOfBounds(prev.x+dir[0], prev.y+dir[1]) {
				continue
			}
			next := garden[prev.x+dir[0]][prev.y+dir[1]]
			if next.i == f.i && !next.clustered {
				next.getCluster(garden, cluster)
			}
		}
		if end {
			break
		}
	}
}

func (c *Cluster) calcPerimeter(garden Garden) {
	// calculate the perimeter
	for _, flower := range c.flowers {
		for _, dir := range [][]int{UP, DOWN, LEFT, RIGHT} {
			if garden.outOfBounds(flower.x+dir[0], flower.y+dir[1]) {
				c.perimeter++
			} else if garden[flower.x+dir[0]][flower.y+dir[1]].i != c.i {
				c.perimeter++
			}
		}
	}
}

// Gets how many sides the cluster has. A side is a number of flowers next to each other
// in a row or column
func (c *Cluster) calcSides(garden Garden) {
	sides := map[string]struct{}{}
	for _, flower := range c.flowers {
		for _, dir := range [][]int{UP, DOWN, LEFT, RIGHT} {
			// We determine the side by where we're coming from
			side := fmt.Sprintf("%v%v%v%v", dir[0], dir[1], flower.x*dir[0], flower.y*dir[1])
			if garden.outOfBounds(flower.x+dir[0], flower.y+dir[1]) {
				// oobSide := fmt.Sprintf("%v%v", dir[0], dir[1])
				sides[side] = struct{}{}
				continue
			}
			next := garden[flower.x+dir[0]][flower.y+dir[1]]
			if next.i != c.i {
				sides[side] = struct{}{}
			}
		}
	}
	slog.Debug("Sides", "sides", sides)
	c.sides = len(sides)
}

func (f *Flower) String() string {
	return f.i
}

type Cluster struct {
	area      int
	perimeter int
	flowers   []*Flower
	i         string
	sides     int
}

func (c *Cluster) farUpLeft() (x, y int) {
	x, y = c.flowers[0].x, c.flowers[0].y
	for _, flower := range c.flowers {
		if flower.x < x {
			x = flower.x
		}
		if flower.y < y {
			y = flower.y
		}
	}
	return
}

func (c *Cluster) farDownLeft() (x, y int) {
	x, y = c.flowers[0].x, c.flowers[0].y
	for _, flower := range c.flowers {
		if flower.x > x {
			x = flower.x
		}
		if flower.y < y {
			y = flower.y
		}
	}
	return
}

func (c *Cluster) farUpRight() (x, y int) {
	x, y = c.flowers[0].x, c.flowers[0].y
	for _, flower := range c.flowers {
		if flower.x < x {
			x = flower.x
		}
		if flower.y > y {
			y = flower.y
		}
	}
	return
}

func (c *Cluster) farDownRight() (x, y int) {
	x, y = c.flowers[0].x, c.flowers[0].y
	for _, flower := range c.flowers {
		if flower.x > x {
			x = flower.x
		}
		if flower.y > y {
			y = flower.y
		}
	}
	return
}

// Checks if a cluster is big enough to contain another cluster inside it
func (c *Cluster) contains(o *Cluster) bool {
	tlx, tly := c.farUpLeft()
	dlx, dly := c.farDownLeft()
	trx, try := c.farUpRight()
	drx, dry := c.farDownRight()
	tlxo, tlyo := o.farUpLeft()
	dlxo, dlyo := o.farDownLeft()
	trxo, tryo := o.farUpRight()
	drxo, dryo := o.farDownRight()
	return tlxo >= tlx &&
		tlyo >= tly &&
		dlxo >= dlx &&
		dlyo <= dly &&
		trxo <= trx &&
		tryo >= try &&
		drxo <= drx &&
		dryo <= dry
}

func (c Cluster) String() string {
	str := strings.Builder{}
	str.WriteString("Cluster: {")
	str.WriteString(fmt.Sprintf("area: %v, perimeter: %v, sides: %v, flowers: [", c.area, c.perimeter, c.sides))
	for i, flower := range c.flowers {
		str.WriteString(flower.String())
		if i < len(c.flowers)-1 {
			str.WriteString(", ")
		}
	}
	str.WriteString("]}")
	return str.String()
}

func part1(input string) (cost int) {
	garden := parseInput(input)
	slog.Debug("Garden", "garden", garden)
	clusters := make([]Cluster, 0)
	for _, row := range garden {
		for _, flower := range row {
			if !flower.clustered {
				cluster := Cluster{}
				flower.getCluster(garden, &cluster)
				cluster.calcPerimeter(garden)
				slog.Debug("Flower is clustered", "flower", flower, "cluster", cluster)
				clusters = append(clusters, cluster)
			}
		}
	}
	slog.Debug("Clusters", "clusters", clusters)
	for i, cluster := range clusters {
		for j, other := range clusters {
			if i == j {
				continue
			}
			if cluster.contains(&other) {
				slog.Debug("Cluster contains another", "cluster", cluster, "other", other)
				cluster.perimeter += other.perimeter
			}
		}
	}
	slog.Debug("Clusters", "clusters", clusters)
	// Calculate the total cost. The cost is obtained by multiplying the area with the perimeter
	for _, cluster := range clusters {
		slog.Debug("Cluster", "cluster", cluster, "cost", cluster.area*cluster.perimeter)
		cost += cluster.area * cluster.perimeter
	}

	return
}

func part2(input string) (cost int) {
	garden := parseInput(input)
	slog.Debug("Garden", "garden", garden)
	clusters := make([]Cluster, 0)
	for _, row := range garden {
		for _, flower := range row {
			if !flower.clustered {
				cluster := Cluster{}
				flower.getCluster(garden, &cluster)
				cluster.calcSides(garden)
				slog.Debug("Flower is clustered", "flower", flower, "cluster", cluster)
				clusters = append(clusters, cluster)
			}
		}
	}
	slog.Debug("Clusters", "clusters", clusters)
	for i, cluster := range clusters {
		for j, other := range clusters {
			if i == j {
				continue
			}
			if cluster.contains(&other) {
				slog.Debug("Cluster contains another", "cluster", cluster, "other", other)
				cluster.sides += other.sides
			}
		}
	}
	slog.Debug("Clusters", "clusters", clusters)
	// Calculate the total cost. The cost is obtained by multiplying the area with the perimeter
	for _, cluster := range clusters {
		slog.Debug("Cluster", "cluster", cluster, "cost", cluster.area*cluster.sides)
		cost += cluster.area * cluster.sides
	}

	return
}

func parseInput(input string) (garden Garden) {
	garden = make(Garden, 0)
	for i, line := range strings.Split(input, "\n") {
		garden = append(garden, make([]*Flower, len(line)))
		for j, c := range line {
			garden[i][j] = &Flower{x: i, y: j, i: string(c)}
		}
	}
	return
}
