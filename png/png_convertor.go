package png

import (
	"bytes"
	"io"

	jpeg "github.com/gen2brain/jpegli"
	"github.com/sunshineplan/imgconv"
)

// converterProcessor converts the card image from any format to PNG
type converterProcessor struct {
	reader  io.Reader
	closer  func() error
	decoded bool
	pngData pngData
	err     error
}

// ScanMode returns the processor itself as it doesn't support scanning
func (p *converterProcessor) ScanMode(scanMode ScanMode) Processor {
	return p
}

// First returns the processor itself as it doesn't support scanning'
func (p *converterProcessor) First() Processor {
	return p.ScanMode(First)
}

// LastVersion returns the processor itself as it doesn't support scanning'
func (p *converterProcessor) LastVersion() Processor {
	return p.ScanMode(LastVersion)
}

// LastLongest returns the processor itself as it doesn't support scanning'
func (p *converterProcessor) LastLongest() Processor {
	return p.ScanMode(LastLongest)
}

// Err returns any error that occurred during processing
func (p *converterProcessor) Err() error {
	return p.err
}

// ImageSize returns the width and height of the converted image
func (p *converterProcessor) ImageSize() (int, int) {
	// Decode the image
	p.decode()
	// If there was an error return -1, -1
	if p.err != nil {
		return -1, -1
	}

	// Return the width and height
	return widthPNG(p.pngData.Header), heightPNG(p.pngData.Header)
}

// Get returns a RawCard from the converted image data
func (p *converterProcessor) Get() (*RawCard, error) {
	// Decode the image
	p.decode()
	if p.err != nil {
		return nil, p.err
	}

	// Return the raw card
	return &RawCard{
		pngData: p.pngData,
	}, nil
}

// Close closes the underlying reader
func (p *converterProcessor) Close() error {
	return p.closer()
}

// decode converts the image data to PNG format if not already decoded
func (p *converterProcessor) decode() {
	// If there is an error or the card was already deocded return
	if p.err != nil || p.decoded {
		return
	}

	// Read all from the input
	data, err := io.ReadAll(p.reader)
	if err != nil {
		p.err = err
		return
	}

	// Decode image
	img, err := imgconv.Decode(bytes.NewReader(data))
	if err != nil {
		// If decoding fails try specialized decoding from jpeg (in case abnormal chrome subsampling)
		img, err = jpeg.Decode(bytes.NewReader(data))
	}
	// If all decoders have failed, return the error
	if err != nil {
		p.err = err
		return
	}

	// Allocate byte buffer
	var buf bytes.Buffer
	// Convert to PNG
	option := imgconv.FormatOption{Format: imgconv.PNG}
	if err = option.Encode(&buf, img); err != nil {
		p.err = err
		return
	}

	// Set a decoded flag to true
	p.decoded = true

	// Set the correct png data
	p.pngData = pngData{
		Header: buf.Next(fullIhdrSize),
		Body:   buf.Bytes(),
	}
}
