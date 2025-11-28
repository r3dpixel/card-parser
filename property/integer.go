package property

import (
	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/spf13/cast"
)

// Integer represents an integer value
type Integer int

// OnValue populates the Integer with the given value (converts to int if possible)
// NOTE: The original value is preserved, if input cannot be converted to int
func (i *Integer) OnValue(value any) {
	if intValue, err := cast.ToIntE(value); err == nil {
		*i = Integer(intValue)
	}
}

// OnNull no-op for Integer, as it cannot be null
// NOTE: The original value is preserved
func (i *Integer) OnNull() {}

// OnComplex is a no-op for Integer, as it is not a complex type
// NOTE: The original value is preserved
func (i *Integer) OnComplex(complex any) {}

// MarshalJSON marshals the Integer to JSON using Sonic
func (i *Integer) MarshalJSON() ([]byte, error) {
	return sonicx.Config.Marshal((*int)(i))
}

// UnmarshalJSON unmarshals JSON data into the Integer using Sonic
func (i *Integer) UnmarshalJSON(data []byte) error {
	return jsonx.HandlePrimitive(data, i)
}

// SetIfPtr updates the Integer if the value is not nil
func (i *Integer) SetIfPtr(value *int) {
	if value != nil {
		*i = Integer(*value)
	}
}

// SetIfPropertyPtr updates the Integer if the value is not nil
func (i *Integer) SetIfPropertyPtr(value *int) {
	if value != nil {
		*i = Integer(*value)
	}
}
