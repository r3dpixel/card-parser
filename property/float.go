package property

import (
	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/spf13/cast"
)

type Float float64

func (f *Float) OnValue(value any) {
	if floatValue, err := cast.ToFloat64E(value); err == nil {
		*f = Float(floatValue)
	}
}

func (f *Float) OnNull() {}

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
