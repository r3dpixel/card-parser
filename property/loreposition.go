package property

import (
	"strings"

	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/stringsx"
	"github.com/r3dpixel/toolkit/symbols"
	"github.com/spf13/cast"
)

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

type LorePosition int

func (l *LorePosition) OnFloat(floatValue float64) {
	*l = lpParser.FromInt(cast.ToInt(floatValue))
}

func (l *LorePosition) OnString(stringValue string) {
	if intValue, err := cast.ToIntE(stringValue); err == nil {
		*l = lpParser.FromInt(intValue)
		return
	}
	*l = lpParser.FromString(stringValue)
}

func (l *LorePosition) OnBool(boolValue bool) {
	*l = lpParser.FromInt(cast.ToInt(boolValue))
}

func (l *LorePosition) OnNull() {
	*l = DefaultLorePosition
}

func (l *LorePosition) OnArray(arrayValue []any) {
	*l = DefaultLorePosition
}

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
	if LorePositionStart <= LorePosition(value) && LorePosition(value) <= LorePositionEnd {
		return LorePosition(value)
	}

	// If the integer value is not a valid lore position value, return the DefaultLorePosition value
	return DefaultLorePosition
}
