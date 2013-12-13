package mnist

import (
	"encoding/binary"
	"image"
	"image/color"
	"io"
)

const (
	IMAGE_MAGIC_NUMBER = int32(0x00000803)
)

type Image struct {
	Width, Height int
	Data          []byte
}

func (i Image) ColorModel() color.Model {
	return MnistModel
}

func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.Width, i.Height)
}

func (i Image) At(x, y int) color.Color {
	return mnistColor(i.Data[y*i.Width+x])
}

type ImageSet []*Image

func NewImageSet(r io.Reader) (ImageSet, error) {
	rs, err := seekify(r)
	if err != nil {
		return nil, err
	}

	ir := &ImageReader{ReadSeeker: rs}
	if err := ir.ValidateHeader(); err != nil {
		return nil, err
	}
	is := make([]*Image, ir.NumImages())
	for i := range is {
		is[i], err = ir.ReadImage(i)
		if err != nil {
			return nil, err
		}
	}
	return is, nil
}

type ImageReader struct {
	headerValidated bool
	width, height   uint32
	numImages       uint32
	io.ReadSeeker
}

func (ir *ImageReader) NumImages() int {
	if !ir.headerValidated {
		panic("Need to call ValidateHeader() first")
	}
	return int(ir.numImages)
}

func (ir *ImageReader) ReadImage(i int) (*Image, error) {
	if !ir.headerValidated {
		panic("Need to call ValidateHeader() first")
	}
	if _, err := ir.Seek(4*4+int64(i)*int64(ir.width)*int64(ir.height), 0); err != nil {
		return nil, err
	}
	img := &Image{
		Width:  int(ir.width),
		Height: int(ir.height),
		Data:   make([]byte, ir.width*ir.height),
	}
	_, err := ir.Read(img.Data)
	return img, err
}

func (ir *ImageReader) ValidateHeader() error {
	if ir.headerValidated {
		return nil
	}
	if _, err := ir.Seek(0, 0); err != nil {
		return err
	}
	if err := validateMagicNumber(ir, IMAGE_MAGIC_NUMBER); err != nil {
		return err
	}
	if err := binary.Read(ir, binary.BigEndian, &ir.numImages); err != nil {
		return err
	}
	if err := binary.Read(ir, binary.BigEndian, &ir.height); err != nil {
		return err
	}
	if err := binary.Read(ir, binary.BigEndian, &ir.width); err != nil {
		return err
	}
	ir.headerValidated = true
	return nil
}
