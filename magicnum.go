package magicnum

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	Unknown    Format = iota // unknown format
	Gzip                     // Gzip compression format; always a tar
	Tar                      // Tar format; normally used
	Tar1                     // Tar1 header format; normalizes to FmtTar
	Tar2                     // Tar1 header format; normalizes to FmtTar
	Zip                      // Zip archive
	ZipEmpty                 // Empty Zip Archive
	ZipSpanned               // Spanned Zip Archive
	Bzip2                    // Bzip2 compression
	LZH                      // LZH compression
	LZW                      // LZW compression
	LZ4                      // LZ4 compression
	RAR                      // RAR 5.0 and later compression
	RAROld                   // Rar pre 1.5 compression
	LZ4                      // LZ4 compression
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
