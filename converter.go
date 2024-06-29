package main

import (
	"os"
	"slices"
)

type Encoder interface {
	Encode(r rune) []byte
	Decode(bytes []byte) []rune
	Bom() []byte
	IsLittleEndian() bool
}

type Converter struct {
	FromEncoding  string
	ToEncoding    string
	InputFilePath string
}

func NewEncoder(encoding string) Encoder {
	switch encoding {
	case "UTF-8":
		return &UTF8{}
	case "UTF-16":
		return &UTF16{littleEndian: false}
	case "UTF-16LE":
		return &UTF16{littleEndian: true}
	case "UTF-32", "UTF-32BE":
		return &UTF32{littleEndian: false}
	case "UTF-32LE":
		return &UTF32{littleEndian: true}
	}
	return nil
}

func (conv Converter) Convert() ([]byte, error) {
	fromEncoder := NewEncoder(conv.FromEncoding)
	toEncoder := NewEncoder(conv.ToEncoding)

	encodedBytes := make([]byte, 0)
	if len(toEncoder.Bom()) > 0 {
		if toEncoder.IsLittleEndian() {
			bom := toEncoder.Bom()
			slices.Reverse(bom)
			encodedBytes = append(encodedBytes, bom...)
		} else {
			encodedBytes = append(encodedBytes, toEncoder.Bom()...)
		}
	}
	bytes, err := os.ReadFile(conv.InputFilePath)
	if err != nil {
		return encodedBytes, err
	}
	for _, r := range fromEncoder.Decode(bytes[len(fromEncoder.Bom()):]) {
		encodedBytes = append(encodedBytes, toEncoder.Encode(r)...)
	}
	return encodedBytes, nil
}
