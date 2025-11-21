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
	SystemRole Role = iota
	UserRole
	AssistantRole
	RoleStart   = SystemRole
	RoleEnd     = AssistantRole
	DefaultRole = SystemRole
)

type Role int

func (r *Role) OnFloat(floatValue float64) {
	*r = rlParser.FromInt(cast.ToInt(floatValue))
}

func (r *Role) OnString(stringValue string) {
	if intValue, err := cast.ToIntE(stringValue); err == nil {
		*r = rlParser.FromInt(intValue)
		return
	}
	*r = rlParser.FromString(stringValue)
}

func (r *Role) OnBool(boolValue bool) {
	*r = rlParser.FromInt(cast.ToInt(boolValue))
}

func (r *Role) OnNull() {
	*r = DefaultRole
}

func (r *Role) OnArray(arrayValue []any) {
	*r = DefaultRole
}

func (r *Role) OnObject(objectValue map[string]any) {
	*r = DefaultRole
}

// MarshalJSON marshals the Role to JSON using Sonic
func (r *Role) MarshalJSON() ([]byte, error) {
	return sonicx.Config.Marshal((*int)(r))
}

// UnmarshalJSON unmarshals JSON data into the Role using Sonic
func (r *Role) UnmarshalJSON(data []byte) error {
	return jsonx.HandleEntity(data, r)
}

// SetIfPtr updates the role if the value is not blank or nil
func (r *Role) SetIfPtr(value *int) {
	if value != nil {
		*r = rlParser.FromInt(*value)
	}
}

// SetIfPropertyPtr updates the Role if the value is not blank or nil
func (r *Role) SetIfPropertyPtr(value *Role) {
	if value != nil {
		*r = *value
	}
}

// RoleParser API to parse string/int into a valid Role
type RoleParser interface {
	FromString(value string) Role
	FromInt(value int) Role
}

type roleParser struct {
	values map[string]Role
}

// rlParser instance of roleParser holding the correct mappings from string to SelectiveLogic
var rlParser = &roleParser{
	values: map[string]Role{
		"system":                 SystemRole,
		"lorebookdepthsystem":    SystemRole,
		"user":                   UserRole,
		"lorebookdepthuser":      UserRole,
		"assistant":              AssistantRole,
		"lorebookdepthchar":      AssistantRole,
		"lorebookdepthassistant": AssistantRole,
	},
}

// RoleProp returns the global RoleParser instance
func RoleProp() RoleParser {
	return rlParser
}

// FromString converts a string value to a Role after sanitization
func (rl *roleParser) FromString(value string) Role {
	// Input value is a string (remove non-ASCII, remove symbols, remove whitespace, lower all characters)
	sanitizedValue := strings.ToLower(stringsx.Remove(value, symbols.NonAlphaNumericWhiteSpaceRegExp))

	// Check if the string input corresponds to any Role value
	if role, exists := rl.values[sanitizedValue]; exists {
		return role
	}
	// Return the DefaultRole value
	return DefaultRole
}

// FromInt converts an integer value to a SelectiveLogic
func (rl *roleParser) FromInt(value int) Role {
	if RoleStart <= Role(value) && Role(value) <= RoleEnd {
		return Role(value)
	}
	return DefaultRole
}
