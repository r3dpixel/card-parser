package png

import (
	"bytes"
	"encoding/binary"
	"io"
	"slices"

	"github.com/r3dpixel/card-parser/character"
	"github.com/r3dpixel/toolkit/bytex"
)

// Sizes in bytes
const (
	headerSize   int = 8  // Size of the PNG standard header in bytes
	ihdrSize     int = 25 // Size of the IHDR header in bytes
	widthSize    int = 4  // Size of the value holding the Width of the PNG in bytes
	heightSize   int = 4  // Size of the value holding the Width of the PNG in bytes
	footerSize   int = 12 // Size of the PNG standard footer in bytes
	fullIhdrSize     = headerSize + ihdrSize

	chunkLengthSize int = 4 // Size of the chunks' length fraction in bytes
	chunkTypeSize   int = 4 // Size of the chunkDetails's type discriminator of in bytes
	chunkCrcSize    int = 4 // Size of the CRC32 checksum in bytes
	chunkHeaderSize     = chunkLengthSize + chunkTypeSize + chunkCrcSize

	charaKeywordSize int = 6 // Size of the 'chara' keyword in bytes
	ccv3KeywordSize  int = 5 // Size of the 'ccv3' keyword in bytes

	minimumSize = headerSize + ihdrSize + footerSize // Minimum size of a PNG in byte

	Extension       string = ".png" // The PNG file extension
	ExtensionLength int    = 4      // PNG extension length

	ihdrWidthOffset  = headerSize + chunkLengthSize + chunkTypeSize
	ihdrHeightOffset = headerSize + chunkLengthSize + chunkTypeSize + widthSize
)

// Byte arrays
var (
	// The standard PNG header (byte array)
	pngHeader = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	// Discriminator 'tEXt' (uint32) - 0x74455874
	chunkTextTypeCode uint32 = 0x74455874
	// Discriminator 'IEND' (uint32) - 0x49454E44
	chunkIENDTypeCode = 0x49454E44
	// 'chara' keyword (byte array)
	charaKeyword = []byte{0x63, 0x68, 0x61, 0x72, 0x61, 0x00}
	// 'ccv3' keyword (byte array)
	ccv3Keyword = []byte{0x63, 0x63, 0x76, 0x33, 0x00}
	// The standard PNG footer (byte array)
	pngFooter = []byte{0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}

	// List of keywords
	keywords = map[character.Revision][]byte{
		character.RevisionV2: charaKeyword,
		character.RevisionV3: ccv3Keyword,
	}
	keywordsLength = map[character.Revision]int{
		character.RevisionV2: len(charaKeyword),
		character.RevisionV3: len(ccv3Keyword),
	}
)

// scanningProcessor implements the Processor interface and is used to scan PNG files for character data
type scanningProcessor struct {
	// Scanner properties
	header   []byte
	reader   io.ReadCloser
	scanMode ScanMode

	// Scanner state and caches
	bodyBuffer   *bytes.Buffer
	chunkDetails chunkDetails
	chunkBuffer  []byte
	rawCard      *RawCard
	err          error
}

// chunkDetails holds the length and discriminator of a PNG chunk
type chunkDetails struct {
	length   uint32
	typeCode uint32
}

// newScanningProcessor creates a new PNG scanner processor
func newScanningProcessor(header []byte, r io.ReadCloser) *scanningProcessor {
	s := &scanningProcessor{
		header:   header,
		reader:   r,
		scanMode: DefaultScanMode,
	}
	return s
}

// ScanMode sets the scan mode for the processor
func (p *scanningProcessor) ScanMode(mode ScanMode) Processor {
	p.scanMode = mode
	return p
}

// First sets the processor to scan for the first chara chunk
func (p *scanningProcessor) First() Processor {
	p.scanMode = First
	return p
}

// LastVersion sets the processor to scan for the latest chara chunk (highest revision)
func (p *scanningProcessor) LastVersion() Processor {
	p.scanMode = LastVersion
	return p
}

// LastLongest sets the processor to scan for the longest chara chunk
func (p *scanningProcessor) LastLongest() Processor {
	p.scanMode = LastLongest
	return p
}

// Err returns any error that occurred during processing
func (p *scanningProcessor) Err() error {
	return p.err
}

// ImageSize returns the width and height of the PNG image
func (p *scanningProcessor) ImageSize() (int, int) {
	if p.err != nil {
		return -1, -1
	}
	return widthPNG(p.header), heightPNG(p.header)
}

// Get processes the PNG and returns a RawCard with extracted character data
func (p *scanningProcessor) Get() (*RawCard, error) {
	defer p.reader.Close()

	// If there is an error return error
	if p.err != nil {
		return nil, p.err
	}

	// Allocate new byte buffers
	p.bodyBuffer = bytes.NewBuffer(make([]byte, 0, 32*bytex.KiB))

	// Set the correct image header
	p.rawCard = &RawCard{
		pngData: pngData{
			Header: p.header,
		},
	}

	// Process PNG chunks
	for {
		// Process the PNG chunk
		err := p.processChunk()
		// If EOF, copy any remaining data, and set the body
		if err == io.EOF {
			// Copy remaining data
			if _, copyErr := io.Copy(p.bodyBuffer, p.reader); copyErr != nil {
				return nil, copyErr
			}
			// Set the body
			p.rawCard.Body = p.bodyBuffer.Bytes()
			// Return the raw card
			return p.rawCard, nil
		}
		// If any other error occurred, return error
		if err != nil {
			return nil, err
		}
	}
}

func (p *scanningProcessor) Close() error {
	return p.reader.Close()
}

// processChunk processes a single PNG chunk and extracts character data if present
func (p *scanningProcessor) processChunk() error {
	// Read the PNG chunk length
	if err := binary.Read(p.reader, binary.BigEndian, &p.chunkDetails.length); err != nil {
		return err
	}

	// Read the PNG chunk discriminator
	if err := binary.Read(p.reader, binary.BigEndian, &p.chunkDetails.typeCode); err != nil {
		return err
	}

	// If the PNG chunk IS NOT a `tEXt` chunk, stream copy it directly to the output
	if p.chunkDetails.typeCode != chunkTextTypeCode {
		return p.streamCopyChunk()
	}

	// Reset the buffer
	p.chunkBuffer = p.chunkBuffer[:0]
	// If the buffer is not large enough, allocate a new one
	if int(p.chunkDetails.length) > cap(p.chunkBuffer) {
		p.chunkBuffer = make([]byte, p.chunkDetails.length)
	}
	// Resize the buffer
	p.chunkBuffer = p.chunkBuffer[:p.chunkDetails.length]

	// Read chunk data
	if _, err := io.ReadFull(p.reader, p.chunkBuffer); err != nil {
		return err
	}

	// Discard the CRC hash
	if _, err := io.CopyN(io.Discard, p.reader, 4); err != nil {
		return err
	}

	// Check if the PNG chunks contains chara data
	revision, isChara := p.isCharaChunk(p.chunkBuffer)
	// If not discard chunk
	if !isChara {
		return nil
	}

	// Check if chara chunk revision is higher than the current revision
	if p.scanMode.criteria(p.rawCard, p.chunkBuffer, revision) {
		p.rawCard.Revision = revision
		p.rawCard.RawCharaData = slices.Clone(p.chunkBuffer[keywordsLength[revision]:])
	}

	// If deep scan is disabled, and we have found a chara chunk return io.EOF so the rest is stream copied
	if !p.scanMode.deepScan && len(p.rawCard.RawCharaData) > 0 {
		return io.EOF
	}

	return nil
}

// streamCopyChunk copies a non-character chunk to the output stream
func (p *scanningProcessor) streamCopyChunk() error {
	// Write the PNG chunk length
	if err := binary.Write(p.bodyBuffer, binary.BigEndian, p.chunkDetails.length); err != nil {
		return err
	}

	// Write the PNG chunk discriminator
	if err := binary.Write(p.bodyBuffer, binary.BigEndian, p.chunkDetails.typeCode); err != nil {
		return err
	}

	// Write the PNG chunk content and the CRC hash
	if _, err := io.CopyN(p.bodyBuffer, p.reader, int64(p.chunkDetails.length)+4); err != nil {
		return err
	}

	// Return
	return nil
}

// isCharaChunk checks if chunk data contains character information and returns the revision
func (p *scanningProcessor) isCharaChunk(chunkData []byte) (character.Revision, bool) {
	// Return false (no chara data)
	if len(chunkData) == 0 {
		return character.RevisionV2, false
	}

	// Detect the correct revision and keyword
	for revision, keyword := range keywords {
		if bytes.HasPrefix(chunkData, keyword) {
			return revision, true
		}
	}

	// No chara keyword detected
	return character.RevisionV2, false
}
