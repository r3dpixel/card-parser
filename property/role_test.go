package property

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/r3dpixel/toolkit/ptr"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/stretchr/testify/assert"
)

type roleTestContainer struct {
	fromString []propertyTestCase[string, Role]
	fromInt    []propertyTestCase[int, Role]
	marshal    []propertyTestCase[Role, string]
}

var roleTests = roleTestContainer{
	fromString: []propertyTestCase[string, Role]{
		{name: "System Lowercase", input: "system", expected: SystemRole},
		{name: "System Uppercase", input: "SYSTEM", expected: SystemRole},
		{name: "System Mixed Case", input: "SyStEm", expected: SystemRole},
		{name: "System With Whitespace", input: "  system  ", expected: SystemRole},
		{name: "System With Symbols", input: "sys-@#$tem", expected: SystemRole},
		{name: "Lorebook Depth System", input: "lorebook_depth_system", expected: SystemRole},
		{name: "User Lowercase", input: "user", expected: UserRole},
		{name: "User Uppercase", input: "USER", expected: UserRole},
		{name: "User Mixed Case", input: "UsEr", expected: UserRole},
		{name: "User With Whitespace", input: "  user  ", expected: UserRole},
		{name: "User With Symbols", input: "u-@#$ser", expected: UserRole},
		{name: "Lorebook Depth User", input: "lorebook_depth_user", expected: UserRole},
		{name: "Assistant Lowercase", input: "assistant", expected: AssistantRole},
		{name: "Assistant Uppercase", input: "ASSISTANT", expected: AssistantRole},
		{name: "Assistant Mixed Case", input: "AsSiStAnT", expected: AssistantRole},
		{name: "Assistant With Whitespace", input: "  assistant  ", expected: AssistantRole},
		{name: "Assistant With Symbols", input: "assis-@#$tant", expected: AssistantRole},
		{name: "Lorebook Depth Char", input: "lorebook_depth_char", expected: AssistantRole},
		{name: "Lorebook Depth Assistant", input: "lorebook_depth_assistant", expected: AssistantRole},
		{name: "Invalid String", input: "not a role", expected: DefaultRole},
		{name: "Empty String", input: "", expected: DefaultRole},
	},
	fromInt: []propertyTestCase[int, Role]{
		{name: "Valid Int 0 (System)", input: 0, expected: SystemRole},
		{name: "Valid Int 1 (User)", input: 1, expected: UserRole},
		{name: "Valid Int 2 (Assistant)", input: 2, expected: AssistantRole},
		{name: "Invalid Int 3", input: 3, expected: DefaultRole},
		{name: "Invalid Int 100", input: 100, expected: DefaultRole},
		{name: "Negative Int -1", input: -1, expected: DefaultRole},
		{name: "Negative Int -2", input: -2, expected: DefaultRole},
	},
	marshal: []propertyTestCase[Role, string]{
		{name: "SystemRole", input: SystemRole, expected: "0"},
		{name: "UserRole", input: UserRole, expected: "1"},
		{name: "AssistantRole", input: AssistantRole, expected: "2"},
	},
}

func TestRole_UnmarshalJSON(t *testing.T) {
	var allTestCases []propertyTestCase[string, Role]

	for _, tc := range roleTests.fromString {
		allTestCases = append(allTestCases, propertyTestCase[string, Role]{
			name:      fmt.Sprintf("From Plain String '%s'", tc.name),
			input:     tc.input,
			shouldErr: true,
			expected:  Role(0),
		})
		allTestCases = append(allTestCases, propertyTestCase[string, Role]{
			name:     fmt.Sprintf("From JSON String '%s'", tc.name),
			input:    fmt.Sprintf(`"%s"`, tc.input),
			expected: tc.expected,
		})
	}

	for _, tc := range roleTests.fromInt {
		allTestCases = append(allTestCases, propertyTestCase[string, Role]{
			name:     fmt.Sprintf("From JSON Number '%s'", tc.name),
			input:    strconv.Itoa(tc.input),
			expected: tc.expected,
		})
		allTestCases = append(allTestCases, propertyTestCase[string, Role]{
			name:     fmt.Sprintf("From JSON String Number '%s'", tc.name),
			input:    fmt.Sprintf(`"%d"`, tc.input),
			expected: tc.expected,
		})
	}

	extraTestCases := []propertyTestCase[string, Role]{
		{name: "JSON Boolean true", input: "true", expected: UserRole},
		{name: "JSON Boolean false", input: "false", expected: SystemRole},
		{name: "JSON Null", input: "null", expected: DefaultRole},
		{name: "Malformed JSON", input: "{", shouldErr: true, expected: Role(0)},
		{name: "JSON Float 1.0", input: "1.0", expected: UserRole},
		{name: "JSON Float 1.9", input: "1.9", expected: UserRole}, // cast truncates
		{name: "JSON Object", input: "{}", expected: DefaultRole},  // casts to 0
		{name: "JSON Array", input: "[]", expected: DefaultRole},   // casts to 0
	}
	allTestCases = append(allTestCases, extraTestCases...)

	for _, tc := range allTestCases {
		t.Run(tc.name, func(t *testing.T) {
			var result Role
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

func TestRole_FromString(t *testing.T) {
	for _, tc := range roleTests.fromString {
		t.Run(tc.name, func(t *testing.T) {
			result := RoleProp().FromString(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRole_FromInt(t *testing.T) {
	for _, tc := range roleTests.fromInt {
		t.Run(tc.name, func(t *testing.T) {
			result := RoleProp().FromInt(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRole_MarshalJSON(t *testing.T) {
	for _, tc := range roleTests.marshal {
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

func TestRole_SetIfPtr(t *testing.T) {
	tests := []struct {
		name     string
		initial  Role
		input    *int
		expected Role
	}{
		{name: "Set valid role 0 (System)", initial: UserRole, input: ptr.Of(0), expected: SystemRole},
		{name: "Set valid role 1 (User)", initial: SystemRole, input: ptr.Of(1), expected: UserRole},
		{name: "Set valid role 2 (Assistant)", initial: SystemRole, input: ptr.Of(2), expected: AssistantRole},
		{name: "Set invalid role 3 (defaults to System)", initial: UserRole, input: ptr.Of(3), expected: DefaultRole},
		{name: "Set invalid negative role (defaults to System)", initial: UserRole, input: ptr.Of(-1), expected: DefaultRole},
		{name: "No change with nil pointer", initial: UserRole, input: nil, expected: UserRole},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIfPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRole_SetIfPropertyPtr(t *testing.T) {
	tests := []struct {
		name     string
		initial  Role
		input    *Role
		expected Role
	}{
		{name: "Set System role", initial: UserRole, input: ptr.Of(SystemRole), expected: SystemRole},
		{name: "Set User role", initial: SystemRole, input: ptr.Of(UserRole), expected: UserRole},
		{name: "Set Assistant role", initial: SystemRole, input: ptr.Of(AssistantRole), expected: AssistantRole},
		{name: "No change with nil pointer", initial: UserRole, input: nil, expected: UserRole},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			result.SetIfPropertyPtr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
