package character

import (
	"github.com/r3dpixel/card-parser/property"
)

const (
	BookNameSeparator        string = " -- "
	BookDescriptionSeparator string = "\n----------------------\n"
	BookNamePlaceholder             = `<<||-@PLACEHOLDER@-||>>`
)

// Book lorebook structure of a V3 chara card
type Book struct {
	Name              property.String  `json:"name"`
	Description       property.String  `json:"description"`
	ScanDepth         property.Integer `json:"scan_depth"`
	TokenBudget       property.Integer `json:"token_budget"`
	RecursiveScanning property.Bool    `json:"recursive_scanning"`
	Extensions        map[string]any   `json:"extensions,omitempty"`
	Entries           []*BookEntry     `json:"entries"`
}

// DefaultBook creates an empty book with an initialized entry list
func DefaultBook() *Book {
	return &Book{}
}

// NormalizeSymbols normalizes the book name and description, and all book entries
func (b *Book) NormalizeSymbols() {
	// Fix Quotes on the book name and description
	b.Name.NormalizeSymbols()
	b.Description.NormalizeSymbols()

	// Fix Quotes on the book entries (name, comment, content)
	// Other fields ARE NOT affected (keywords, secondary keywords, etc.)
	for _, entry := range b.Entries {
		entry.MirrorNameAndComment()
		entry.Name.NormalizeSymbols()
		entry.Comment.NormalizeSymbols()
		entry.Content.NormalizeSymbols()
	}
}
