package property

import (
	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/spf13/cast"
)

// StringArray represents an array of strings
type StringArray []string

// OnFloat populates the StringArray with a single string containing the float value
func (s *StringArray) OnFloat(floatValue float64) {
	*s = StringArray{cast.ToString(floatValue)}
}

// OnString populates the StringArray with a single string containing the string value
func (s *StringArray) OnString(stringValue string) {
	*s = StringArray{stringValue}
}

// OnBool populates the StringArray with a single string containing the bool value
func (s *StringArray) OnBool(boolValue bool) {
	*s = StringArray{cast.ToString(boolValue)}
}

// OnNull populates the StringArray with a nil array
func (s *StringArray) OnNull() {
	*s = make([]string, 0)
}

// OnObject populates the StringArray with a single string containing the JSON representation of the object
func (s *StringArray) OnObject(objectValue map[string]any) {
	*s = StringArray{jsonx.String(objectValue)}
}

// OnArray populates the StringArray with the array values converted to strings
func (s *StringArray) OnArray(arrayValue []any) {
	// Create a new array of strings
	stringItems := make([]string, len(arrayValue))
	// Iterate over the array and convert each item to a string
	for index, item := range arrayValue {
		// Convert the item to a string
		var stringItem string
		switch v := item.(type) {
		case []any, map[string]any:
			// If the item is an array or object, convert it to JSON string
			stringItem = jsonx.String(v)
		default:
			// Otherwise, convert the primitive value to a string
			stringItem = cast.ToString(item)
		}
		// Save the string item in the array
		stringItems[index] = stringItem
	}
	// Set the StringArray to the new array
	*s = stringItems
}

// MarshalJSON marshals the StringArray to JSON using Sonic
func (s *StringArray) MarshalJSON() ([]byte, error) {
	return sonicx.Config.Marshal((*[]string)(s))
}

// UnmarshalJSON unmarshals JSON data into the StringArray using Sonic
func (s *StringArray) UnmarshalJSON(data []byte) error {
	return jsonx.HandleEntity(data, s)
}
