package character

import (
	"cmp"
	"io"

	gcmp "github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/r3dpixel/card-parser/property"
	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/stringsx"
)

// cmpOptions are used to compare Sheets
var cmpOptions = []gcmp.Option{
	cmpopts.EquateEmpty(),
	cmpopts.SortSlices(comparator[string]),
	cmpopts.SortSlices(comparator[int]),
	cmpopts.SortSlices(comparator[float64]),
	cmpopts.SortSlices(comparator[property.String]),
	cmpopts.SortSlices(comparator[property.Integer]),
	cmpopts.SortSlices(comparator[property.Float]),
}

const (
	// CreatorNotesSeparator is the separator between different parts of the merged creator notes
	CreatorNotesSeparator = "\n\n"
	// AnonymousCreator is the name used for anonymous creators
	AnonymousCreator = "Anonymous"
)

// sheetWrapper is used to wrap the Sheet content in a JSON object for marshaling and unmarshalling
type sheetWrapper struct {
	Spec    Spec     `json:"spec"`
	Version Version  `json:"spec_version"`
	Content *Content `json:"data"`
}

// Sheet structure of a V3 chara card
type Sheet struct {
	Spec     Spec
	Version  Version
	Revision Revision
	Content
}

// DefaultSheet returns an empty chara sheet with the given Revision
func DefaultSheet(revision Revision) *Sheet {
	// Create a new sheet
	sheet := &Sheet{}
	// Set the revision
	sheet.SetRevision(revision)
	// Return the sheet
	return sheet
}

// MarshalJSON marshals Sheet into JSON format with Content wrapped under "data" using Sonic
func (s *Sheet) MarshalJSON() ([]byte, error) {
	// Wrap the content in a JSON object
	wrapper := sheetWrapper{
		Spec:    s.Spec,
		Version: s.Version,
		Content: &s.Content,
	}
	// Encode the JSON object using Sonic
	return sonicx.Config.Marshal(&wrapper)
}

// UnmarshalJSON decode a chara sheet from JSON using Sonic
func (s *Sheet) UnmarshalJSON(data []byte) error {
	// Decode the JSON object using Sonic
	wrap, err := sonicx.GetFromString(stringsx.FromBytes(data))
	if err != nil {
		return err
	}

	// Extract metadata without copying
	spec := wrap.GetByPath("spec").String()
	version := wrap.GetByPath("spec_version").String()
	rawData := wrap.GetByPath("data").Raw()
	if err := sonicx.Config.UnmarshalFromString(rawData, &s.Content); err != nil {
		return err
	}

	// Set the correct revision, spec and version
	revision := RevisionV2
	if spec == string(SpecV3) || version == string(V3) {
		revision = RevisionV3
	}
	s.SetRevision(revision)

	// Decoding complete
	return nil
}

// SetRevision sets the sheet revision, spec and version
func (s *Sheet) SetRevision(revision Revision) {
	// Get the correct stamp
	stamp := Stamps[revision]

	// Set the revision, spec and version
	s.Revision = revision
	s.Spec = stamp.Spec
	s.Version = stamp.Version
}

// ToJSON converts the sheet to its JSON representation and writes it to the given output io.Writer using Sonic streaming
func (s *Sheet) ToJSON(w io.Writer, opts ...jsonx.Options) error {
	return jsonx.ToJSON(s, w, opts...)
}

// ToFile converts the sheet to its JSON representation and writes it to the given output file destination
func (s *Sheet) ToFile(path string, opts ...jsonx.Options) error {
	return jsonx.ToFile(s, path, opts...)
}

// ToBytes converts the sheet to its JSON representation and returns the JSON byte slice
func (s *Sheet) ToBytes(opts ...jsonx.Options) ([]byte, error) {
	return jsonx.ToBytes(s, opts...)
}

// DeepEquals returns true if the two sheets are deeply equal
func (s *Sheet) DeepEquals(other *Sheet) bool {
	return gcmp.Equal(s, other, cmpOptions...)
}

// FromJSON decodes the JSON from the given input io.Reader and returns the decoded sheet using Sonic streaming
func FromJSON(r io.Reader) (*Sheet, error) {
	return jsonx.FromJSON[*Sheet](r)
}

// FromFile decodes the JSON from the given input file and returns the decoded sheet
func FromFile(path string) (*Sheet, error) {
	return jsonx.FromFile[*Sheet](path)
}

// FromBytes decodes the JSON from the given input byte slice and returns the decoded sheet
func FromBytes(b []byte) (*Sheet, error) {
	return jsonx.FromBytes[*Sheet](b)
}

// comparator is used to compare slices of any type
func comparator[T cmp.Ordered](a, b T) bool {
	return a < b
}
