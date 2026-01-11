package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAttribute(t *testing.T) {
	attr := NewAttribute()
	assert.NotNil(t, attr)
	assert.NotNil(t, attr.vals)
	assert.Len(t, attr.vals, 0)
}

func TestAttributeSet(t *testing.T) {
	attr := NewAttribute()

	attr.Set("key1", "value1")
	assert.Equal(t, "value1", attr.vals["key1"])

	attr.Set("key2", 123)
	assert.Equal(t, 123, attr.vals["key2"])

	attr.Set("key1", "new_value")
	assert.Equal(t, "new_value", attr.vals["key1"])
	assert.Len(t, attr.vals, 2)
}

func TestAttributeGet(t *testing.T) {
	attr := NewAttribute()
	attr.Set("string_val", "hello")
	attr.Set("int_val", 42)
	attr.Set("float_val", 3.14)
	attr.Set("bool_val", true)

	tests := []struct {
		name     string
		key      string
		expected any
	}{
		{"string value", "string_val", "hello"},
		{"int value", "int_val", 42},
		{"float value", "float_val", 3.14},
		{"bool value", "bool_val", true},
		{"missing key", "missing", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := attr.Get(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAttributeGetAfterSet(t *testing.T) {
	attr := NewAttribute()

	attr.Set("key", "initial")
	assert.Equal(t, "initial", attr.Get("key"))

	attr.Set("key", "updated")
	assert.Equal(t, "updated", attr.Get("key"))

	attr.Set("key", nil)
	assert.Nil(t, attr.Get("key"))
}

func TestAttributeWithComplexTypes(t *testing.T) {
	attr := NewAttribute()

	slice := []int{1, 2, 3}
	attr.Set("slice", slice)
	assert.Equal(t, slice, attr.Get("slice"))

	mapVal := map[string]int{"a": 1, "b": 2}
	attr.Set("map", mapVal)
	assert.Equal(t, mapVal, attr.Get("map"))

	structVal := struct{ Name string }{"test"}
	attr.Set("struct", structVal)
	assert.Equal(t, structVal, attr.Get("struct"))
}
