package mnist

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

const (
	MAGIC_NUMBER = int32(0x00000801)
)

var (
	ErrInvalidMagicNumber = fmt.Errorf("Invalid magic number")
)

type Label uint8

type LabelSet []Label

func NewLabelSet(r io.Reader) (LabelSet, error) {
	var err error
	var rs io.ReadSeeker
	if nrs, ok := r.(io.ReadSeeker); ok {
		rs = nrs
	} else {
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		rs = bytes.NewReader(buf)
	}

	lr := &LabelReader{ReadSeeker: rs}
	if err := lr.ValidateHeader(); err != nil {
		return nil, err
	}
	ls := make([]Label, lr.Len())
	for i := range ls {
		ls[i], err = lr.ReadLabel(i)
		if err != nil {
			return nil, err
		}
	}
	return ls, nil
}

type LabelReader struct {
	headerValidated bool
	numLabels       uint32
	io.ReadSeeker
}

func (lr *LabelReader) ReadLabel(i int) (Label, error) {
	if !lr.headerValidated {
		panic("Need to call ValidateHeader() first")
	}
	if _, err := lr.Seek(int64(8+i), 0); err != nil {
		return Label(0), err
	}
	buf := []byte{0}
	_, err := lr.Read(buf)
	return Label(uint8(buf[0])), err
}

func (lr *LabelReader) Len() int {
	if !lr.headerValidated {
		panic("Need to call ValidateHeader() first")
	}
	return int(lr.numLabels)
}

func (lr *LabelReader) ValidateHeader() error {
	if lr.headerValidated {
		return nil
	}
	if _, err := lr.Seek(0, 0); err != nil {
		return err
	}
	if err := validateMagicNumber(lr); err != nil {
		return err
	}
	if err := binary.Read(lr, binary.BigEndian, &lr.numLabels); err != nil {
		return err
	}
	lr.headerValidated = true
	return nil
}

func validateMagicNumber(r io.Reader) error {
	var magic int32
	if err := binary.Read(r, binary.BigEndian, &magic); err != nil {
		return err
	}
	if magic != MAGIC_NUMBER {
		return ErrInvalidMagicNumber
	}
	return nil
}
