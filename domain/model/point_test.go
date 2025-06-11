package model

import (
	"testing"
)

func TestPoint_Add(t *testing.T) {
	tests := []struct {
		name     string
		point    Point
		other    Point
		expected Point
	}{
		{
			name:     "正の値同士の足し算",
			point:    Point{X: 1, Y: 2},
			other:    Point{X: 3, Y: 4},
			expected: Point{X: 4, Y: 6},
		},
		{
			name:     "負の値を含む足し算",
			point:    Point{X: 5, Y: 3},
			other:    Point{X: -2, Y: -1},
			expected: Point{X: 3, Y: 2},
		},
		{
			name:     "ゼロとの足し算",
			point:    Point{X: 7, Y: 8},
			other:    Point{X: 0, Y: 0},
			expected: Point{X: 7, Y: 8},
		},
		{
			name:     "負の結果となる足し算",
			point:    Point{X: 2, Y: 1},
			other:    Point{X: -5, Y: -3},
			expected: Point{X: -3, Y: -2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.point.Add(tt.other)
			if result != tt.expected {
				t.Errorf("Point.Add() = %v, want %v", result, tt.expected)
			}
		})
	}
}
