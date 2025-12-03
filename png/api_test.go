package png

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/r3dpixel/card-parser/character"
	"github.com/r3dpixel/card-parser/property"
	"github.com/r3dpixel/toolkit/reqx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	minimalIHDR = []byte{
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89,
	}

	testCards = struct {
		smallV2  *character.Sheet
		largeV3  *character.Sheet
		tinyV2   *character.Sheet
		firstV2  *character.Sheet
		secondV2 *character.Sheet
	}{
		smallV2:  createSheet(character.RevisionV2, "Small V2"),
		largeV3:  createSheet(character.RevisionV3, "Much larger V3 card with extra data to make it clearly bigger than the V2 card for testing longest chunk mode"),
		tinyV2:   createSheet(character.RevisionV2, "Tiny"),
		firstV2:  createSheet(character.RevisionV2, "First Card"),
		secondV2: createSheet(character.RevisionV2, "Last  Card"), // Same length as firstV2
	}
)

func createSheet(revision character.Revision, name string) *character.Sheet {
	sheet := &character.Sheet{
		Revision: character.RevisionV2,
		Spec:     character.SpecV2,
		Version:  character.V2,
		Content:  character.Content{Name: property.String(name)},
	}
	if revision == character.RevisionV3 {
		sheet.Revision = character.RevisionV3
		sheet.Spec = character.SpecV3
		sheet.Version = character.V3
	}
	return sheet
}

func encodeCardData(t *testing.T, sheet *character.Sheet) []byte {
	t.Helper()
	cardJSON, err := sheet.ToBytes()
	require.NoError(t, err)
	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(cardJSON)))
	base64.StdEncoding.Encode(b64, cardJSON)
	return b64
}

func createTestPNG(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(width-1, height-1, color.RGBA{B: 255, A: 255})
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	require.NoError(t, err)
	return buf.Bytes()
}

func createTestJPG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	require.NoError(t, err)
	return buf.Bytes()
}

// injectSingleChunk creates a PNG with a single character chunk
func injectSingleChunk(t *testing.T, basePNG []byte, sheet *character.Sheet, atEnd bool) []byte {
	t.Helper()
	data := encodeCardData(t, sheet)
	return injectChunk(t, basePNG, sheet.Revision, data, atEnd)
}

// injectDoubleChunk creates a PNG with two character chunks
func injectDoubleChunk(t *testing.T, basePNG []byte, first, second *character.Sheet) []byte {
	t.Helper()
	firstData := encodeCardData(t, first)
	secondData := encodeCardData(t, second)
	withFirst := injectChunk(t, basePNG, first.Revision, firstData, false)
	return injectChunk(t, withFirst, second.Revision, secondData, true)
}

func injectChunk(t *testing.T, pngBytes []byte, version character.Revision, data []byte, atEnd bool) []byte {
	t.Helper()
	keyword := keywords[version]
	require.NotNil(t, keyword)

	// Use streaming approach like production code
	buf := new(bytes.Buffer)
	chunkDataLen := uint32(len(keyword) + len(data))

	// Write chunk header directly
	require.NoError(t, binary.Write(buf, binary.BigEndian, chunkDataLen))
	require.NoError(t, binary.Write(buf, binary.BigEndian, chunkTextTypeCode))

	// Write data with CRC calculation
	crcHasher := crc32.NewIEEE()
	multiWriter := io.MultiWriter(buf, crcHasher)
	_, err := multiWriter.Write(keyword)
	require.NoError(t, err)
	_, err = multiWriter.Write(data)
	require.NoError(t, err)

	// Write CRC
	require.NoError(t, binary.Write(buf, binary.BigEndian, crcHasher.Sum32()))

	charaChunk := buf.Bytes()

	if atEnd {
		iendStart := len(pngBytes) - footerSize
		return slices.Concat(pngBytes[:iendStart], charaChunk, pngBytes[iendStart:])
	}
	injectionPoint := headerSize + ihdrSize
	return slices.Concat(pngBytes[:injectionPoint], charaChunk, pngBytes[injectionPoint:])
}

func TestProcessor_Constructors(t *testing.T) {
	pngBytes := createTestPNG(t, 4, 4)
	jpgBytes := createTestJPG(t)
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.png")
	err := os.WriteFile(filePath, pngBytes, 0644)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ok" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(pngBytes)
			return
		}
		if r.URL.Path == "/404" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	t.Run("Happy Paths", func(t *testing.T) {
		testCases := []struct {
			name      string
			processor Processor
		}{
			{"FromFile", FromFile(filePath)},
			{"FromURL", FromURL(reqx.NewClient(reqx.Options{}), server.URL+"/ok")},
			{"FromImagePNG", FromBytes(pngBytes)},
			{"FromImageJPG", FromBytes(jpgBytes)},
			{"FromImage", FromBytes(pngBytes)},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := tc.processor.Get()
				require.NoError(t, err)
				require.NoError(t, tc.processor.Err())
			})
		}
	})

	t.Run("Error Paths", func(t *testing.T) {
		unreadableFile := filepath.Join(tempDir, "unreadable.png")
		require.NoError(t, os.WriteFile(unreadableFile, []byte{}, 0000))

		testCases := []struct {
			name      string
			processor Processor
		}{
			{"FromFile non-existent", FromFile("non-existent.png")},
			{"FromFile unreadable", FromFile(unreadableFile)},
			{"FromURL network error", FromURL(reqx.NewClient(reqx.Options{}), "http://localhost:99999")},
			{"FromURL non-200 status", FromURL(reqx.NewClient(reqx.Options{}), server.URL+"/404")},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				assert.Error(t, tc.processor.Err())
			})
		}
	})
}

func TestFromURL_MultipleURLs(t *testing.T) {
	pngBytes := createTestPNG(t, 4, 4)
	// Disable retries for cleaner test (1 = try once, no retries)
	client := reqx.NewClient(reqx.Options{
		RetryCount: 1,
	})

	// Track which URLs were accessed
	var accessLog []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessLog = append(accessLog, r.URL.Path)

		switch r.URL.Path {
		case "/success":
			w.Header().Set("Content-Type", "image/png")
			w.WriteHeader(http.StatusOK)
			w.Write(pngBytes)
		case "/failure":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	t.Run("first URL succeeds - should not try subsequent URLs", func(t *testing.T) {
		accessLog = nil
		processor := FromURL(client, server.URL+"/success", server.URL+"/failure")
		_, err := processor.Get()

		require.NoError(t, err)
		require.NoError(t, processor.Err())
		assert.Equal(t, []string{"/success"}, accessLog, "should only access first URL")
	})

	t.Run("first URL fails, second succeeds", func(t *testing.T) {
		accessLog = nil
		processor := FromURL(client, server.URL+"/failure", server.URL+"/success")
		_, err := processor.Get()

		require.NoError(t, err)
		require.NoError(t, processor.Err())
		// Client retries failed requests once, so first URL hit twice, then success
		assert.Equal(t, []string{"/failure", "/failure", "/success"}, accessLog, "should try first URL twice then succeed on second")
	})

	t.Run("all URLs fail - returns last error", func(t *testing.T) {
		accessLog = nil
		processor := FromURL(client, server.URL+"/failure", server.URL+"/404", server.URL+"/500")

		assert.Error(t, processor.Err())
		assert.GreaterOrEqual(t, len(accessLog), 1, "should have tried at least one URL")
	})

	t.Run("multiple failures then success", func(t *testing.T) {
		accessLog = nil
		processor := FromURL(client,
			server.URL+"/404",
			server.URL+"/failure",
			server.URL+"/success",
			server.URL+"/never-reached",
		)
		_, err := processor.Get()

		require.NoError(t, err)
		require.NoError(t, processor.Err())
		// Each failed URL gets 2 attempts (1 retry), then success on third URL
		assert.Equal(t, []string{"/404", "/404", "/failure", "/failure", "/success"}, accessLog, "should stop at first success")
	})
}

func TestProcessor_ScanModes(t *testing.T) {
	basePNG := createTestPNG(t, 4, 4)

	tests := []struct {
		name     string
		data     []byte
		scanMode ScanMode
		want     *character.Sheet
	}{
		{
			name:     "First - single V2",
			data:     injectSingleChunk(t, basePNG, testCards.smallV2, false),
			scanMode: First,
			want:     testCards.smallV2,
		},
		{
			name:     "First - single V3",
			data:     injectSingleChunk(t, basePNG, testCards.largeV3, false),
			scanMode: First,
			want:     testCards.largeV3,
		},
		{
			name:     "First - V2 then V3 (takes first)",
			data:     injectDoubleChunk(t, basePNG, testCards.smallV2, testCards.largeV3),
			scanMode: First,
			want:     testCards.smallV2,
		},
		{
			name:     "HighestVersion - V3 over V2",
			data:     injectDoubleChunk(t, basePNG, testCards.smallV2, testCards.largeV3),
			scanMode: LastVersion,
			want:     testCards.largeV3,
		},
		{
			name:     "HighestVersion - V3 with V2 data",
			data:     injectSingleChunk(t, basePNG, testCards.largeV3, false),
			scanMode: LastVersion,
			want:     testCards.largeV3,
		},
		{
			name:     "LongestChunk - large V3 beats tiny V2",
			data:     injectDoubleChunk(t, basePNG, testCards.tinyV2, testCards.largeV3),
			scanMode: LastLongest,
			want:     testCards.largeV3,
		},
		{
			name:     "LongestChunk - small V2 beats tiny V2",
			data:     injectDoubleChunk(t, basePNG, testCards.tinyV2, testCards.smallV2),
			scanMode: LastLongest,
			want:     testCards.smallV2,
		},
		{
			name:     "LongestChunk - same size prefers last",
			data:     injectDoubleChunk(t, basePNG, testCards.firstV2, testCards.secondV2),
			scanMode: LastLongest,
			want:     testCards.secondV2,
		},
		{
			name:     "LongestChunk - larger V2 first beats smaller V3",
			data:     injectDoubleChunk(t, basePNG, testCards.largeV3, testCards.tinyV2),
			scanMode: LastLongest,
			want:     testCards.largeV3,
		},
		{
			name:     "No metadata - PNG",
			data:     basePNG,
			scanMode: First,
			want:     nil,
		},
		{
			name:     "No metadata - JPG",
			data:     createTestJPG(t),
			scanMode: First,
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := FromBytes(tt.data).ScanMode(tt.scanMode)
			rawCard, err := processor.Get()
			require.NoError(t, err)

			if tt.want == nil {
				assert.Empty(t, rawCard.RawCharaData)
				assert.Equal(t, character.Revision(0), rawCard.Revision)
			} else {
				assert.Equal(t, tt.want.Revision, rawCard.Revision)
				decodedCard, err := rawCard.Decode()
				require.NoError(t, err)
				require.NotNil(t, decodedCard.Sheet)
				assert.Equal(t, string(tt.want.Content.Name), string(decodedCard.Sheet.Content.Name))
			}
		})
	}

	t.Run("malformed chunk error", func(t *testing.T) {
		incompleteChunk := []byte{0x00, 0x00, 0x01}
		malformed := slices.Concat(pngHeader, minimalIHDR, incompleteChunk)
		_, err := FromBytes(malformed).Get()
		assert.Error(t, err)
	})
}
