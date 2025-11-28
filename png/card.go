package png

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/png"
	"io"
	"os"

	"github.com/r3dpixel/card-parser/character"
	"github.com/r3dpixel/toolkit/filex"
)

// RawCard encoded chara PNG card
type RawCard struct {
	pngData
	RawCharaData []byte
	Revision     character.Revision
}

// RawJsonCard encoded chara PNG card with JSON data
type RawJsonCard struct {
	pngData
	RawJsonData []byte
	Revision    character.Revision
}

// CharacterCard decoded chara PNG card
type CharacterCard struct {
	pngData
	*character.Sheet
}

// PlaceholderCharacterCard returns a placeholder character card of the given size (black PNG image)
func PlaceholderCharacterCard(size int) (*RawCard, error) {
	// Create a new black image
	img := image.NewGray(image.Rect(0, 0, size, size))

	// Encode to PNG bytes
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	// Return the RawCard
	return FromImage(io.NopCloser(&buf)).First().Get()
}

// ToRawJson converts a RawCard to a RawJsonCard by decoding the base64 data
func (rc *RawCard) ToRawJson() (*RawJsonCard, error) {
	// Create a new RawJsonCard
	rawJsonCard := &RawJsonCard{
		pngData:  rc.pngData,
		Revision: rc.Revision,
	}

	// If there is no chara data, return the RawJsonCard as is
	if len(rc.RawCharaData) == 0 {
		return rawJsonCard, nil
	}

	// Decode chara data from base64
	decodedJSON := make([]byte, base64.StdEncoding.DecodedLen(len(rc.RawCharaData)))
	n, err := base64.StdEncoding.Decode(decodedJSON, rc.RawCharaData)
	if err != nil {
		return nil, err
	}

	// Set the JSON data in the RawJsonCard
	rawJsonCard.RawJsonData = decodedJSON[:n]

	// Return the RawJsonCard
	return rawJsonCard, nil
}

// ToCharacter converts a RawJsonCard to a CharacterCard by parsing the JSON data
func (rjc *RawJsonCard) ToCharacter() (*CharacterCard, error) {
	// Create a new CharacterCard
	characterCard := &CharacterCard{
		pngData: rjc.pngData,
	}

	// If there is no JSON data, return a default sheet
	if len(rjc.RawJsonData) == 0 {
		characterCard.Sheet = character.DefaultSheet(character.RevisionV2)
		return characterCard, nil
	}

	// Decode chara data from JSON into a Sheet
	sheet, err := character.FromBytes(rjc.RawJsonData)
	if err != nil {
		return nil, err
	}

	// Set the correct spec/version
	stamp := character.Stamps[rjc.Revision]
	sheet.Revision = rjc.Revision
	sheet.Spec = stamp.Spec
	sheet.Version = stamp.Version

	// Set the sheet in the CharacterCard
	characterCard.Sheet = sheet

	// Return the CharacterCard
	return characterCard, nil
}

// ToRawJson converts a CharacterCard to a RawJsonCard by serializing the Sheet to JSON
func (cc *CharacterCard) ToRawJson() (*RawJsonCard, error) {
	rawJsonCard := &RawJsonCard{
		pngData: cc.pngData,
	}

	if cc.Sheet == nil {
		return rawJsonCard, nil
	}

	// Encode the sheet to JSON byte slice
	jsonData, err := cc.Sheet.ToBytes()
	if err != nil {
		return nil, err
	}

	rawJsonCard.RawJsonData = jsonData
	rawJsonCard.Revision = cc.Sheet.Revision

	return rawJsonCard, nil
}

// ToRaw converts a RawJsonCard to a RawCard by encoding the JSON data as base64
func (rjc *RawJsonCard) ToRaw() *RawCard {
	// Create a new RawCard
	rawCard := &RawCard{
		pngData:  rjc.pngData,
		Revision: rjc.Revision,
	}

	// If there is no JSON data, return the RawCard as is
	if len(rjc.RawJsonData) == 0 {
		return rawCard
	}

	// Encode the JSON byte slice to base64
	encodedJSON := make([]byte, base64.StdEncoding.EncodedLen(len(rjc.RawJsonData)))
	base64.StdEncoding.Encode(encodedJSON, rjc.RawJsonData)

	// Set the base64 data in the RawCard
	rawCard.RawCharaData = encodedJSON

	// Return the RawCard
	return rawCard
}

// Decode converts a RawCard to a CharacterCard by decoding the base64 character data
func (rc *RawCard) Decode() (*CharacterCard, error) {
	// Decode the character data from base64
	rjc, err := rc.ToRawJson()
	if err != nil {
		return nil, err
	}
	// Decode the JSON data into a Sheet
	return rjc.ToCharacter()
}

// Encode converts a CharacterCard to a RawCard by encoding the character data as base64
func (cc *CharacterCard) Encode() (*RawCard, error) {
	// Encode the JSON data into a RawJsonCard
	rjc, err := cc.ToRawJson()
	if err != nil {
		return nil, err
	}
	// Encode the RawJsonCard to a RawCard
	return rjc.ToRaw(), nil
}

// ToImage writes the RawCard as a PNG image to the provided writer
func (rc *RawCard) ToImage(w io.Writer) error {
	// Write the header of the image first
	if _, err := w.Write(rc.Header); err != nil {
		return err
	}

	// Write the chara chunk
	if err := rc.streamCharaChunk(w, rc.Revision); err != nil {
		return err
	}

	// Write the image body
	_, err := w.Write(rc.Body)

	// Return
	return err
}

// ToFile saves the RawCard as a PNG image file at the specified path
func (rc *RawCard) ToFile(path string) error {
	// Open a file io.Writer
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filex.FilePermission)
	if err != nil {
		return err
	}
	// Defer the writer closing
	defer file.Close()

	// Write the image to the file
	return rc.ToImage(file)
}

// ToBytes returns the RawCard as a PNG image byte slice
func (rc *RawCard) ToBytes() ([]byte, error) {
	// Create a byte buffer
	buf := new(bytes.Buffer)
	// Write the image to the byte buffer
	if err := rc.ToImage(buf); err != nil {
		return nil, err
	}
	// Return the byte slice
	return buf.Bytes(), nil
}

// streamCharaChunk writes the character data chunk to the PNG stream
func (rc *RawCard) streamCharaChunk(w io.Writer, revision character.Revision) error {
	// If there is no chara data return empty byte slice
	if len(rc.RawCharaData) == 0 {
		return nil
	}

	// Write the correct chara keyword (fallback to V2)
	keyword := keywords[revision]
	if keyword == nil {
		keyword = keywords[character.RevisionV2]
	}

	// Write the correct PNG chunk length
	chunkDataLen := uint32(len(keyword) + len(rc.RawCharaData))
	if err := binary.Write(w, binary.BigEndian, chunkDataLen); err != nil {
		return err
	}

	// Create a new crc hasher
	crcHasher := crc32.NewIEEE()
	// Stream the writings to the output, as well as to the crc hasher
	multiWriter := io.MultiWriter(w, crcHasher)

	// Write the PNG chunk `tEXt` type
	if err := binary.Write(multiWriter, binary.BigEndian, chunkTextTypeCode); err != nil {
		return err
	}

	// Write the chara keyword
	if _, err := multiWriter.Write(keyword); err != nil {
		return err
	}

	// Write the chara data
	if _, err := multiWriter.Write(rc.RawCharaData); err != nil {
		return err
	}

	// Write the crc hash
	return binary.Write(w, binary.BigEndian, crcHasher.Sum32())
}
