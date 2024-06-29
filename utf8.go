package main

type UTF8 struct{}

func (enc *UTF8) Decode(bytes []byte) []rune {
	codePoints := make([]rune, 0)
	var lower byte = 0x80
	var upper byte = 0xBF
	var codePoint rune = 0x0
	bytesSeen := 0
	bytesNeeded := 0
	for index := 0; index < len(bytes); index++ {
		if bytesNeeded == 0 {
			if bytes[index] > 0x00 && bytes[index] < 0x80 {
				codePoints = append(codePoints, rune(bytes[index]))
			} else if bytes[index] >= 0xC2 && bytes[index] <= 0xDF {
				codePoint = rune(bytes[index] & 0x1F)
				bytesNeeded = 1
			} else if bytes[index] >= 0xE0 && bytes[index] <= 0xEF {
				codePoint = rune(bytes[index] & 0xF)
				bytesNeeded = 2
				if bytes[index] == 0xE0 {
					lower = 0xA0
				}
				if bytes[index] == 0xED {
					upper = 0x9F
				}
			} else if bytes[index] >= 0xF0 && bytes[index] <= 0xF4 {
				codePoint = rune(bytes[index] & 0x7)
				bytesNeeded = 3
				if bytes[index] == 0xF0 {
					lower = 0x90
				}
				if bytes[index] == 0xF4 {
					upper = 0x8F
				}
			} else {
				// error case so replace error with replacement character
				codePoints = append(codePoints, 0xFFFD)
			}
			continue
		}

		if !(bytes[index] >= lower && bytes[index] <= upper) {
			// error case so replace error with replacement character
			codePoints = append(codePoints, 0xFFFD)
			// reset all other variables
			codePoint = 0x0
			bytesNeeded = 0
			bytesSeen = 0
			lower = 0x80
			upper = 0xBF
			continue
		}

		codePoint = codePoint<<6 | rune(bytes[index]&0x3f)
		lower = 0x80
		upper = 0xBF
		bytesSeen++

		if bytesSeen != bytesNeeded {
			continue
		}

		codePoints = append(codePoints, codePoint)

		// reset all other variables
		codePoint = 0x0
		bytesNeeded = 0
		bytesSeen = 0
	}
	//fmt.Printf("% x", codePoints)

	return codePoints
}

func (enc *UTF8) Encode(r rune) []byte {
	var offset byte = 0x0
	count := 0
	if r >= 0x0000 && r <= 0x007F {
		// ascii => return one byte
		return []byte{byte(r)}
	} else if r >= 0x0080 && r <= 0x07FF {
		offset = 0xC0
		count = 1
	} else if r >= 0x0800 && r <= 0xFFFF {
		offset = 0xE0
		count = 2
	} else if r >= 0x010000 && r <= 0x10FFFF {
		offset = 0xF0
		count = 3
	}

	bytes := make([]byte, 0, count)
	bytes = append(bytes, byte(r>>(6*count))+offset)

	for count > 0 {
		temp := byte(r >> (6 * (count - 1)))
		bytes = append(bytes, 0x80|(temp&0x3F))
		count--
	}

	return bytes
}

func (enc *UTF8) Bom() []byte {
	return []byte{}
}

func (enc *UTF8) IsLittleEndian() bool {
	return false
}
