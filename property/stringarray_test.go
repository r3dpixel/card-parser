package property

import (
	"testing"

	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/stretchr/testify/assert"
)

type stringArrayTestContainer struct {
	unmarshal []propertyTestCase[string, StringArray]
	marshal   []propertyTestCase[StringArray, string]
}

var stringArrayTests = stringArrayTestContainer{
	unmarshal: []propertyTestCase[string, StringArray]{
		{
			name:     "JSON Array of Strings",
			input:    `["a", "b", "c"]`,
			expected: StringArray{"a", "b", "c"},
		},
		{
			name:     "Single JSON String",
			input:    `"hello"`,
			expected: StringArray{"hello"},
		},
		{
			name:     "Single JSON Number",
			input:    "123",
			expected: StringArray{"123"},
		},
		{
			name:     "JSON Boolean true",
			input:    "true",
			expected: StringArray{"true"},
		},
		{
			name:     "JSON Boolean false",
			input:    "false",
			expected: StringArray{"false"},
		},
		{
			name:     "JSON Object",
			input:    `{"key":"value","number":42}`,
			expected: StringArray{`{"key":"value","number":42}`},
		},
		{
			name:     "JSON Object Empty",
			input:    `{}`,
			expected: StringArray{`{}`},
		},
		{
			name:      "Plain String (Invalid JSON)",
			shouldErr: true,
			input:     "plain text",
			expected:  nil,
		},
		{
			name:      "Malformed JSON",
			shouldErr: true,
			input:     "{malformed",
			expected:  nil,
		},

		{
			name:     "Array with Mixed Types",
			input:    `["a", 1, true, {"key":"value"}]`,
			expected: StringArray{"a", "1", "true", `{"key":"value"}`},
		},
		{
			name:     "JSON Null",
			input:    "null",
			expected: []string{},
		},
		{
			name:     "Empty JSON Array",
			input:    "[]",
			expected: StringArray{},
		},
	},
	marshal: []propertyTestCase[StringArray, string]{
		{
			name:     "Valid StringArray",
			input:    StringArray{"a", "b", "c"},
			expected: `["a","b","c"]`,
		},
		{
			name:     "Empty StringArray",
			input:    StringArray{},
			expected: `[]`,
		},
		{
			name:     "Nil StringArray",
			input:    nil,
			expected: `[]`,
		},
	},
}

func TestStringArray_UnmarshalJSON(t *testing.T) {
	originalConfig := sonicx.Config
	defer func() { sonicx.Config = originalConfig }()
	sonicx.Config = sonicx.StableSort

	for _, tc := range stringArrayTests.unmarshal {
		t.Run(tc.name, func(t *testing.T) {
			var result StringArray
			err := sonicx.Config.UnmarshalFromString(tc.input, &result)
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestStringArray_MarshalJSON(t *testing.T) {
	for _, tc := range stringArrayTests.marshal {
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
