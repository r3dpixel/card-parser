package property

import (
	"testing"

	"github.com/r3dpixel/toolkit/ptr"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/stretchr/testify/assert"
)

type intPropertyTestContainer struct {
	unmarshal []propertyTestCase[string, int]
	marshal   []propertyTestCase[Integer, string]
}

var intPropertyTests = intPropertyTestContainer{
	unmarshal: []propertyTestCase[string, int]{
		{name: "JSON Number Positive", input: "123", expected: 123},
		{name: "JSON Number Negative", input: "-50", expected: -50},
		{name: "JSON Number Zero", input: "0", expected: 0},
		{name: "JSON Float", input: "123.45", expected: 123},

		{name: "JSON String with Number", input: `"456"`, expected: 456},
		{name: "JSON String with Negative Number", input: `"-78"`, expected: -78},

		{name: "Plain String Number", input: "99", expected: 99},
		{name: "Plain String with Whitespace", input: "  101  ", expected: 101},

		{name: "JSON Boolean true", input: "true", expected: 1},
		{name: "JSON Boolean false", input: "false", expected: 0},

		{name: "JSON Null", input: "null", expected: 0},
		{name: "Empty Input", shouldErr: true, input: "", expected: 0},
		{name: "Empty JSON String", input: `""`, expected: 0},
		{name: "Non-numeric JSON String", input: `"hello"`, expected: 0},
		{name: "Non-numeric Plain String", shouldErr: true, input: "world", expected: 0},
		{name: "JSON Object", input: "{}", expected: 0},
		{name: "JSON Array", input: "[]", expected: 0},
		{name: "Malformed JSON", shouldErr: true, input: "{", expected: 0},
	},
	marshal: []propertyTestCase[Integer, string]{
		{name: "Positive Value", input: 1000, expected: "1000"},
		{name: "Negative Value", input: -250, expected: "-250"},
		{name: "Zero Value", input: 0, expected: "0"},
	},
}

func TestIntProperty_UnmarshalJSON(t *testing.T) {
	for _, tc := range intPropertyTests.unmarshal {
		t.Run(tc.name, func(t *testing.T) {
			var result Integer
			err := sonicx.Config.UnmarshalFromString(tc.input, &result)
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, int(result))
		})
	}
}

func TestIntProperty_MarshalJSON(t *testing.T) {
	for _, tc := range intPropertyTests.marshal {
		t.Run(tc.name, func(t *testing.T) {
			bytes, err := sonicx.Config.Marshal(&tc.input)
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.JSONEq(t, tc.expected, string(bytes))
		})
	}
}

func TestInteger_SetIfPtr(t *testing.T) {
	tests := []struct {
		name     string
		initial  Integer
		input    *int
		expected Integer
	}{
		{name: "Set positive value with valid pointer", initial: Integer(0), input: ptr.Of(123), expected: Integer(123)},
		{name: "Set negative value with valid pointer", initial: Integer(0), input: ptr.Of(-456), expected: Integer(-456)},
		{name: "Set zero with valid pointer", initial: Integer(123), input: ptr.Of(0), expected: Integer(0)},
		{name: "No change with nil pointer", initial: Integer(123), input: nil, expected: Integer(123)},
		{name: "No change with nil pointer from zero", initial: Integer(0), input: nil, expected: Integer(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIfPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInteger_SetIfPropertyPtr(t *testing.T) {
	tests := []struct {
		name     string
		initial  Integer
		input    *int
		expected Integer
	}{
		{name: "Set positive value with valid pointer", initial: Integer(0), input: ptr.Of(123), expected: Integer(123)},
		{name: "Set negative value with valid pointer", initial: Integer(0), input: ptr.Of(-456), expected: Integer(-456)},
		{name: "Set zero with valid pointer", initial: Integer(123), input: ptr.Of(0), expected: Integer(0)},
		{name: "No change with nil pointer", initial: Integer(123), input: nil, expected: Integer(123)},
		{name: "No change with nil pointer from zero", initial: Integer(0), input: nil, expected: Integer(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIfPropertyPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
