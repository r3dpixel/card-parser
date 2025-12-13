package property

import (
	"testing"

	"github.com/r3dpixel/toolkit/ptr"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/stringsx"
	"github.com/stretchr/testify/assert"
)

func TestBool_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Boolean values
		{name: "Bool true", input: "true", expected: true},
		{name: "Bool false", input: "false", expected: false},

		// Float values (cast library: 0 = false, non-zero = true)
		{name: "Float zero", input: "0.0", expected: false},
		{name: "Float positive", input: "1.0", expected: true},
		{name: "Float negative", input: "-1.0", expected: true},
		{name: "Float decimal", input: "0.1", expected: true},

		// String values (cast library: "false", "0", "" = false, others = true)
		{name: "String true", input: `"true"`, expected: true},
		{name: "String false", input: `"false"`, expected: false},
		{name: "String 1", input: `"1"`, expected: true},
		{name: "String 0", input: `"0"`, expected: false},
		{name: "String empty", input: `""`, expected: false},
		{name: "String yes", input: `"yes"`, expected: false},
		{name: "String no", input: `"no"`, expected: false},
		{name: "String random", input: `"hello"`, expected: false},

		// Null value
		{name: "Null", input: "null", expected: false},

		// Array values (default to false)
		{name: "Array empty", input: "[]", expected: false},
		{name: "Array simple", input: `["a", "b", "c"]`, expected: false},
		{name: "Array mixed", input: `[1, "two", true, null]`, expected: false},

		// Object values (default to false)
		{name: "Object empty", input: "{}", expected: false},
		{name: "Object simple", input: `{"key": "value"}`, expected: false},
		{name: "Object complex", input: `{"num": 42, "str": "hello"}`, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Bool
			err := sonicx.Config.UnmarshalFromString(tt.input, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, bool(result))
		})
	}
}

func TestBool_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    Bool
		expected string
	}{
		{name: "Bool true", input: Bool(true), expected: "true"},
		{name: "Bool false", input: Bool(false), expected: "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sonicx.Config.Marshal(&tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestBool_RoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		original bool
	}{
		{name: "True value", original: true},
		{name: "False value", original: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := sonicx.Config.Marshal(&tt.original)
			assert.NoError(t, err)

			// Unmarshal back
			var result Bool
			err = sonicx.Config.UnmarshalFromString(stringsx.FromBytes(jsonData), &result)
			assert.NoError(t, err)

			// Should be equal
			assert.Equal(t, tt.original, bool(result))
		})
	}
}

func TestBool_TypeConversions(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{name: "int zero to bool", input: 0, expected: false},
		{name: "int positive to bool", input: 1, expected: true},
		{name: "int negative to bool", input: -1, expected: true},
		{name: "float zero to bool", input: 0.0, expected: false},
		{name: "float positive to bool", input: 1.5, expected: true},
		{name: "bool true to bool", input: true, expected: true},
		{name: "bool false to bool", input: false, expected: false},
		{name: "nil to bool", input: nil, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := sonicx.Config.Marshal(tt.input)
			assert.NoError(t, err)

			var result Bool
			err = sonicx.Config.UnmarshalFromString(stringsx.FromBytes(jsonData), &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, bool(result))
		})
	}
}

func TestBool_CastLibraryBehavior(t *testing.T) {
	// Test specific cast library behavior for edge cases
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// String cases that cast library handles specifically
		{name: "String f", input: `"f"`, expected: false},
		{name: "String F", input: `"F"`, expected: false},
		{name: "String FALSE", input: `"FALSE"`, expected: false},
		{name: "String True", input: `"True"`, expected: true},
		{name: "String TRUE", input: `"TRUE"`, expected: true},
		{name: "String t", input: `"t"`, expected: true},
		{name: "String T", input: `"T"`, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Bool
			err := sonicx.Config.UnmarshalFromString(tt.input, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, bool(result))
		})
	}
}

func TestBool_ErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "Malformed JSON", input: "{"},
		{name: "Invalid JSON", input: "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Bool
			err := sonicx.Config.UnmarshalFromString(tt.input, &result)
			assert.Error(t, err)
		})
	}
}

func TestBool_SetIfPtr(t *testing.T) {
	tests := []struct {
		name     string
		initial  Bool
		input    *bool
		expected Bool
	}{
		{name: "Set true with valid pointer", initial: Bool(false), input: ptr.Of(true), expected: Bool(true)},
		{name: "Set false with valid pointer", initial: Bool(true), input: ptr.Of(false), expected: Bool(false)},
		{name: "No change with nil pointer", initial: Bool(true), input: nil, expected: Bool(true)},
		{name: "No change with nil pointer from false", initial: Bool(false), input: nil, expected: Bool(false)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIfPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBool_SetIfPropertyPtr(t *testing.T) {
	tests := []struct {
		name     string
		initial  Bool
		input    *Bool
		expected Bool
	}{
		{name: "Set true with valid pointer", initial: Bool(false), input: ptr.Of(Bool(true)), expected: Bool(true)},
		{name: "Set false with valid pointer", initial: Bool(true), input: ptr.Of(Bool(false)), expected: Bool(false)},
		{name: "No change with nil pointer", initial: Bool(true), input: nil, expected: Bool(true)},
		{name: "No change with nil pointer from false", initial: Bool(false), input: nil, expected: Bool(false)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIfPropertyPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
