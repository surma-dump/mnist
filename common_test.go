package mnist

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestValidateMagicNumber_WrongMagicNumber(t *testing.T) {
	data := "NO_THE_LABEL_MAGIC_NUMBER"
	err := validateMagicNumber(bytes.NewReader([]byte(data)), LABEL_MAGIC_NUMBER)
	if err != ErrInvalidMagicNumber {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestValidateMagicNumber_InsufficientData(t *testing.T) {
	data := "X"
	err := validateMagicNumber(bytes.NewReader([]byte(data)), LABEL_MAGIC_NUMBER)
	if err == nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestValidateMagicNumber_CorrectMagicNumber(t *testing.T) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, LABEL_MAGIC_NUMBER)
	err := validateMagicNumber(buf, LABEL_MAGIC_NUMBER)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

type unseekableReader struct {
	r *bytes.Reader
}

func (ub *unseekableReader) Read(v []byte) (n int, err error) {
	return ub.r.Read(v)
}

var (
	seekBuf = []byte("01234567890123456789")
)

func TestSeekify_Seekable(t *testing.T) {
	sr, err := seekify(bytes.NewReader(seekBuf))
	if err != nil {
		t.Fatalf("Could not seekify reader: %s", err)
	}

	sr.Seek(7, 0)
	buf := []byte{0}
	if _, err := sr.Read(buf); err != nil {
		t.Fatalf("Could not read after seeking: %s", err)
	}
	if buf[0] != '7' {
		t.Fatalf("Unexpected content: %#v", buf)
	}
}

func TestSeekify_Unseekable(t *testing.T) {
	sr, err := seekify(&unseekableReader{bytes.NewReader(seekBuf)})
	if err != nil {
		t.Fatalf("Could not seekify reader: %s", err)
	}

	sr.Seek(7, 0)
	buf := []byte{0}
	if _, err := sr.Read(buf); err != nil {
		t.Fatalf("Could not read after seeking: %s", err)
	}
	if buf[0] != '7' {
		t.Fatalf("Unexpected content: %#v", buf)
	}
}
