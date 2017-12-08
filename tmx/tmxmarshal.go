package tmx

import (
	"io"
	"errors"
)

const (
	wireTypeVarint      = 0
	wireType64          = 1
	wireTypeLengthDelim = 2
	wireType32          = 5
)

const (
	mapFieldVersion      = iota
	mapFieldTiledVersion
	mapFieldOrientation
	mapFieldRenderOrder
	mapFieldWidth
	mapFieldHeight
	mapFieldTileWidth
	mapFieldTileHeight
	mapFieldHexSideLength
	mapFieldStaggerAxis
	mapFieldStaggerIndex
	mapFieldBackgroundColor
	mapFieldNextObjectID
	mapFieldProperties
	mapFieldTileSets
	mapFieldLayers
)

func writeUVarint(bw io.ByteWriter, num uint64, n *int) (err error) {
	for num >= 0x80 {
		if err = bw.WriteByte(byte(num) | 0x80); err != nil {
			return
		}
		num >>= 7
		*n++
	}
	if err = bw.WriteByte(byte(num)); err != nil {
		return
	}
	*n++
	return
}

func writeString(bw io.ByteWriter, s string, n *int) (err error) {
	if err = writeUVarint(bw, uint64(len([]byte(s))), n); err != nil {
		return
	}
	for _, c := range []byte(s) {
		if err = bw.WriteByte(c); err != nil {
			return
		}
		*n++
	}
	return
}

func writeUVarintField(bw io.ByteWriter, fieldNum byte, num uint64, n *int) (err error) {
	if err = bw.WriteByte(fieldNum<<3 | wireTypeLengthDelim); err != nil {
		return
	}
	*n++
	err = writeUVarint(bw, num, n)
	return
}

func writeStringField(bw io.ByteWriter, fieldNum byte, s string, n *int) (err error) {
	if err = bw.WriteByte(fieldNum<<3 | wireTypeLengthDelim); err != nil {
		return
	}
	*n++
	err = writeString(bw, s, n)
	return
}


func (m *Map) MarshalGen(bw io.ByteWriter) (n int, err error) {
	// Version
	if m.Version != "" {
		if err = writeStringField(bw, mapFieldVersion, m.Version, &n); err != nil {
			return
		}
	}
	if m.TiledVersion != "" {
		if err = writeStringField(bw, mapFieldTiledVersion, m.TiledVersion, &n); err != nil {
			return
		}
	}
	if m.Orientation != "" {
		if err = writeStringField(bw, mapFieldOrientation, m.Orientation, &n); err != nil {
			return
		}
	}
	if m.RenderOrder != nil {
		if err = writeStringField(bw, mapFieldRenderOrder, *m.RenderOrder, &n); err != nil {
			return
		}
	}
	if m.Width != 0 {
		if err = writeUVarintField(bw, mapFieldWidth, uint64(m.Width), &n); err != nil {
			return
		}
	}






	return
}


func readUVarint64(br io.ByteReader) (num uint64, n int, err error) {
	var b byte
	var s uint
	for {
		b, err = br.ReadByte()
		if err != nil {
			return
		}
		n++
		if b < 0x80 {
			num |= uint64(b) << s
			return
		}
		num |= uint64(b&0x7F) << s
		s += 7
		if s > 64 {
			err = errors.New("invalid encoding, varint overflow")
			return
		}
	}
}

func (m *Map) UnmarshalGen(br io.ByteReader) (n int, err error) {

	// Version

	return
}
