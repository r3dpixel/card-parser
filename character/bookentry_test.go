package character

import (
	"encoding/json"
	"testing"

	"github.com/r3dpixel/card-parser/property"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/stringsx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyBookEntry(t *testing.T) {
	entry := DefaultBookEntry()

	// Basic structure
	assert.NotNil(t, entry)
	assert.Equal(t, property.Union{}, entry.ID)
	assert.Equal(t, []string{}, []string(entry.Keys))
	assert.Equal(t, []string{}, []string(entry.SecondaryKeys))
	assert.Equal(t, "", string(entry.Name))
	assert.Equal(t, "", string(entry.Comment))
	assert.Equal(t, "", string(entry.Content))
	assert.Equal(t, false, bool(entry.Constant))
	assert.Equal(t, false, bool(entry.Selective))
	assert.Equal(t, 10, int(entry.InsertionOrder))
	assert.Equal(t, true, bool(entry.Enabled))
	assert.Equal(t, true, bool(entry.UseRegex))
	assert.Nil(t, entry.RawExtensions)
}

func TestDefaultBookEntry(t *testing.T) {
	entry := FilledBookEntry("Test Name", "Test Content")

	assert.NotNil(t, entry)
	assert.Equal(t, []string{"Test Name"}, []string(entry.Keys))
	assert.Equal(t, "Test Content", string(entry.Content))
	assert.Equal(t, "Test Name", string(entry.Name))
	assert.Equal(t, "Test Name", string(entry.Comment))
	assert.Equal(t, true, bool(entry.Enabled))
	assert.Equal(t, 10, int(entry.InsertionOrder))

	// Should inherit all the same extension defaults as DefaultBookEntry
	assert.Equal(t, property.DefaultLorePosition, entry.Extensions.LorePosition)
	assert.Equal(t, 100.00, float64(entry.Extensions.Probability))
	assert.Equal(t, property.DefaultSelectiveLogic, entry.Extensions.SelectiveLogic)
}

func TestBookEntry_MirrorNameAndComment(t *testing.T) {
	testCases := []struct {
		name            string
		initialName     string
		initialComment  string
		expectedName    string
		expectedComment string
	}{
		{
			name:            "Comment is mirrored from Name",
			initialName:     "entry name",
			initialComment:  "",
			expectedName:    "entry name",
			expectedComment: "entry name",
		},
		{
			name:            "Name is mirrored from Comment",
			initialName:     "",
			initialComment:  "entry comment",
			expectedName:    "entry comment",
			expectedComment: "entry comment",
		},
		{
			name:            "Comment is mirrored from Name (when comment is blank)",
			initialName:     "entry name",
			initialComment:  "  ",
			expectedName:    "entry name",
			expectedComment: "entry name",
		},
		{
			name:            "Name is mirrored from Comment (when name is blank)",
			initialName:     "\t",
			initialComment:  "entry comment",
			expectedName:    "entry comment",
			expectedComment: "entry comment",
		},
		{
			name:            "No change when both exist and are different",
			initialName:     "name",
			initialComment:  "comment",
			expectedName:    "name",
			expectedComment: "comment",
		},
		{
			name:            "No change when both exist and are the same",
			initialName:     "same",
			initialComment:  "same",
			expectedName:    "same",
			expectedComment: "same",
		},
		{
			name:            "No change when both are empty",
			initialName:     "",
			initialComment:  "",
			expectedName:    "",
			expectedComment: "",
		},
		{
			name:            "No change when both are blank",
			initialName:     "  ",
			initialComment:  "\t\n",
			expectedName:    "  ",
			expectedComment: "\t\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry := &BookEntry{
				BookEntryCore: BookEntryCore{
					Name:    property.String(tc.initialName),
					Comment: property.String(tc.initialComment),
				},
			}

			entry.MirrorNameAndComment()

			assert.Equal(t, tc.expectedName, string(entry.Name))
			assert.Equal(t, tc.expectedComment, string(entry.Comment))
		})
	}
}

func TestBookEntry_UnmarshalJSON(t *testing.T) {
	t.Run("Basic unmarshal with no extensions", func(t *testing.T) {
		jsonData := `{
			"id": 123,
			"keys": ["test"],
			"name": "Test Entry",
			"content": "Test content",
			"enabled": true
		}`

		var entry BookEntry
		err := sonicx.Config.UnmarshalFromString(jsonData, &entry)

		require.NoError(t, err)
		assert.Equal(t, "Test Entry", string(entry.Name))
		assert.Equal(t, "Test content", string(entry.Content))
		assert.Equal(t, true, bool(entry.Enabled))
	})

	t.Run("Unmarshal with extensions", func(t *testing.T) {
		jsonData := `{
			"id": 456,
			"keys": ["test"],
			"name": "Test Entry",
			"extensions": {
				"position": 1,
				"probability": 75.5,
				"depth": 5,
				"selectiveLogic": 2,
				"match_whole_words": true,
				"case_sensitive": false,
				"role": 1,
				"sticky": 2,
				"cooldown": 10,
				"delay": 5,
				"unknown_field": "should remain"
			}
		}`

		var entry BookEntry
		err := sonicx.Config.UnmarshalFromString(jsonData, &entry)

		require.NoError(t, err)
		extensions := entry.Extensions
		// Check that known extensions were moved to struct fields
		assert.Equal(t, 1, int(extensions.LorePosition))
		assert.Equal(t, 75.5, float64(extensions.Probability))
		assert.Equal(t, 5, int(extensions.Depth))
		assert.Equal(t, 2, int(extensions.SelectiveLogic))
		assert.Equal(t, true, bool(extensions.MatchWholeWords))
		assert.Equal(t, false, bool(extensions.CaseSensitive))
		assert.Equal(t, 1, int(extensions.Role))
		assert.Equal(t, 2, int(extensions.Sticky))
		assert.Equal(t, 10, int(extensions.Cooldown))
		assert.Equal(t, 5, int(extensions.Delay))

		// Check that known extensions were removed from map
		assert.NotContains(t, entry.RawExtensions, EntryPosition)
		assert.NotContains(t, entry.RawExtensions, EntryProbability)
		assert.NotContains(t, entry.RawExtensions, EntryDepth)
		assert.NotContains(t, entry.RawExtensions, EntrySelectiveLogic)
		assert.NotContains(t, entry.RawExtensions, EntryMatchWholeWords)
		assert.NotContains(t, entry.RawExtensions, EntryCaseSensitive)
		assert.NotContains(t, entry.RawExtensions, EntryRole)
		assert.NotContains(t, entry.RawExtensions, EntrySticky)
		assert.NotContains(t, entry.RawExtensions, EntryCooldown)
		assert.NotContains(t, entry.RawExtensions, EntryDelay)

		// Check that unknown extensions remain in map
		assert.Contains(t, entry.RawExtensions, "unknown_field")
		assert.Equal(t, "should remain", entry.RawExtensions["unknown_field"])
	})

	t.Run("Unmarshal with empty extensions", func(t *testing.T) {
		jsonData := `{
			"id": 789,
			"keys": ["test"],
			"extensions": {}
		}`

		var entry BookEntry
		err := sonicx.Config.UnmarshalFromString(jsonData, &entry)

		require.NoError(t, err)
		assert.NotNil(t, entry.RawExtensions)
		assert.Len(t, entry.RawExtensions, 0)
	})

	t.Run("Unmarshal with straggler extensions", func(t *testing.T) {
		jsonData := `{
			"selectiveLogic": "NOT_ANY",
			"position": 3,
			"probability": 13.05,
			"case_sensitive": true
		}`

		var entry BookEntry
		err := sonicx.Config.UnmarshalFromString(jsonData, &entry)
		require.NoError(t, err)

		assert.Equal(t, 2, int(entry.Extensions.SelectiveLogic))
		assert.Equal(t, 3, int(entry.Extensions.LorePosition))
		assert.Equal(t, 13.05, float64(entry.Extensions.Probability))
		assert.Equal(t, true, bool(entry.Extensions.CaseSensitive))
	})

	t.Run("Unmarshal with straggler and normal extensions", func(t *testing.T) {
		jsonData := `{
			"selectiveLogic": "NOT_ANY",
			"position": 3,
			"probability": 13.05,
			"case_sensitive": true,
			"extensions": {
				"selectiveLogic": "and__all",
				"position": 5,
				"probability": 98.05,
				"case_sensitive": false
			}
		}`

		var entry BookEntry
		err := sonicx.Config.UnmarshalFromString(jsonData, &entry)
		require.NoError(t, err)

		assert.Equal(t, 3, int(entry.Extensions.SelectiveLogic))
		assert.Equal(t, 5, int(entry.Extensions.LorePosition))
		assert.Equal(t, 98.05, float64(entry.Extensions.Probability))
		assert.Equal(t, false, bool(entry.Extensions.CaseSensitive))
	})

	t.Run("Unmarshal with null extensions", func(t *testing.T) {
		jsonData := `{
			"id": 999,
			"keys": ["test"],
			"extensions": null
		}`

		var entry BookEntry
		err := sonicx.Config.UnmarshalFromString(jsonData, &entry)

		require.NoError(t, err)
		// RawExtensions should be nil, not processed
		assert.Nil(t, entry.RawExtensions)
	})

	t.Run("Unmarshal with invalid JSON", func(t *testing.T) {
		jsonData := `{invalid json`

		var entry BookEntry
		err := sonicx.Config.UnmarshalFromString(jsonData, &entry)

		assert.Error(t, err)
	})

	t.Run("Unmarshal with extension type conversion", func(t *testing.T) {
		jsonData := `{
			"id": 111,
			"keys": ["test"],
			"extensions": {
				"position": "2",
				"probability": "85",
				"depth": true,
				"match_whole_words": "true",
				"case_sensitive": 1
			}
		}`

		var entry BookEntry
		err := sonicx.Config.UnmarshalFromString(jsonData, &entry)

		require.NoError(t, err)
		extensions := entry.Extensions
		// Check type conversions work through Entity interface
		assert.Equal(t, 2, int(extensions.LorePosition))
		assert.Equal(t, 85.00, float64(extensions.Probability))
		assert.Equal(t, 1, int(extensions.Depth)) // true -> 1
		assert.Equal(t, true, bool(extensions.MatchWholeWords))
		assert.Equal(t, true, bool(extensions.CaseSensitive)) // 1 -> true
	})
}

func TestBookEntry_MarshalJSON(t *testing.T) {
	t.Run("Basic marshal with no extensions", func(t *testing.T) {
		entry := &BookEntry{
			BookEntryCore: BookEntryCore{
				ID:      property.Union{},
				Keys:    property.StringArray{"test"},
				Name:    "Test Entry",
				Content: "Test content",
				Enabled: true,
			},
		}

		data, err := sonicx.Config.Marshal(entry)
		require.NoError(t, err)

		var result map[string]any
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		assert.Equal(t, "Test Entry", result["name"])
		assert.Equal(t, "Test content", result["content"])
		assert.Equal(t, true, result["enabled"])

		// RawExtensions should be present with struct field values
		extensions, ok := result["extensions"].(map[string]any)
		assert.True(t, ok)
		assert.NotNil(t, extensions)
	})

	t.Run("Marshal with existing extensions", func(t *testing.T) {
		entry := &BookEntry{
			BookEntryCore: BookEntryCore{
				ID:   property.Union{},
				Keys: property.StringArray{"test"},
			},
			RawExtensions: map[string]any{
				"custom_field": "custom_value",
			},
			Extensions: BookEntryExtensions{
				LorePosition: property.LorePosition(3),
				Probability:  90.5,
				Depth:        7,
			},
		}

		data, err := sonicx.Config.Marshal(entry)
		require.NoError(t, err)

		var result map[string]any
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		extensions, ok := result["extensions"].(map[string]any)
		assert.True(t, ok)

		// Check that custom extensions are preserved
		assert.Equal(t, "custom_value", extensions["custom_field"])

		// Check that struct fields are added to extensions
		assert.Contains(t, extensions, EntryPosition)
		assert.Contains(t, extensions, EntryProbability)
		assert.Contains(t, extensions, EntryDepth)
	})

	t.Run("Marshal idempotent - original not modified", func(t *testing.T) {
		originalExtensions := map[string]any{
			"original_field": "original_value",
		}

		entry := &BookEntry{
			BookEntryCore: BookEntryCore{
				ID:   property.Union{},
				Keys: property.StringArray{"test"},
			},
			RawExtensions: originalExtensions,
		}

		data, err := sonicx.Config.Marshal(entry)
		require.NoError(t, err)

		// Check that original extensions map is unchanged
		// (We verify this by checking the content hasn't been modified)
		assert.Len(t, entry.RawExtensions, 1)
		assert.Equal(t, "original_value", entry.RawExtensions["original_field"])
		assert.NotContains(t, entry.RawExtensions, EntryPosition) // Should not be added to original

		// But marshaled data should contain all extensions
		var result map[string]any
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		extensions := result["extensions"].(map[string]any)
		assert.Contains(t, extensions, "original_field")
		assert.Contains(t, extensions, EntryPosition)
	})
}

func TestBookEntry_MarshalUnmarshalRoundTrip(t *testing.T) {
	t.Run("Round trip preserves all data", func(t *testing.T) {
		original := &BookEntry{
			BookEntryCore: BookEntryCore{
				ID:             property.Union{IntValue: &[]int{42}[0]},
				Keys:           property.StringArray{"key1", "key2"},
				SecondaryKeys:  property.StringArray{"sec1"},
				Name:           "Test Name",
				Comment:        "Test Comment",
				Content:        "Test Content",
				Constant:       true,
				Selective:      false,
				InsertionOrder: 15,
				Enabled:        true,
				UseRegex:       false,
			},
			RawExtensions: map[string]any{
				"custom_field": "custom_value",
				"number_field": 123,
			},
			Extensions: BookEntryExtensions{
				LorePosition:    property.LorePosition(2),
				Probability:     88.8,
				Depth:           6,
				SelectiveLogic:  property.SelectiveLogic(1),
				MatchWholeWords: true,
				CaseSensitive:   false,
				Role:            2,
				Sticky:          4,
				Cooldown:        20,
				Delay:           10,
			},
		}

		// Marshal
		data, err := sonicx.Config.Marshal(original)
		require.NoError(t, err)

		// Unmarshal
		var restored BookEntry
		err = sonicx.Config.UnmarshalFromString(stringsx.FromBytes(data), &restored)
		require.NoError(t, err)

		// Compare all fields
		assert.Equal(t, original.Name, restored.Name)
		assert.Equal(t, original.Comment, restored.Comment)
		assert.Equal(t, original.Content, restored.Content)
		assert.Equal(t, original.Constant, restored.Constant)
		assert.Equal(t, original.Selective, restored.Selective)
		assert.Equal(t, original.InsertionOrder, restored.InsertionOrder)
		assert.Equal(t, original.Enabled, restored.Enabled)
		assert.Equal(t, original.UseRegex, restored.UseRegex)

		assertBookEntryExtensions(t, original.Extensions, restored.Extensions)

		assert.Contains(t, restored.RawExtensions, "custom_field")
		assert.Equal(t, "custom_value", restored.RawExtensions["custom_field"])
		assert.Contains(t, restored.RawExtensions, "number_field")
		assert.Equal(t, float64(123), restored.RawExtensions["number_field"])
	})
}

func TestExtensionConstants(t *testing.T) {
	t.Run("Constants match JSON tags", func(t *testing.T) {
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
	})
}

func TestBookEntry_EdgeCases(t *testing.T) {
	t.Run("Empty JSON object", func(t *testing.T) {
		jsonData := `{}`

		var entry BookEntry
		err := sonicx.Config.UnmarshalFromString(jsonData, &entry)

		require.NoError(t, err)
		assert.Equal(t, "", string(entry.Name))
		assert.Equal(t, true, bool(entry.Enabled))
	})

	t.Run("Partial extension overlap", func(t *testing.T) {
		jsonData := `{
			"extensions": {
				"position": 1,
				"unknown": "value"
			}
		}`

		var entry BookEntry
		err := sonicx.Config.UnmarshalFromString(jsonData, &entry)

		require.NoError(t, err)
		assert.Equal(t, 1, int(entry.Extensions.LorePosition))
		assert.Contains(t, entry.RawExtensions, "unknown")
		assert.NotContains(t, entry.RawExtensions, EntryPosition)
	})

	t.Run("RawExtensions with null values", func(t *testing.T) {
		jsonData := `{
			"extensions": {
				"position": null,
				"probability": null
			}
		}`

		var entry BookEntry
		err := sonicx.Config.UnmarshalFromString(jsonData, &entry)

		require.NoError(t, err)
		assert.NotContains(t, entry.RawExtensions, EntryPosition)
		assert.NotContains(t, entry.RawExtensions, EntryProbability)
	})
}
