package property

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/r3dpixel/toolkit/ptr"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/stretchr/testify/assert"
)

type propertyTestCase[T, V any] struct {
	name      string
	shouldErr bool
	input     T
	expected  V
}

type lorePositionTestContainer struct {
	fromString []propertyTestCase[string, LorePosition]
	fromInt    []propertyTestCase[int, LorePosition]
	marshal    []propertyTestCase[LorePosition, string]
}

var lorePositionTests = lorePositionTestContainer{
	fromString: []propertyTestCase[string, LorePosition]{
		{name: "Lowercase", input: "before_char", expected: BeforeCharPosition},
		{name: "Uppercase", input: "AFTER_CHAR", expected: AfterCharPosition},
		{name: "Mixed Case", input: "BeFoRe_ChAr", expected: BeforeCharPosition},
		{name: "With Whitespace AfterChar", input: "  after _char  ", expected: AfterCharPosition},
		{name: "With Whitespace BeforeChar", input: "  before    _char  ", expected: BeforeCharPosition},
		{name: "With Symbols", input: "before-@#$char", expected: BeforeCharPosition},
		{name: "Lorebook Entry Before", input: "lorebook_entry_before_char", expected: BeforeCharPosition},
		{name: "Lorebook Entry After", input: "lorebook_entry_after_char", expected: AfterCharPosition},
		{name: "Before Author Notes", input: "before_an", expected: BeforeAuthorNotes},
		{name: "Lorebook Entry Before AN", input: "lorebook_entry_before_an", expected: BeforeAuthorNotes},
		{name: "After Author Notes", input: "after_an", expected: AfterAuthorNotes},
		{name: "Lorebook Entry After AN", input: "lorebook_entry_after_an", expected: AfterAuthorNotes},
		{name: "At Depth", input: "at_depth", expected: AtDepth},
		{name: "Lorebook Entry Depth", input: "lorebook_entry_depth", expected: AtDepth},
		{name: "Lorebook In Chat 1", input: "in _$%$#^chat", expected: AtDepth},
		{name: "Lorebook In Chat 2", input: "inchat", expected: AtDepth},
		{name: "Before Example Messages", input: "before_em", expected: BeforeExampleMessages},
		{name: "Lorebook Entry Before EM", input: "lorebook_entry_before_em", expected: BeforeExampleMessages},
		{name: "After Example Messages", input: "after_em", expected: AfterExampleMessages},
		{name: "Lorebook Entry After EM", input: "lorebook_entry_after_em", expected: AfterExampleMessages},
		{name: "Invalid String", input: "not a position", expected: DefaultLorePosition},
		{name: "Empty String", input: "", expected: DefaultLorePosition},
	},
	fromInt: []propertyTestCase[int, LorePosition]{
		{name: "Valid Int 0 (BeforeChar)", input: 0, expected: BeforeCharPosition},
		{name: "Valid Int 1 (AfterChar)", input: 1, expected: AfterCharPosition},
		{name: "Valid Int 2 (BeforeAuthorNotes)", input: 2, expected: BeforeAuthorNotes},
		{name: "Valid Int 3 (AfterAuthorNotes)", input: 3, expected: AfterAuthorNotes},
		{name: "Valid Int 4 (AtDepth)", input: 4, expected: AtDepth},
		{name: "Valid Int 5 (BeforeExampleMessages)", input: 5, expected: BeforeExampleMessages},
		{name: "Valid Int 6 (AfterExampleMessages)", input: 6, expected: AfterExampleMessages},
		{name: "Invalid Int 7", input: 7, expected: DefaultLorePosition},
		{name: "Invalid Int 100", input: 100, expected: DefaultLorePosition},
		{name: "Negative Int -1", input: -1, expected: DefaultLorePosition},
		{name: "Negative Int -2", input: -2, expected: DefaultLorePosition},
	},
	marshal: []propertyTestCase[LorePosition, string]{
		{name: "BeforeCharPosition", input: BeforeCharPosition, expected: "0"},
		{name: "AfterCharPosition", input: AfterCharPosition, expected: "1"},
		{name: "BeforeAuthorNotes", input: BeforeAuthorNotes, expected: "2"},
		{name: "AfterAuthorNotes", input: AfterAuthorNotes, expected: "3"},
		{name: "AtDepth", input: AtDepth, expected: "4"},
		{name: "BeforeExampleMessages", input: BeforeExampleMessages, expected: "5"},
		{name: "AfterExampleMessages", input: AfterExampleMessages, expected: "6"},
	},
}

func TestLorePosition_UnmarshalJSON(t *testing.T) {
	var allTestCases []propertyTestCase[string, LorePosition]

	for _, tc := range lorePositionTests.fromString {
		allTestCases = append(allTestCases, propertyTestCase[string, LorePosition]{
			name:      fmt.Sprintf("From Plain String '%s'", tc.name),
			input:     tc.input,
			shouldErr: true,
			expected:  LorePosition(0),
		})
		allTestCases = append(allTestCases, propertyTestCase[string, LorePosition]{
			name:     fmt.Sprintf("From JSON String '%s'", tc.name),
			input:    fmt.Sprintf(`"%s"`, tc.input),
			expected: tc.expected,
		})
	}

	for _, tc := range lorePositionTests.fromInt {
		allTestCases = append(allTestCases, propertyTestCase[string, LorePosition]{
			name:     fmt.Sprintf("From JSON Number '%s'", tc.name),
			input:    strconv.Itoa(tc.input),
			expected: tc.expected,
		})
		allTestCases = append(allTestCases, propertyTestCase[string, LorePosition]{
			name:     fmt.Sprintf("From JSON String Number '%s'", tc.name),
			input:    fmt.Sprintf(`"%d"`, tc.input),
			expected: tc.expected,
		})
	}

	extraTestCases := []propertyTestCase[string, LorePosition]{
		{name: "JSON Boolean true", input: "true", expected: AfterCharPosition},
		{name: "JSON Boolean false", input: "false", expected: BeforeCharPosition},
		{name: "JSON Null", input: "null", expected: DefaultLorePosition},
		{name: "Malformed JSON", input: "{", shouldErr: true, expected: LorePosition(0)},
		{name: "JSON Float 1.0", input: "1.0", expected: AfterCharPosition},
		{name: "JSON Float 1.9", input: "1.9", expected: AfterCharPosition}, // cast truncates
		{name: "JSON Object", input: "{}", expected: DefaultLorePosition},   // casts to 0
		{name: "JSON Array", input: "[]", expected: DefaultLorePosition},    // casts to 0
	}
	allTestCases = append(allTestCases, extraTestCases...)

	for _, tc := range allTestCases {
		t.Run(tc.name, func(t *testing.T) {
			var result LorePosition
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

func TestLorePosition_FromString(t *testing.T) {
	for _, tc := range lorePositionTests.fromString {
		t.Run(tc.name, func(t *testing.T) {
			result := LorePositionProp().FromString(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestLorePosition_FromInt(t *testing.T) {
	for _, tc := range lorePositionTests.fromInt {
		t.Run(tc.name, func(t *testing.T) {
			result := LorePositionProp().FromInt(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestLorePosition_MarshalJSON(t *testing.T) {
	for _, tc := range lorePositionTests.marshal {
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

func TestLorePosition_SetIfPtr(t *testing.T) {
	tests := []struct {
		name     string
		initial  LorePosition
		input    *int
		expected LorePosition
	}{
		{name: "Set valid position 0 (BeforeChar)", initial: AfterCharPosition, input: ptr.Of(0), expected: BeforeCharPosition},
		{name: "Set valid position 1 (AfterChar)", initial: BeforeCharPosition, input: ptr.Of(1), expected: AfterCharPosition},
		{name: "Set valid position 2 (BeforeAuthorNotes)", initial: BeforeCharPosition, input: ptr.Of(2), expected: BeforeAuthorNotes},
		{name: "Set valid position 3 (AfterAuthorNotes)", initial: BeforeCharPosition, input: ptr.Of(3), expected: AfterAuthorNotes},
		{name: "Set valid position 4 (AtDepth)", initial: BeforeCharPosition, input: ptr.Of(4), expected: AtDepth},
		{name: "Set valid position 5 (BeforeExampleMessages)", initial: BeforeCharPosition, input: ptr.Of(5), expected: BeforeExampleMessages},
		{name: "Set valid position 6 (AfterExampleMessages)", initial: BeforeCharPosition, input: ptr.Of(6), expected: AfterExampleMessages},
		{name: "Set invalid position 7 (defaults to BeforeChar)", initial: AfterCharPosition, input: ptr.Of(7), expected: DefaultLorePosition},
		{name: "Set invalid negative position (defaults to BeforeChar)", initial: AfterCharPosition, input: ptr.Of(-1), expected: DefaultLorePosition},
		{name: "No change with nil pointer", initial: AfterCharPosition, input: nil, expected: AfterCharPosition},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIfPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLorePosition_SetIfPropertyPtr(t *testing.T) {
	tests := []struct {
		name     string
		initial  LorePosition
		input    *LorePosition
		expected LorePosition
	}{
		{name: "Set BeforeChar position", initial: AfterCharPosition, input: ptr.Of(BeforeCharPosition), expected: BeforeCharPosition},
		{name: "Set AfterChar position", initial: BeforeCharPosition, input: ptr.Of(AfterCharPosition), expected: AfterCharPosition},
		{name: "Set BeforeAuthorNotes position", initial: BeforeCharPosition, input: ptr.Of(BeforeAuthorNotes), expected: BeforeAuthorNotes},
		{name: "Set AfterAuthorNotes position", initial: BeforeCharPosition, input: ptr.Of(AfterAuthorNotes), expected: AfterAuthorNotes},
		{name: "Set AtDepth position", initial: BeforeCharPosition, input: ptr.Of(AtDepth), expected: AtDepth},
		{name: "Set BeforeExampleMessages position", initial: BeforeCharPosition, input: ptr.Of(BeforeExampleMessages), expected: BeforeExampleMessages},
		{name: "Set AfterExampleMessages position", initial: BeforeCharPosition, input: ptr.Of(AfterExampleMessages), expected: AfterExampleMessages},
		{name: "No change with nil pointer", initial: AfterCharPosition, input: nil, expected: AfterCharPosition},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIfPropertyPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
