package mnist

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

var (
	ErrInvalidMagicNumber = fmt.Errorf("Invalid magic number")
)

func validateMagicNumber(r io.Reader, magic int32) error {
	var readMagic int32
	if err := binary.Read(r, binary.BigEndian, &readMagic); err != nil {
		return err
	}
	if readMagic != magic {
		return ErrInvalidMagicNumber
	}
	return nil
}

func seekify(r io.Reader) (io.ReadSeeker, error) {
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
	return rs, nil
}
