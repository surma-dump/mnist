package mnist

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestValidateMagicNumber_WrongMagicNumber(t *testing.T) {
	data := "NO_THE_MAGIC_NUMBER"
	err := validateMagicNumber(bytes.NewReader([]byte(data)))
	if err != ErrInvalidMagicNumber {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestValidateMagicNumber_InsufficientData(t *testing.T) {
	data := "X"
	err := validateMagicNumber(bytes.NewReader([]byte(data)))
	if err == nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestValidateMagicNumber_CorrectMagicNumber(t *testing.T) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, MAGIC_NUMBER)
	err := validateMagicNumber(buf)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestLabelReader_ReadLabel(t *testing.T) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, MAGIC_NUMBER)
	binary.Write(buf, binary.BigEndian, int32(4))
	binary.Write(buf, binary.BigEndian, []int8{0, 1, 2, 3})

	lr := &LabelReader{
		ReadSeeker: bytes.NewReader(buf.Bytes()),
	}
	if err := lr.ValidateHeader(); err != nil {
		t.Fatalf("Could not parse data: %s", err)
	}
	if n := lr.Len(); n != 4 {
		t.Fatalf("Unexpected number of labels: %d", n)
	}
	for i := 0; i < 4; i++ {
		if l, err := lr.ReadLabel(i); err != nil {
			t.Fatalf("Could not read label #%d: %s", i, err)
		} else if int(l) != i {
			t.Fatalf("Read unexpected label #%d: %d", i, l)
		}
	}
}

func TestLabelReader_Len_Unvalidated(t *testing.T) {
	defer func() {
		if x := recover(); x == nil {
			t.Fatalf("Unvalidated read did not panic")
		}
	}()
	lr := &LabelReader{}
	lr.Len()
}

func TestLabelReader_ReadLabel_Unvalidated(t *testing.T) {
	defer func() {
		if x := recover(); x == nil {
			t.Fatalf("Unvalidated read did not panic")
		}
	}()
	lr := &LabelReader{}
	lr.ReadLabel(0)
}

func TestNewLabelReader_Seekable(t *testing.T) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, MAGIC_NUMBER)
	binary.Write(buf, binary.BigEndian, int32(4))
	binary.Write(buf, binary.BigEndian, []int8{0, 1, 2, 3})

	ls, err := NewLabelSet(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("Could not parse data: %s", err)
	}
	if n := len(ls); n != 4 {
		t.Fatalf("Unexpected number of labels: %d", n)
	}
	for i := 0; i < 4; i++ {
		if int(ls[i]) != i {
			t.Fatalf("Read unexpected label #%d: %d", i, ls[i])
		}
	}
}

type unseekableReader struct {
	r *bytes.Reader
}

func (ub *unseekableReader) Read(v []byte) (n int, err error) {
	return ub.r.Read(v)
}

func TestNewLabelReader_Unseekable(t *testing.T) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, MAGIC_NUMBER)
	binary.Write(buf, binary.BigEndian, int32(4))
	binary.Write(buf, binary.BigEndian, []int8{0, 1, 2, 3})

	ls, err := NewLabelSet(&unseekableReader{bytes.NewReader(buf.Bytes())})
	if err != nil {
		t.Fatalf("Could not parse data: %s", err)
	}
	if n := len(ls); n != 4 {
		t.Fatalf("Unexpected number of labels: %d", n)
	}
	for i := 0; i < 4; i++ {
		if int(ls[i]) != i {
			t.Fatalf("Read unexpected label #%d: %d", i, ls[i])
		}
	}
}
