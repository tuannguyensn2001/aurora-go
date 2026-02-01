package core

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tuannguyensn2001/aurora-go/core/evaluator"
)

func TestEqualOperator(t *testing.T) {
	tests := []struct {
		name     string
		a, b     any
		expected bool
	}{
		{"nil both", nil, nil, true},
		{"nil a", nil, "value", false},
		{"nil b", "value", nil, false},
		{"int equal", 5, 5, true},
		{"int not equal", 5, 10, false},
		{"int vs float equal", 5, 5.0, true},
		{"int vs float not equal", 5, 5.5, false},
		{"float32 equal", float32(5.5), float32(5.5), true},
		{"float64 equal", 5.5, 5.5, true},
		{"uint equal", uint(10), uint(10), true},
		{"string equal", "hello", "hello", true},
		{"string not equal", "hello", "world", false},
		{"bool true", true, true, true},
		{"bool false", false, false, true},
		{"bool mixed", true, false, false},
		{"slice equal", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"slice not equal", []int{1, 2, 3}, []int{1, 2, 4}, false},
		{"map equal", map[string]int{"a": 1}, map[string]int{"a": 1}, true},
		{"map not equal", map[string]int{"a": 1}, map[string]int{"a": 2}, false},
		{"int types mixed", int8(5), int64(5), true},
		{"uint types mixed", uint32(10), uint64(10), true},
		{"negative int", -5, -5, true},
		{"zero", 0, 0, true},
		{"empty string", "", "", true},
		{"string vs int", "5", 5, false},
		{"bool vs int", true, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluator.EqualOp(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNotEqualOperator(t *testing.T) {
	tests := []struct {
		name     string
		a, b     any
		expected bool
	}{
		{"nil both", nil, nil, false},
		{"equal values", 5, 5, false},
		{"not equal", 5, 10, true},
		{"string equal", "hello", "hello", false},
		{"string not equal", "hello", "world", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluator.NotEqualOp(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGreaterThanOperator(t *testing.T) {
	tests := []struct {
		name     string
		a, b     any
		expected bool
	}{
		{"nil a", nil, 5, false},
		{"nil b", 5, nil, false},
		{"int greater", 10, 5, true},
		{"int less", 5, 10, false},
		{"int equal", 5, 5, false},
		{"float greater", 5.5, 5.0, true},
		{"float less", 5.0, 5.5, false},
		{"string greater", "z", "a", true},
		{"string less", "a", "z", false},
		{"string equal", "abc", "abc", false},
		{"int vs float greater", 10, 5.5, true},
		{"negative numbers", -5, -10, true},
		{"zero vs positive", 0, 1, false},
		{"mixed int types", int64(100), int32(50), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluator.GreaterThanOp(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLessThanOperator(t *testing.T) {
	tests := []struct {
		name     string
		a, b     any
		expected bool
	}{
		{"nil a", nil, 5, false},
		{"nil b", 5, nil, false},
		{"int less", 5, 10, true},
		{"int greater", 10, 5, false},
		{"int equal", 5, 5, false},
		{"float less", 5.0, 5.5, true},
		{"float greater", 5.5, 5.0, false},
		{"string less", "a", "z", true},
		{"string greater", "z", "a", false},
		{"string equal", "abc", "abc", false},
		{"negative numbers", -10, -5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluator.LessThanOp(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGreaterThanOrEqualOperator(t *testing.T) {
	tests := []struct {
		name     string
		a, b     any
		expected bool
	}{
		{"equal", 5, 5, true},
		{"greater", 10, 5, true},
		{"less", 5, 10, false},
		{"float equal", 5.5, 5.5, true},
		{"float greater", 6.0, 5.5, true},
		{"string equal", "abc", "abc", true},
		{"string greater", "z", "a", true},
		{"string less", "a", "z", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluator.GreaterThanOrEqualOp(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLessThanOrEqualOperator(t *testing.T) {
	tests := []struct {
		name     string
		a, b     any
		expected bool
	}{
		{"equal", 5, 5, true},
		{"greater", 10, 5, false},
		{"less", 5, 10, true},
		{"float equal", 5.5, 5.5, true},
		{"float less", 5.0, 5.5, true},
		{"string equal", "abc", "abc", true},
		{"string less", "a", "z", true},
		{"string greater", "z", "a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluator.LessThanOrEqualOp(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContainsOperator(t *testing.T) {
	tests := []struct {
		name     string
		a, b     any
		expected bool
	}{
		{"nil a", nil, "value", false},
		{"nil b", "value", nil, false},
		{"string contains", "hello world", "world", true},
		{"string does not contain", "hello world", "xyz", false},
		{"empty substring", "hello", "", true},
		{"substring at start", "hello", "hel", true},
		{"substring at end", "hello", "llo", true},
		{"slice contains int", []int{1, 2, 3}, 2, true},
		{"slice does not contain int", []int{1, 2, 3}, 5, false},
		{"slice contains string", []string{"a", "b", "c"}, "b", true},
		{"empty slice", []int{}, 1, false},
		{"slice with duplicates", []int{1, 2, 2, 3}, 2, true},
		{"int vs int in slice", 5, 5, false},
		{"string vs int slice", "2", []int{1, 2, 3}, false},
		{"slice of slices", [][]int{{1, 2}, {3, 4}}, []int{1, 2}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluator.ContainsOp(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInOperator(t *testing.T) {
	tests := []struct {
		name     string
		a, b     any
		expected bool
	}{
		{"nil a", nil, []int{1, 2, 3}, false},
		{"nil b", 1, nil, false},
		{"int in slice", 2, []int{1, 2, 3}, true},
		{"int not in slice", 5, []int{1, 2, 3}, false},
		{"string in slice", "b", []string{"a", "b", "c"}, true},
		{"string not in slice", "x", []string{"a", "b", "c"}, false},
		{"empty slice", 1, []int{}, false},
		{"first element", 1, []int{1, 2, 3}, true},
		{"last element", 3, []int{1, 2, 3}, true},
		{"int in float slice", 5, []float64{1.0, 5.0, 10.0}, true},
		{"float in int slice", 5.0, []int{1, 5, 10}, true},
		{"not slice or array", 1, 5, false},
		{"slice with duplicates", 2, []int{1, 2, 2, 3}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluator.InOp(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNotInOperator(t *testing.T) {
	tests := []struct {
		name     string
		a, b     any
		expected bool
	}{
		{"int in slice", 2, []int{1, 2, 3}, false},
		{"int not in slice", 5, []int{1, 2, 3}, true},
		{"string in slice", "b", []string{"a", "b", "c"}, false},
		{"string not in slice", "x", []string{"a", "b", "c"}, true},
		{"empty slice", 1, []int{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluator.NotInOp(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUTF8StringContains(t *testing.T) {
	tests := []struct {
		name      string
		s, substr string
		expected  bool
	}{
		{"UTF-8 with accents", "café", "afé", true},
		{"UTF-8 emoji", "hello👋world", "👋", true},
		{"UTF-8 Chinese", "你好世界", "世界", true},
		{"UTF-8 mixed", "hello世界", "世界", true},
		{"UTF-8 not found", "café", "xyz", false},
		{"empty substring", "café", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.Contains(tt.s, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}
