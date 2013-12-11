package mnist

import (
	"image/color"
	"reflect"
	"testing"
)

func TestRemapColor(t *testing.T) {
	in, expected := uint8(0x00), uint32(0x00000000)
	if out := remapColor(in); out != expected {
		t.Fatalf("0x%02x got remapped to 0x%04x, expected 0x%04x", in, out, expected)
	}

	in, expected = uint8(0xFF), uint32(0x0000FFFF)
	if out := remapColor(in); out != expected {
		t.Fatalf("0x%02x got remapped to 0x%04x, expected 0x%04x", in, out, expected)
	}
}

func TestMnistModel_Convert(t *testing.T) {
	mnBlack := MnistModel.Convert(color.Black)
	mnWhite := MnistModel.Convert(color.White)
	if c, ok := mnWhite.(mnistColor); !ok || uint8(c) != 0 {
		t.Fatalf("White did not get convorted to MNIST color properly")
	}
	if c, ok := mnBlack.(mnistColor); !ok || uint8(c) != 255 {
		t.Fatalf("White did not get convorted to MNIST color properly")
	}
	if r, g, b, a := mnWhite.RGBA(); !reflect.DeepEqual([]uint32{r, g, b, a}, []uint32{0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF}) {
		t.Fatalf("White did not get converted back to RGBA properly")
	}
	if r, g, b, a := mnBlack.RGBA(); !reflect.DeepEqual([]uint32{r, g, b, a}, []uint32{0x0000, 0x0000, 0x0000, 0xFFFF}) {
		t.Fatalf("White did not get converted back to RGBA properly")
	}
}
