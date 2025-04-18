package tinysa

import (
	"fmt"
	"image"
	"image/color"
)

// Capture returns the current screen of the device as image.Image type.
func (d *Device) Capture() (image.Image, error) {
	d.logger.Info("capturing image")

	imgRaw, err := d.sendCommandBinary("capture")
	if err != nil {
		return nil, err
	}

	// Convert bin to RGBA image
	img, err := convertBinCaptureToImage(imgRaw, d.width, d.height)
	if err != nil {
		d.logger.Error("failed to convert binary image", "err", err)
		return nil, fmt.Errorf("%w: failed to convert binary image: %v", ErrCommandFailed, err)
	}

	return img, nil
}

// convertBinCaptureToImage converts a binary image from the device to an image.Image type.
func convertBinCaptureToImage(data []byte, width int, height int) (image.Image, error) {
	expectedSize := width * height * 2
	if len(data) != expectedSize {
		return nil, fmt.Errorf("expected size %d (%d * %d * 2), got %d", expectedSize, width, height, len(data))
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := (y*width + x) * 2
			pixel := uint16(data[i])<<8 | uint16(data[i+1]) // big-endian
			img.Set(x, y, convertRGB565PixelToRGBA(pixel))
		}
	}

	return img, nil
}

// convertRGB565PixelToRGBA converts a RGBA565 pixel to RGBA.
func convertRGB565PixelToRGBA(pixel uint16) color.RGBA {
	// Extract RGB components from RGB565 format (RRRRRGGG GGGBBBBB)
	r := (pixel >> 11) & 0x1F
	g := (pixel >> 5) & 0x3F
	b := pixel & 0x1F

	// Convert to 8-bit per channel with proper rounding
	return color.RGBA{
		R: uint8((r*255 + 15) / 31),
		G: uint8((g*255 + 31) / 63),
		B: uint8((b*255 + 15) / 31),
		A: 255,
	}
}
