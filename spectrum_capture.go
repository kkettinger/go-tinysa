package tinysa

import (
	"fmt"
	"image"
	"image/color"
)

// Capture returns the current screen of the device as image.Image type. The expected width and height is based
// on the detected model.
// TODO: refactor
func (d *Device) Capture() (image.Image, error) {
	d.logger.Info("capturing image")

	imgRaw, err := d.sendCommandBinary("capture")
	if err != nil {
		return nil, err
	}

	// we expect width * height * 2 bytes
	expectedLen := d.width * d.height * 2
	if expectedLen != len(imgRaw) {
		d.logger.Error("capture length mismatch (was the correct model detected?)",
			"got", len(imgRaw), "expected", expectedLen)
		return nil, fmt.Errorf("%w: capture length mismatch: expected %d, got %d (was the correct model detected?)",
			ErrCommandFailed, expectedLen, len(imgRaw))
	}

	// convert bin to rgba image
	img, err := convertBinToRGBA(imgRaw, d.width, d.height)
	if err != nil {
		d.logger.Error("failed to convert binary image to rgba", "err", err)
		return nil, fmt.Errorf("%w: failed to convert binary image to rgba: %w", ErrCommandFailed, err)
	}

	return img, nil
}

// convertBinToRGBA converts a binary image from the device to an image.Image type.
// TODO: refactor
func convertBinToRGBA(data []byte, width int, height int) (image.Image, error) {
	expectedSize := width * height * 2
	dataLen := len(data)

	if dataLen != expectedSize {
		newData := make([]byte, expectedSize)
		copy(newData, data) // Handles both padding and truncation
		data = newData
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset := ((y/2)*(2*width) + (y%2)*width + x) * 2

			if offset+1 >= len(data) {
				// Handle out-of-bounds access, which should be rare but possible due to padding
				img.Set(x, y, color.RGBA{A: 255}) // Default to black in case of error
				continue
			}

			pixel := uint16(data[offset])<<8 | uint16(data[offset+1])
			img.Set(x, y, convertRGB565ToRGBA(pixel))
		}
	}

	return img, nil
}

// convertRGB565ToRGBA converts a RGBA565 pixel to RGBA.
// TODO: refactor
func convertRGB565ToRGBA(pixel uint16) color.RGBA {
	// Extract RGB components from RGB565 format (RRRRRGGG GGGBBBBB)
	r := uint8((pixel >> 11) & 0x1F)
	g := uint8((pixel >> 5) & 0x3F)
	b := uint8(pixel & 0x1F)

	// Convert to 8-bit per channel with proper rounding
	return color.RGBA{
		R: uint8((uint16(r)*255 + 15) / 31),
		G: uint8((uint16(g)*255 + 31) / 63),
		B: uint8((uint16(b)*255 + 15) / 31),
		A: 255,
	}
}
