package mnist

import (
	"encoding/binary"
	"io"
)

const (
	LABEL_MAGIC_NUMBER = int32(0x00000801)
)

type Label uint8

type LabelSet []Label

func NewLabelSet(r io.Reader) (LabelSet, error) {
	rs, err := seekify(r)
	if err != nil {
		return nil, err
	}

	lr := &LabelReader{ReadSeeker: rs}
	if err := lr.ValidateHeader(); err != nil {
		return nil, err
	}
	ls := make([]Label, lr.NumLabels())
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
	if _, err := lr.Seek(int64(2*4+i), 0); err != nil {
		return Label(0), err
	}
	buf := []byte{0}
	_, err := lr.Read(buf)
	return Label(uint8(buf[0])), err
}

func (lr *LabelReader) NumLabels() int {
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
	if err := validateMagicNumber(lr, LABEL_MAGIC_NUMBER); err != nil {
		return err
	}
	if err := binary.Read(lr, binary.BigEndian, &lr.numLabels); err != nil {
		return err
	}
	lr.headerValidated = true
	return nil
}
