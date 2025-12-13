package png

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"slices"

	"github.com/r3dpixel/card-parser/character"
	"github.com/r3dpixel/toolkit/reqx"
)

// criteria defines the conditions for a chunk to be considered a valid PNG chara chunk
type criteria func(rawCard *RawCard, chunk []byte, revision character.Revision) bool

// isLarger checks if the chunk is larger than the raw chara data
func isLarger(rawCard *RawCard, chunk []byte, revision character.Revision) bool {
	return len(chunk)-keywordsLength[revision] >= len(rawCard.RawCharaData)
}

// isHigherVersion checks if the chunk revision is higher than the raw card revision
func isHigherVersion(rawCard *RawCard, chunk []byte, revision character.Revision) bool {
	return revision >= rawCard.Revision
}

// ScanMode defines the scan mode for PNG card decoding
type ScanMode struct {
	deepScan bool
	criteria criteria
}

// ScanMode values
var (
	First = ScanMode{
		deepScan: false,
		criteria: isLarger,
	}
	LastVersion = ScanMode{
		deepScan: true,
		criteria: isHigherVersion,
	}
	LastLongest = ScanMode{
		deepScan: true,
		criteria: isLarger,
	}
	DefaultScanMode = First
)

// Processor API for decoding chara PNG cards
type Processor interface {
	ScanMode(scanMode ScanMode) Processor
	First() Processor
	LastVersion() Processor
	LastLongest() Processor
	Err() error
	ImageSize() (int, int)
	Get() (*RawCard, error)
	Close() error
}

// FromImage creates a Processor from an io.Reader containing PNG image data
func FromImage(r io.ReadCloser) Processor {
	// Read the PNG header
	header := make([]byte, fullIhdrSize)
	// If the header cannot be read or is not long enough, return a converter processor
	if _, err := io.ReadFull(r, header); err != nil {
		return &converterProcessor{reader: io.MultiReader(bytes.NewReader(header), r), closer: r.Close}
	}
	// If the header does not match the PNG header, return a converter processor
	if !slices.Equal(header[0:headerSize], pngHeader) {
		return &converterProcessor{reader: io.MultiReader(bytes.NewReader(header), r), closer: r.Close}
	}
	// Return a scanning processor
	return newScanningProcessor(header, r)
}

// FromFile creates a Processor from a PNG file at the given path
func FromFile(path string) Processor {
	// Open the PNG file
	f, err := os.Open(path)
	if err != nil {
		return &converterProcessor{err: err}
	}
	// Return a processor from the file
	return FromImage(f)
}

// FromBytes creates a Processor from a byte slice containing PNG image data
func FromBytes(data []byte) Processor {
	// Return a processor from the byte slice
	return FromImage(io.NopCloser(bytes.NewReader(data)))
}

// FromURL creates a Processor by fetching a PNG image from the given URL
func FromURL(c *reqx.Client, urls ...string) Processor {
	// fetchErr will be the final error
	var fetchErr error

	// Loop through the URLs and fetch the image
	for _, url := range urls {
		// Fetch the image from the URL
		response, err := c.R().SetHeader("Accept", "image/png").Get(url)
		if err == nil {
			// Return a processor from the image
			return FromImage(response.Body)
		}
		// If there was an error, set it
		fetchErr = err
	}

	// Return a converter processor with the final error
	return &converterProcessor{err: fetchErr}
}

// widthPNG extracts the width from PNG header bytes
func widthPNG(bytes []byte) int {
	return int(binary.BigEndian.Uint32(bytes[ihdrWidthOffset : ihdrWidthOffset+widthSize]))
}

// heightPNG extracts the height from PNG header bytes
func heightPNG(bytes []byte) int {
	return int(binary.BigEndian.Uint32(bytes[ihdrHeightOffset : ihdrHeightOffset+heightSize]))
}
