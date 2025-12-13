package property

import (
	"testing"

	"github.com/r3dpixel/toolkit/ptr"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/stringsx"
	"github.com/stretchr/testify/assert"
)

func TestFloat_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		// Float values
		{name: "Float positive", input: "123.45", expected: 123.45},
		{name: "Float negative", input: "-67.89", expected: -67.89},
		{name: "Float zero", input: "0.0", expected: 0.0},
		{name: "Float integer", input: "42", expected: 42.0},
		{name: "Float scientific", input: "1.23e-4", expected: 1.23e-4},

		// String values
		{name: "String number", input: `"123.45"`, expected: 123.45},
		{name: "String negative", input: `"-67.89"`, expected: -67.89},
		{name: "String zero", input: `"0"`, expected: 0.0},
		{name: "String integer", input: `"42"`, expected: 42.0},
		{name: "String invalid", input: `"not a number"`, expected: 0.0},
		{name: "String empty", input: `""`, expected: 0.0},

		// Boolean values
		{name: "Bool true", input: "true", expected: 1.0},
		{name: "Bool false", input: "false", expected: 0.0},

		// Null value
		{name: "Null", input: "null", expected: 0.0},

		// Array values (default to 0.0)
		{name: "Array empty", input: "[]", expected: 0.0},
		{name: "Array simple", input: `["a", "b", "c"]`, expected: 0.0},
		{name: "Array mixed", input: `[1, "two", true, null]`, expected: 0.0},

		// Object values (default to 0.0)
		{name: "Object empty", input: "{}", expected: 0.0},
		{name: "Object simple", input: `{"key": "value"}`, expected: 0.0},
		{name: "Object complex", input: `{"num": 42, "str": "hello"}`, expected: 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Float
			err := sonicx.Config.UnmarshalFromString(tt.input, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, float64(result))
		})
	}
}

func TestFloat_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    Float
		expected string
	}{
		{name: "Float positive", input: 123.45, expected: "123.45"},
		{name: "Float negative", input: -67.89, expected: "-67.89"},
		{name: "Float zero", input: 0.0, expected: "0"},
		{name: "Float integer", input: 42.0, expected: "42"},
		{name: "Float scientific small", input: 1.23e-4, expected: "0.000123"},
		{name: "Float scientific large", input: 1.23e6, expected: "1230000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sonicx.Config.Marshal(&tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestFloat_RoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		original Float
	}{
		{name: "Positive float", original: 123.45},
		{name: "Negative float", original: -67.89},
		{name: "Zero", original: 0.0},
		{name: "Integer float", original: 42.0},
		{name: "Small decimal", original: 0.123},
		{name: "Large number", original: 123456.789},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := sonicx.Config.Marshal(&tt.original)
			assert.NoError(t, err)

			// Unmarshal back
			var result Float
			err = sonicx.Config.UnmarshalFromString(stringsx.FromBytes(jsonData), &result)
			assert.NoError(t, err)

			// Should be equal
			assert.Equal(t, tt.original, result)
		})
	}
}

func TestFloat_TypeConversions(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected float64
	}{
		{name: "int to float", input: 42, expected: 42.0},
		{name: "float to float", input: 3.14, expected: 3.14},
		{name: "bool true to float", input: true, expected: 1.0},
		{name: "bool false to float", input: false, expected: 0.0},
		{name: "nil to float", input: nil, expected: 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := sonicx.Config.Marshal(tt.input)
			assert.NoError(t, err)

			var result Float
			err = sonicx.Config.UnmarshalFromString(stringsx.FromBytes(jsonData), &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, float64(result))
		})
	}
}

func TestFloat_ErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "Malformed JSON", input: "{"},
		{name: "Invalid JSON", input: "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Float
			err := sonicx.Config.UnmarshalFromString(tt.input, &result)
			assert.Error(t, err)
		})
	}
}

func TestFloat_SetIfPtr(t *testing.T) {
	tests := []struct {
		name     string
		initial  Float
		input    *float64
		expected Float
	}{
		{name: "Set positive value with valid pointer", initial: Float(0.0), input: ptr.Of(123.45), expected: Float(123.45)},
		{name: "Set negative value with valid pointer", initial: Float(0.0), input: ptr.Of(-67.89), expected: Float(-67.89)},
		{name: "Set zero with valid pointer", initial: Float(123.45), input: ptr.Of(0.0), expected: Float(0.0)},
		{name: "No change with nil pointer", initial: Float(123.45), input: nil, expected: Float(123.45)},
		{name: "No change with nil pointer from zero", initial: Float(0.0), input: nil, expected: Float(0.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIfPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFloat_SetIfPropertyPtr(t *testing.T) {
	tests := []struct {
		name     string
		initial  Float
		input    *Float
		expected Float
	}{
		{name: "Set positive value with valid pointer", initial: Float(0.0), input: ptr.Of(Float(123.45)), expected: Float(123.45)},
		{name: "Set negative value with valid pointer", initial: Float(0.0), input: ptr.Of(Float(-67.89)), expected: Float(-67.89)},
		{name: "Set zero with valid pointer", initial: Float(123.45), input: ptr.Of(Float(0.0)), expected: Float(0.0)},
		{name: "No change with nil pointer", initial: Float(123.45), input: nil, expected: Float(123.45)},
		{name: "No change with nil pointer from zero", initial: Float(0.0), input: nil, expected: Float(0.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIfPropertyPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
