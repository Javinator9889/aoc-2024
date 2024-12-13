package main

import (
	"log/slog"
	"testing"
)

var example = `RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE`

var example2 = `OOOOO
OXOXO
OOOOO
OXOXO
OOOOO`

var simple = `AAAA
BBCD
BBCC
EEEC`

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "containing example",
			input: example2,
			want:  772,
		},
		{
			name:  "example",
			input: example,
			want:  1930,
		},
		{
			name:  "actual",
			input: input,
			want:  1371306,
		},
	}
	slog.SetLogLoggerLevel(slog.LevelDebug)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part1(tt.input); got != tt.want {
				t.Errorf("part1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_part2(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "simple example",
			input: simple,
			want:  80,
		},
		{
			name:  "containing example",
			input: example2,
			want:  436,
		},
		{
			name:  "example",
			input: example,
			want:  1206,
		},
		{
			name:  "actual",
			input: input,
			want:  805880,
		},
	}
	slog.SetLogLoggerLevel(slog.LevelDebug)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part2(tt.input); got != tt.want {
				t.Errorf("part2() = %v, want %v", got, tt.want)
			}
		})
	}
}
