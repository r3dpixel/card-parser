package character

import (
	"testing"

	"github.com/r3dpixel/card-parser/property"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/stringsx"
	"github.com/r3dpixel/toolkit/timestamp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContent_MarshalJSONTo(t *testing.T) {
	tests := []struct {
		name     string
		content  *Content
		contains []string
	}{
		{
			name: "marshal without depth prompt",
			content: &Content{
				Title:       property.String("Test Character"),
				Name:        property.String("TestChar"),
				Description: property.String("A test character"),
			},
			contains: []string{`"title":"Test Character"`, `"name":"TestChar"`, `"description":"A test character"`},
		},
		{
			name: "marshal with depth prompt",
			content: &Content{
				Title:       property.String("Test Character"),
				Name:        property.String("TestChar"),
				Description: property.String("A test character"),
				DepthPrompt: DepthPrompt{
					Prompt: "Test depth prompt",
					Depth:  5,
				},
			},
			contains: []string{`"title":"Test Character"`, `"depth_prompt"`, `"prompt":"Test depth prompt"`, `"depth":5`},
		},
		{
			name: "marshal with existing extensions and depth prompt",
			content: &Content{
				Title:       property.String("Test Character"),
				Name:        property.String("TestChar"),
				Description: property.String("A test character"),
				Extensions: map[string]any{
					"existing_key": "existing_value",
				},
				DepthPrompt: DepthPrompt{
					Prompt: "Test depth prompt",
					Depth:  3,
				},
			},
			contains: []string{`"existing_key":"existing_value"`, `"depth_prompt"`, `"prompt":"Test depth prompt"`, `"depth":3`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := sonicx.Config.Marshal(tt.content)
			require.NoError(t, err)

			result := string(data)
			for _, expected := range tt.contains {
				assert.Contains(t, result, expected)
			}

			if stringsx.IsNotBlank(tt.content.DepthPrompt.Prompt) {
				assert.Contains(t, result, `"depth_prompt"`)
				if tt.content.Extensions != nil {
					assert.NotContains(t, tt.content.Extensions, DepthPromptKey)
				}
			}
		})
	}
}

func TestContent_UnmarshalJSONFrom(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected func(t *testing.T, content *Content)
	}{
		{
			name:     "unmarshal without depth prompt",
			jsonData: `{"title":"Test Character","name":"TestChar","description":"A test character"}`,
			expected: func(t *testing.T, content *Content) {
				assert.Equal(t, "Test Character", string(content.Title))
				assert.Equal(t, "TestChar", string(content.Name))
				assert.Equal(t, "A test character", string(content.Description))
				assert.Empty(t, content.DepthPrompt.Prompt)
				assert.Zero(t, content.DepthPrompt.Depth)
				assert.Empty(t, content.Extensions)
			},
		},
		{
			name:     "unmarshal with depth prompt",
			jsonData: `{"title":"Test Character","name":"TestChar","description":"A test character","extensions":{"depth_prompt":{"prompt":"Test depth prompt","depth":5}}}`,
			expected: func(t *testing.T, content *Content) {
				assert.Equal(t, "Test Character", string(content.Title))
				assert.Equal(t, "TestChar", string(content.Name))
				assert.Equal(t, "A test character", string(content.Description))
				assert.Equal(t, "Test depth prompt", content.DepthPrompt.Prompt)
				assert.Equal(t, 5, content.DepthPrompt.Depth)
				assert.Empty(t, content.Extensions)
			},
		},
		{
			name:     "unmarshal with depth prompt and default level",
			jsonData: `{"title":"Test Character","name":"TestChar","description":"A test character","extensions":{"depth_prompt":{"prompt":"Test depth prompt"}}}`,
			expected: func(t *testing.T, content *Content) {
				assert.Equal(t, "Test Character", string(content.Title))
				assert.Equal(t, "Test depth prompt", content.DepthPrompt.Prompt)
				assert.Equal(t, DefaultDepth, content.DepthPrompt.Depth)
				assert.Empty(t, content.Extensions)
			},
		},
		{
			name:     "unmarshal with depth prompt and other extensions",
			jsonData: `{"title":"Test Character","name":"TestChar","description":"A test character","extensions":{"depth_prompt":{"prompt":"Test depth prompt","depth":3},"other_ext":"other_value"}}`,
			expected: func(t *testing.T, content *Content) {
				assert.Equal(t, "Test Character", string(content.Title))
				assert.Equal(t, "Test depth prompt", content.DepthPrompt.Prompt)
				assert.Equal(t, 3, content.DepthPrompt.Depth)
				expected := map[string]any{"other_ext": "other_value"}
				assert.Equal(t, expected, content.Extensions)
				assert.NotContains(t, content.Extensions, DepthPromptKey)
			},
		},
		{
			name:     "unmarshal with depth prompt containing other keys",
			jsonData: `{"title":"Test Character","name":"TestChar","description":"A test character","extensions":{"depth_prompt":{"prompt":"Test depth prompt","depth":2,"other_key":"other_value"}}}`,
			expected: func(t *testing.T, content *Content) {
				assert.Equal(t, "Test Character", string(content.Title))
				assert.Equal(t, "Test depth prompt", content.DepthPrompt.Prompt)
				assert.Equal(t, 2, content.DepthPrompt.Depth)
				expected := map[string]any{
					DepthPromptKey: map[string]any{
						"other_key": "other_value",
					},
				}
				assert.Equal(t, expected, content.Extensions)
			},
		},
		{
			name:     "unmarshal with empty depth prompt",
			jsonData: `{"title":"Test Character","name":"TestChar","description":"A test character","extensions":{"depth_prompt":{"prompt":"","depth":5}}}`,
			expected: func(t *testing.T, content *Content) {
				assert.Equal(t, "Test Character", string(content.Title))
				assert.Empty(t, content.DepthPrompt.Prompt)
				assert.Zero(t, content.DepthPrompt.Depth)
				assert.Empty(t, content.Extensions)
			},
		},
		{
			name:     "unmarshal with whitespace-only depth prompt",
			jsonData: `{"title":"Test Character","name":"TestChar","description":"A test character","extensions":{"depth_prompt":{"prompt":"   \t\n   ","depth":5}}}`,
			expected: func(t *testing.T, content *Content) {
				assert.Equal(t, "Test Character", string(content.Title))
				assert.Empty(t, content.DepthPrompt.Prompt)
				assert.Zero(t, content.DepthPrompt.Depth)
				assert.Empty(t, content.Extensions)
			},
		},
		{
			name:     "unmarshal with invalid depth prompt format",
			jsonData: `{"title":"Test Character","name":"TestChar","description":"A test character","extensions":{"depth_prompt":"invalid_format"}}`,
			expected: func(t *testing.T, content *Content) {
				assert.Equal(t, "Test Character", string(content.Title))
				assert.Equal(t, "invalid_format", content.DepthPrompt.Prompt)
				assert.Equal(t, DefaultDepth, content.DepthPrompt.Depth)
				assert.Empty(t, content.Extensions)
			},
		},
		{
			name:     "unmarshal with invalid depth level",
			jsonData: `{"title":"Test Character","name":"TestChar","description":"A test character","extensions":{"depth_prompt":{"prompt":"Test depth prompt","depth":"invalid"}}}`,
			expected: func(t *testing.T, content *Content) {
				assert.Equal(t, "Test Character", string(content.Title))
				assert.Equal(t, "Test depth prompt", content.DepthPrompt.Prompt)
				assert.Equal(t, DefaultDepth, content.DepthPrompt.Depth)
				assert.Empty(t, content.Extensions)
			},
		},
		{
			name:     "unmarshal with depth prompt as array",
			jsonData: `{"title":"Test Character","name":"TestChar","description":"A test character","extensions":{"depth_prompt":["some","array","data"]}}`,
			expected: func(t *testing.T, content *Content) {
				assert.Equal(t, "Test Character", string(content.Title))
				assert.Equal(t, `["some","array","data"]`, content.DepthPrompt.Prompt) // Array should be stringified via jsonx.String
				assert.Equal(t, DefaultDepth, content.DepthPrompt.Depth)
				assert.Empty(t, content.Extensions)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var content Content
			err := sonicx.Config.UnmarshalFromString(tt.jsonData, &content)
			require.NoError(t, err)

			tt.expected(t, &content)
		})
	}
}

func TestContent_UnmarshalJSON_ErrorCases(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{
			name:     "invalid JSON syntax",
			jsonData: `{"title":"Test Character","invalid json}`,
		},
		{
			name:     "malformed JSON",
			jsonData: `{title:"Test Character"}`, // Missing quotes around key
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var content Content
			err := sonicx.Config.UnmarshalFromString(tt.jsonData, &content)
			assert.Error(t, err) // Should return an error for invalid JSON
		})
	}
}

func TestContent_MarshalUnmarshal_Roundtrip(t *testing.T) {
	original := Content{
		Title:       property.String("Roundtrip Test"),
		Name:        property.String("RoundtripChar"),
		Description: property.String("A character for roundtrip testing"),
		Creator:     property.String("Test Creator"),
		Tags:        property.StringArray{"test", "roundtrip"},
		Extensions: map[string]any{
			"custom_ext": "custom_value",
		},
		DepthPrompt: DepthPrompt{
			Prompt: "Roundtrip depth prompt",
			Depth:  7,
		},
		CreationDate:     timestamp.Seconds(1234567890),
		ModificationDate: timestamp.Seconds(1234567999),
	}

	jsonData, err := sonicx.Config.Marshal(&original)
	require.NoError(t, err)

	var unmarshaled Content
	err = sonicx.Config.UnmarshalFromString(stringsx.FromBytes(jsonData), &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, original.Title, unmarshaled.Title)
	assert.Equal(t, original.Name, unmarshaled.Name)
	assert.Equal(t, original.Description, unmarshaled.Description)
	assert.Equal(t, original.Creator, unmarshaled.Creator)
	assert.Equal(t, original.Tags, unmarshaled.Tags)
	assert.Equal(t, original.CreationDate, unmarshaled.CreationDate)
	assert.Equal(t, original.ModificationDate, unmarshaled.ModificationDate)

	expectedExtensions := map[string]any{
		"custom_ext": "custom_value",
	}
	assert.Equal(t, expectedExtensions, unmarshaled.Extensions)

	assert.Equal(t, original.DepthPrompt.Prompt, unmarshaled.DepthPrompt.Prompt)
	assert.Equal(t, original.DepthPrompt.Depth, unmarshaled.DepthPrompt.Depth)
}

func TestContent_NormalizeSymbols_NameAndComment(t *testing.T) {
	tests := []struct {
		name     string
		content  *Content
		expected func(t *testing.T, content *Content)
	}{
		{
			name: "content with no character book",
			content: &Content{
				Title: property.String("Test"),
			},
			expected: func(t *testing.T, content *Content) {
				assert.Nil(t, content.CharacterBook)
			},
		},
		{
			name: "content with character book having entries with only name",
			content: &Content{
				Title: property.String("Test"),
				CharacterBook: &Book{
					Entries: []*BookEntry{
						{BookEntryCore: BookEntryCore{Name: "entry1"}},
						{BookEntryCore: BookEntryCore{Name: "entry2"}},
					},
				},
			},
			expected: func(t *testing.T, content *Content) {
				require.NotNil(t, content.CharacterBook)
				require.Len(t, content.CharacterBook.Entries, 2)
				assert.Equal(t, "entry1", string(content.CharacterBook.Entries[0].Name))
				assert.Equal(t, "entry1", string(content.CharacterBook.Entries[0].Comment)) // Should mirror name to comment
				assert.Equal(t, "entry2", string(content.CharacterBook.Entries[1].Name))
				assert.Equal(t, "entry2", string(content.CharacterBook.Entries[1].Comment)) // Should mirror name to comment
			},
		},
		{
			name: "content with character book having entries with only comment",
			content: &Content{
				Title: property.String("Test"),
				CharacterBook: &Book{
					Entries: []*BookEntry{
						{BookEntryCore: BookEntryCore{Comment: "comment1"}},
						{BookEntryCore: BookEntryCore{Comment: "comment2"}},
					},
				},
			},
			expected: func(t *testing.T, content *Content) {
				require.NotNil(t, content.CharacterBook)
				require.Len(t, content.CharacterBook.Entries, 2)
				assert.Equal(t, "comment1", string(content.CharacterBook.Entries[0].Name)) // Should mirror comment to name
				assert.Equal(t, "comment1", string(content.CharacterBook.Entries[0].Comment))
				assert.Equal(t, "comment2", string(content.CharacterBook.Entries[1].Name)) // Should mirror comment to name
				assert.Equal(t, "comment2", string(content.CharacterBook.Entries[1].Comment))
			},
		},
		{
			name: "content with character book having mixed entries",
			content: &Content{
				Title: property.String("Test"),
				CharacterBook: &Book{
					Entries: []*BookEntry{
						{BookEntryCore: BookEntryCore{Name: "name1"}},
						{BookEntryCore: BookEntryCore{Comment: "comment2"}},
						{BookEntryCore: BookEntryCore{Name: "name3", Comment: "comment3"}},
						{},
					},
				},
			},
			expected: func(t *testing.T, content *Content) {
				require.NotNil(t, content.CharacterBook)
				require.Len(t, content.CharacterBook.Entries, 4)

				// Entry with only name - comment should be mirrored
				assert.Equal(t, "name1", string(content.CharacterBook.Entries[0].Name))
				assert.Equal(t, "name1", string(content.CharacterBook.Entries[0].Comment))

				// Entry with only comment - name should be mirrored
				assert.Equal(t, "comment2", string(content.CharacterBook.Entries[1].Name))
				assert.Equal(t, "comment2", string(content.CharacterBook.Entries[1].Comment))

				// Entry with both - should remain unchanged
				assert.Equal(t, "name3", string(content.CharacterBook.Entries[2].Name))
				assert.Equal(t, "comment3", string(content.CharacterBook.Entries[2].Comment))

				// Empty entry - should remain empty
				assert.Empty(t, content.CharacterBook.Entries[3].Name)
				assert.Empty(t, content.CharacterBook.Entries[3].Comment)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.content.NormalizeSymbols()
			tt.expected(t, tt.content)
		})
	}
}

func TestContent_NormalizeSymbols(t *testing.T) {
	content := &Content{
		Title:                   property.String(`Test's 《Title》`),
		Name:                    property.String("Character's 'Name'"),
		Description:             property.String(`Description with „quotes"`),
		Personality:             property.String(`Personality's "text"`),
		Scenario:                property.String("Scenario's 'text'"),
		FirstMessage:            property.String(`First's "message"`),
		MessageExamples:         property.String("Message's 'examples'"),
		CreatorNotes:            property.String(`Creator's 「notes」`),
		SystemPrompt:            property.String(`System's 'prompt'`),
		PostHistoryInstructions: property.String(`Post's «instructions»`),
		AlternateGreetings:      property.StringArray{`Greeting's 〈one〉`, `Greeting's ‹two›`},
		CharacterBook: &Book{
			Entries: []*BookEntry{
				{
					BookEntryCore: BookEntryCore{
						Name:    `Entry's ‚name'`,
						Comment: `Entry's ‛comment'`,
						Content: `Entry's 〝content〞`,
					},
				},
			},
		},
		DepthPrompt: DepthPrompt{
			Prompt: `Depth's 〈prompt〉`, // 〈〉 should become ""
			Depth:  5,
		},
	}

	originalTitle := string(content.Title)
	originalDescription := string(content.Description)
	originalCreatorNotes := string(content.CreatorNotes)
	originalPostHistory := string(content.PostHistoryInstructions)
	originalGreeting1 := content.AlternateGreetings[0]
	originalGreeting2 := content.AlternateGreetings[1]
	originalEntryName := content.CharacterBook.Entries[0].Name
	originalEntryComment := content.CharacterBook.Entries[0].Comment
	originalEntryContent := content.CharacterBook.Entries[0].Content
	originalDepthPrompt := content.DepthPrompt.Prompt

	assert.Contains(t, originalTitle, "《")
	assert.Contains(t, originalDescription, "„")
	assert.Contains(t, originalCreatorNotes, "「")
	assert.Contains(t, originalPostHistory, "«")
	assert.Contains(t, originalGreeting1, "〈")
	assert.Contains(t, originalGreeting2, "‹")
	assert.Contains(t, originalEntryName, "‚")
	assert.Contains(t, originalEntryComment, "‛")
	assert.Contains(t, originalEntryContent, "〝")
	assert.Contains(t, originalDepthPrompt, "〈")

	content.NormalizeSymbols()

	assert.Equal(t, `Test's 《Title》`, string(content.Title))
	assert.Equal(t, "Character's 'Name'", string(content.Name)) // Should remain unchanged
	assert.Equal(t, `Description with "quotes"`, string(content.Description))
	assert.Equal(t, `Personality's "text"`, string(content.Personality))     // Should remain unchanged
	assert.Equal(t, "Scenario's 'text'", string(content.Scenario))           // Should remain unchanged
	assert.Equal(t, `First's "message"`, string(content.FirstMessage))       // Should remain unchanged
	assert.Equal(t, "Message's 'examples'", string(content.MessageExamples)) // Should remain unchanged
	assert.Equal(t, `Creator's "notes"`, string(content.CreatorNotes))
	assert.Equal(t, `System's 'prompt'`, string(content.SystemPrompt)) // Should remain unchanged
	assert.Equal(t, `Post's "instructions"`, string(content.PostHistoryInstructions))
	assert.Equal(t, `Greeting's "one"`, content.AlternateGreetings[0])
	assert.Equal(t, `Greeting's "two"`, content.AlternateGreetings[1])
	assert.Equal(t, `Entry's ,name'`, string(content.CharacterBook.Entries[0].Name))
	assert.Equal(t, `Entry's 'comment'`, string(content.CharacterBook.Entries[0].Comment))
	assert.Equal(t, `Entry's "content"`, string(content.CharacterBook.Entries[0].Content))
	assert.Equal(t, `Depth's "prompt"`, content.DepthPrompt.Prompt)

	assert.Contains(t, string(content.Title), "《")
	assert.NotContains(t, string(content.Description), "„")
	assert.NotContains(t, string(content.CreatorNotes), "「")
	assert.NotContains(t, string(content.PostHistoryInstructions), "«")
	assert.NotContains(t, content.AlternateGreetings[0], "〈")
	assert.NotContains(t, content.AlternateGreetings[1], "‹")
	assert.NotContains(t, content.CharacterBook.Entries[0].Name, "‚")
	assert.NotContains(t, content.CharacterBook.Entries[0].Comment, "‛")
	assert.NotContains(t, content.CharacterBook.Entries[0].Content, "〝")
	assert.NotContains(t, content.DepthPrompt.Prompt, "〈")
}

func TestContent_Integrity(t *testing.T) {
	tests := []struct {
		name     string
		content  *Content
		expected bool
	}{
		{
			name: "well-formed content",
			content: &Content{
				Title:            property.String("Valid Title"),
				Name:             property.String("Valid Name"),
				Description:      property.String("Valid Description"),
				FirstMessage:     property.String("Valid First Message"),
				Creator:          property.String("Valid Creator"),
				Nickname:         property.String("Valid Nickname"),
				SourceID:         property.String("Valid Source ID"),
				CreationDate:     timestamp.Seconds(1234567890),
				ModificationDate: timestamp.Seconds(1234567999),
			},
			expected: true,
		},
		{
			name: "missing title",
			content: &Content{
				Name:             property.String("Valid Name"),
				Description:      property.String("Valid Description"),
				FirstMessage:     property.String("Valid First Message"),
				Creator:          property.String("Valid Creator"),
				Nickname:         property.String("Valid Nickname"),
				SourceID:         property.String("Valid Source ID"),
				CreationDate:     timestamp.Seconds(1234567890),
				ModificationDate: timestamp.Seconds(1234567999),
			},
			expected: false,
		},
		{
			name: "missing creation date",
			content: &Content{
				Title:            property.String("Valid Title"),
				Name:             property.String("Valid Name"),
				Description:      property.String("Valid Description"),
				FirstMessage:     property.String("Valid First Message"),
				Creator:          property.String("Valid Creator"),
				Nickname:         property.String("Valid Nickname"),
				SourceID:         property.String("Valid Source ID"),
				ModificationDate: timestamp.Seconds(1234567999),
			},
			expected: false,
		},
		{
			name: "empty strings are considered blank",
			content: &Content{
				Title:            property.String(""),
				Name:             property.String("Valid Name"),
				Description:      property.String("Valid Description"),
				FirstMessage:     property.String("Valid First Message"),
				Creator:          property.String("Valid Creator"),
				Nickname:         property.String("Valid Nickname"),
				SourceID:         property.String("Valid Source ID"),
				CreationDate:     timestamp.Seconds(1234567890),
				ModificationDate: timestamp.Seconds(1234567999),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.content.Integrity()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContent_fixUserCharTemplate(t *testing.T) {
	c := &Content{}
	tests := []struct {
		input string
		want  string
	}{
		{"{char}", "{{char}}"},
		{"{{{char}}}", "{{char}}"},
		{"{{{{char}}}}", "{{char}}"},
		{"{user}", "{{user}}"},
		{"{{{user}}}", "{{user}}"},
		{"{{{{user}}}}", "{{user}}"},
		{"{char} and {user}", "{{char}} and {{user}}"},
		{"{{char}} meets {{{user}}}", "{{char}} meets {{user}}"},
		{"Hello {char}, I'm {user}!", "Hello {{char}}, I'm {{user}}!"},
		{"No templates here", "No templates here"},
		{"", ""},
		{"{char}{user}", "{{char}}{{user}}"},
		{"Multiple {char} and {user} in {char} one {user} string", "Multiple {{char}} and {{user}} in {{char}} one {{user}} string"},
		{"{CHAR} and {USER}", "{CHAR} and {USER}"},
		{"{chars} and {users}", "{chars} and {users}"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := c.fixUserCharTemplate(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContent_FixUserCharTemplates(t *testing.T) {
	tests := []struct {
		name   string
		input  *Content
		expect *Content
	}{
		{
			name: "all fields with templates",
			input: &Content{
				Description:             property.String("{char} is brave, {user} is clever"),
				Personality:             property.String("{{{char}}} is friendly to {user}"),
				Scenario:                property.String("{char} meets {{{user}}} in the forest"),
				FirstMessage:            property.String("Hello {user}, I'm {char}!"),
				MessageExamples:         property.String("{char}: Hi!\n{user}: Hello!"),
				SystemPrompt:            property.String("You are {char}, talking to {user}"),
				PostHistoryInstructions: property.String("Remember {char} and {user} context"),
				AlternateGreetings:      property.StringArray{"Hey {user}!", "{char} waves at {{{user}}}"},
				DepthPrompt:             DepthPrompt{Prompt: "{char} depth with {user}", Depth: 4},
			},
			expect: &Content{
				Description:             property.String("{{char}} is brave, {{user}} is clever"),
				Personality:             property.String("{{char}} is friendly to {{user}}"),
				Scenario:                property.String("{{char}} meets {{user}} in the forest"),
				FirstMessage:            property.String("Hello {{user}}, I'm {{char}}!"),
				MessageExamples:         property.String("{{char}}: Hi!\n{{user}}: Hello!"),
				SystemPrompt:            property.String("You are {{char}}, talking to {{user}}"),
				PostHistoryInstructions: property.String("Remember {{char}} and {{user}} context"),
				AlternateGreetings:      property.StringArray{"Hey {{user}}!", "{{char}} waves at {{user}}"},
				DepthPrompt:             DepthPrompt{Prompt: "{{char}} depth with {{user}}", Depth: 4},
			},
		},
		{
			name: "mixed valid and invalid templates",
			input: &Content{
				Description:  property.String("{char} and {CHAR} and {chars}"),
				Personality:  property.String("{user} but not {USER} or {users}"),
				FirstMessage: property.String("{{char}} already correct, {user} needs fix"),
			},
			expect: &Content{
				Description:  property.String("{{char}} and {CHAR} and {chars}"),
				Personality:  property.String("{{user}} but not {USER} or {users}"),
				FirstMessage: property.String("{{char}} already correct, {{user}} needs fix"),
			},
		},
		{
			name: "empty fields",
			input: &Content{
				Description: property.String(""),
				Personality: property.String(""),
			},
			expect: &Content{
				Description: property.String(""),
				Personality: property.String(""),
			},
		},
		{
			name: "no templates",
			input: &Content{
				Description:  property.String("Just plain text"),
				Personality:  property.String("No special markers"),
				Scenario:     property.String("Normal scenario"),
				FirstMessage: property.String("Hello world"),
			},
			expect: &Content{
				Description:  property.String("Just plain text"),
				Personality:  property.String("No special markers"),
				Scenario:     property.String("Normal scenario"),
				FirstMessage: property.String("Hello world"),
			},
		},
		{
			name: "excessive braces",
			input: &Content{
				Description:        property.String("{{{{char}}}} and {{{{{user}}}}}"),
				AlternateGreetings: property.StringArray{"{{{{{{char}}}}}}", "{user}"},
				DepthPrompt:        DepthPrompt{Prompt: "{{{{{{{{char}}}}}}}}", Depth: 5},
			},
			expect: &Content{
				Description:        property.String("{{char}} and {{user}}"),
				AlternateGreetings: property.StringArray{"{{char}}", "{{user}}"},
				DepthPrompt:        DepthPrompt{Prompt: "{{char}}", Depth: 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.FixUserCharTemplates()
			assert.Equal(t, tt.expect.Description, tt.input.Description)
			assert.Equal(t, tt.expect.Personality, tt.input.Personality)
			assert.Equal(t, tt.expect.Scenario, tt.input.Scenario)
			assert.Equal(t, tt.expect.FirstMessage, tt.input.FirstMessage)
			assert.Equal(t, tt.expect.MessageExamples, tt.input.MessageExamples)
			assert.Equal(t, tt.expect.SystemPrompt, tt.input.SystemPrompt)
			assert.Equal(t, tt.expect.PostHistoryInstructions, tt.input.PostHistoryInstructions)
			assert.Equal(t, tt.expect.AlternateGreetings, tt.input.AlternateGreetings)
			assert.Equal(t, tt.expect.DepthPrompt.Prompt, tt.input.DepthPrompt.Prompt)
			assert.Equal(t, tt.expect.DepthPrompt.Depth, tt.input.DepthPrompt.Depth)
		})
	}
}
