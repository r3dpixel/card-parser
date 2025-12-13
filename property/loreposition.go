package property

import (
	"strings"

	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/stringsx"
	"github.com/r3dpixel/toolkit/symbols"
	"github.com/spf13/cast"
)

// LorePosition constants
const (
	BeforeCharPosition LorePosition = iota
	AfterCharPosition
	BeforeAuthorNotes
	AfterAuthorNotes
	AtDepth
	BeforeExampleMessages
	AfterExampleMessages

	LorePositionStart   = BeforeCharPosition
	LorePositionEnd     = AfterExampleMessages
	DefaultLorePosition = BeforeCharPosition // DefaultLorePosition is BeforeCharPosition
)

// LorePosition represents the position of a book entry in the Lorebook
type LorePosition int

// OnFloat converts the float value to an integer and sets the LorePosition to the corresponding value
func (l *LorePosition) OnFloat(floatValue float64) {
	*l = lpParser.FromInt(cast.ToInt(floatValue))
}

// OnString converts the string value to an integer and sets the LorePosition to the corresponding value
// If the conversion fails, the LorePosition is set to the parsed string value (numeric values have priority over string values)
func (l *LorePosition) OnString(stringValue string) {
	if intValue, err := cast.ToIntE(stringValue); err == nil {
		*l = lpParser.FromInt(intValue)
		return
	}
	*l = lpParser.FromString(stringValue)
}

// OnBool converts the bool value to an integer and sets the LorePosition to the corresponding value
func (l *LorePosition) OnBool(boolValue bool) {
	*l = lpParser.FromInt(cast.ToInt(boolValue))
}

// OnNull sets the LorePosition to the default value
func (l *LorePosition) OnNull() {
	*l = DefaultLorePosition
}

// OnArray is a no-op for LorePosition, as it is not a complex type (sets default value)
func (l *LorePosition) OnArray(arrayValue []any) {
	*l = DefaultLorePosition
}

// OnObject is a no-op for LorePosition, as it is not a complex type (sets default value)
func (l *LorePosition) OnObject(objectValue map[string]any) {
	*l = DefaultLorePosition
}

// MarshalJSON marshals the LorePosition to JSON using Sonic
func (l *LorePosition) MarshalJSON() ([]byte, error) {
	return sonicx.Config.Marshal((*int)(l))
}

// UnmarshalJSON unmarshals JSON data into the LorePosition using Sonic
func (l *LorePosition) UnmarshalJSON(data []byte) error {
	return jsonx.HandleEntity(data, l)
}

// SetIfPtr updates the LorePosition if the value is not nil
func (l *LorePosition) SetIfPtr(value *int) {
	if value != nil {
		*l = lpParser.FromInt(*value)
	}
}

// SetIfPropertyPtr updates the LorePosition if the value is not nil
func (l *LorePosition) SetIfPropertyPtr(value *LorePosition) {
	if value != nil {
		*l = *value
	}
}

// LorePositionParser API to parse string/int into a valid LorePosition
type LorePositionParser interface {
	FromString(value string) LorePosition
	FromInt(value int) LorePosition
}

// lorePositionParser API to parse string into a valid LorePosition
type lorePositionParser struct {
	strs map[string]LorePosition
}

// lpParser instance of lorePositionParser holding the correct mappings from int to LorePosition
var lpParser = &lorePositionParser{
	strs: map[string]LorePosition{
		"beforechar":              BeforeCharPosition,
		"lorebookentrybeforechar": BeforeCharPosition,
		"afterchar":               AfterCharPosition,
		"lorebookentryafterchar":  AfterCharPosition,
		"beforean":                BeforeAuthorNotes,
		"lorebookentrybeforean":   BeforeAuthorNotes,
		"afteran":                 AfterAuthorNotes,
		"lorebookentryafteran":    AfterAuthorNotes,
		"atdepth":                 AtDepth,
		"lorebookentrydepth":      AtDepth,
		"inchat":                  AtDepth,
		"beforeem":                BeforeExampleMessages,
		"lorebookentrybeforeem":   BeforeExampleMessages,
		"afterem":                 AfterExampleMessages,
		"lorebookentryafterem":    AfterExampleMessages,
	},
}

// LorePositionProp returns the global LorePositionParser instance
func LorePositionProp() LorePositionParser {
	return lpParser
}

// FromString converts a string value to a LorePosition after sanitization
func (lp *lorePositionParser) FromString(value string) LorePosition {
	sanitizedValue := strings.ToLower(stringsx.Remove(value, symbols.NonAlphaNumericWhiteSpaceRegExp))

	// Check string sets for before position
	if position, ok := lp.strs[sanitizedValue]; ok {
		return position
	}

	return DefaultLorePosition
}

// FromInt converts an integer value to a LorePosition
func (lp *lorePositionParser) FromInt(value int) LorePosition {
	// Check if the integer value is within the valid range
	if LorePositionStart <= LorePosition(value) && LorePosition(value) <= LorePositionEnd {
		return LorePosition(value)
	}

	// Return the DefaultLorePosition value, otherwise
	return DefaultLorePosition
}
