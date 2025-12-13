package property

import (
	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/stringsx"
	"github.com/spf13/cast"
)

// String represents a string value
type String string

// NormalizeSymbols normalizes the symbols in the String
func (s *String) NormalizeSymbols() {
	*s = String(stringsx.NormalizeSymbols(string(*s)))
}

// OnValue populates the String with the value converted to a string
func (s *String) OnValue(value any) {
	if stringValue, err := cast.ToStringE(value); err == nil {
		*s = String(stringValue)
	}
}

// OnNull populates the String with an empty string
func (s *String) OnNull() {
	*s = String("")
}

// OnComplex populates the String with the JSON representation of the complex value
func (s *String) OnComplex(complex any) {
	*s = String(jsonx.String(complex))
}

// MarshalJSON marshals the String to JSON using Sonic
func (s *String) MarshalJSON() ([]byte, error) {
	return sonicx.Config.Marshal((*string)(s))
}

// UnmarshalJSON unmarshals JSON data into the String using Sonic
func (s *String) UnmarshalJSON(data []byte) error {
	return jsonx.HandlePrimitive(data, s)
}

// SetIf updates the String if the value is not blank
func (s *String) SetIf(value string) {
	if stringsx.IsNotBlank(value) {
		*s = String(value)
	}
}

// SetIfPtr updates the String if the value is not blank or nil
func (s *String) SetIfPtr(value *string) {
	if stringsx.IsNotBlankPtr(value) {
		*s = String(*value)
	}
}

// SetIfProperty updates the String if the value is not blank
func (s *String) SetIfProperty(value String) {
	if stringsx.IsNotBlank(string(value)) {
		*s = value
	}
}

// SetIfPropertyPtr updates the String if the value is not blank or nil
func (s *String) SetIfPropertyPtr(value *String) {
	if stringsx.IsNotBlankPtr((*string)(value)) {
		*s = *value
	}
}
