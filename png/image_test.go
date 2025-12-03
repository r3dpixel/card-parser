package png

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPngDataTest(t *testing.T) *pngData {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 200, 100))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(199, 99, color.RGBA{B: 255, A: 255})

	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	require.NoError(t, err)
	pngBytes := buf.Bytes()

	headerEnd := headerSize + ihdrSize
	return &pngData{
		Header: pngBytes[:headerEnd],
		Body:   pngBytes[headerEnd:],
	}
}

func TestPngData_Dimensions(t *testing.T) {
	pd := setupPngDataTest(t)
	assert.Equal(t, 200, pd.Width(), "Width should be correctly read from the IHDR chunkDetails")
	assert.Equal(t, 100, pd.Height(), "Height should be correctly read from the IHDR chunkDetails")
}

func TestPngData_Image(t *testing.T) {
	pd := setupPngDataTest(t)
	img, err := pd.Image()
	require.NoError(t, err)
	require.NotNil(t, img)

	assert.Equal(t, 200, img.Bounds().Dx())
	assert.Equal(t, 100, img.Bounds().Dy())

	r, _, _, a := img.At(0, 0).RGBA()
	assert.Equal(t, uint32(0xffff), r, "Pixel at (0,0) should be red")
	assert.Equal(t, uint32(0xffff), a, "Pixel at (0,0) should be opaque")
}

func TestPngData_Thumbnail(t *testing.T) {
	pd := setupPngDataTest(t)
	thumbnailSize := 50

	thumb, err := pd.Thumbnail(thumbnailSize)
	require.NoError(t, err)
	require.NotNil(t, thumb)

	bounds := thumb.Bounds()
	assert.Equal(t, thumbnailSize, bounds.Dx(), "Thumbnail width should be the target size")
	assert.Equal(t, 25, bounds.Dy(), "Thumbnail height should be scaled proportionally")
}

func TestPngData_ScaleDown(t *testing.T) {
	pd := setupPngDataTest(t)
	scaleDownSize := 40

	assert.Equal(t, 200, pd.Width())
	assert.Equal(t, 100, pd.Height())

	err := pd.ScaleDown(scaleDownSize)
	require.NoError(t, err)

	assert.Equal(t, 40, pd.Width(), "Width should be updated to the new scaled size")
	assert.Equal(t, 20, pd.Height(), "Height should be updated proportionally")

	img, err := pd.Image()
	require.NoError(t, err)
	assert.Equal(t, 40, img.Bounds().Dx())
	assert.Equal(t, 20, img.Bounds().Dy())
}
