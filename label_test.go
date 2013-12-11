package mnist

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestNewLabelSet_WrongMagicNumber(t *testing.T) {
	data := "NO_THE_MAGIC_NUMBER"
	_, err := NewLabelSet(bytes.NewReader([]byte(data)))
	if err != ErrInvalidMagicNumber {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestNewLabelSet_InsufficientData(t *testing.T) {
	data := "X"
	_, err := NewLabelSet(bytes.NewReader([]byte(data)))
	if err == nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestNewLabelSet_CorrectMagicNumber(t *testing.T) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.BigEndian, MAGIC_NUMBER)
	if err != nil {
		t.Fatalf("Could not write to buffer: %s", err)
	}
	_, err = NewLabelSet(buf)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}
