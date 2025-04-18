package tinysa

import (
	_ "golang.org/x/image/bmp"
	"image"
	"image/color"
	"os"
	"testing"
)

func TestConvertRGB565ToRGBA(t *testing.T) {
	// RGB565: Red = 0b11111, Green = 0b000000, Blue = 0b00000 (pure red)
	pixel := uint16(0b1111100000000000)
	rgba := convertRGB565PixelToRGBA(pixel)

	expected := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	if rgba != expected {
		t.Fatalf("expected %v, got %v", expected, rgba)
	}
}

func TestConvertBinToRGBA(t *testing.T) {
	width, height := 2, 2
	// 2x2 image = 4 pixels -> 8 bytes (2 bytes per pixel)
	// Fill with 0xF800 (pure red in RGB565) -> bytes: [0xF8, 0x00] * 4
	data := []byte{
		0xF8, 0x00,
		0xF8, 0x00,
		0xF8, 0x00,
		0xF8, 0x00,
	}

	img, err := convertBinCaptureToImage(data, width, height)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 || a>>8 != 255 {
				t.Errorf("pixel at (%d, %d) unexpected value: R=%d, G=%d, B=%d, A=%d",
					x, y, r>>8, g>>8, b>>8, a>>8)
			}
		}
	}
}

func TestCompareWithGoldenBMP(t *testing.T) {
	width, height := 480, 320
	bin, err := os.ReadFile("testdata/capture_480x320.bin")
	if err != nil {
		t.Fatal(err)
	}

	refImgFile, err := os.Open("testdata/capture_480x320.bmp")
	if err != nil {
		t.Fatal(err)
	}
	defer refImgFile.Close()

	refImg, _, err := image.Decode(refImgFile)
	if err != nil {
		t.Fatal(err)
	}

	testImg, err := convertBinCaptureToImage(bin, width, height)
	if err != nil {
		t.Fatal(err)
	}

	// Pixel-by-pixel comparison
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r1, g1, b1, _ := refImg.At(x, y).RGBA()
			r2, g2, b2, _ := testImg.At(x, y).RGBA()

			if r1 != r2 || g1 != g2 || b1 != b2 {
				t.Errorf("pixel mismatch at (%d,%d): BMP=(%d,%d,%d), got=(%d,%d,%d)",
					x, y, r1>>8, g1>>8, b1>>8, r2>>8, g2>>8, b2>>8)
			}
		}
	}
}
