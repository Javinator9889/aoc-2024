package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"slices"
	"sort"
	"strings"

	"github.com/Javinator9889/aoc-2024/cast"
	"github.com/Javinator9889/aoc-2024/util"
)

type Rule struct {
	before []int
	after  []int
}

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

func part1(input string) (centerSum int) {
	rules, pages := parseInput(input)
	for i := range pages {
		slog.Debug("Page", "i", i, "pages", pages[i])
		// For each page, check if the rules are satisfied
		// If they are, add the page to the center sum
		valid := true
		for j := range pages[i] {
			rule := rules[pages[i][j]]
			before := pages[i][:j]
			after := pages[i][j+1:]

			// With the intersection, verify if any element that should go after is not before
			// and vice versa
			if len(intersection(rule.after, after)) != 0 ||
				len(intersection(rule.before, before)) != 0 {
				slog.Debug(
					"Intersection",
					"after", intersection(rule.after, after),
					"before", intersection(rule.before, before),
					"rule", rule,
				)
				valid = false
				break
			}
		}
		if valid {
			centerSum += pages[i][len(pages[i])/2]
		}
	}

	return
}

func With(rules map[int]*Rule) Sorter {
	w := func(i, j int) bool {
		return slices.Contains(rules[i].before, j)
	}
	return Sorter{with: w, changed: false}
}

type Sorter struct {
	with    func(i, j int) bool
	page    []int
	changed bool
}

func (s *Sorter) Len() int {
	return len(s.page)
}

func (s *Sorter) Swap(i, j int) {
	s.changed = true
	s.page[i], s.page[j] = s.page[j], s.page[i]
}

func (s *Sorter) Less(i, j int) bool {
	return s.with(s.page[i], s.page[j])
}

func (s *Sorter) Sort(page []int) {
	s.page = page
	sort.Sort(s)
}

func (s *Sorter) Changed() bool {
	return s.changed
}

func part2(input string) (centerSum int) {
	rules, pages := parseInput(input)
	// For this part, the algorithm is the same but we have to consider the incorrectly-ordered
	// pages and sort them according to the rules.
	for i := range pages {
		slog.Debug("Page", "i", i, "pages", pages[i])
		sorter := With(rules)
		sorter.Sort(pages[i])
		slog.Debug("Sorted", "i", i, "pages", pages[i])
		if sorter.Changed() {
			centerSum += pages[i][len(pages[i])/2]
		}
	}
	return
}

func parseInput(input string) (rules map[int]*Rule, pages [][]int) {
	rules = make(map[int]*Rule)
	pages = make([][]int, 0)
	rulesDone := false

	// The input is divided in two parts: The set of rules and the set of pages
	// Rules are in the form "X|Y" where X is the page that must be before Y,
	// and pages are in the form "X,Y,Z" where X, Y and Z are the pages that are to be be
	// printed in that order.
	// The rules are stored in a map of int to Rule, where Rule is a struct with two slices of int
	// that represent the pages that must be before and after the page that is the key of the map.
	// The pages are stored in a slice of slices of int, where each slice of int
	// represents the pages that are to be printed in that order.
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			rulesDone = true
			continue
		}
		if !rulesDone {
			parts := strings.Split(line, "|")
			page1 := cast.ToInt(parts[0])
			page2 := cast.ToInt(parts[1])
			if _, ok := rules[page1]; !ok {
				rules[page1] = &Rule{before: []int{}, after: []int{}}
			}
			if _, ok := rules[page2]; !ok {
				rules[page2] = &Rule{before: []int{}, after: []int{}}
			}
			rules[page1].before = append(rules[page1].before, page2)
			rules[page2].after = append(rules[page2].after, page1)
		} else {
			var current []int
			for _, val := range strings.Split(line, ",") {
				current = append(current, cast.ToInt(val))
			}
			pages = append(pages, current)
		}
	}
	return rules, pages
}
