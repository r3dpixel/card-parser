package property

import (
	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/spf13/cast"
)

type Bool bool

func (b *Bool) OnValue(value any) {
	if boolValue, err := cast.ToBoolE(value); err == nil {
		*b = Bool(boolValue)
	}
}
func (b *Bool) OnNull()               {}
func (b *Bool) OnComplex(complex any) {}

// MarshalJSON marshals the Bool to JSON using Sonic
func (b *Bool) MarshalJSON() ([]byte, error) {
	return sonicx.Config.Marshal((*bool)(b))
}

// UnmarshalJSON unmarshals JSON data into the Bool using Sonic
func (b *Bool) UnmarshalJSON(data []byte) error {
	return jsonx.HandlePrimitive(data, b)
}

// SetIfPtr updates the Bool if the value is not nil
func (b *Bool) SetIfPtr(value *bool) {
	if value != nil {
		*b = Bool(*value)
	}
}

// SetIfPropertyPtr updates the Bool if the value is not nil
func (b *Bool) SetIfPropertyPtr(value *Bool) {
	if value != nil {
		*b = *value
	}
}
