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
	LZH                      // LZH compression
	LZW                      // LZW compression
	LZ4                      // LZ4 compression
	RAR                      // RAR 5.0 and later compression
	RAROld                   // Rar pre 1.5 compression
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
	case LZH:
		return "lzh"
	case LZW:
		return "lzw"
	case LZ4:
		return "lz4"
	case RAR:
		return "rar post 5.0"
	case RAROld:
		return "rar pre 1.5"
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
	case LZH:
		return ".lzh"
	case LZW:
		return ".Z"
	case LZ4:
		return "lz4"
	case RAR, RAROld:
		return ".rar"
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
	case "lzh":
		return LZH
	case "lzw", "Z":
		return LZW
	case "lz4":
		return LZ4
	case "rar":
		return RAR
	}
	return Unknown
}

// Magic numbers for supported formats
var (
	headerGzip       = []byte{0x1f, 0x8b}
	headerTar1       = []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00, 0x30, 0x30} // offset: 257
	headerTar2       = []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00, 0x20, 0x00} // offset: 257
	headerZip        = []byte{0x50, 0x4b, 0x03, 0x04}
	headerZipEmpty   = []byte{0x50, 0x4b, 0x05, 0x06}
	headerZipSpanned = []byte{0x50, 0x4b, 0x07, 0x08}
	headerBzip2      = []byte{0x42, 0x5a, 0x68}
	headerLZH        = []byte{0x1F, 0xa0}
	headerLZW        = []byte{0x1F, 0x9d}
	headerLZ4        = []byte{0x18, 0x4d, 0x22, 0x04}
	headerRAR        = []byte{0x52, 0x61, 0x72, 0x21, 0x1a, 0x07, 0x01, 0x00}
	headerRAROld     = []byte{0x52, 0x61, 0x72, 0x21, 0x1a, 0x07, 0x00}
)

// GetFormat tries to match up the data in the Reader to a supported
// magic number, if a match isn't found, UnsupportedFmt is returned
func GetFormat(r io.ReaderAt) (Format, error) {
	ok, err := IsLZ4(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return LZ4, nil
	}
	return Unknown, errors.New("unsupported format: input format is not known")
}

func IsLZ4(r io.ReaderAt) (bool, error) {
	h := make([]byte, 4)
	// Reat the first 4 bytes
	n, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	if n < 4 {
		return false, fmt.Errorf("magic number error: short read, expected to read 4 bytes, read %d bytes", n)
	}
	var h32 uint32
	// check for lz4
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.LittleEndian, &h32)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched LZ4's magic number: %s", err)
	}
	var m32 uint32
	mbuf := bytes.NewBuffer(headerLZ4)
	err = binary.Read(mbuf, binary.BigEndian, &m32)
	if err != nil {
		return false, fmt.Errorf("error while converting LZ4 magic number for comparison: %s", err)
	}
	if h32 == m32 {
		return true, nil
	}
	return false, nil
}

func IsLZW(r io.ReaderAt) (bool, error) {
	h := make([]byte, 2)
	// Reat the first 8 bytes since that's where most magic numbers are
	n, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	if n < 2 {
		return false, fmt.Errorf("magic number error: short read, expected to read 2 bytes, read %d bytes", n)
	}
	var h16 uint16
	// check for lzw
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.LittleEndian, &h16)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched LZW's magic number: %s", err)
	}
	var m16 uint16
	mbuf := bytes.NewBuffer(headerLZW)
	err = binary.Read(mbuf, binary.BigEndian, &m16)
	fmt.Printf("%x %x %x\n", h16, m16, headerLZW)
	if err != nil {
		return false, fmt.Errorf("error while converting LZW magic number for comparison: %s", err)
	}
	if h16 == m16 {
		return true, nil
	}
	return false, nil
}
