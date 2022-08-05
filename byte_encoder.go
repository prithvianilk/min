package main

import (
	"os"
)

type ByteEncoder struct {
	bufferSize int
	buffer     byte
	file       *os.File
}

func createNewByteEncoder(file *os.File) ByteEncoder {
	return ByteEncoder{file: file}
}

func (be *ByteEncoder) WriteInt32(num int32) {
	byteSlice := make([]byte, 4)
	for i := 3; i >= 0; i-- {
		index := 3 - i
		byteSlice[index] = byte(num >> (8 * i))
	}
	_, err := be.file.Write(byteSlice)
	check(err)
}

func (be *ByteEncoder) WriteTokenEncodingPair(token rune, encoding string) {
	be.WriteInt32(token)
	be.WriteEncoding(encoding)
}

func (be *ByteEncoder) WriteEncoding(encoding string) {
	be.WriteInt32(int32(len(encoding)))
	for _, bitRune := range encoding {
		if bitRune == '1' {
			be.buffer |= (1 << (7 - be.bufferSize))
		}
		be.bufferSize++
		if be.bufferSize == 8 {
			be.flush()
		}
	}
	be.flush()
}

func (be *ByteEncoder) flush() {
	if be.bufferSize == 0 {
		return
	}
	_, err := be.file.Write([]byte{be.buffer})
	check(err)
	be.buffer = 0
	be.bufferSize = 0
}
