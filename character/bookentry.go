package character

import (
	"maps"

	"github.com/r3dpixel/card-parser/property"
	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/stringsx"
)

// BookEntry lorebook entry structure
type bookEntryAlias BookEntry
type BookEntry struct {
	BookEntryCore
	RawExtensions map[string]any      `json:"-"`
	Extensions    BookEntryExtensions `json:"extensions"`
}

type bookEntryWrapper struct {
	BookEntryCore
	Extensions map[string]any `json:"extensions"`
}

type BookEntryCore struct {
	ID             property.Union       `json:"id"`
	Keys           property.StringArray `json:"keys"`
	SecondaryKeys  property.StringArray `json:"secondary_keys"`
	Name           property.String      `json:"name"`
	Comment        property.String      `json:"comment"`
	Content        property.String      `json:"content"`
	Constant       property.Bool        `json:"constant"`
	Selective      property.Bool        `json:"selective"`
	InsertionOrder property.Integer     `json:"insertion_order"`
	Enabled        property.Bool        `json:"enabled"`
	UseRegex       property.Bool        `json:"use_regex"`
}

// DefaultBookEntry returns and empty lorebook entry with ZERO value fields
func DefaultBookEntry() *BookEntry {
	return &BookEntry{
		BookEntryCore: BookEntryCore{
			ID:             property.Union{},
			Keys:           []string{},
			SecondaryKeys:  []string{},
			Name:           property.String(stringsx.Empty),
			Comment:        property.String(stringsx.Empty),
			Content:        property.String(stringsx.Empty),
			Constant:       false,
			Selective:      false,
			InsertionOrder: 10,
			Enabled:        true,
			UseRegex:       true,
		},
		RawExtensions: nil,
		Extensions:    DefaultBookEntryExtensions(),
	}
}

// FilledBookEntry returns a lorebook entry with default configuration
func FilledBookEntry(entryName, entryDescription string) *BookEntry {
	entry := DefaultBookEntry()
	entry.Keys = property.StringArray{entryName}
	entry.Name = property.String(entryName)
	entry.Comment = property.String(entryName)
	entry.Content = property.String(entryDescription)
	return entry
}

// MirrorNameAndComment assures that the comment/name of the entry are consistent
func (e *BookEntry) MirrorNameAndComment() {
	// Check the blank status of each field
	nameIsBlank := stringsx.IsBlank(string(e.Name))
	commentIsBlank := stringsx.IsBlank(string(e.Comment))

	switch {
	case nameIsBlank && !commentIsBlank:
		// If the name is blank and the comment is not, copy the comment to name
		e.Name = e.Comment
	case !nameIsBlank && commentIsBlank:
		// If the comment is blank and the name is not, copy the name to comment.
		e.Comment = e.Name
	}
}

func (e *BookEntry) MarshalJSON() ([]byte, error) {
	temp := bookEntryWrapper{}
	temp.BookEntryCore = e.BookEntryCore

	knownExtensions, err := jsonx.StructToMap(e.Extensions)
	if err != nil {
		return nil, err
	}
	if e.RawExtensions != nil {
		temp.Extensions = maps.Clone(e.RawExtensions)
		for k, v := range knownExtensions {
			temp.Extensions[k] = v
		}
	} else {
		temp.Extensions = knownExtensions
	}

	return sonicx.Config.Marshal(&temp)
}

func (e *BookEntry) UnmarshalJSON(data []byte) error {
	*e = *DefaultBookEntry()
	ref := stringsx.FromBytes(data)
	if err := sonicx.Config.UnmarshalFromString(ref, (*bookEntryAlias)(e)); err != nil {
		return err
	}
	var rawMap map[string]any
	if err := sonicx.Config.UnmarshalFromString(ref, &rawMap); err != nil {
		return err
	}
	extensionsMap, ok := rawMap["extensions"].(map[string]any)

	if caseSensitive, straggler := stragglerKey(EntryCaseSensitive, rawMap, extensionsMap); straggler {
		jsonx.HandlePrimitiveValue(caseSensitive, &e.Extensions.CaseSensitive)
	}
	if lorePosition, straggler := stragglerKey(EntryPosition, rawMap, extensionsMap); straggler {
		jsonx.HandleEntityValue(lorePosition, &e.Extensions.LorePosition)
	}
	if probability, straggler := stragglerKey(EntryProbability, rawMap, extensionsMap); straggler {
		jsonx.HandlePrimitiveValue(probability, &e.Extensions.Probability)
	}
	if selectiveLogic, straggler := stragglerKey(EntrySelectiveLogic, rawMap, extensionsMap); straggler {
		jsonx.HandleEntityValue(selectiveLogic, &e.Extensions.SelectiveLogic)
	}
	if role, straggler := stragglerKey(EntryRole, rawMap, extensionsMap); straggler {
		jsonx.HandleEntityValue(role, &e.Extensions.Role)
	}

	if ok {
		for _, fieldName := range bookEntryExtensionFields {
			delete(extensionsMap, fieldName)
		}
		e.RawExtensions = extensionsMap
	}

	return nil
}

func stragglerKey(key BookEntryExtension, entryMap map[BookEntryExtension]any, extensionsMap map[BookEntryExtension]any) (any, bool) {
	topLevelValue, isTopLevel := entryMap[key]
	_, isExtension := extensionsMap[key]
	if !isExtension && isTopLevel {
		return topLevelValue, true
	}
	return nil, false
}
