package tmx

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
	"strconv"
)

// TileInstance is about to get completely refactored.
type TileInstance struct {
	GID                 uint32
	FlippedHorizontally bool
	FlippedVertically   bool
	FlippedDiagonally   bool
}

// Constants for parsing GID data
const (
	FlippedHorizontallyFlag uint32 = 0x80000000
	FlippedVerticallyFlag   uint32 = 0x40000000
	FlippedDiagonallyFlag   uint32 = 0x20000000
	GIDMask                 uint32 = ^(FlippedHorizontallyFlag | FlippedVerticallyFlag | FlippedDiagonallyFlag)
)

// TileIterator provides a generic access method for TMX Tiles with different
// encoding methods.
type TileIterator interface {
	Next() bool
	Error() error
	Get() TileInstance
	GetIndex() uint32
}

type xmlIterator struct {
	d *Data
	i uint32
}

func (xi *xmlIterator) Next() bool {
	xi.i++
	return int(xi.i) >= len(xi.d.TileData)
}

func (xi *xmlIterator) Error() error {
	return nil
}

func (xi *xmlIterator) Get() TileInstance {
	return TileInstance{
		GID:                 GIDMask & xi.d.TileData[xi.i].GID,
		FlippedHorizontally: xi.d.TileData[xi.i].GID&FlippedHorizontallyFlag > 0,
		FlippedVertically:   xi.d.TileData[xi.i].GID&FlippedVerticallyFlag > 0,
		FlippedDiagonally:   xi.d.TileData[xi.i].GID&FlippedDiagonallyFlag > 0,
	}
}

func (xi *xmlIterator) GetIndex() uint32 {
	return xi.i - 1
}

type csvIterator struct {
	d     *Data
	start int
	end   int
	i     uint32
	err   error
	done  bool
}

func (ci *csvIterator) Next() bool {
	if ci.done {
		return false
	}
	ci.i++
	ci.start = ci.end
	i := bytes.IndexByte(ci.d.Data[ci.start:], ',')
	if i == -1 {
		ci.done = true
	}
	ci.end += i
	return true
}

func (ci *csvIterator) Error() error {
	return ci.err
}

func (ci *csvIterator) Get() TileInstance {
	g, err := strconv.ParseUint(string(ci.d.Data[ci.start:ci.end]), 10, 32)
	if err != nil {
		ci.err = err
		return TileInstance{}
	}
	gid := uint32(g)
	return TileInstance{
		GID:                 GIDMask & gid,
		FlippedHorizontally: gid&FlippedHorizontallyFlag > 0,
		FlippedVertically:   gid&FlippedVerticallyFlag > 0,
		FlippedDiagonally:   gid&FlippedDiagonallyFlag > 0,
	}
}

func (ci *csvIterator) GetIndex() uint32 {
	return ci.i - 1
}

type b64Iterator struct {
	r   io.Reader
	tok TileInstance
	err error
	i   uint32
}

func (bi *b64Iterator) Next() bool {
	var gid uint32
	bi.err = binary.Read(bi.r, binary.LittleEndian, &gid)
	if bi.err != nil {
		if bi.err == io.EOF {
			bi.err = nil
		}
		return false
	}
	bi.i++
	bi.tok = TileInstance{
		GID:                 GIDMask & gid,
		FlippedHorizontally: gid&FlippedHorizontallyFlag > 0,
		FlippedVertically:   gid&FlippedVerticallyFlag > 0,
		FlippedDiagonally:   gid&FlippedDiagonallyFlag > 0,
	}
	return true
}

func (bi *b64Iterator) Error() error {
	return bi.err
}

func (bi *b64Iterator) Get() TileInstance {
	return bi.tok
}

func (bi *b64Iterator) GetIndex() uint32 {
	return bi.i - 1
}

func (d *Data) Iter() (TileIterator, error) {
	switch {
	case d.Encoding == nil && d.Compression != nil:
		return nil, errors.New("compression without encoding is not possible")
	case d.Encoding == nil && d.Compression == nil:
		return &xmlIterator{d: d}, nil
	case *d.Encoding == "csv":
		return &csvIterator{}, nil
	case *d.Encoding == "base64":
		var r io.Reader
		var err error
		r = bytes.NewReader(bytes.TrimSpace(d.Data))
		r = base64.NewDecoder(base64.StdEncoding, r)
		switch {
		case d.Compression == nil, *d.Compression == "":
			// Do nothing
		case *d.Compression == "gzip":
			r, err = gzip.NewReader(r)
		case *d.Compression == "zlib":
			r, err = zlib.NewReader(r)
		default:
			err = errors.New("invalid encoding")
		}
		if err != nil {
			return nil, errors.Wrap(err, "could not load base64 tile data")
		}
		return &b64Iterator{r: r}, nil

	default:
		return nil, errors.Errorf("invalid encoding: %s", *d.Encoding)
	}
}

func (d *Data) Tiles() ([]TileInstance, error) {
	iter, err := d.Iter()
	if err != nil {
		return nil, errors.Wrap(err, "bad iterator")
	}

	var tis []TileInstance
	for iter.Next() {
		tis = append(tis, iter.Get())
	}
	return tis, errors.Wrap(err, "error reading iterator")
}
