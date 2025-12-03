package property

import (
	"strings"

	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/stringsx"
	"github.com/r3dpixel/toolkit/symbols"
	"github.com/spf13/cast"
)

// SelectiveLogic values
const (
	SelectiveAndAny SelectiveLogic = iota // Starts at 0
	SelectiveNotAll
	SelectiveNotAny
	SelectiveAndAll

	SelectiveLogicStart   = SelectiveAndAny // SelectiveLogicStart - Start Value for SelectiveLogic (inclusive)
	SelectiveLogicEnd     = SelectiveAndAll // SelectiveLogicEnd - End Value for SelectiveLogic (inclusive)
	DefaultSelectiveLogic = SelectiveAndAny // DefaultSelectiveLogic is SelectiveAndAny
)

// SelectiveLogic represents the selective logic of a book entry
type SelectiveLogic int

// OnFloat converts the float value to an integer and sets the SelectiveLogic to the corresponding value
func (s *SelectiveLogic) OnFloat(floatValue float64) {
	*s = slParser.FromInt(cast.ToInt(floatValue))
}

// OnString converts the string value to an integer and sets the SelectiveLogic to the corresponding value
// If the conversion fails, the SelectiveLogic is set to the parsed string value (numeric values have priority over string values)
func (s *SelectiveLogic) OnString(stringValue string) {
	if intValue, err := cast.ToIntE(stringValue); err == nil {
		*s = slParser.FromInt(intValue)
		return
	}
	*s = slParser.FromString(stringValue)
}

// OnBool converts the bool value to an integer and sets the SelectiveLogic to the corresponding value
func (s *SelectiveLogic) OnBool(boolValue bool) {
	*s = slParser.FromInt(cast.ToInt(boolValue))
}

// OnNull sets the SelectiveLogic to the default value
func (s *SelectiveLogic) OnNull() {
	*s = DefaultSelectiveLogic
}

// OnArray is a no-op for SelectiveLogic, as it is not a complex type (sets default value)
func (s *SelectiveLogic) OnArray(arrayValue []any) {
	*s = DefaultSelectiveLogic
}

// OnObject is a no-op for SelectiveLogic, as it is not a complex type (sets default value)
func (s *SelectiveLogic) OnObject(objectValue map[string]any) {
	*s = DefaultSelectiveLogic
}

// MarshalJSON marshals the SelectiveLogic to JSON using Sonic
func (s *SelectiveLogic) MarshalJSON() ([]byte, error) {
	return sonicx.Config.Marshal((*int)(s))
}

// UnmarshalJSON unmarshals JSON data into the SelectiveLogic using Sonic
func (s *SelectiveLogic) UnmarshalJSON(data []byte) error {
	return jsonx.HandleEntity(data, s)
}

// SetIfPtr updates the selectivr logic if the value is not blank or nil
func (s *SelectiveLogic) SetIfPtr(value *int) {
	if value != nil {
		*s = slParser.FromInt(*value)
	}
}

// SetIfPropertyPtr updates the SelectiveLogic if the value is not blank or nil
func (s *SelectiveLogic) SetIfPropertyPtr(value *SelectiveLogic) {
	if value != nil {
		*s = *value
	}
}

// SelectiveLogicParser API to parse string/int into a valid SelectiveLogic
type SelectiveLogicParser interface {
	FromString(value string) SelectiveLogic
	FromInt(value int) SelectiveLogic
}

// selectiveLogicParser API to parse string into a valid SelectiveLogic
type selectiveLogicParser struct {
	values map[string]SelectiveLogic
}

// slParser instance of selectiveLogicParser holding the correct mappings from string to SelectiveLogic
var slParser = &selectiveLogicParser{
	values: map[string]SelectiveLogic{
		"andany": SelectiveAndAny,
		"notall": SelectiveNotAll,
		"notany": SelectiveNotAny,
		"andall": SelectiveAndAll,
	},
}

// SelectiveLogicProp returns the global SelectiveLogicParser instance
func SelectiveLogicProp() SelectiveLogicParser {
	return slParser
}

// FromString converts a string value to a SelectiveLogic after sanitization
func (sl *selectiveLogicParser) FromString(value string) SelectiveLogic {
	// Input value is a string (remove non-ASCII, remove symbols, remove whitespace, lower all characters)
	sanitizedValue := strings.ToLower(stringsx.Remove(value, symbols.NonAlphaNumericWhiteSpaceRegExp))

	// Check if the string input corresponds to any SelectiveLogic value
	if selectiveValue, exists := sl.values[sanitizedValue]; exists {
		return selectiveValue
	}

	// Return the DefaultSelectiveLogic value, otherwise
	return DefaultSelectiveLogic
}

// FromInt converts an integer value to a SelectiveLogic
func (sl *selectiveLogicParser) FromInt(value int) SelectiveLogic {
	// Check if the integer value is within the valid range
	if SelectiveLogicStart <= SelectiveLogic(value) && SelectiveLogic(value) <= SelectiveLogicEnd {
		return SelectiveLogic(value)
	}

	// Return the DefaultSelectiveLogic value, otherwise
	return DefaultSelectiveLogic
}
