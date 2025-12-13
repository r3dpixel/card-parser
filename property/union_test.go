package property

import (
	"testing"

	"github.com/r3dpixel/toolkit/ptr"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/stretchr/testify/assert"
)

type bookEntryIDTestContainer struct {
	unmarshal []propertyTestCase[string, Union]
	marshal   []propertyTestCase[Union, string]
}

var bookEntryIDTests = bookEntryIDTestContainer{
	unmarshal: []propertyTestCase[string, Union]{
		{name: "JSON Number", input: "123", expected: Union{IntValue: ptr.Of(123)}},
		{name: "JSON String Number", input: `"456"`, expected: Union{IntValue: ptr.Of(456)}},
		{name: "JSON Boolean true", input: "true", expected: Union{IntValue: ptr.Of(1)}},
		{name: "JSON Boolean false", input: "false", expected: Union{IntValue: ptr.Of(0)}},
		{name: "JSON Null", input: "null", expected: Union{IntValue: ptr.Of(0)}},

		{name: "JSON String", input: `"hello"`, expected: Union{StringValue: ptr.Of("hello")}},
		{name: "Plain String (Invalid JSON)", shouldErr: true, input: "world", expected: Union{}},
		{name: "Plain String Number", input: "99", expected: Union{IntValue: ptr.Of(99)}},
		{name: "Empty JSON String", input: `""`, expected: Union{IntValue: ptr.Of(0)}},
		{name: "Empty Input", shouldErr: true, input: "", expected: Union{}},
		{name: "JSON Object", input: `{"a":"prop"}`, expected: Union{StringValue: ptr.Of(`{"a":"prop"}`)}},
		{name: "JSON Array", input: "[]", expected: Union{StringValue: ptr.Of("[]")}},
		{name: "Malformed JSON", shouldErr: true, input: "{", expected: Union{}},
	},
	marshal: []propertyTestCase[Union, string]{
		{name: "With IntValue", input: Union{IntValue: ptr.Of(9999)}, expected: "9999"},
		{name: "With StringValue", input: Union{StringValue: ptr.Of("random_id")}, expected: `"random_id"`},
		{name: "With Both Values (IntValue takes precedence)", input: Union{IntValue: ptr.Of(123), StringValue: ptr.Of("abc")}, expected: "123"},
		{name: "With No Values", input: Union{}, expected: "null"},
		//{name: "With Nil Pointer", input: *(*Union)(nil), expected: "null"},
	},
}

func assertBookEntryIDEqual(t *testing.T, expected, actual Union, msgAndArgs ...interface{}) {
	t.Helper()
	assert.EqualValues(t, expected, actual, msgAndArgs...)
}

func TestBookEntryID_UnmarshalJSON(t *testing.T) {
	for _, tc := range bookEntryIDTests.unmarshal {
		t.Run(tc.name, func(t *testing.T) {
			var result Union
			err := sonicx.Config.UnmarshalFromString(tc.input, &result)

			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assertBookEntryIDEqual(t, tc.expected, result)
		})
	}
}

func TestBookEntryID_MarshalJSON(t *testing.T) {
	for _, tc := range bookEntryIDTests.marshal {
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
