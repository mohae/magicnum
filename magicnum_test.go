package magicnum

import (
	"bytes"
	//"compress/lzw"
	//"io"
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
//func TestIsBzip2(t *testing.T) {

//}

func TestIsLZ4(t *testing.T) {
	var buf bytes.Buffer
	lw := lz4.NewWriter(&buf)
	n, err := lw.Write(testVal)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if n != 452 {
		t.Errorf("Expected 452 bytes to be copied; %d were", n)
	}
	lw.Close()
	r := bytes.NewReader(buf.Bytes())
	ok, err := IsLZ4(r)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	format, err := GetFormat(r)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if format != LZ4 {
		t.Errorf("Expected format to be LZ4 got %s", format)
	}
}

// TODO: commented out because LZW output doesn't have the magic number.
// Figure out why and resolve.
// Another example http://play.golang.org/p/zGLAj1ruoh
/*
func TestIsLZW(t *testing.T) {
	//r := bytes.NewReader(testVal)
	var buf bytes.Buffer
	lw := lzw.NewWriter(&buf, lzw.LSB, 8)
	n, err := lw.Write(testVal)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if n != 452 {
		t.Errorf("Expected 452 bytes to be copied; %d were", n)
	}
	lw.Close()
	t.Errorf("%x", buf.Bytes())
	rr := bytes.NewReader(buf.Bytes())
	ok, err := IsLZW(rr)
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
	if format != LZW {
		t.Errorf("Expected format to be LZW got %s", format)
	}
}
*/
