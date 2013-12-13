package mnist

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"reflect"
	"testing"
)

var (
	testImage = Image{
		Width:  2,
		Height: 2,
		Data:   []byte{0, 255, 255, 0},
	}

	testImageFile []byte

	black, white []uint32
)

func init() {
	r, g, b, a := color.Black.RGBA()
	black = []uint32{r, g, b, a}

	r, g, b, a = color.White.RGBA()
	white = []uint32{r, g, b, a}

	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, IMAGE_MAGIC_NUMBER)
	binary.Write(buf, binary.BigEndian, int32(2))
	binary.Write(buf, binary.BigEndian, int32(2))
	binary.Write(buf, binary.BigEndian, int32(3))
	binary.Write(buf, binary.BigEndian, []uint8{0, 0, 0, 0, 0, 0})
	binary.Write(buf, binary.BigEndian, []uint8{255, 255, 255, 255, 255, 255})
	testImageFile = buf.Bytes()
}

func TestImage_ColorModel(t *testing.T) {
	if cm := testImage.ColorModel(); cm != MnistModel {
		t.Fatalf("Unexpected color model on image")
	}
}

func TestImage_Bounds(t *testing.T) {
	if b := testImage.Bounds(); !b.Eq(image.Rect(0, 0, 2, 2)) {
		t.Fatalf("Unexpected bounds: %s", b)
	}
}

func TestImage_At(t *testing.T) {
	if r, g, b, a := testImage.At(0, 0).RGBA(); !reflect.DeepEqual([]uint32{r, g, b, a}, white) {
		t.Fatalf("Expected white pixel, got %#v", []uint32{r, g, b, a})
	}
	if r, g, b, a := testImage.At(1, 1).RGBA(); !reflect.DeepEqual([]uint32{r, g, b, a}, white) {
		t.Fatalf("Expected white pixel, got %#v", []uint32{r, g, b, a})
	}
	if r, g, b, a := testImage.At(1, 0).RGBA(); !reflect.DeepEqual([]uint32{r, g, b, a}, black) {
		t.Fatalf("Expected black pixel, got %#v", []uint32{r, g, b, a})
	}
	if r, g, b, a := testImage.At(0, 1).RGBA(); !reflect.DeepEqual([]uint32{r, g, b, a}, black) {
		t.Fatalf("Expected black pixel, got %#v", []uint32{r, g, b, a})
	}
}

func TestImageReader_ReadImage(t *testing.T) {
	ir := &ImageReader{
		ReadSeeker: bytes.NewReader(testImageFile),
	}
	if err := ir.ValidateHeader(); err != nil {
		t.Fatalf("Could not validate header: %s", err)
	}
	if ir.numImages != 2 {
		t.Fatalf("Unexpected number of images: %d", ir.numImages)
	}
	if ir.width != 3 {
		t.Fatalf("Unexpected width: %d", ir.width)
	}
	if ir.height != 2 {
		t.Fatalf("Unexpected height: %d", ir.width)
	}

	if n := ir.NumImages(); n != 2 {
		t.Fatalf("Unexpected number of images: %d", n)
	}

	img, err := ir.ReadImage(0)
	if err != nil {
		t.Fatalf("Could not read image #0: %s", err)
	}
	if r, g, b, a := img.At(0, 0).RGBA(); !reflect.DeepEqual([]uint32{r, g, b, a}, white) {
		t.Fatalf("Expected white pixel, got %#v", []uint32{r, g, b, a})
	}
	if r, g, b, a := img.At(2, 1).RGBA(); !reflect.DeepEqual([]uint32{r, g, b, a}, white) {
		t.Fatalf("Expected white pixel, got %#v", []uint32{r, g, b, a})
	}

	img, err = ir.ReadImage(1)
	if err != nil {
		t.Fatalf("Could not read image #1: %s", err)
	}
	if r, g, b, a := img.At(0, 0).RGBA(); !reflect.DeepEqual([]uint32{r, g, b, a}, black) {
		t.Fatalf("Expected black pixel, got %#v", []uint32{r, g, b, a})
	}
	if r, g, b, a := img.At(2, 1).RGBA(); !reflect.DeepEqual([]uint32{r, g, b, a}, black) {
		t.Fatalf("Expected black pixel, got %#v", []uint32{r, g, b, a})
	}
}

func TestImageReader_NumImages_Unvalidated(t *testing.T) {
	defer func() {
		if x := recover(); x == nil {
			t.Fatalf("Unvalidated read did not panic")
		}
	}()
	ir := &ImageReader{}
	ir.NumImages()
}

func TestImageReader_ReadImage_Unvalidated(t *testing.T) {
	defer func() {
		if x := recover(); x == nil {
			t.Fatalf("Unvalidated read did not panic")
		}
	}()
	ir := &ImageReader{}
	ir.ReadImage(0)
}

func TestNewImageSet(t *testing.T) {
	is, err := NewImageSet(bytes.NewReader(testImageFile))
	if err != nil {
		t.Fatalf("Could not parse data: %s", err)
	}
	if n := len(is); n != 2 {
		t.Fatalf("Unexpected number of labels: %d", n)
	}
	if b := is[0].Bounds(); b.Dx() != 3 || b.Dy() != 2 {
		t.Fatalf("Unexpected dimensions: %#v", b)
	}
}
