package character

import (
	"testing"

	"github.com/r3dpixel/card-parser/property"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/stringsx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyBook(t *testing.T) {
	book := DefaultBook()

	assert.NotNil(t, book)
	assert.Nil(t, book.Entries)
	assert.Equal(t, 0, len(book.Entries))
	assert.Empty(t, book.Name)
	assert.Empty(t, book.Description)
	assert.False(t, bool(book.RecursiveScanning))
}

func TestBook_NormalizeSymbols(t *testing.T) {
	tests := []struct {
		name     string
		book     *Book
		expected *Book
	}{
		{
			name: "normalize book name and description",
			book: &Book{
				Name:        `Book "Name"`,
				Description: "Description 'with' quotes",
				Entries:     []*BookEntry{},
			},
			expected: &Book{
				Name:        `Book "Name"`,
				Description: "Description 'with' quotes",
				Entries:     []*BookEntry{},
			},
		},
		{
			name: "normalize book entries",
			book: &Book{
				Name:        "Book Name",
				Description: "Description",
				Entries: []*BookEntry{
					{
						BookEntryCore: BookEntryCore{
							Name:    `Entry "Name"`,
							Comment: "Comment 'with' quotes",
							Content: `Content "text"`,
						},
					},
				},
			},
			expected: &Book{
				Name:        "Book Name",
				Description: "Description",
				Entries: []*BookEntry{
					{
						BookEntryCore: BookEntryCore{
							Name:    `Entry "Name"`,
							Comment: "Comment 'with' quotes",
							Content: `Content "text"`,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.book.NormalizeSymbols()
			assert.Equal(t, tt.expected.Name, tt.book.Name)
			assert.Equal(t, tt.expected.Description, tt.book.Description)
			for i, entry := range tt.book.Entries {
				assert.Equal(t, tt.expected.Entries[i].Name, entry.Name)
				assert.Equal(t, tt.expected.Entries[i].Comment, entry.Comment)
				assert.Equal(t, tt.expected.Entries[i].Content, entry.Content)
			}
		})
	}
}

func TestBook_JSONMarshal(t *testing.T) {
	tests := []struct {
		name     string
		book     *Book
		expected string
	}{
		{
			name: "empty book",
			book: DefaultBook(),
			expected: `{
				"name": "",
				"description": "",
				"scan_depth": 0,
				"token_budget": 0,
				"recursive_scanning": false,
				"entries": []
			}`,
		},
		{
			name: "book with basic fields",
			book: &Book{
				Name:              "Test Book",
				Description:       "Test Description",
				ScanDepth:         5,
				TokenBudget:       1000,
				RecursiveScanning: true,
				Entries:           []*BookEntry{},
			},
			expected: `{
				"name": "Test Book",
				"description": "Test Description",
				"scan_depth": 5,
				"token_budget": 1000,
				"recursive_scanning": true,
				"entries": []
			}`,
		},
		{
			name: "book with extensions",
			book: &Book{
				Name:        "Test Book",
				Description: "Test Description",
				Extensions: map[string]any{
					"custom_field": "custom_value",
					"number_field": 42,
				},
				Entries: []*BookEntry{},
			},
			expected: `{
				"name": "Test Book",
				"description": "Test Description",
				"scan_depth": 0,
				"token_budget": 0,
				"recursive_scanning": false,
				"extensions": {
					"custom_field": "custom_value",
					"number_field": 42
				},
				"entries": []
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := sonicx.Config.Marshal(tt.book)
			require.NoError(t, err)

			var expectedMap, actualMap map[string]any
			err = sonicx.Config.UnmarshalFromString(tt.expected, &expectedMap)
			require.NoError(t, err)
			err = sonicx.Config.UnmarshalFromString(stringsx.FromBytes(data), &actualMap)
			require.NoError(t, err)

			assert.Equal(t, expectedMap, actualMap)
		})
	}
}

func TestBook_JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected *Book
	}{
		{
			name: "basic book",
			jsonData: `{
				"name": "Test Book",
				"description": "Test Description",
				"scan_depth": 5,
				"token_budget": 1000,
				"recursive_scanning": true,
				"entries": []
			}`,
			expected: &Book{
				Name:              "Test Book",
				Description:       "Test Description",
				ScanDepth:         5,
				TokenBudget:       1000,
				RecursiveScanning: true,
				Entries:           []*BookEntry{},
			},
		},
		{
			name: "book with extensions",
			jsonData: `{
				"name": "Test Book",
				"description": "Test Description",
				"scan_depth": 0,
				"token_budget": 0,
				"recursive_scanning": false,
				"extensions": {
					"custom_field": "custom_value",
					"number_field": 42
				},
				"entries": []
			}`,
			expected: &Book{
				Name:        "Test Book",
				Description: "Test Description",
				Extensions: map[string]any{
					"custom_field": "custom_value",
					"number_field": float64(42), // JSON numbers unmarshal to float64
				},
				Entries: []*BookEntry{},
			},
		},
		{
			name: "minimal book",
			jsonData: `{
				"name": "",
				"description": "",
				"scan_depth": 0,
				"token_budget": 0,
				"recursive_scanning": false,
				"entries": []
			}`,
			expected: &Book{
				Name:              property.String(""),
				Description:       property.String(""),
				ScanDepth:         0,
				TokenBudget:       0,
				RecursiveScanning: false,
				Entries:           []*BookEntry{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var book Book
			err := sonicx.Config.UnmarshalFromString(tt.jsonData, &book)
			require.NoError(t, err)

			assert.Equal(t, tt.expected.Name, book.Name)
			assert.Equal(t, tt.expected.Description, book.Description)
			assert.Equal(t, tt.expected.ScanDepth, book.ScanDepth)
			assert.Equal(t, tt.expected.TokenBudget, book.TokenBudget)
			assert.Equal(t, tt.expected.RecursiveScanning, book.RecursiveScanning)
			assert.Equal(t, tt.expected.Extensions, book.Extensions)
			assert.Equal(t, len(tt.expected.Entries), len(book.Entries))
		})
	}
}

func TestBook_JSONRoundTrip(t *testing.T) {
	original := &Book{
		Name:              "Round Trip Book",
		Description:       "Testing round trip",
		ScanDepth:         10,
		TokenBudget:       2000,
		RecursiveScanning: true,
		Extensions: map[string]any{
			"test_field": "test_value",
			"numeric":    123,
		},
		Entries: []*BookEntry{},
	}

	// Marshal to JSON
	data, err := sonicx.Config.Marshal(original)
	require.NoError(t, err)

	// Unmarshal back to struct
	var roundTrip Book
	err = sonicx.Config.UnmarshalFromString(stringsx.FromBytes(data), &roundTrip)
	require.NoError(t, err)

	// Compare
	assert.Equal(t, original.Name, roundTrip.Name)
	assert.Equal(t, original.Description, roundTrip.Description)
	assert.Equal(t, original.ScanDepth, roundTrip.ScanDepth)
	assert.Equal(t, original.TokenBudget, roundTrip.TokenBudget)
	assert.Equal(t, original.RecursiveScanning, roundTrip.RecursiveScanning)
	assert.Equal(t, len(original.Entries), len(roundTrip.Entries))

	// RawExtensions comparison (accounting for JSON number conversion)
	assert.Equal(t, original.Extensions["test_field"], roundTrip.Extensions["test_field"])
	assert.Equal(t, float64(123), roundTrip.Extensions["numeric"]) // JSON converts numbers to float64
}

func TestBook_JSONInvalidData(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{
			name:     "wrong type for scan_depth",
			jsonData: `{"scan_depth": "not_a_number"}`,
		},
		{
			name:     "wrong type for recursive_scanning",
			jsonData: `{"recursive_scanning": "not_a_bool"}`,
		},
		{
			name:     "wrong type for token_budget",
			jsonData: `{"token_budget": "random"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book := DefaultBook()
			err := sonicx.Config.UnmarshalFromString(tt.jsonData, &book)
			assert.NoError(t, err)
		})
	}
}

func TestBookConstants(t *testing.T) {
	assert.Equal(t, " -- ", BookNameSeparator)
	assert.Equal(t, "\n----------------------\n", BookDescriptionSeparator)
	assert.Equal(t, "<<||-@PLACEHOLDER@-||>>", BookNamePlaceholder)
}
