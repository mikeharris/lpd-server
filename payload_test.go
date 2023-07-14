package lpd

import (
	"bytes"
	"testing"
)

// Test Byte to KB conversion
func TestPayload_GetFileSizeInKB(t *testing.T) {
	// Test correct conversion
	exp1 := 2.048
	p := Payload{FileSizeInBytes: 2048}
	r1 := p.GetFileSizeInKB()

	if r1 != exp1 {
		t.Fatalf("Expected %f but got %f for input %d.\n", exp1, r1, p.FileSizeInBytes)
	}

	exp1 = 1.055
	p = Payload{FileSizeInBytes: 1055}
	r1 = p.GetFileSizeInKB()

	if r1 != exp1 {
		t.Fatalf("Expected %f but got %f for input %d.\n", exp1, r1, p.FileSizeInBytes)
	}

	// Test empty payload
	p = Payload{}
	r2 := p.GetFileSizeInKB()
	if r2 != 0.0 {
		t.Fatalf("Expected %f but got %f for input %d.\n", exp1, r1, p.FileSizeInBytes)
	}
}

// Test Payload Unmarshal
func TestPayload_Unmarshal(t *testing.T) {
	// Read typical file
	b := []byte{01, 01, 01, 01, 01, 01, 01, 01, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 00}
	p := Payload{}
	r := bytes.NewReader(b)
	err := p.unmarshal(r)
	if err != nil || p.FileSizeInBytes != 20 {
		t.Fatalf("Expected file size of %d but got %d.\n", len(b), p.FileSizeInBytes)
	}

	if !bytes.Equal(p.PrintFile, b) {
		t.Fatalf("File unmarshaled %b does not match file read %b.", p.PrintFile, b)
	}

	// Read file stream without an empty byte at the end
	p = Payload{}
	b = []byte{01, 01, 01, 01, 01, 01, 01, 01, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}
	r = bytes.NewReader(b)
	err = p.unmarshal(r)
	if err != nil {
		t.Fatal("Received an unexpected error when reading byte stream.", err)
	}

	if err != nil || p.FileSizeInBytes != 19 {
		t.Fatalf("Expected file size of %d but got %d.\n", len(b), p.FileSizeInBytes)
	}

	if !bytes.Equal(p.PrintFile, b) {
		t.Fatalf("File unmarshaled %b does not match file read %b.", p.PrintFile, b)
	}
}
