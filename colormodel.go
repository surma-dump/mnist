package mnist

import (
	"image/color"
)

var MnistModel = color.ModelFunc(mnistModel)

type mnistColor uint8

func (mc mnistColor) RGBA() (r, g, b, a uint32) {
	return remapColor(255 - uint8(mc)), remapColor(255 - uint8(mc)), remapColor(255 - uint8(mc)), remapColor(255)
}

func mnistModel(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	avg := (r + g + b) / 3
	return mnistColor(255 - avg)
}

func remapColor(c uint8) uint32 {
	return uint32(float64(c) / float64(0xFF) * float64(0xFFFF))
}
