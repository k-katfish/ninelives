package stego

import (
	"errors"
	"fmt"
	"image"
	"io"
)

func ExtractLSB(reader io.Reader) ([]byte, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	bits := make([]byte, 0, width*height*3)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			bits = append(bits, byte(r)&0x01)
			bits = append(bits, byte(g)&0x01)
			bits = append(bits, byte(b)&0x01)
		}
	}

	if len(bits) < 32 {
		return nil, errors.New("image too small or no data found")
	}
	length := 0
	for i := 0; i < 32; i++ {
		length = (length << 1) | int(bits[i])
	}

	totalBits := length * 8
	if len(bits) < totalBits {
		return nil, errors.New("not enough data in image for expected payload length")
	}

	payload := make([]byte, length)
	for i := 0; i < length; i++ {
		for j := 0; j < 8; j++ {
			payload[i] = (payload[i] << 1) | bits[32+i*8+j]
		}
	}
	return payload, nil
}
