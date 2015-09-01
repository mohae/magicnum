# magicnum
provides functions for checking magic numbers; currently supports a limited number of file formats

## About
The functions within this package currently accepts and `io.ReaderAt` and trys to determine its format by comparing the data to the magic numbers of various formats. 

It was created to determine the compression or archive format. Support for other magic numbers will be added as needed. This may also mean support for other formats, e.g. images.

## Usage
Get:

    go get github.com/mohae/magicnum
Import:

    import github.com/mohae/magicnum
    
Get a reader. 
To check to see if any format is matched:

    format, err := magicnum.GetFormat(r)
    
If `magicnum` isn't able to determine the format the `format` will be: `Unsupported` and an `NotSupportedErr` will be returned. Otherwise a `magicnum.Format` will be returned. 

To check to see if it is of a specific format:

    ok, err := magicnum.IsLZ4(r)
    
If the reader has data in the specified format, `lz4` in this case, `ok == true`, otherwise `ok == false`. If an error occurs during processing, that will be returned.

### Format
The `magicnum.Format` fulfills the `Stringer` interface so the common abbreviation for the format can be obtained with `format.String()`.

Since formats can have different valid extensions, an extension can be matched to a `Format` by using `magicnum.GetFormatFromString(ext)`.

## Endianness
This library currently supports checking the magic number using _little endian_ only. If there is need for _big endian_ support, please file an issue or make a pull request. Ofc, the latter is preferred!

## Copyright
Copyright 2015 by Joel Scoble.
This is provided under the __MIT License__. Please check the included LICENSE file for more details.
