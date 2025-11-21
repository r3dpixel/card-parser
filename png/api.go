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

type criteria func(rawCard *RawCard, chunk []byte, revision character.Revision) bool

func isLarger(rawCard *RawCard, chunk []byte, revision character.Revision) bool {
	return len(chunk)-keywordsLength[revision] >= len(rawCard.RawCharaData)
}

func isHigherVersion(rawCard *RawCard, chunk []byte, revision character.Revision) bool {
	return revision >= rawCard.Revision
}

type ScanMode struct {
	deepScan bool
	criteria criteria
}

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
	header := make([]byte, fullIhdrSize)
	if _, err := io.ReadFull(r, header); err != nil {
		return &converterProcessor{reader: io.MultiReader(bytes.NewReader(header), r), closer: r.Close}
	}

	if !slices.Equal(header[0:headerSize], pngHeader) {
		return &converterProcessor{reader: io.MultiReader(bytes.NewReader(header), r), closer: r.Close}
	}

	return newScanningProcessor(header, r)
}

// FromFile creates a Processor from a PNG file at the given path
func FromFile(path string) Processor {
	f, err := os.Open(path)
	if err != nil {
		return &converterProcessor{err: err}
	}
	return FromImage(f)
}

// FromBytes creates a Processor from a byte slice containing PNG image data
func FromBytes(data []byte) Processor {
	return FromImage(io.NopCloser(bytes.NewReader(data)))
}

// FromURL creates a Processor by fetching a PNG image from the given URL
func FromURL(c *reqx.Client, urls ...string) Processor {
	var fetchErr error

	for _, url := range urls {
		response, err := c.R().SetHeader("Accept", "image/png").Get(url)
		if err == nil {
			return FromImage(response.Body)
		}
		fetchErr = err
	}

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
