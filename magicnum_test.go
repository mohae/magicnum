package magicnum

import (
	"bytes"
	"compress/lzw"
	"io"
	"testing"

	"github.com/pierrec/lz4"
)

var testVal = []byte(`
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
 incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
 nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
 Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore
 eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt
 in culpa qui officia deserunt mollit anim id est laborum.
`)

// The tests to check format involve creating a compressed version using the
// desired algorithm and then checking its header.
//
// All algorithm specific tests will also call GetFormat() to validate its
// behavior for that algorithm.
func TestIsLZ4(t *testing.T) {
	b := make([]byte, 0, 512)
	r := bytes.NewReader(testVal)
	w := bytes.NewBuffer(b)
	lw := lz4.NewWriter(w)
	n, err := io.Copy(lw, r)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if n != 452 {
		t.Errorf("Expected 452 bytes to be copied; %d were", n)
	}
	lw.Close()
	rr := bytes.NewReader(w.Bytes())
	ok, err := IsLZ4(rr)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	format, err := GetFormat(rr)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if format != LZ4 {
		t.Errorf("Expected format to be LZ4 got %s", format)
	}
}
