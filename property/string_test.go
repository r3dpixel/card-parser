package property

import (
	"testing"

	"github.com/r3dpixel/toolkit/ptr"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/stringsx"
	"github.com/stretchr/testify/assert"
)

func TestString_NormalizeSymbols(t *testing.T) {
	tests := []struct {
		name     string
		input    String
		expected String
	}{
		{name: "No symbols", input: "hello world", expected: "hello world"},
		{name: "Empty string", input: "", expected: ""},
		{name: "Unicode symbols", input: String("héllo—world"), expected: String("héllo—world")},
		{name: "Mixed symbols", input: String(`test‘s 《quotes》`), expected: String("test's \"quotes\"")},
		{name: "Already normalized", input: String("normal text"), expected: String("normal text")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.NormalizeSymbols()
			assert.Equal(t, tt.expected, tt.input)
		})
	}
}

func TestString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Float values
		{name: "Float positive", input: "123.45", expected: "123.45"},
		{name: "Float negative", input: "-67.89", expected: "-67.89"},
		{name: "Float zero", input: "0.0", expected: "0"},
		{name: "Float integer", input: "42.0", expected: "42"},

		// String values
		{name: "String simple", input: `"hello"`, expected: "hello"},
		{name: "String empty", input: `""`, expected: ""},
		{name: "String with spaces", input: `"hello world"`, expected: "hello world"},
		{name: "String with quotes", input: `"say \"hello\""`, expected: `say "hello"`},
		{name: "String unicode", input: `"héllo 世界"`, expected: "héllo 世界"},

		// Boolean values
		{name: "Bool true", input: "true", expected: "true"},
		{name: "Bool false", input: "false", expected: "false"},

		// Null value
		{name: "Null", input: "null", expected: ""},

		// Array values
		{name: "Array empty", input: "[]", expected: "[]"},
		{name: "Array simple", input: `["a", "b", "c"]`, expected: `["a","b","c"]`},
		{name: "Array mixed", input: `[1, "two", true, null]`, expected: `[1,"two",true,null]`},
		{name: "Array nested", input: `[[1, 2], ["a", "b"]]`, expected: `[[1,2],["a","b"]]`},

		// Object values
		{name: "Object empty", input: "{}", expected: "{}"},
		{name: "Object simple", input: `{"key": "value"}`, expected: `{"key":"value"}`},
		{name: "Object nested", input: `{"outer": {"inner": "value"}}`, expected: `{"outer":{"inner":"value"}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result String
			err := sonicx.Config.UnmarshalFromString(tt.input, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "String simple", input: "hello", expected: `"hello"`},
		{name: "String empty", input: "", expected: `""`},
		{name: "String with quotes", input: `say "hello"`, expected: `"say \"hello\""`},
		{name: "String unicode", input: "héllo 世界", expected: `"héllo 世界"`},
		{name: "String number", input: "123", expected: `"123"`},
		{name: "String boolean", input: "true", expected: `"true"`},
		{name: "String JSON", input: `{"key":"value"}`, expected: `"{\"key\":\"value\"}"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sonicx.Config.Marshal(&tt.input)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(result))
		})
	}
}

func TestString_RoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		original String
	}{
		{name: "Simple string", original: "hello world"},
		{name: "Empty string", original: ""},
		{name: "Unicode string", original: "héllo 世界"},
		{name: "JSON-like string", original: `{"key":"value"}`},
		{name: "Number string", original: "123.45"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := sonicx.Config.Marshal(&tt.original)
			assert.NoError(t, err)

			// Unmarshal back
			var result String
			err = sonicx.Config.UnmarshalFromString(stringsx.FromBytes(jsonData), &result)
			assert.NoError(t, err)

			// Should be equal
			assert.Equal(t, tt.original, result)
		})
	}
}

func TestString_TypeConversions(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{name: "int to string", input: 42, expected: "42"},
		{name: "float to string", input: 3.14, expected: "3.14"},
		{name: "bool true to string", input: true, expected: "true"},
		{name: "bool false to string", input: false, expected: "false"},
		{name: "nil to string", input: nil, expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := sonicx.Config.Marshal(tt.input)
			assert.NoError(t, err)

			var result String
			err = sonicx.Config.UnmarshalFromString(stringsx.FromBytes(jsonData), &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestString_ComplexObject(t *testing.T) {
	// Test complex object separately to handle key ordering
	input := `{"num": 42, "str": "hello", "bool": true, "null": null}`
	var result String
	err := sonicx.Config.UnmarshalFromString(input, &result)
	assert.NoError(t, err)

	// Verify it's valid JSON and contains all expected keys
	resultStr := string(result)
	assert.Contains(t, resultStr, `"num":42`)
	assert.Contains(t, resultStr, `"str":"hello"`)
	assert.Contains(t, resultStr, `"bool":true`)
	assert.Contains(t, resultStr, `"null":null`)
}

func TestString_SetIf(t *testing.T) {
	tests := []struct {
		name     string
		initial  String
		value    string
		expected String
	}{
		{name: "Update with non-blank value", initial: "old", value: "new", expected: "new"},
		{name: "Update empty string with value", initial: "", value: "value", expected: "value"},
		{name: "No update with blank value", initial: "original", value: "", expected: "original"},
		{name: "No update with whitespace only", initial: "original", value: "   ", expected: "original"},
		{name: "No update with tab only", initial: "original", value: "\t", expected: "original"},
		{name: "No update with newline only", initial: "original", value: "\n", expected: "original"},
		{name: "Update with valid whitespace content", initial: "old", value: "hello world", expected: "hello world"},
		{name: "Update with unicode", initial: "old", value: "héllo 世界", expected: "héllo 世界"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIf(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestString_ErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "Malformed JSON", input: "{"},
		{name: "Invalid JSON", input: "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result String
			err := sonicx.Config.UnmarshalFromString(tt.input, &result)
			assert.Error(t, err)
		})
	}
}

func TestString_SetIfPtr(t *testing.T) {
	tests := []struct {
		name     string
		initial  String
		input    *string
		expected String
	}{
		{name: "Set with non-blank value", initial: "old", input: ptr.Of("new"), expected: "new"},
		{name: "Set empty string with value", initial: "", input: ptr.Of("value"), expected: "value"},
		{name: "No update with blank value", initial: "original", input: ptr.Of(""), expected: "original"},
		{name: "No update with whitespace only", initial: "original", input: ptr.Of("   "), expected: "original"},
		{name: "No update with nil pointer", initial: "original", input: nil, expected: "original"},
		{name: "Update with unicode", initial: "old", input: ptr.Of("héllo 世界"), expected: "héllo 世界"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIfPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestString_SetIfProperty(t *testing.T) {
	tests := []struct {
		name     string
		initial  String
		input    String
		expected String
	}{
		{name: "Set with non-blank value", initial: "old", input: "new", expected: "new"},
		{name: "Set empty string with value", initial: "", input: "value", expected: "value"},
		{name: "No update with blank value", initial: "original", input: "", expected: "original"},
		{name: "No update with whitespace only", initial: "original", input: "   ", expected: "original"},
		{name: "Update with unicode", initial: "old", input: "héllo 世界", expected: "héllo 世界"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIfProperty(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestString_SetIfPropertyPtr(t *testing.T) {
	tests := []struct {
		name     string
		initial  String
		input    *String
		expected String
	}{
		{name: "Set with non-blank value", initial: "old", input: ptr.Of(String("new")), expected: "new"},
		{name: "Set empty string with value", initial: "", input: ptr.Of(String("value")), expected: "value"},
		{name: "No update with blank value", initial: "original", input: ptr.Of(String("")), expected: "original"},
		{name: "No update with whitespace only", initial: "original", input: ptr.Of(String("   ")), expected: "original"},
		{name: "No update with nil pointer", initial: "original", input: nil, expected: "original"},
		{name: "Update with unicode", initial: "old", input: ptr.Of(String("héllo 世界")), expected: "héllo 世界"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIfPropertyPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
