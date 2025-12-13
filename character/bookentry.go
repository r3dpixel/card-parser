package character

import (
	"maps"

	"github.com/r3dpixel/card-parser/property"
	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/stringsx"
)

// bookEntryAlias is used to avoid circular references
type bookEntryAlias BookEntry

// BookEntry lorebook entry structure
type BookEntry struct {
	BookEntryCore
	RawExtensions map[string]any      `json:"-"`
	Extensions    BookEntryExtensions `json:"extensions"`
}

// bookEntryWrapper is used to marshal/unmarshal the BookEntry struct with the extension map
type bookEntryWrapper struct {
	BookEntryCore
	Extensions map[string]any `json:"extensions"`
}

// BookEntryCore lorebook entry core structure
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

// DefaultBookEntry returns and empty lorebook entry with default value fields
func DefaultBookEntry() *BookEntry {
	return &BookEntry{
		BookEntryCore: BookEntryCore{
			ID:             property.Union{},
			Keys:           []string{},
			SecondaryKeys:  []string{},
			Name:           property.String(""),
			Comment:        property.String(""),
			Content:        property.String(""),
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

// FilledBookEntry returns a lorebook entry with the given name and description (default values for other fields)
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

	// Mirror the comment to the name if the name is blank, and vice versa
	switch {
	case nameIsBlank && !commentIsBlank:
		// If the name is blank and the comment is not, copy the comment to name
		e.Name = e.Comment
	case !nameIsBlank && commentIsBlank:
		// If the comment is blank and the name is not, copy the name to comment.
		e.Comment = e.Name
	}
}

// MarshalJSON marshals the BookEntry struct to JSON
func (e *BookEntry) MarshalJSON() ([]byte, error) {
	// Copy the BookEntryCore struct to avoid circular references
	temp := bookEntryWrapper{}
	temp.BookEntryCore = e.BookEntryCore

	// Extract the typed extensions
	knownExtensions, err := jsonx.StructToMap(e.Extensions)
	if err != nil {
		return nil, err
	}

	// Merge the dynamic extension map with the known extensions
	if e.RawExtensions != nil {
		// Clone the raw extensions map to avoid modifying the original
		temp.Extensions = maps.Clone(e.RawExtensions)
		// Add the known extensions to the raw extensions map
		for k, v := range knownExtensions {
			temp.Extensions[k] = v
		}
	} else {
		// If the raw extensions map is nil, use the known extensions directly
		temp.Extensions = knownExtensions
	}

	// Marshal the BookEntryWrapper struct to JSON
	return sonicx.Config.Marshal(&temp)
}

// UnmarshalJSON unmarshals JSON data into the BookEntry struct
func (e *BookEntry) UnmarshalJSON(data []byte) error {
	// Initialize the BookEntry struct with default values
	*e = *DefaultBookEntry()

	// Convert to string without copying the underlying array
	ref := stringsx.FromBytes(data)

	// Unmarshal the alias book entry struct from JSON
	if err := sonicx.Config.UnmarshalFromString(ref, (*bookEntryAlias)(e)); err != nil {
		return err
	}

	// Unmarshal to a raw map as well (double unmarshalling necessary, unfortunately)
	var rawMap map[string]any
	if err := sonicx.Config.UnmarshalFromString(ref, &rawMap); err != nil {
		return err
	}

	// Extract the extension map
	extensionsMap, ok := rawMap["extensions"].(map[string]any)

	// Extract straggler keyed extensions
	// A straggler key extension is an extension that is mapped outside the extension map itself
	// Examples (in this case case_sensitive is a straggler since it's outside the extension map):
	// {
	//   "extensions": {
	//     "probability": 100
	//   }
	//   "case_sensitive": true
	// }
	// Extract case_sensitive from the top level map, if it exists
	if caseSensitive, straggler := stragglerKey(EntryCaseSensitive, rawMap, extensionsMap); straggler {
		jsonx.HandlePrimitiveValue(caseSensitive, &e.Extensions.CaseSensitive)
	}
	// Extract lore position from the top level map, if it exists
	if lorePosition, straggler := stragglerKey(EntryPosition, rawMap, extensionsMap); straggler {
		jsonx.HandleEntityValue(lorePosition, &e.Extensions.LorePosition)
	}
	// Extract probability from the top level map, if it exists
	if probability, straggler := stragglerKey(EntryProbability, rawMap, extensionsMap); straggler {
		jsonx.HandlePrimitiveValue(probability, &e.Extensions.Probability)
	}
	// Extract selective logic from the top level map, if it exists
	if selectiveLogic, straggler := stragglerKey(EntrySelectiveLogic, rawMap, extensionsMap); straggler {
		jsonx.HandleEntityValue(selectiveLogic, &e.Extensions.SelectiveLogic)
	}
	// Extract role from the top level map, if it exists
	if role, straggler := stragglerKey(EntryRole, rawMap, extensionsMap); straggler {
		jsonx.HandleEntityValue(role, &e.Extensions.Role)
	}

	// If an extension map was found, remove the typed extensions from the raw extensions map
	if ok {
		// Remove the typed extensions from the raw extensions map
		for _, fieldName := range bookEntryExtensionFields {
			delete(extensionsMap, fieldName)
		}
		// Set the raw extensions map on the BookEntry struct
		e.RawExtensions = extensionsMap
	}

	// Return nil (success)
	return nil
}

// stragglerKey returns the value of a key if it exists outside the extension map, and whether the key is a straggler
func stragglerKey(key BookEntryExtension, entryMap map[BookEntryExtension]any, extensionsMap map[BookEntryExtension]any) (any, bool) {
	// Check if the key exists in the top level map
	topLevelValue, isTopLevel := entryMap[key]
	// Check if the key exists in the extension map
	_, isExtension := extensionsMap[key]
	// If the key is not an extension and is at the top level then it is a straggler
	if !isExtension && isTopLevel {
		return topLevelValue, true
	}
	// Otherwise return nil and false
	return nil, false
}
