package character

import (
	"testing"

	"github.com/r3dpixel/card-parser/property"
	"github.com/r3dpixel/toolkit/ptr"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/stretchr/testify/assert"
)

func TestBookEntryExtensions_Constants(t *testing.T) {
	assert.Equal(t, "position", EntryPosition)
	assert.Equal(t, "probability", EntryProbability)
	assert.Equal(t, "depth", EntryDepth)
	assert.Equal(t, "selectiveLogic", EntrySelectiveLogic)
	assert.Equal(t, "match_whole_words", EntryMatchWholeWords)
	assert.Equal(t, "case_sensitive", EntryCaseSensitive)
	assert.Equal(t, "role", EntryRole)
	assert.Equal(t, "sticky", EntrySticky)
	assert.Equal(t, "cooldown", EntryCooldown)
	assert.Equal(t, "delay", EntryDelay)
}

func TestBookEntryExtensions_DefaultMissing(t *testing.T) {
	type Container struct {
		Role           property.Role
		SelectiveLogic property.SelectiveLogic
		LorePosition   *property.LorePosition
		Probability    *property.Float
	}
	container := Container{}

	err := sonicx.Config.UnmarshalFromString(`{"Role":1}`, &container)
	assert.NoError(t, err, "unmarshaling should not produce an error")
	assert.Equal(t, property.UserRole, container.Role, "Role should be set to UserRole")
	assert.Equal(t, property.SelectiveAndAny, container.SelectiveLogic, "SelectiveLogic should default to SelectiveAndAny")
	assert.Equal(t, (*property.LorePosition)(nil), container.LorePosition, "LorePosition should remain nil")
	assert.Equal(t, (*property.Float)(nil), container.Probability, "Probability should remain nil")

	result := Container{
		Role:           property.AssistantRole,
		SelectiveLogic: property.SelectiveNotAll,
		LorePosition:   ptr.Of(property.AfterExampleMessages),
		Probability:    ptr.Of(property.Float(87.50)),
	}
	result.LorePosition.SetIfPropertyPtr(container.LorePosition)
	result.Probability.SetIfPropertyPtr(container.Probability)
	assert.Equal(t, property.AfterExampleMessages, *result.LorePosition, "LorePosition should be preserved when not set")
	assert.Equal(t, 87.50, float64(*result.Probability), "Probability should be preserved when not set")
}

func TestBookEntryExtensions_Default(t *testing.T) {
	defaults := DefaultBookEntryExtensions()

	assert.Equal(t, property.DefaultLorePosition, defaults.LorePosition)
	assert.Equal(t, 100.00, float64(defaults.Probability))
	assert.Equal(t, DefaultDepth, int(defaults.Depth))
	assert.Equal(t, property.DefaultSelectiveLogic, defaults.SelectiveLogic)
	assert.Equal(t, false, bool(defaults.MatchWholeWords))
	assert.Equal(t, false, bool(defaults.CaseSensitive))
	assert.Equal(t, 0, int(defaults.Role))
	assert.Equal(t, 0, int(defaults.Sticky))
	assert.Equal(t, 0, int(defaults.Cooldown))
	assert.Equal(t, 0, int(defaults.Delay))
}

// assertBookEntryExtensions is a helper function that asserts BookEntryExtensions values
// against expected primitive values, handling the property-to-primitive conversion.
func assertBookEntryExtensions(t *testing.T, expected BookEntryExtensions, actual BookEntryExtensions) {
	assert.Equal(t, int(expected.LorePosition), int(actual.LorePosition))
	assert.Equal(t, float64(expected.Probability), float64(actual.Probability))
	assert.Equal(t, int(expected.Depth), int(actual.Depth))
	assert.Equal(t, int(expected.SelectiveLogic), int(actual.SelectiveLogic))
	assert.Equal(t, bool(expected.MatchWholeWords), bool(actual.MatchWholeWords))
	assert.Equal(t, bool(expected.CaseSensitive), bool(actual.CaseSensitive))
	assert.Equal(t, int(expected.Role), int(actual.Role))
	assert.Equal(t, int(expected.Sticky), int(actual.Sticky))
	assert.Equal(t, int(expected.Cooldown), int(actual.Cooldown))
	assert.Equal(t, int(expected.Delay), int(actual.Delay))
}

// assertBookEntryExtensionsFromMap is a helper function that asserts BookEntryExtensions values
// from a map[string]any, handling the property-to-primitive conversion.
func assertBookEntryExtensionsFromMap(
	t *testing.T,
	assertFunc func(t assert.TestingT, expected, actual any, msgAndArgs ...any) bool,
	expected BookEntryExtensions,
	actualMap map[string]any,
) {
	assertFunc(t, int(expected.LorePosition), int(actualMap[EntryPosition].(property.LorePosition)))
	assertFunc(t, float64(expected.Probability), float64(actualMap[EntryProbability].(property.Float)))
	assertFunc(t, int(expected.Depth), int(actualMap[EntryDepth].(property.Integer)))
	assertFunc(t, int(expected.SelectiveLogic), int(actualMap[EntrySelectiveLogic].(property.SelectiveLogic)))
	assertFunc(t, bool(expected.MatchWholeWords), bool(actualMap[EntryMatchWholeWords].(property.Bool)))
	assertFunc(t, bool(expected.CaseSensitive), bool(actualMap[EntryCaseSensitive].(property.Bool)))
	assertFunc(t, int(expected.Role), int(actualMap[EntryRole].(property.Role)))
	assertFunc(t, int(expected.Sticky), int(actualMap[EntrySticky].(property.Integer)))
	assertFunc(t, int(expected.Cooldown), int(actualMap[EntryCooldown].(property.Integer)))
	assertFunc(t, int(expected.Delay), int(actualMap[EntryDelay].(property.Integer)))
}
