package property

import (
	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/ptr"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/spf13/cast"
)

// Union represents a union of integer and string values
type Union struct {
	IntValue    *int
	StringValue *string
}

// OnFloat populates the Union with an integer value from a float64
func (u *Union) OnFloat(floatValue float64) {
	// If float value is detected, convert to integer and save it in the integer field
	u.IntValue = ptr.Of(cast.ToInt(floatValue))
	u.StringValue = nil
}

// OnString populates the Union with a string value from a string
func (u *Union) OnString(stringValue string) {
	// If string value is detected, convert to integer, and save it in the integer field
	if intValue, err := cast.ToIntE(stringValue); err == nil {
		u.IntValue = &intValue
		u.StringValue = nil
		return
	}
	// Fallback to string value, and save it to string field
	u.IntValue = nil
	u.StringValue = &stringValue
}

// OnBool populates the Union with an integer value from a bool
func (u *Union) OnBool(boolValue bool) {
	// If bool value is detected, convert to integer and save it in the integer field
	u.IntValue = ptr.Of(cast.ToInt(boolValue))
	u.StringValue = nil
}

// OnNull populates the Union with a null value (zero value)
func (u *Union) OnNull() {
	// If null is detected, save 0 in the integer field
	u.IntValue = ptr.Of(0)
	u.StringValue = nil
}

// OnArray populates the Union with a string value from an array
func (u *Union) OnArray(arrayValue []any) {
	// If array is detected convert to json string and save it in the string field
	u.StringValue = ptr.Of(jsonx.String(arrayValue))
	u.IntValue = nil
}

// OnObject populates the Union with a string value from an object
func (u *Union) OnObject(objectValue map[string]any) {
	// If map is detected convert to json string and save it in the string field
	u.StringValue = ptr.Of(jsonx.String(objectValue))
	u.IntValue = nil
}

// MarshalJSON marshals the Union to JSON using the provided encoder
func (u *Union) MarshalJSON() ([]byte, error) {
	switch {
	case u.IntValue != nil:
		// Integer values have priority (marshal integer value if it exists)
		return sonicx.Config.Marshal(*u.IntValue)
	case u.StringValue != nil:
		// Fallback to marshalling the string value
		return sonicx.Config.Marshal(*u.StringValue)
	default:
		// If nothing exists marshall nil
		return sonicx.Config.Marshal(nil)
	}
}

// UnmarshalJSON unmarshals JSON data into the Union using the provided decoder
func (u *Union) UnmarshalJSON(data []byte) error {
	return jsonx.HandleEntity(data, u)
}
