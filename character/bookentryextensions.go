package character

import (
	"github.com/r3dpixel/card-parser/property"
	"github.com/r3dpixel/toolkit/jsonx"
)

// BookEntryExtension is a string alias for the extension keys of a book entry
type BookEntryExtension = string

const (
	EntryPosition        BookEntryExtension = "position"
	EntryProbability     BookEntryExtension = "probability"
	EntryDepth           BookEntryExtension = "depth"
	EntrySelectiveLogic  BookEntryExtension = "selectiveLogic"
	EntryMatchWholeWords BookEntryExtension = "match_whole_words"
	EntryCaseSensitive   BookEntryExtension = "case_sensitive"
	EntryRole            BookEntryExtension = "role"
	EntrySticky          BookEntryExtension = "sticky"
	EntryCooldown        BookEntryExtension = "cooldown"
	EntryDelay           BookEntryExtension = "delay"
)

const (
	DefaultEntryProbability float64 = 100.00 // Default probability for entries
	DefaultEntryDepth       int     = 4      // Default depth for entries
)

// bookEntryExtensionFields is a helper variable that extracts the field names from BookEntryExtensions (typed extension struct)
var bookEntryExtensionFields = jsonx.ExtractJsonFieldNames(BookEntryExtensions{})

// BookEntryExtensions is a typed struct for extensions that can be added to a BookEntry
type BookEntryExtensions struct {
	LorePosition    property.LorePosition   `json:"position"`
	Probability     property.Float          `json:"probability"`
	Depth           property.Integer        `json:"depth"`
	SelectiveLogic  property.SelectiveLogic `json:"selectiveLogic"`
	MatchWholeWords property.Bool           `json:"match_whole_words"`
	CaseSensitive   property.Bool           `json:"case_sensitive"`
	Role            property.Role           `json:"role"`
	Sticky          property.Integer        `json:"sticky"`
	Cooldown        property.Integer        `json:"cooldown"`
	Delay           property.Integer        `json:"delay"`
}

// DefaultBookEntryExtensions returns an initialized BookEntryExtensions struct with default values
func DefaultBookEntryExtensions() BookEntryExtensions {
	return BookEntryExtensions{
		LorePosition:    property.DefaultLorePosition,
		Probability:     property.Float(DefaultEntryProbability),
		Depth:           property.Integer(DefaultEntryDepth),
		SelectiveLogic:  property.DefaultSelectiveLogic,
		MatchWholeWords: false,
		CaseSensitive:   false,
		Role:            property.DefaultRole,
		Sticky:          0,
		Cooldown:        0,
		Delay:           0,
	}
}
