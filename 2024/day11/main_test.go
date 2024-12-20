package main

import (
	"log/slog"
	"testing"
)

var example = `125 17`

func Test_part1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "example",
			input: example,
			want:  55312,
		},
		{
			name:  "actual",
			input: input,
			want:  186175,
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
			name:  "actual",
			input: input,
			want:  220566831337810,
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

func recursiveImpl(stones *Stone, n int) int {
	return blink(n, stones)
}

func cachedImpl(stones *Stone, n int) (count int) {
	for st := stones; st != nil; st = st.next {
		count += st.Stones(n)
	}
	return
}

func BenchmarkRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stones := parseInput(input)
		recursiveImpl(stones, 25)
	}
}

func BenchmarkCached(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stones := parseInput(input)
		cachedImpl(stones, 25)
	}
}
