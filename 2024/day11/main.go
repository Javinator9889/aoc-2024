package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Javinator9889/aoc-2024/2024/day11/cache"
	"github.com/Javinator9889/aoc-2024/cast"
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

type Stone struct {
	value int
	next  *Stone
}

func (s Stone) Len() int {
	return len(cast.ToString(s.value))
}

func (s *Stone) Blink() (next *Stone) {
	switch {
	case s.value == 0:
		s.value = 1
		next = s.next
	case s.Len()%2 == 0:
		sv := cast.ToString(s.value)
		first, second := sv[:len(sv)/2], sv[len(sv)/2:]
		s.value = cast.ToInt(first)
		newStone := &Stone{value: cast.ToInt(second)}
		s.next, newStone.next = newStone, s.next
		next = newStone.next // Skip the new stone
	default:
		s.value *= 2024
		next = s.next
	}
	return
}

func (s *Stone) Size() (size int) {
	for stone := s; stone != nil; stone = stone.next {
		size++
	}
	return
}

func (s *Stone) String() string {
	var sb strings.Builder
	for stone := s; stone != nil; stone = stone.next {
		sb.WriteString(fmt.Sprintf("%d ", stone.value))
	}
	return sb.String()
}

func digits(n int) (count int) {
	for n > 0 {
		n /= 10
		count++
	}
	return
}

func (s *Stone) Stones(n int) int {
	var stonesFn cache.Int2Func
	stonesFn = cache.Cached(func(value, n int) (count int) {
		prev := value
		for i := 0; i < n; i++ {
			switch {
			case prev == 0:
				prev = 1
			case digits(prev)%2 == 0:
				sv := cast.ToString(prev)
				first, second := sv[:len(sv)/2], sv[len(sv)/2:]
				prev = cast.ToInt(first)
				if n-i > 0 {
					count += stonesFn(cast.ToInt(second), n-i-1)
				}
			default:
				prev *= 2024
			}
		}
		count++ // Add the current stone
		return
	})
	return stonesFn(s.value, n)
}

func blink(times int, ref *Stone) int {
	for i := 0; i < times; i++ {
		start := time.Now()
		for st := ref; st != nil; st = st.Blink() {
		}
		slog.Debug("Blink", "i", i, "elapsed", time.Since(start), "size", ref.Size(), "stones", ref)
	}
	return ref.Size()
}

func part1(input string) int {
	stones := parseInput(input)
	slog.Debug("Stones:", "stones", stones)
	return blink(25, stones)
}

func part2(input string) (count int) {
	stones := parseInput(input)
	slog.Debug("Stones:", "stones", stones)
	start := time.Now()
	for st := stones; st != nil; st = st.next {
		istart := time.Now()
		count += st.Stones(75)
		slog.Debug("Stones", "count", count, "elapsed", time.Since(istart))
	}
	slog.Debug("Elapsed", "elapsed", time.Since(start))
	return
}

func parseInput(input string) (stones *Stone) {
	var prev *Stone
	for _, line := range strings.Split(input, "\n") {
		for _, num := range strings.Fields(line) {
			if prev == nil {
				prev = &Stone{value: cast.ToInt(num)}
				stones = prev
				continue
			}
			newStone := &Stone{value: cast.ToInt(num)}
			prev.next = newStone
			prev = newStone
		}
	}
	return
}
