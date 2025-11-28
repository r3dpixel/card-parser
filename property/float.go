package property

import (
	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/spf13/cast"
)

// Float represents a float value
type Float float64

// OnValue populates the Float with the given value (converts to float64 if possible)
// NOTE: The original value is preserved, if input cannot be converted to float64
func (f *Float) OnValue(value any) {
	if floatValue, err := cast.ToFloat64E(value); err == nil {
		*f = Float(floatValue)
	}
}

// OnNull no-op for Float, as it cannot be null
// NOTE: The original value is preserved
func (f *Float) OnNull() {}

// OnComplex is a no-op for Float, as it is not a complex type
// NOTE: The original value is preserved
func (f *Float) OnComplex(complex any) {}

// MarshalJSON marshals the Float to JSON using Sonic
func (f *Float) MarshalJSON() ([]byte, error) {
	return sonicx.Config.Marshal((*float64)(f))
}

// UnmarshalJSON unmarshals JSON data into the Float using Sonic
func (f *Float) UnmarshalJSON(data []byte) error {
	return jsonx.HandlePrimitive(data, f)
}

// SetIfPtr updates the Float if the value is not nil
func (f *Float) SetIfPtr(value *float64) {
	if value != nil {
		*f = Float(*value)
	}
}

// SetIfPropertyPtr updates the Float if the value is not nil
func (f *Float) SetIfPropertyPtr(value *Float) {
	if value != nil {
		*f = *value
	}
}
