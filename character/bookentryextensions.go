package character

import (
	"github.com/r3dpixel/card-parser/property"
	"github.com/r3dpixel/toolkit/jsonx"
)

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
	DefaultEntryProbability float64 = 100.00
	DefaultEntryDepth       int     = 4
)

var bookEntryExtensionFields = jsonx.ExtractJsonFieldNames(BookEntryExtensions{})

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
