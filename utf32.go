package main

import "encoding/binary"

type UTF32 struct {
	littleEndian bool
}

func (enc *UTF32) Decode(bytes []byte) []rune {
	codePoints := make([]rune, 0)
	for index := 0; index < len(bytes); index += 4 {
		var u uint32

		if enc.littleEndian {
			u = binary.LittleEndian.Uint32(bytes[index : index+4])
		} else {
			u = binary.BigEndian.Uint32(bytes[index : index+4])
		}
		codePoints = append(codePoints, rune(u))
	}
	return codePoints
}

func (enc *UTF32) Encode(r rune) []byte {
	if r >= 0x00 && r <= 0x10FFFF {
		hh := byte((r >> 24) & 0xFF)
		hl := byte((r >> 16) & 0xFF)
		lh := byte((r >> 8) & 0xFF)
		ll := byte(r & 0xFF)
		if enc.littleEndian {
			return []byte{ll, lh, hl, hh}
		}
		return []byte{hh, hl, lh, ll}
	}
	return []byte{}
}

func (enc *UTF32) Bom() []byte {
	return []byte{0x00, 0x00, 0xFE, 0xFF}
}

func (enc *UTF32) IsLittleEndian() bool {
	return enc.littleEndian
}
