package property

import (
	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/spf13/cast"
)

type Integer int

func (i *Integer) OnValue(value any) {
	if intValue, err := cast.ToIntE(value); err == nil {
		*i = Integer(intValue)
	}
}

func (i *Integer) OnNull() {}

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
