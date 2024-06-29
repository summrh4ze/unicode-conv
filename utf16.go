package main

import "encoding/binary"

type UTF16 struct {
	littleEndian bool
}

func (enc *UTF16) Decode(bytes []byte) []rune {
	codePoints := make([]rune, 0)

	index := 0
	for index < len(bytes) {
		var u uint16

		if enc.littleEndian {
			u = binary.LittleEndian.Uint16(bytes[index : index+2])
		} else {
			u = binary.BigEndian.Uint16(bytes[index : index+2])
		}

		if u <= 0xD7FF || (u >= 0xE000 && u <= 0xFFFF) {
			codePoints = append(codePoints, rune(u))
			index += 2
		} else {
			if u >= 0xD800 && u <= 0xDBFF {
				// high surrogate
				var codePoint rune = 0x0
				codePoint = rune((u - 0xD800) * 0x400)
				// read next 2 bytes
				var uu uint16

				if enc.littleEndian {
					uu = binary.LittleEndian.Uint16(bytes[index+2 : index+4])
				} else {
					uu = binary.BigEndian.Uint16(bytes[index+2 : index+4])
				}
				if uu >= 0xDC00 && uu <= 0xDFFF {
					// low surrogate
					codePoint = codePoint + rune(uu-0xDC00) + 0x10000
					codePoints = append(codePoints, codePoint)
					index += 4
					continue
				} else {
					// error return replacement character
					codePoints = append(codePoints, 0xFFFD)
					index += 4
					continue
				}
			} else {
				// error return replacement character
				codePoints = append(codePoints, 0xFFFD)
				index += 2
				continue
			}
		}
	}
	return codePoints
}

func (enc *UTF16) Encode(r rune) []byte {
	if (r >= 0x0000 && r <= 0xD7FF) || (r >= 0xE000 && r <= 0xFFFF) {
		high := byte((r >> 8) & 0xFF)
		low := byte(r & 0xFF)
		if enc.littleEndian {
			return []byte{low, high}
		}
		return []byte{high, low}
	} else if r >= 0x010000 && r <= 0x10FFFF {
		s := r - 0x10000
		highSurrogate := (s >> 10) + 0xD800
		lowSurrogate := (s & 0x3FF) + 0xDC00
		highSH := byte((highSurrogate >> 8) & 0xFF)
		highSL := byte(highSurrogate & 0xFF)
		lowSH := byte((lowSurrogate >> 8) & 0xFF)
		lowSL := byte(lowSurrogate & 0xFF)
		if enc.littleEndian {
			return []byte{lowSL, lowSH, highSL, highSH}
		}
		return []byte{highSH, highSL, lowSH, lowSL}
	}
	return []byte{}
}

func (enc *UTF16) Bom() []byte {
	return []byte{0xFE, 0xFF}
}

func (enc *UTF16) IsLittleEndian() bool {
	return enc.littleEndian
}
