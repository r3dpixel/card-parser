package property

import (
	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/spf13/cast"
)

type StringArray []string

func (s *StringArray) OnFloat(floatValue float64) {
	*s = StringArray{cast.ToString(floatValue)}
}

func (s *StringArray) OnString(stringValue string) {
	*s = StringArray{stringValue}
}

func (s *StringArray) OnBool(boolValue bool) {
	*s = StringArray{cast.ToString(boolValue)}
}

func (s *StringArray) OnNull() {
	*s = make([]string, 0)
}

func (s *StringArray) OnObject(objectValue map[string]any) {
	*s = StringArray{jsonx.String(objectValue)}
}

func (s *StringArray) OnArray(arrayValue []any) {
	stringItems := make([]string, len(arrayValue))
	for index, item := range arrayValue {
		var stringItem string
		switch v := item.(type) {
		case []any, map[string]any:
			stringItem = jsonx.String(v)
		default:
			stringItem = cast.ToString(item)
		}
		stringItems[index] = stringItem
	}
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
