package main

import (
	"log/slog"
	"testing"
)

var example = `p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3`

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "example",
			input: example,
			want:  12,
		},
		// {
		// 	name:  "actual",
		// 	input: input,
		// 	want:  0,
		// },
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
			name:  "example",
			input: example,
			want:  0,
		},
		// {
		// 	name:  "actual",
		// 	input: input,
		// 	want:  0,
		// },
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
