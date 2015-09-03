package magicnum

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"os"
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

// this test uses a tarball compressed with bzip2 because compress/bzip2
// doesn't have a compressor.
func TestIsBzip2(t *testing.T) {
	f, err := os.Open("test_files/test.bz2")
	if err != nil {
		t.Errorf("open test.bz2: expected no error, got %s", err)
		return
	}
	defer f.Close()
	ok, err := IsBzip2(f)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !ok {
		t.Error("expected ok to be true for bzip2, got false")
	}
	format, err := GetFormat(f)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if format != Bzip2 {
		t.Errorf("expected format to be bzip2, got %s", format)
	}
}

func TestIsGzip(t *testing.T) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	n, err := w.Write(testVal)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if n != 452 {
		t.Errorf("Expected 452 bytes to be written; %d were", n)
	}
	w.Close()
	r := bytes.NewReader(buf.Bytes())
	ok, err := IsGzip(r)
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
	if format != Gzip {
		t.Errorf("Expected format to be gzip got %s", format)
	}
}

func TestIsLZ4(t *testing.T) {
	var buf bytes.Buffer
	lw := lz4.NewWriter(&buf)
	n, err := lw.Write(testVal)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if n != 452 {
		t.Errorf("Expected 452 bytes to be written; %d were", n)
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

// a file is used for the non-empty test because creating one using the test
// data in this func resulted in the zip empty header...probably an error on
// my part.
func TestIsZip(t *testing.T) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	err := w.Close()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	r := bytes.NewReader(buf.Bytes())
	ok, err := IsZip(r)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	f, err := os.Open("test_files/test.zip")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	defer f.Close()
	ok, err = IsZip(f)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	format, err := GetFormat(f)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if format != Zip {
		t.Errorf("Expected format to be gzip got %s", format)
	}
}

// TODO: commented out because LZW output doesn't have the magic number.
// Figure out why and resolve.
// Another example http://play.golang.org/p/zGLAj1ruoh
/*
func TestIsLZW(t *testing.T) {
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
