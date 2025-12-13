package character

import (
	"strings"

	"github.com/r3dpixel/card-parser/property"
	"github.com/r3dpixel/toolkit/ptr"
	"github.com/r3dpixel/toolkit/stringsx"
)

// BookMerger merges multiple lorebooks through a safe API
type BookMerger struct {
	book               *Book
	nameBuilder        *tokenAppender
	descriptionBuilder *tokenAppender
	entryIndex         int
}

// NewBookMerger creates a new lorebook merger
func NewBookMerger() *BookMerger {
	merger := &BookMerger{
		book:               DefaultBook(),
		nameBuilder:        newTokenAppender(BookNameSeparator),
		descriptionBuilder: newTokenAppender(BookDescriptionSeparator),
		entryIndex:         0,
	}
	return merger
}

// AppendBook appends the given lorebook
func (bm *BookMerger) AppendBook(book *Book) {
	// If the book is nil, return (NO-OP)
	if book == nil {
		return
	}

	// If the book has no entries, move the description as the only entry
	if len(book.Entries) == 0 {
		// Set the only entry to the book description
		book.Entries = []*BookEntry{FilledBookEntry(string(book.Name), string(book.Description))}
		// Reset the book description
		book.Description = property.String("")
	}

	// Append book properties
	bm.AppendProperties(int(book.ScanDepth), int(book.TokenBudget), bool(book.RecursiveScanning))

	// Append book name and description
	bm.AppendNameAndDescription(string(book.Name), string(book.Description))

	// Append book extensions
	bm.AppendMapExtensions(book.Extensions)

	// Append book entries
	bm.AppendEntries(book.Entries)
}

// AppendProperties compute new properties of the merged book
func (bm *BookMerger) AppendProperties(scanDepth int, tokenBudget int, recursiveScanning bool) {
	// Scan depth will always be the maximum scan depth found
	bm.book.ScanDepth = max(bm.book.ScanDepth, property.Integer(scanDepth))
	// Token Budget will always have the maximum value found
	bm.book.TokenBudget = max(bm.book.TokenBudget, property.Integer(tokenBudget))
	// Recursive scanning will be turned on, if at least one book has recursive scanning
	bm.book.RecursiveScanning = bm.book.RecursiveScanning || property.Bool(recursiveScanning)
}

// AppendNameAndDescription appends the name and description to the merged book
func (bm *BookMerger) AppendNameAndDescription(name string, description string) {
	bm.nameBuilder.appendToken(name)
	bm.descriptionBuilder.appendToken(description)
}

// AppendEntries appends the given entries
func (bm *BookMerger) AppendEntries(entries []*BookEntry) {
	for _, entry := range entries {
		bm.AppendEntry(entry)
	}
}

// AppendEntry appends the given entry
func (bm *BookMerger) AppendEntry(entry *BookEntry) {
	// Mirror the name and comment for SillyTavern
	entry.MirrorNameAndComment()
	// Assign the entryIndex as the ID of the entry
	entry.ID = property.Union{IntValue: ptr.Of(bm.entryIndex)}
	// Append the entry to the merged book
	bm.book.Entries = append(bm.book.Entries, entry)
	// Increment the entry index for the next entry
	bm.entryIndex++
}

// AppendMapExtensions Append extension map
func (bm *BookMerger) AppendMapExtensions(extensions map[string]any) {
	// If the extensions map is empty, return (NO-OP)
	if len(extensions) == 0 {
		return
	}

	// Create a merged book extensions map (if it doesn't exist)
	if bm.book.Extensions == nil {
		bm.book.Extensions = make(map[string]any)
	}

	// Copy extensions into accumulator
	for k, v := range extensions {
		if _, duplicate := bm.book.Extensions[k]; !duplicate {
			bm.book.Extensions[k] = v
		}
	}
}

// Build builds the merged book
func (bm *BookMerger) Build() *Book {
	// If there are no entries, return nil (no book needed)
	if len(bm.book.Entries) == 0 {
		return nil
	}

	// Assign book name
	bm.book.Name = property.String(strings.TrimSpace(bm.nameBuilder.get()))

	// Assign book description
	bm.book.Description = property.String(strings.TrimSpace(bm.descriptionBuilder.get()))

	// Return merged book
	return bm.book
}

// tokenAppender token appender that handles adding separators between tokens automatically
type tokenAppender struct {
	stringBuilder      strings.Builder
	separator          string
	tokenIndex         int
	nonEmptyTokenIndex int
}

// newTokenAppender returns a new token appender
func newTokenAppender(separator string) *tokenAppender {
	return &tokenAppender{
		stringBuilder:      strings.Builder{},
		separator:          separator,
		tokenIndex:         0,
		nonEmptyTokenIndex: 0,
	}
}

// appendToken appends the given token to the current string, with respect to the separator
func (t *tokenAppender) appendToken(token string) {
	parsedToken := strings.TrimSpace(token)
	t.tokenIndex++

	// If the book name is empty, return
	if stringsx.IsBlank(parsedToken) {
		return
	}

	// If not first non-empty token adds separator
	if t.nonEmptyTokenIndex != 0 {
		t.stringBuilder.WriteString(t.separator)
	}

	// Add token
	t.stringBuilder.WriteString(parsedToken)
	t.nonEmptyTokenIndex++
}

// get returns the built string
func (t *tokenAppender) get() string {
	return t.stringBuilder.String()
}
