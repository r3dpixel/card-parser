package character

import (
	"testing"

	"github.com/r3dpixel/card-parser/property"
	"github.com/stretchr/testify/assert"
)

func TestNewBookMerger(t *testing.T) {
	merger := NewBookMerger()
	assert.NotNil(t, merger)
	assert.NotNil(t, merger.book)
	assert.NotNil(t, merger.nameBuilder)
	assert.NotNil(t, merger.descriptionBuilder)
	assert.Zero(t, merger.entryIndex)
}

func TestBookMerger_AppendProperties(t *testing.T) {
	testCases := []struct {
		name                string
		initialScanDepth    int
		initialTokenBudget  int
		initialRecursive    bool
		appendScanDepth     int
		appendTokenBudget   int
		appendRecursive     bool
		expectedScanDepth   int
		expectedTokenBudget int
		expectedRecursive   bool
	}{
		{"Higher values win", 100, 200, false, 150, 250, true, 150, 250, true},
		{"Lower values do not change", 150, 250, true, 100, 200, false, 150, 250, true},
		{"Recursive stays true", 0, 0, true, 0, 0, true, 0, 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			merger := NewBookMerger()
			merger.book.ScanDepth = property.Integer(tc.initialScanDepth)
			merger.book.TokenBudget = property.Integer(tc.initialTokenBudget)
			merger.book.RecursiveScanning = property.Bool(tc.initialRecursive)

			merger.AppendProperties(tc.appendScanDepth, tc.appendTokenBudget, tc.appendRecursive)

			assert.Equal(t, tc.expectedScanDepth, int(merger.book.ScanDepth))
			assert.Equal(t, tc.expectedTokenBudget, int(merger.book.TokenBudget))
			assert.Equal(t, tc.expectedRecursive, bool(merger.book.RecursiveScanning))
		})
	}
}

func TestBookMerger_AppendNameAndDescription(t *testing.T) {
	merger := NewBookMerger()
	merger.AppendNameAndDescription("First Name", "First Desc")
	merger.AppendNameAndDescription("  ", "  ") // Blank should be ignored
	merger.AppendNameAndDescription("Second Name", "Second Desc")
	merger.AppendEntry(FilledBookEntry("one", "one"))
	book := merger.Build()
	assert.NotNil(t, book)
	assert.Equal(t, "First Name -- Second Name", string(book.Name))
	assert.Equal(t, "First Desc\n----------------------\nSecond Desc", string(book.Description))
}

func TestBookMerger_AppendEntry(t *testing.T) {
	merger := NewBookMerger()
	entry1 := &BookEntry{BookEntryCore: BookEntryCore{Name: "entry1"}}
	entry2 := &BookEntry{BookEntryCore: BookEntryCore{Name: "entry2"}}

	merger.AppendEntry(entry1)
	merger.AppendEntry(entry2)

	assert.Len(t, merger.book.Entries, 2)
	assert.Equal(t, 0, *merger.book.Entries[0].ID.IntValue)
	assert.Equal(t, 1, *merger.book.Entries[1].ID.IntValue)
	// Check if mirroring was called (comment should equal name)
	assert.Equal(t, merger.book.Entries[0].Name, merger.book.Entries[0].Comment)
}

func TestBookMerger_AppendMapExtensions(t *testing.T) {
	merger := NewBookMerger()
	merger.AppendMapExtensions(map[string]any{"key1": "val1"})
	merger.AppendMapExtensions(map[string]any{"key1": "SHOULD_NOT_OVERWRITE", "key2": "val2"})

	assert.Len(t, merger.book.Extensions, 2)
	assert.Equal(t, "val1", merger.book.Extensions["key1"])
	assert.Equal(t, "val2", merger.book.Extensions["key2"])
}

func TestBookMerger_AppendBook(t *testing.T) {
	t.Run("Appends book with existing entries", func(t *testing.T) {
		merger := NewBookMerger()
		bookToAppend := DefaultBook()
		bookToAppend.Name = "Test Book"
		bookToAppend.Entries = append(bookToAppend.Entries, &BookEntry{BookEntryCore: BookEntryCore{Content: "entry content"}})

		merger.AppendBook(bookToAppend)

		assert.Len(t, merger.book.Entries, 1)
		assert.Equal(t, "entry content", string(merger.book.Entries[0].Content))
		book := merger.Build()
		assert.Equal(t, "Test Book", string(book.Name))
	})

	t.Run("Appends book with no entries, creating a default entry", func(t *testing.T) {
		merger := NewBookMerger()
		bookToAppend := DefaultBook()
		bookToAppend.Name = "Book With No Entries"
		bookToAppend.Description = "Desc With No Entries"

		merger.AppendBook(bookToAppend)

		assert.Len(t, merger.book.Entries, 1)
		defaultEntry := merger.book.Entries[0]
		assert.Equal(t, "Book With No Entries", string(defaultEntry.Name))
		assert.Equal(t, "Desc With No Entries", string(defaultEntry.Content))
	})
}

func TestBookMerger_Build(t *testing.T) {
	t.Run("Build returns nil if no entries", func(t *testing.T) {
		merger := NewBookMerger()
		book := merger.Build()
		assert.Nil(t, book)
	})

	t.Run("Build returns a complete book", func(t *testing.T) {
		merger := NewBookMerger()
		merger.AppendNameAndDescription("Final Name", "Final Desc")
		merger.AppendEntry(&BookEntry{})

		book := merger.Build()

		assert.NotNil(t, book)
		assert.Equal(t, "Final Name", string(book.Name))
		assert.Equal(t, "Final Desc", string(book.Description))
		assert.Len(t, book.Entries, 1)
	})
}

func TestBookMerger_AppendEntries(t *testing.T) {
	merger := NewBookMerger()
	entries := []*BookEntry{
		{BookEntryCore: BookEntryCore{Name: "entry1", Content: "content1"}},
		{BookEntryCore: BookEntryCore{Name: "entry2", Content: "content2"}},
	}

	merger.AppendEntries(entries)

	assert.Len(t, merger.book.Entries, 2)
	assert.Equal(t, "content1", string(merger.book.Entries[0].Content))
	assert.Equal(t, "content2", string(merger.book.Entries[1].Content))
	assert.Equal(t, 0, *merger.book.Entries[0].ID.IntValue)
	assert.Equal(t, 1, *merger.book.Entries[1].ID.IntValue)
}

func TestBookMerger_EdgeCases(t *testing.T) {
	t.Run("Build with empty name and description", func(t *testing.T) {
		merger := NewBookMerger()
		merger.AppendEntry(&BookEntry{BookEntryCore: BookEntryCore{Content: "test"}})

		book := merger.Build()
		assert.NotNil(t, book)
		assert.Empty(t, book.Name)
		assert.Empty(t, book.Description)
	})

	t.Run("Build with whitespace-only name and description", func(t *testing.T) {
		merger := NewBookMerger()
		merger.AppendNameAndDescription("   ", "\t\n ")
		merger.AppendEntry(&BookEntry{BookEntryCore: BookEntryCore{Content: "test"}})

		book := merger.Build()
		assert.NotNil(t, book)
		assert.Empty(t, book.Name)
		assert.Empty(t, book.Description)
	})

	t.Run("AppendMapExtensions with nil extensions", func(t *testing.T) {
		merger := NewBookMerger()
		merger.AppendMapExtensions(nil)
		assert.Nil(t, merger.book.Extensions)
	})

	t.Run("AppendMapExtensions with empty map", func(t *testing.T) {
		merger := NewBookMerger()
		merger.AppendMapExtensions(map[string]any{})
		assert.Nil(t, merger.book.Extensions)
	})
}

func TestTokenAppender(t *testing.T) {
	t.Run("newTokenAppender creates correct instance", func(t *testing.T) {
		appender := newTokenAppender(" | ")
		assert.Equal(t, " | ", appender.separator)
		assert.Equal(t, 0, appender.tokenIndex)
		assert.Equal(t, 0, appender.nonEmptyTokenIndex)
		assert.Empty(t, appender.get())
	})

	t.Run("appendToken with single token", func(t *testing.T) {
		appender := newTokenAppender(" | ")
		appender.appendToken("first")
		assert.Equal(t, "first", appender.get())
		assert.Equal(t, 1, appender.tokenIndex)
		assert.Equal(t, 1, appender.nonEmptyTokenIndex)
	})

	t.Run("appendToken with multiple tokens", func(t *testing.T) {
		appender := newTokenAppender(" | ")
		appender.appendToken("first")
		appender.appendToken("second")
		appender.appendToken("third")
		assert.Equal(t, "first | second | third", appender.get())
		assert.Equal(t, 3, appender.tokenIndex)
		assert.Equal(t, 3, appender.nonEmptyTokenIndex)
	})

	t.Run("appendToken ignores empty and whitespace tokens", func(t *testing.T) {
		appender := newTokenAppender(" -- ")
		appender.appendToken("first")
		appender.appendToken("")      // empty
		appender.appendToken("   ")   // whitespace only
		appender.appendToken("\t\n ") // whitespace only
		appender.appendToken("second")
		assert.Equal(t, "first -- second", appender.get())
		assert.Equal(t, 5, appender.tokenIndex)
		assert.Equal(t, 2, appender.nonEmptyTokenIndex)
	})

	t.Run("appendToken trims whitespace", func(t *testing.T) {
		appender := newTokenAppender(" | ")
		appender.appendToken("  first  ")
		appender.appendToken("\tsecond\n")
		assert.Equal(t, "first | second", appender.get())
	})

	t.Run("appendToken with custom separator", func(t *testing.T) {
		appender := newTokenAppender("\n----------------------\n")
		appender.appendToken("paragraph 1")
		appender.appendToken("paragraph 2")
		expected := "paragraph 1\n----------------------\nparagraph 2"
		assert.Equal(t, expected, appender.get())
	})

	t.Run("appendToken with all empty tokens", func(t *testing.T) {
		appender := newTokenAppender(" | ")
		appender.appendToken("")
		appender.appendToken("   ")
		appender.appendToken("\t")
		assert.Empty(t, appender.get())
		assert.Equal(t, 3, appender.tokenIndex)
		assert.Equal(t, 0, appender.nonEmptyTokenIndex)
	})
}
