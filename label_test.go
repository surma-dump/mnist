package mnist

import (
	"bytes"
	"encoding/binary"
	"testing"
)

var testLabelFile []byte

func init() {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, LABEL_MAGIC_NUMBER)
	binary.Write(buf, binary.BigEndian, int32(4))
	binary.Write(buf, binary.BigEndian, []uint8{0, 1, 2, 3})
	testLabelFile = buf.Bytes()
}

func TestLabelReader_ReadLabel(t *testing.T) {
	lr := &LabelReader{
		ReadSeeker: bytes.NewReader(testLabelFile),
	}
	if err := lr.ValidateHeader(); err != nil {
		t.Fatalf("Could not parse data: %s", err)
	}
	if n := lr.NumLabels(); n != 4 {
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

func TestLabelReader_NumLabels_Unvalidated(t *testing.T) {
	defer func() {
		if x := recover(); x == nil {
			t.Fatalf("Unvalidated read did not panic")
		}
	}()
	lr := &LabelReader{}
	lr.NumLabels()
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

func TestNewLabelSet(t *testing.T) {
	ls, err := NewLabelSet(bytes.NewReader(testLabelFile))
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
