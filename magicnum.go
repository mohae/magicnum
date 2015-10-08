package magicnum

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	Unknown    Format = iota // unknown format
	Gzip                     // Gzip compression format
	Tar                      // Tar format; normally used
	Tar1                     // Tar1 header format; normalizes to Tar
	Tar2                     // Tar1 header format; normalizes to Tar
	Zip                      // Zip archive
	ZipEmpty                 // Empty Zip Archive
	ZipSpanned               // Spanned Zip Archive
	Bzip2                    // Bzip2 compression
	//LZW                      // LZW compression
	LZ4 // LZ4 compression
)

// Magic numbers for magicnum for compression and archive formats
var (
	headerGzip       = []byte{0x1f, 0x8b}
	headerTar1       = []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00, 0x30, 0x30} // offset: 257
	headerTar2       = []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00, 0x20, 0x00} // offset: 257
	headerZip        = []byte{0x50, 0x4b, 0x03, 0x04}
	headerZipEmpty   = []byte{0x50, 0x4b, 0x05, 0x06}
	headerZipSpanned = []byte{0x50, 0x4b, 0x07, 0x08}
	headerBzip2      = []byte{0x42, 0x5a, 0x68}
	//headerLZW        = []byte{0x1F, 0x9d}
	headerLZ4 = []byte{0x18, 0x4d, 0x22, 0x04}
)

// TODO: should Format be more specific? e.g. CompressionFormat, MediaFormat, etc.
type Format int

func (f Format) String() string {
	switch f {
	case Gzip:
		return "gzip"
	case Tar, Tar1, Tar2:
		return "tar"
	case Zip:
		return "zip"
	case ZipEmpty:
		return "empty zip archive"
	case ZipSpanned:
		return "spanned zip archive"
	case Bzip2:
		return "bzip2"
	//case LZW:
	//	return "lzw"
	case LZ4:
		return "lz4"
	}
	return "unknown"
}

// Ext returns the extension for the format. Formats may have more than one
// accepted extension; alternate extensiona are not supported.
func (f Format) Ext() string {
	switch f {
	case Gzip:
		return ".gz"
	case Tar, Tar1, Tar2:
		return ".tar"
	case Zip, ZipEmpty, ZipSpanned:
		return ".zip"
	case Bzip2:
		return ".bz2"
	//case LZW:
	//	return ".Z"
	case LZ4:
		return ".lz4"
	}
	return "unknown"
}

func FormatFromString(s string) Format {
	s = strings.ToLower(s)
	switch s {
	case "gzip", "gz":
		return Gzip
	case "tar":
		return Tar
	case "zip":
		return Zip
	case "bzip2", "bz2":
		return Bzip2
	//case "lzw", "Z":
	//	return LZW
	case "lz4":
		return LZ4
	}
	return Unknown
}

// ParseFormat takes a string and returns the format or unknown. Any compressed
// tar extensions are returned as the compression format and not tar.
//
// If the passed string starts with a '.', it is removed.
// All strings are lowercased
func ParseFormat(s string) Format {
	if s[0] == '.' {
		s = s[1:]
	}
	s = strings.ToLower(s)
	switch s {
	case "gzip", "tar.gz", "tgz":
		return Gzip
	case "tar":
		return Tar
	case "bz2", "tbz", "tb2", "tbz2", "tar.bz2":
		return Bzip2
	case "lz4", "tar.lz4", "tz4":
		return LZ4
	case "zip":
		return Zip
	}
	return Unknown
}

// GetFormat tries to match up the data in the Reader to a supported
// magic number, if a match isn't found, UnsupportedFmt is returned
//
// For zips, this will also match on files with empty zip or spanned zip magic
// numbers.  If you need to distinguich between the various zip formats, use
// something else.
func GetFormat(r io.ReaderAt) (Format, error) {
	ok, err := IsLZ4(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return LZ4, nil
	}
	ok, err = IsGzip(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return Gzip, nil
	}
	ok, err = IsZip(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return Zip, nil
	}
	ok, err = IsTar(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return Tar, nil
	}
	ok, err = IsBzip2(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return Bzip2, nil
	}
	//ok, err = IsLZW(r)
	//if err != nil {
	//	return Unknown, err
	//}
	//if ok {
	//	return LZW, nil
	//}
	return Unknown, errors.New("unsupported format: input format is not known")
}

// IsBzip2 checks to see if the received reader's contents are in bzip2 format
// by checking the magic numbers.
func IsBzip2(r io.ReaderAt) (bool, error) {
	h := make([]byte, 3)
	// Read the first 3 bytes
	_, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	var hb [3]byte
	// check for bzip2
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.LittleEndian, &hb)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched bzip2's magic number: %s", err)
	}
	var cb [3]byte
	cbuf := bytes.NewBuffer(headerBzip2)
	err = binary.Read(cbuf, binary.BigEndian, &cb)
	if err != nil {
		return false, fmt.Errorf("error while converting bzip2 magic number for comparison: %s", err)
	}
	if hb == cb {
		return true, nil
	}
	return false, nil
}

// IsGzip checks to see if the received reader's contents are in gzip format
// by checking the magic numbers.
func IsGzip(r io.ReaderAt) (bool, error) {
	h := make([]byte, 2)
	// Read the first 2 bytes
	_, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	var h16 uint16
	// check for gzip
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.BigEndian, &h16)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched bzip2's magic number: %s", err)
	}
	var c16 uint16
	cbuf := bytes.NewBuffer(headerGzip)
	err = binary.Read(cbuf, binary.BigEndian, &c16)
	if err != nil {
		return false, fmt.Errorf("error while converting bzip2 magic number for comparison: %s", err)
	}
	if h16 == c16 {
		return true, nil
	}
	return false, nil
}

// IsLZ4 checks to see if the received reader's contents are in LZ4 foramt by
// checking the magic numbers.
func IsLZ4(r io.ReaderAt) (bool, error) {
	h := make([]byte, 4)
	// Read the first 4 bytes
	_, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	var h32 uint32
	// check for lz4
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.LittleEndian, &h32)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched LZ4's magic number: %s", err)
	}
	var c32 uint32
	cbuf := bytes.NewBuffer(headerLZ4)
	err = binary.Read(cbuf, binary.BigEndian, &c32)
	if err != nil {
		return false, fmt.Errorf("error while converting LZ4 magic number for comparison: %s", err)
	}
	if h32 == c32 {
		return true, nil
	}
	return false, nil
}

// IsLZW checks to see if the received reader's contents are in LZ4 format by
// checking the magic numbers.
//
// TODO: unsupported until I have a better understanding of how to handle LZW
/*
func IsLZW(r io.ReaderAt) (bool, error) {
	h := make([]byte, 2)
	// Reat the first 8 bytes since that's where most magic numbers are
	_, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	var h16 uint16
	// check for lzw
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.LittleEndian, &h16)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched LZW's magic number: %s", err)
	}
	var c16 uint16
	cbuf := bytes.NewBuffer(headerLZW)
	err = binary.Read(cbuf, binary.BigEndian, &c16)
	if err != nil {
		return false, fmt.Errorf("error while converting LZW magic number for comparison: %s", err)
	}
	if h16 == c16 {
		return true, nil
	}
	return false, nil
}
*/

// IsTar checks to see if the received reader's contents are in the tar format
// by checking the magic numbers. This evaluates using both tar1 and tar2 magic
// numbers.
func IsTar(r io.ReaderAt) (bool, error) {
	h := make([]byte, 8)
	// Read the first 8 bytes at offset 257
	_, err := r.ReadAt(h, 257)
	if err != nil {
		return false, err
	}
	var h64 uint64
	// check for Zip
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.BigEndian, &h64)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched tar's magic number: %s", err)
	}
	var c64 uint64
	cbuf := bytes.NewBuffer(headerTar1)
	err = binary.Read(cbuf, binary.BigEndian, &c64)
	if err != nil {
		return false, fmt.Errorf("error while converting the tar magic number for comparison: %s", err)
	}
	if h64 == c64 {
		return true, nil
	}
	cbuf = bytes.NewBuffer(headerTar2)
	err = binary.Read(cbuf, binary.BigEndian, &c64)
	if err != nil {
		return false, fmt.Errorf("error while converting the empty tar magic number for comparison: %s", err)
	}
	if h64 == c64 {
		return true, nil
	}
	return false, nil
}

// IsZip checks to see if the received reader's contents are in the zip format
// by checking the magic numbers. This will match on zip, empty zip and spanned
// zip magic numbers. If you need to distinguish between those, use something
// else.
func IsZip(r io.ReaderAt) (bool, error) {
	h := make([]byte, 4)
	// Read the first 4 bytes
	_, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	var h32 uint32
	// check for Zip
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.BigEndian, &h32)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched zip's magic number: %s", err)
	}
	var c32 uint32
	cbuf := bytes.NewBuffer(headerZip)
	err = binary.Read(cbuf, binary.BigEndian, &c32)
	if err != nil {
		return false, fmt.Errorf("error while converting the zip magic number for comparison: %s", err)
	}
	if h32 == c32 {
		return true, nil
	}
	cbuf = bytes.NewBuffer(headerZipEmpty)
	err = binary.Read(cbuf, binary.BigEndian, &c32)
	if err != nil {
		return false, fmt.Errorf("error while converting the empty zip magic number for comparison: %s", err)
	}
	if h32 == c32 {
		return true, nil
	}
	cbuf = bytes.NewBuffer(headerZipSpanned)
	err = binary.Read(cbuf, binary.BigEndian, &c32)
	if err != nil {
		return false, fmt.Errorf("error while converting the spanned zip magic number for comparison: %s", err)
	}
	if h32 == c32 {
		return true, nil
	}
	return false, nil
}
