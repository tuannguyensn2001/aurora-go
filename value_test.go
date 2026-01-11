package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResolvedValue(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		matched  bool
		expected *resolvedValue
	}{
		{"matched boolean", true, true, &resolvedValue{value: true, matched: true}},
		{"matched string", "hello", true, &resolvedValue{value: "hello", matched: true}},
		{"matched int", 42, true, &resolvedValue{value: 42, matched: true}},
		{"matched float", 3.14, true, &resolvedValue{value: 3.14, matched: true}},
		{"unmatched nil", nil, false, &resolvedValue{value: nil, matched: false}},
		{"matched nil value", nil, true, &resolvedValue{value: nil, matched: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewResolvedValue(tt.value, tt.matched)
			assert.Equal(t, tt.expected.value, result.value)
			assert.Equal(t, tt.expected.matched, result.matched)
		})
	}
}

func TestResolvedValueBoolean(t *testing.T) {
	tests := []struct {
		name       string
		value      any
		matched    bool
		defaultVal bool
		expected   bool
	}{
		{"matched true", true, true, false, true},
		{"matched false", false, true, true, false},
		{"unmatched uses default", true, false, true, true},
		{"matched nil uses default", nil, true, true, true},
		{"wrong type uses default", "true", true, false, false},
		{"int 1 uses default", 1, true, false, false},
		{"int 0 uses default", 0, true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rv := NewResolvedValue(tt.value, tt.matched)
			result := rv.Boolean(tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResolvedValueString(t *testing.T) {
	tests := []struct {
		name       string
		value      any
		matched    bool
		defaultVal string
		expected   string
	}{
		{"matched string", "hello", true, "default", "hello"},
		{"matched empty string", "", true, "default", ""},
		{"unmatched uses default", "hello", false, "default", "default"},
		{"matched nil uses default", nil, true, "default", "default"},
		{"wrong type uses default", 123, true, "default", "default"},
		{"matched bool uses default", true, true, "default", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rv := NewResolvedValue(tt.value, tt.matched)
			result := rv.String(tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResolvedValueInt(t *testing.T) {
	tests := []struct {
		name       string
		value      any
		matched    bool
		defaultVal int
		expected   int
	}{
		{"matched int", 42, true, 0, 42},
		{"matched zero", 0, true, 99, 0},
		{"unmatched uses default", 42, false, 99, 99},
		{"matched nil uses default", nil, true, 99, 99},
		{"wrong type uses default", "42", true, 99, 99},
		{"matched bool uses default", true, true, 99, 99},
		{"matched float uses default", 42.5, true, 99, 99},
		{"matched int8 uses default", int8(10), true, 0, 0},
		{"matched int16 uses default", int16(20), true, 0, 0},
		{"matched int32 uses default", int32(30), true, 0, 0},
		{"matched int64 uses default", int64(40), true, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rv := NewResolvedValue(tt.value, tt.matched)
			result := rv.Int(tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResolvedValueFloat(t *testing.T) {
	tests := []struct {
		name       string
		value      any
		matched    bool
		defaultVal float64
		expected   float64
	}{
		{"matched float", 3.14, true, 0.0, 3.14},
		{"matched zero", 0.0, true, 99.0, 0.0},
		{"unmatched uses default", 3.14, false, 99.0, 99.0},
		{"matched nil uses default", nil, true, 99.0, 99.0},
		{"wrong type uses default", "3.14", true, 99.0, 99.0},
		{"matched bool uses default", true, true, 99.0, 99.0},
		{"matched int uses default", 42, true, 99.0, 99.0},
		{"matched float32 uses default", float32(2.5), true, 0.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rv := NewResolvedValue(tt.value, tt.matched)
			result := rv.Float(tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResolvedValueTypeSafety(t *testing.T) {
	t.Run("matched false with default false", func(t *testing.T) {
		rv := NewResolvedValue(false, true)
		assert.False(t, rv.Boolean(true))
	})

	t.Run("matched empty string with default", func(t *testing.T) {
		rv := NewResolvedValue("", true)
		assert.Equal(t, "", rv.String("default"))
	})

	t.Run("matched zero int with default", func(t *testing.T) {
		rv := NewResolvedValue(0, true)
		assert.Equal(t, 0, rv.Int(99))
	})

	t.Run("matched zero float with default", func(t *testing.T) {
		rv := NewResolvedValue(0.0, true)
		assert.Equal(t, 0.0, rv.Float(99.0))
	})

	t.Run("unmatched returns default regardless of value type", func(t *testing.T) {
		rv := NewResolvedValue("not an int", false)
		assert.Equal(t, 42, rv.Int(42))

		rv = NewResolvedValue(true, false)
		assert.Equal(t, 3.14, rv.Float(3.14))
	})
}
