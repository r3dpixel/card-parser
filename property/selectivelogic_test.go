package property

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/r3dpixel/toolkit/ptr"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/stretchr/testify/assert"
)

// selectiveLogicTestContainer holds all test data in a single namespace.
type selectiveLogicTestContainer struct {
	fromString []propertyTestCase[string, SelectiveLogic]
	fromInt    []propertyTestCase[int, SelectiveLogic]
	marshal    []propertyTestCase[SelectiveLogic, string]
}

// selectiveLogicTests holds all defined test cases.
var selectiveLogicTests = func() selectiveLogicTestContainer {
	container := selectiveLogicTestContainer{
		fromInt: []propertyTestCase[int, SelectiveLogic]{
			{name: "Valid Int 0", input: 0, expected: SelectiveAndAny},
			{name: "Valid Int 1", input: 1, expected: SelectiveNotAll},
			{name: "Valid Int 2", input: 2, expected: SelectiveNotAny},
			{name: "Valid Int 3", input: 3, expected: SelectiveAndAll},
			{name: "Invalid Int 4", input: 4, expected: DefaultSelectiveLogic},
			{name: "Negative Int", input: -1, expected: DefaultSelectiveLogic},
		},
		marshal: []propertyTestCase[SelectiveLogic, string]{
			{name: "Marshal AndAny", input: SelectiveAndAny, expected: "0"},
			{name: "Marshal NotAll", input: SelectiveNotAll, expected: "1"},
			{name: "Marshal NotAny", input: SelectiveNotAny, expected: "2"},
			{name: "Marshal AndAll", input: SelectiveAndAll, expected: "3"},
			{name: "Marshal Default", input: DefaultSelectiveLogic, expected: "0"},
		},
	}

	// Add extensive FromString test cases
	fromStringFormats := make([][]string, SelectiveLogicEnd+1)
	fromStringFormats[SelectiveAndAny] = []string{"and_any", "AND_ANY", "AnD_Any", "A_n_D_a_n_y", "an_d_A_N_Y"}
	fromStringFormats[SelectiveNotAll] = []string{"not_all", "NOT_ALL", "NoT_All", "N_o_T_A_l_l", "no_t_A_L_L"}
	fromStringFormats[SelectiveNotAny] = []string{"not_any", "NOT_ANY", "NoT_Any", "N_o_T_A_n_y", "no_t_A_N_Y"}
	fromStringFormats[SelectiveAndAll] = []string{"and_all", "AND_ALL", "AnD_All", "A_n_d_A_l_l", "an_d_A_L_L"}

	for i, formats := range fromStringFormats {
		for _, format := range formats {
			container.fromString = append(container.fromString, propertyTestCase[string, SelectiveLogic]{
				name:     fmt.Sprintf("FromString format '%s'", format),
				input:    format,
				expected: SelectiveLogic(i),
			})
		}
	}
	// Add other FromString cases
	container.fromString = append(container.fromString,
		propertyTestCase[string, SelectiveLogic]{name: "With Whitespace AndAll", input: "  and all  ", expected: SelectiveAndAll},
		propertyTestCase[string, SelectiveLogic]{name: "With Symbols notall", input: "not-@#$all", expected: SelectiveNotAll},
		propertyTestCase[string, SelectiveLogic]{name: "Invalid String", input: "nix nox", expected: DefaultSelectiveLogic},
	)

	return container
}()

func TestSelectiveLogic_UnmarshalJSON(t *testing.T) {
	var allTestCases = []propertyTestCase[string, SelectiveLogic]{
		{name: "JSON String 'andany'", input: `"andany"`, expected: SelectiveAndAny},
		{name: "Plain String 'notall'", input: "notall", shouldErr: true, expected: DefaultSelectiveLogic},
		{name: "JSON Boolean true", input: "true", expected: SelectiveNotAll},   // true -> 1
		{name: "JSON Boolean false", input: "false", expected: SelectiveAndAny}, // false -> 0
		{name: "JSON Null", input: "null", expected: DefaultSelectiveLogic},
		{name: "JSON Array Empty", input: "[]", expected: DefaultSelectiveLogic},
		{name: "JSON Array Simple", input: `["a", "b", "c"]`, expected: DefaultSelectiveLogic},
		{name: "JSON Array Mixed", input: `[1, "two", true, null]`, expected: DefaultSelectiveLogic},
		{name: "JSON Object Empty", input: "{}", expected: DefaultSelectiveLogic},
		{name: "JSON Object Simple", input: `{"key": "value"}`, expected: DefaultSelectiveLogic},
		{name: "JSON Object Complex", input: `{"num": 42, "str": "hello", "bool": true}`, expected: DefaultSelectiveLogic},
		{name: "Malformed JSON", input: "{", shouldErr: true, expected: DefaultSelectiveLogic},
	}

	const offset = 3000
	for i := -offset; i <= offset; i++ {
		expected := slParser.FromInt(i)

		allTestCases = append(allTestCases, propertyTestCase[string, SelectiveLogic]{
			name:     fmt.Sprintf("JSON Number %d", i),
			input:    strconv.Itoa(i),
			expected: expected,
		})
		allTestCases = append(allTestCases, propertyTestCase[string, SelectiveLogic]{
			name:     fmt.Sprintf("JSON String Number %d", i),
			input:    fmt.Sprintf(`"%d"`, i),
			expected: expected,
		})
	}

	for _, tc := range allTestCases {
		t.Run(tc.name, func(t *testing.T) {
			var result SelectiveLogic
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

func TestSelectiveLogic_FromString(t *testing.T) {
	for _, tc := range selectiveLogicTests.fromString {
		t.Run(tc.name, func(t *testing.T) {
			result := SelectiveLogicProp().FromString(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSelectiveLogic_FromInt(t *testing.T) {
	for _, tc := range selectiveLogicTests.fromInt {
		t.Run(tc.name, func(t *testing.T) {
			result := SelectiveLogicProp().FromInt(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSelectiveLogic_MarshalJSON(t *testing.T) {
	for _, tc := range selectiveLogicTests.marshal {
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

func TestSelectiveLogic_SetIfPtr(t *testing.T) {
	tests := []struct {
		name     string
		initial  SelectiveLogic
		input    *int
		expected SelectiveLogic
	}{
		{name: "Set valid logic 0 (AndAny)", initial: SelectiveNotAll, input: ptr.Of(0), expected: SelectiveAndAny},
		{name: "Set valid logic 1 (NotAll)", initial: SelectiveAndAny, input: ptr.Of(1), expected: SelectiveNotAll},
		{name: "Set valid logic 2 (NotAny)", initial: SelectiveAndAny, input: ptr.Of(2), expected: SelectiveNotAny},
		{name: "Set valid logic 3 (AndAll)", initial: SelectiveAndAny, input: ptr.Of(3), expected: SelectiveAndAll},
		{name: "Set invalid logic 4 (defaults to AndAny)", initial: SelectiveNotAll, input: ptr.Of(4), expected: DefaultSelectiveLogic},
		{name: "Set invalid negative logic (defaults to AndAny)", initial: SelectiveNotAll, input: ptr.Of(-1), expected: DefaultSelectiveLogic},
		{name: "No change with nil pointer", initial: SelectiveNotAll, input: nil, expected: SelectiveNotAll},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIfPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSelectiveLogic_SetIfPropertyPtr(t *testing.T) {
	tests := []struct {
		name     string
		initial  SelectiveLogic
		input    *SelectiveLogic
		expected SelectiveLogic
	}{
		{name: "Set AndAny logic", initial: SelectiveNotAll, input: ptr.Of(SelectiveAndAny), expected: SelectiveAndAny},
		{name: "Set NotAll logic", initial: SelectiveAndAny, input: ptr.Of(SelectiveNotAll), expected: SelectiveNotAll},
		{name: "Set NotAny logic", initial: SelectiveAndAny, input: ptr.Of(SelectiveNotAny), expected: SelectiveNotAny},
		{name: "Set AndAll logic", initial: SelectiveAndAny, input: ptr.Of(SelectiveAndAll), expected: SelectiveAndAll},
		{name: "No change with nil pointer", initial: SelectiveNotAll, input: nil, expected: SelectiveNotAll},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIfPropertyPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
