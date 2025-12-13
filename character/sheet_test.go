package character

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/r3dpixel/card-parser/property"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/timestamp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultSheet(t *testing.T) {
	tests := []struct {
		name     string
		revision Revision
		expected Stamp
	}{
		{
			name:     "create V2 sheet",
			revision: RevisionV2,
			expected: Stamps[RevisionV2],
		},
		{
			name:     "create V3 sheet",
			revision: RevisionV3,
			expected: Stamps[RevisionV3],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sheet := DefaultSheet(tt.revision)

			assert.NotNil(t, sheet)
			assert.Equal(t, tt.expected.Revision, sheet.Revision)
			assert.Equal(t, tt.expected.Spec, sheet.Spec)
			assert.Equal(t, tt.expected.Version, sheet.Version)
		})
	}
}

func TestSheet_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		sheet    *Sheet
		contains []string
	}{
		{
			name: "marshal V3 sheet with content",
			sheet: &Sheet{
				Spec:    SpecV3,
				Version: V3,
				Content: Content{
					Title:       property.String("Test Character"),
					Name:        property.String("TestChar"),
					Description: property.String("A test character"),
					Creator:     property.String("Test Creator"),
				},
			},
			contains: []string{
				`"spec":"chara_card_v3"`,
				`"spec_version":"3.0"`,
				`"data":{`,
				`"title":"Test Character"`,
				`"name":"TestChar"`,
			},
		},
		{
			name: "marshal V2 sheet",
			sheet: &Sheet{
				Spec:    SpecV2,
				Version: V2,
				Content: Content{
					Title: property.String("V2 Character"),
					Name:  property.String("V2Char"),
				},
			},
			contains: []string{
				`"spec":"chara_card_v2"`,
				`"spec_version":"2.0"`,
				`"title":"V2 Character"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := sonicx.Config.Marshal(tt.sheet)
			require.NoError(t, err)

			result := string(data)
			for _, expected := range tt.contains {
				assert.Contains(t, result, expected)
			}
		})
	}
}

func TestSheet_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected func(t *testing.T, sheet *Sheet)
	}{
		{
			name:     "unmarshal V3 sheet",
			jsonData: `{"spec":"chara_card_v3","spec_version":"3.0","data":{"title":"Test Character","name":"TestChar","description":"A test character"}}`,
			expected: func(t *testing.T, sheet *Sheet) {
				assert.Equal(t, SpecV3, sheet.Spec)
				assert.Equal(t, V3, sheet.Version)
				assert.Equal(t, RevisionV3, sheet.Revision)
				assert.Equal(t, "Test Character", string(sheet.Content.Title))
				assert.Equal(t, "TestChar", string(sheet.Content.Name))
				assert.Equal(t, "A test character", string(sheet.Content.Description))
			},
		},
		{
			name:     "unmarshal V2 sheet",
			jsonData: `{"spec":"chara_card_v2","spec_version":"2.0","data":{"title":"V2 Character","name":"V2Char"}}`,
			expected: func(t *testing.T, sheet *Sheet) {
				assert.Equal(t, SpecV2, sheet.Spec)
				assert.Equal(t, V2, sheet.Version)
				assert.Equal(t, RevisionV2, sheet.Revision)
				assert.Equal(t, "V2 Character", string(sheet.Content.Title))
				assert.Equal(t, "V2Char", string(sheet.Content.Name))
			},
		},
		{
			name:     "unmarshal with alternate greetings from string",
			jsonData: `{"spec":"chara_card_v3","spec_version":"3.0","data":{"title":"Test","name":"Test","alternate_greetings":"single greeting"}}`,
			expected: func(t *testing.T, sheet *Sheet) {
				assert.Equal(t, SpecV3, sheet.Spec)
				assert.Len(t, sheet.Content.AlternateGreetings, 1)
				assert.Equal(t, "single greeting", sheet.Content.AlternateGreetings[0])
			},
		},
		{
			name:     "unmarshal with alternate greetings from array",
			jsonData: `{"spec":"chara_card_v3","spec_version":"3.0","data":{"title":"Test","name":"Test","alternate_greetings":["greeting1","greeting2","greeting3"]}}`,
			expected: func(t *testing.T, sheet *Sheet) {
				assert.Equal(t, SpecV3, sheet.Spec)
				assert.Len(t, sheet.Content.AlternateGreetings, 3)
				assert.Equal(t, []string{"greeting1", "greeting2", "greeting3"}, []string(sheet.Content.AlternateGreetings))
			},
		},
		{
			name:     "unmarshal with alternate greetings from numbers",
			jsonData: `{"spec":"chara_card_v3","spec_version":"3.0","data":{"title":"Test","name":"Test","alternate_greetings":[324,325,326]}}`,
			expected: func(t *testing.T, sheet *Sheet) {
				assert.Equal(t, SpecV3, sheet.Spec)
				assert.Len(t, sheet.Content.AlternateGreetings, 3)
				assert.Equal(t, []string{"324", "325", "326"}, []string(sheet.Content.AlternateGreetings))
			},
		},
		{
			name:     "unmarshal with depth prompt",
			jsonData: `{"spec":"chara_card_v3","spec_version":"3.0","data":{"title":"Test","name":"Test","extensions":{"depth_prompt":{"prompt":"test prompt","depth":5}}}}`,
			expected: func(t *testing.T, sheet *Sheet) {
				assert.Equal(t, "test prompt", sheet.Content.DepthPrompt.Prompt)
				assert.Equal(t, 5, sheet.Content.DepthPrompt.Depth)
			},
		},
		{
			name:     "unmarshal with depth prompt string depth",
			jsonData: `{"spec":"chara_card_v3","spec_version":"3.0","data":{"title":"Test","name":"Test","extensions":{"depth_prompt":{"prompt":"test prompt","depth":"10"}}}}`,
			expected: func(t *testing.T, sheet *Sheet) {
				assert.Equal(t, "test prompt", sheet.Content.DepthPrompt.Prompt)
				assert.Equal(t, 10, sheet.Content.DepthPrompt.Depth)
			},
		},
		{
			name:     "unmarshal with depth prompt missing depth",
			jsonData: `{"spec":"chara_card_v3","spec_version":"3.0","data":{"title":"Test","name":"Test","extensions":{"depth_prompt":{"prompt":"test prompt"}}}}`,
			expected: func(t *testing.T, sheet *Sheet) {
				assert.Equal(t, "test prompt", sheet.Content.DepthPrompt.Prompt)
				assert.Equal(t, DefaultDepth, sheet.Content.DepthPrompt.Depth)
			},
		},
		{
			name:     "unmarshal with empty depth prompt",
			jsonData: `{"spec":"chara_card_v3","spec_version":"3.0","data":{"title":"Test","name":"Test","extensions":{"depth_prompt":{"prompt":"","depth":10}}}}`,
			expected: func(t *testing.T, sheet *Sheet) {
				assert.Empty(t, sheet.Content.DepthPrompt.Prompt)
				assert.Zero(t, sheet.Content.DepthPrompt.Depth)
			},
		},
		{
			name:     "unmarshal legacy sheet without spec defaults to V2",
			jsonData: `{"data":{"title":"Legacy Character","name":"LegacyChar"}}`,
			expected: func(t *testing.T, sheet *Sheet) {
				assert.Equal(t, SpecV2, sheet.Spec)
				assert.Equal(t, V2, sheet.Version)
				assert.Equal(t, RevisionV2, sheet.Revision)
				assert.Equal(t, "Legacy Character", string(sheet.Content.Title))
			},
		},
		{
			name:     "unmarshal with empty spec defaults to V2",
			jsonData: `{"spec":"","spec_version":"","data":{}}`,
			expected: func(t *testing.T, sheet *Sheet) {
				assert.Equal(t, SpecV2, sheet.Spec)
				assert.Equal(t, V2, sheet.Version)
				assert.Equal(t, RevisionV2, sheet.Revision)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sheet Sheet
			err := sonicx.Config.UnmarshalFromString(tt.jsonData, &sheet)
			require.NoError(t, err)

			tt.expected(t, &sheet)
		})
	}
}

func TestSheet_UnmarshalJSON_ErrorCases(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{
			name:     "invalid JSON syntax",
			jsonData: `{"spec":"chara_card_v3","invalid json}`,
		},
		{
			name:     "malformed JSON",
			jsonData: `{spec:"chara_card_v3"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sheet Sheet
			err := sonicx.Config.UnmarshalFromString(tt.jsonData, &sheet)
			assert.Error(t, err)
		})
	}
}

func TestSheet_ToJSON(t *testing.T) {
	sheet := &Sheet{
		Spec:    SpecV3,
		Version: V3,
		Content: Content{
			Title:       property.String("JSON Test"),
			Name:        property.String("JSONChar"),
			Description: property.String("Testing JSON output"),
		},
	}

	var buf bytes.Buffer
	err := sheet.ToJSON(&buf)
	require.NoError(t, err)

	result := buf.String()
	assert.Contains(t, result, `"spec":"chara_card_v3"`)
	assert.Contains(t, result, `"title":"JSON Test"`)
	assert.Contains(t, result, `"name":"JSONChar"`)
}

func TestSheet_ToBytes(t *testing.T) {
	sheet := &Sheet{
		Spec:    SpecV3,
		Version: V3,
		Content: Content{
			Title: property.String("Bytes Test"),
			Name:  property.String("BytesChar"),
		},
	}

	data, err := sheet.ToBytes()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	result := string(data)
	assert.Contains(t, result, `"title":"Bytes Test"`)
	assert.Contains(t, result, `"name":"BytesChar"`)
}

func TestSheet_ToFile(t *testing.T) {
	sheet := &Sheet{
		Spec:    SpecV3,
		Version: V3,
		Content: Content{
			Title: property.String("File Test"),
			Name:  property.String("FileChar"),
		},
	}

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "sheet_test_*.json")
	require.NoError(t, err)
	_ = tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Write to file
	err = sheet.ToFile(tmpFile.Name())
	require.NoError(t, err)

	// Read back and verify
	data, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)

	result := string(data)
	assert.Contains(t, result, `"title":"File Test"`)
	assert.Contains(t, result, `"name":"FileChar"`)
}

func TestFromJSON(t *testing.T) {
	jsonData := `{"spec":"chara_card_v3","spec_version":"3.0","data":{"title":"Reader Test","name":"ReaderChar"}}`
	reader := strings.NewReader(jsonData)

	sheet, err := FromJSON(reader)
	require.NoError(t, err)
	require.NotNil(t, sheet)

	assert.Equal(t, SpecV3, sheet.Spec)
	assert.Equal(t, V3, sheet.Version)
	assert.Equal(t, "Reader Test", string(sheet.Content.Title))
	assert.Equal(t, "ReaderChar", string(sheet.Content.Name))
}

func TestFromBytes(t *testing.T) {
	jsonData := []byte(`{"spec":"chara_card_v3","spec_version":"3.0","data":{"title":"Bytes Reader Test","name":"BytesReaderChar"}}`)

	sheet, err := FromBytes(jsonData)
	require.NoError(t, err)
	require.NotNil(t, sheet)

	assert.Equal(t, SpecV3, sheet.Spec)
	assert.Equal(t, "Bytes Reader Test", string(sheet.Content.Title))
}

func TestFromFile(t *testing.T) {
	// Create a temporary file with JSON data
	tmpFile, err := os.CreateTemp("", "sheet_read_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	jsonData := `{"spec":"chara_card_v3","spec_version":"3.0","data":{"title":"File Reader Test","name":"FileReaderChar"}}`
	_, err = tmpFile.WriteString(jsonData)
	require.NoError(t, err)
	_ = tmpFile.Close()

	sheet, err := FromFile(tmpFile.Name())
	require.NoError(t, err)
	require.NotNil(t, sheet)

	assert.Equal(t, SpecV3, sheet.Spec)
	assert.Equal(t, "File Reader Test", string(sheet.Content.Title))
	assert.Equal(t, "FileReaderChar", string(sheet.Content.Name))
}

func TestFromFile_Error(t *testing.T) {
	_, err := FromFile("nonexistent_file.json")
	assert.Error(t, err)
}

func TestSheet_MarshalDepthPromptNonDestructively(t *testing.T) {
	t.Run("depth prompt with existing extensions", func(t *testing.T) {
		sheet := DefaultSheet(RevisionV3)
		sheet.Content.DepthPrompt = DepthPrompt{
			Prompt: "test prompt",
			Depth:  5,
		}
		sheet.Content.Extensions = map[string]any{
			"role": "user",
			DepthPromptKey: map[string]any{
				"other_prop": "should be preserved",
			},
		}

		jsonBytes, err := sheet.ToBytes()
		require.NoError(t, err)
		unmarshaledSheet, err := FromBytes(jsonBytes)
		require.NoError(t, err)

		assert.Equal(t, "user", unmarshaledSheet.Content.Extensions["role"])
		assert.Equal(t, "test prompt", unmarshaledSheet.Content.DepthPrompt.Prompt)
		assert.Equal(t, 5, unmarshaledSheet.Content.DepthPrompt.Depth)

		depthPromptMap, ok := unmarshaledSheet.Content.Extensions[DepthPromptKey].(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "should be preserved", depthPromptMap["other_prop"])
	})

	t.Run("depth prompt without other keys", func(t *testing.T) {
		sheet := DefaultSheet(RevisionV3)
		sheet.Content.DepthPrompt = DepthPrompt{
			Prompt: "test prompt",
			Depth:  5,
		}
		sheet.Content.Extensions = map[string]any{
			"role": "user",
		}

		jsonBytes, err := sheet.ToBytes()
		require.NoError(t, err)
		unmarshaledSheet, err := FromBytes(jsonBytes)
		require.NoError(t, err)

		assert.Equal(t, "user", unmarshaledSheet.Content.Extensions["role"])
		assert.Equal(t, "test prompt", unmarshaledSheet.Content.DepthPrompt.Prompt)
		assert.Equal(t, 5, unmarshaledSheet.Content.DepthPrompt.Depth)

		// depth_prompt should be completely removed since it only had prompt/depth
		_, ok := unmarshaledSheet.Content.Extensions[DepthPromptKey]
		assert.False(t, ok)
		assert.Len(t, unmarshaledSheet.Content.Extensions, 1)
	})
}

func TestSheet_NormalizeSymbols(t *testing.T) {
	abnormalQuotesJSON := `{
		"spec": "chara_card_v3",
		"spec_version": "3.0",
		"data": {
			"title": "Test",
			"name": "Test",
			"description": "'description'",
			"first_mes": "\"first message\"",
			"alternate_greetings": ["«greeting 1»"],
			"character_book": {
				"name": "„book name\"",
				"entries": [{"name": "「entry name」", "content": "《entry content》"}]
			},
			"extensions": {"depth_prompt": {"prompt": "'depth prompt'"}}
		}
	}`

	sheet, err := FromBytes([]byte(abnormalQuotesJSON))
	require.NoError(t, err)

	sheet.Content.NormalizeSymbols()

	assert.Equal(t, "'description'", string(sheet.Content.Description))
	assert.Equal(t, `"first message"`, string(sheet.Content.FirstMessage))
	assert.Equal(t, `"greeting 1"`, sheet.Content.AlternateGreetings[0])
	assert.Equal(t, `"book name"`, string(sheet.Content.CharacterBook.Name))
	assert.Equal(t, `"entry name"`, string(sheet.Content.CharacterBook.Entries[0].Name))
	assert.Equal(t, `"entry content"`, string(sheet.Content.CharacterBook.Entries[0].Content))
	assert.Equal(t, `'depth prompt'`, sheet.Content.DepthPrompt.Prompt)
}

func TestSheet_NormalizeSymbols_NameAndComment(t *testing.T) {
	sheet := DefaultSheet(RevisionV3)
	sheet.Content.CharacterBook = &Book{
		Entries: []*BookEntry{
			{BookEntryCore: BookEntryCore{Name: "entry name"}},
			{BookEntryCore: BookEntryCore{Comment: "entry comment"}},
		},
	}

	sheet.Content.NormalizeSymbols()

	require.Len(t, sheet.Content.CharacterBook.Entries, 2)
	entry1 := sheet.Content.CharacterBook.Entries[0]
	entry2 := sheet.Content.CharacterBook.Entries[1]

	assert.Equal(t, "entry name", string(entry1.Comment))
	assert.Equal(t, "entry comment", string(entry2.Name))
}

func TestSheet_Integrity(t *testing.T) {
	tests := []struct {
		name     string
		content  Content
		expected bool
	}{
		{
			name: "well-formed sheet",
			content: Content{
				Title:            property.String("Test Sheet"),
				Name:             property.String("Test Name"),
				Description:      property.String("A description."),
				FirstMessage:     property.String("A first message."),
				Creator:          property.String("A creator."),
				Nickname:         property.String("A Nickname"),
				CreationDate:     timestamp.Seconds(12345),
				ModificationDate: timestamp.Seconds(12345),
				SourceID:         property.String("A Source ID"),
			},
			expected: true,
		},
		{
			name: "malformed due to blank title",
			content: Content{
				Title:            property.String(" "),
				Name:             property.String("Test Name"),
				Description:      property.String("A description."),
				FirstMessage:     property.String("A first message."),
				Creator:          property.String("A creator."),
				Nickname:         property.String("A Nickname"),
				CreationDate:     timestamp.Seconds(12345),
				ModificationDate: timestamp.Seconds(12345),
				SourceID:         property.String("A Source ID"),
			},
			expected: false,
		},
		{
			name: "malformed due to zero creation date",
			content: Content{
				Title:            property.String("Test Sheet"),
				Name:             property.String("Test Name"),
				Description:      property.String("A description."),
				FirstMessage:     property.String("A first message."),
				Creator:          property.String("A creator."),
				Nickname:         property.String("A Nickname"),
				CreationDate:     timestamp.Seconds(0),
				ModificationDate: timestamp.Seconds(12345),
				SourceID:         property.String("A Source ID"),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sheet := &Sheet{Content: tt.content}
			assert.Equal(t, tt.expected, sheet.Content.Integrity())
		})
	}
}

func TestSheet_RoundTripIntegrity(t *testing.T) {
	originalJSON := `{
		"spec": "chara_card_v3", 
		"spec_version": "3.0",
		"data": {
			"title": "Test Character", 
			"name": "TestChar",
			"description": "A test sheet.",
			"alternate_greetings": ["Hi", "Hello"],
			"modification_date": 100,
			"extensions": {
				"depth_prompt": {"prompt": "A deep prompt.", "depth": 8, "other_data": "preserved"},
				"misc": "some data"
			},
			"tags": []
		}
	}`

	originalSheet, err := FromBytes([]byte(originalJSON))
	require.NoError(t, err)

	encodedBytes, err := originalSheet.ToBytes()
	require.NoError(t, err)

	secondarySheet, err := FromBytes(encodedBytes)
	require.NoError(t, err)

	assert.Equal(t, originalSheet.Spec, secondarySheet.Spec)
	assert.Equal(t, originalSheet.Version, secondarySheet.Version)
	assert.Equal(t, originalSheet.Revision, secondarySheet.Revision)
	assert.Equal(t, originalSheet.Content.Title, secondarySheet.Content.Title)
	assert.Equal(t, originalSheet.Content.Name, secondarySheet.Content.Name)
	assert.Equal(t, originalSheet.Content.Description, secondarySheet.Content.Description)
	assert.Equal(t, originalSheet.Content.AlternateGreetings, secondarySheet.Content.AlternateGreetings)
	assert.Equal(t, originalSheet.Content.ModificationDate, secondarySheet.Content.ModificationDate)
	assert.Equal(t, originalSheet.Content.DepthPrompt.Prompt, secondarySheet.Content.DepthPrompt.Prompt)
	assert.Equal(t, originalSheet.Content.DepthPrompt.Depth, secondarySheet.Content.DepthPrompt.Depth)

	// Check that extensions are preserved
	assert.Equal(t, "some data", secondarySheet.Content.Extensions["misc"])
	depthPromptMap, ok := secondarySheet.Content.Extensions[DepthPromptKey].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "preserved", depthPromptMap["other_data"])
}

func TestSheet_ComprehensiveRoundTrip(t *testing.T) {
	// Create a comprehensive sheet JSON with every possible field populated
	comprehensiveJSON := `{
		"spec": "chara_card_v3",
		"spec_version": "3.0",
		"data": {
			"title": "Comprehensive Test Character",
			"name": "ComprehensiveChar",
			"description": "A character with every possible field populated for testing.",
			"personality": "Friendly, outgoing, and comprehensive",
			"scenario": "Testing scenario with detailed background",
			"first_mes": "Hello! I'm a comprehensive test character.",
			"mes_example": "<START>\n{{user}}: Hello\n{{char}}: Hi there!\n<START>\n{{user}}: How are you?\n{{char}}: I'm doing great!",
			"creator_notes": "Created for comprehensive testing purposes\n\nIncludes all possible fields",
			"system_prompt": "You are a helpful assistant for testing.",
			"post_history_instructions": "Remember to stay in character.",
			"alternate_greetings": [
				"Hi there! Ready for some comprehensive testing?",
				"Greetings! I have all the fields populated.",
				"Hey! Testing every possible property."
			],
			"character_book": {
				"name": "Comprehensive Lorebook",
				"description": "A lorebook with all possible configurations",
				"scan_depth": 100,
				"token_budget": 2048,
				"recursive_scanning": true,
				"extensions": {
					"custom_book_field": "custom_book_value",
					"book_metadata": {
						"version": "1.0",
						"author": "Test Suite"
					}
				},
				"entries": [
					{
						"id": 1,
						"keys": ["comprehensive", "test", "character"],
						"secondary_keys": ["comp", "test"],
						"name": "Comprehensive Entry",
						"comment": "Main character entry",
						"content": "This is comprehensive test content for the character.",
						"constant": true,
						"selective": true,
						"insertion_order": 100,
						"enabled": true,
						"use_regex": true,
						"extensions": {
							"position": 2,
							"probability": 85.00,
							"depth": 3,
							"selectiveLogic": 3,
							"match_whole_words": true,
							"case_sensitive": false,
							"role": 1,
							"sticky": 2,
							"cooldown": 5,
							"delay": 1,
							"entry_custom": "entry_value"
						}
					},
					{
						"id": 2,
						"keys": ["c", "t", "cc"],
						"secondary_keys": ["cc", "tt"],
						"name": "Comprehensive Entry2",
						"comment": "Main character entry2",
						"content": "This is comprehensive test content for the character2.",
						"constant": false,
						"selective": false,
						"insertion_order": 85,
						"enabled": false,
						"use_regex": false,
						"extensions": {
							"position": 3,
							"probability": 95.00,
							"depth": 2,
							"selectiveLogic": 1,
							"match_whole_words": false,
							"case_sensitive": true,
							"role": 2,
							"sticky": 3,
							"cooldown": 5,
							"delay": 2,
							"entry_custom2": "entry_value2"
						}
					}
				]
			},
			"tags": ["comprehensive", "test", "full-featured", "roundtrip"],
			"creator": "Test Suite Author",
			"character_version": "2.1.0",
			"creation_date": 1640995200,
			"modification_date": 1672531200,
			"nickname": "CompChar",
			"extensions": {
				"depth_prompt": {
					"prompt": "Think deeply about this comprehensive character.",
					"depth": 10,
					"custom_depth_field": "custom_value"
				},
				"custom_extension_1": "value1",
				"custom_extension_2": {
					"nested": "data",
					"number": 42,
					"boolean": true,
					"array": ["item1", "item2", "item3"]
				},
				"character_metadata": {
					"test_version": "1.0",
					"features": ["comprehensive", "roundtrip", "validation"]
				}
			},
			"source_id": "comprehensive_test_001",
			"character_id": "comprehensive_id_001",
			"platform_id": "comprehensive_pt_id_001",
			"direct_link": "https://example.com/comprehensive_test_001"
		}
	}`

	originalSheet, err := FromBytes([]byte(comprehensiveJSON))
	require.NoError(t, err)

	marshaledBytes, err := originalSheet.ToBytes()
	require.NoError(t, err)

	roundtripSheet, err := FromBytes(marshaledBytes)
	require.NoError(t, err)
	println(cmp.Diff(originalSheet, roundtripSheet))
	assert.True(t, cmp.Equal(originalSheet, roundtripSheet, cmpopts.EquateEmpty()))
}

func TestSheet_DeepEquals(t *testing.T) {
	tests := []struct {
		name     string
		sheet1   *Sheet
		sheet2   *Sheet
		expected bool
	}{
		{
			name: "identical sheets",
			sheet1: &Sheet{
				Spec:    SpecV3,
				Version: V3,
				Content: Content{
					Title:       property.String("Test"),
					Name:        property.String("TestChar"),
					Description: property.String("A test character"),
				},
			},
			sheet2: &Sheet{
				Spec:    SpecV3,
				Version: V3,
				Content: Content{
					Title:       property.String("Test"),
					Name:        property.String("TestChar"),
					Description: property.String("A test character"),
				},
			},
			expected: true,
		},
		{
			name: "different titles",
			sheet1: &Sheet{
				Spec:    SpecV3,
				Version: V3,
				Content: Content{
					Title: property.String("Test1"),
					Name:  property.String("TestChar"),
				},
			},
			sheet2: &Sheet{
				Spec:    SpecV3,
				Version: V3,
				Content: Content{
					Title: property.String("Test2"),
					Name:  property.String("TestChar"),
				},
			},
			expected: false,
		},
		{
			name: "different specs",
			sheet1: &Sheet{
				Spec:    SpecV2,
				Version: V2,
				Content: Content{
					Title: property.String("Test"),
					Name:  property.String("TestChar"),
				},
			},
			sheet2: &Sheet{
				Spec:    SpecV3,
				Version: V3,
				Content: Content{
					Title: property.String("Test"),
					Name:  property.String("TestChar"),
				},
			},
			expected: false,
		},
		{
			name: "empty sheets",
			sheet1: &Sheet{
				Spec:    SpecV3,
				Version: V3,
			},
			sheet2: &Sheet{
				Spec:    SpecV3,
				Version: V3,
			},
			expected: true,
		},
		{
			name: "sheets with different alternate greetings order",
			sheet1: &Sheet{
				Spec:    SpecV3,
				Version: V3,
				Content: Content{
					Title:              property.String("Test"),
					Name:               property.String("TestChar"),
					AlternateGreetings: []string{"Hi", "Hello", "Hey"},
				},
			},
			sheet2: &Sheet{
				Spec:    SpecV3,
				Version: V3,
				Content: Content{
					Title:              property.String("Test"),
					Name:               property.String("TestChar"),
					AlternateGreetings: []string{"Hey", "Hello", "Hi"},
				},
			},
			expected: true,
		},
		{
			name: "sheets with different tags order",
			sheet1: &Sheet{
				Spec:    SpecV3,
				Version: V3,
				Content: Content{
					Title: property.String("Test"),
					Name:  property.String("TestChar"),
					Tags:  []string{"tag1", "tag2", "tag3"},
				},
			},
			sheet2: &Sheet{
				Spec:    SpecV3,
				Version: V3,
				Content: Content{
					Title: property.String("Test"),
					Name:  property.String("TestChar"),
					Tags:  []string{"tag3", "tag1", "tag2"},
				},
			},
			expected: true,
		},
		{
			name: "one sheet nil, other not nil",
			sheet1: &Sheet{
				Spec:    SpecV3,
				Version: V3,
				Content: Content{
					Title: property.String("Test"),
				},
			},
			sheet2:   nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.sheet2 == nil {
				assert.Equal(t, tt.expected, tt.sheet1.DeepEquals(tt.sheet2))
			} else {
				assert.Equal(t, tt.expected, tt.sheet1.DeepEquals(tt.sheet2))
			}
		})
	}
}

func TestSheet_SetRevision(t *testing.T) {
	tests := []struct {
		name     string
		revision Revision
		expected Stamp
	}{
		{
			name:     "set V2 revision",
			revision: RevisionV2,
			expected: Stamps[RevisionV2],
		},
		{
			name:     "set V3 revision",
			revision: RevisionV3,
			expected: Stamps[RevisionV3],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sheet := &Sheet{}
			sheet.SetRevision(tt.revision)

			assert.Equal(t, tt.expected.Revision, sheet.Revision)
			assert.Equal(t, tt.expected.Spec, sheet.Spec)
			assert.Equal(t, tt.expected.Version, sheet.Version)
		})
	}
}

func TestSheetConstants(t *testing.T) {
	assert.Equal(t, "\n\n", CreatorNotesSeparator)
	assert.Equal(t, "Anonymous", AnonymousCreator)
}
