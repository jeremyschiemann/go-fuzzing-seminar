package main

import (
	oraclert "gfuzz/pkg/oraclert"
	"testing"
)

func TestSumSquares(t *testing.T) {
	oracleEntry := oraclert.BeforeRun()
	defer oraclert.AfterRun(oracleEntry)
	tests := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{"Empty slice", []int{}, 0},
		{"Single number", []int{4}, 16},
		{"Multiple numbers", []int{1, 2, 3, 4}, 30},
		{"Negative numbers", []int{-1, -2, -3}, 14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SumSquares(tt.numbers)
			if result != tt.expected {
				t.Errorf("got %d, want %d", result, tt.expected)
			}
		})
	}
}
