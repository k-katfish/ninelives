package stego

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func EmbedLSB(coverFile, outFile string, payload []byte) error {
	file, err := os.Open(coverFile)
	if err != nil {
		return fmt.Errorf("%w: %s", err, coverFile)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("%w: %s", err, coverFile)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	totalPixels := width * height
	maxBytes := totalPixels * 3 / 8

	if len(payload)+4 > maxBytes {
		return fmt.Errorf("input too large for output. max %d bytes, got %d bytes", maxBytes-4, len(payload))
	}

	payloadWithLength := make([]byte, len(payload)+4)
	payloadWithLength[0] = byte(len(payload) >> 24)
	payloadWithLength[1] = byte(len(payload) >> 16)
	payloadWithLength[2] = byte(len(payload) >> 8)
	payloadWithLength[3] = byte(len(payload))
	copy(payloadWithLength[4:], payload)

	outImg := image.NewRGBA(bounds)
	bitIndex := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			if bitIndex < len(payloadWithLength)*8 {
				r8 = (r8 & 0xFE) | (payloadWithLength[bitIndex/8] >> (7 - bitIndex%8) & 0x01)
				bitIndex++
			}
			if bitIndex < len(payloadWithLength)*8 {
				g8 = (g8 & 0xFE) | (payloadWithLength[bitIndex/8] >> (7 - bitIndex%8) & 0x01)
				bitIndex++
			}
			if bitIndex < len(payloadWithLength)*8 {
				b8 = (b8 & 0xFE) | (payloadWithLength[bitIndex/8] >> (7 - bitIndex%8) & 0x01)
				bitIndex++
			}
			outImg.Set(x, y, color.RGBA{r8, g8, b8, a8})
		}
	}

	outputFile, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("%w: %s", err, outFile)
	}
	defer outputFile.Close()

	if strings.ToLower(format) == "png" {
		return png.Encode(outputFile, outImg)
	}
	return jpeg.Encode(outputFile, outImg, &jpeg.Options{Quality: 100})
}
