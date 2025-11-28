package png

import (
	"bytes"
	"image"
	"image/png"
	"io"

	"github.com/sunshineplan/imgconv"
)

// pngData PNG image data
type pngData struct {
	Header []byte
	Body   []byte
}

// Width returns the width in pixels of the PNG
func (p *pngData) Width() int {
	return widthPNG(p.Header)
}

// Height returns the height in pixels of the PNG
func (p *pngData) Height() int {
	return heightPNG(p.Header)
}

// Thumbnail Create a thumbnail from the image of the raw context
func (p *pngData) Thumbnail(size int) (image.Image, error) {
	// FromBytes the image from the raw p
	imageSource, err := p.Image()
	if err != nil {
		return nil, err
	}
	// Return the scaled-down image (to the down scale size)
	return resizeImage(imageSource, size), nil
}

// ScaleDown Scale down the png image
func (p *pngData) ScaleDown(size int) error {
	// Decode the image
	imageSource, err := p.Image()
	if err != nil {
		return err
	}

	// Scale down the image
	downScaledImageSource := resizeImage(imageSource, size)

	// Encode the scaled-down image to PNG bytes
	writer := new(bytes.Buffer)
	err = png.Encode(writer, downScaledImageSource)
	if err != nil {
		return err
	}

	// Extract the header and body from the writer
	p.Header = writer.Next(headerSize + ihdrSize)
	p.Body = writer.Bytes()

	// Return nil (success)
	return nil
}

// Image FromBytes just the image from the raw context
func (p *pngData) Image() (image.Image, error) {
	// Use the prefix data and suffix data to reconstruct the image bytes (eliminates all the metadata)
	imageByteReader := io.MultiReader(bytes.NewReader(p.Header), bytes.NewReader(p.Body))
	// Decode the image from the image bytes
	imageSource, _, err := image.Decode(imageByteReader)

	// If the extraction failed, return nil
	if err != nil {
		return nil, err
	}

	// Return the decoded image
	return imageSource, nil
}

// resizeImage Resize the image to fit a square based on a given size
func resizeImage(image image.Image, size int) image.Image {
	// The scaled-down image should always fit into a square of size
	resizeOption := imgconv.ResizeOption{}
	width := image.Bounds().Dx()
	height := image.Bounds().Dy()
	// Scale down on either width or height
	if width > height {
		// If the width is larger than the height, scale down is based on width
		resizeOption.Width = size
	} else {
		// If the height is larger than the width, scale down is based on height
		resizeOption.Height = size
	}

	// Return thumbnail
	return imgconv.Resize(image, &resizeOption)
}
