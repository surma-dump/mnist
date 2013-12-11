package mnist

import (
	"encoding/binary"
	"fmt"
	"io"
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
	if err := validateMagicNumber(r); err != nil {
		return nil, err
	}
	return nil, nil
}

func validateMagicNumber(r io.Reader) error {
	var magic int32
	err := binary.Read(r, binary.BigEndian, &magic)
	if err != nil {
		return err
	}
	if magic != MAGIC_NUMBER {
		return ErrInvalidMagicNumber
	}
	return nil
}
