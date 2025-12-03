package property

import (
	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/spf13/cast"
)

// Bool represents a boolean value
type Bool bool

// OnValue populates the Bool with the given value (converts to bool if possible)
// NOTE: The original value is preserved, if input cannot be converted to bool
func (b *Bool) OnValue(value any) {
	if boolValue, err := cast.ToBoolE(value); err == nil {
		*b = Bool(boolValue)
	}
}

// OnNull no-op for Bool, as it cannot be null
// NOTE: The original value is preserved
func (b *Bool) OnNull() {}

// OnComplex is a no-op for Bool, as it is not a complex type
// NOTE: The original value is preserved
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
