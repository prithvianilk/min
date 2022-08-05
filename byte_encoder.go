package main

import (
	"os"
)

type ByteEncoder struct {
	bufferSize int
	buffer     byte
	file       *os.File
}

func createNewByteEncoder(zipfpath string) ByteEncoder {
	file, err := os.Create(zipfpath)
	check(err)
	return ByteEncoder{file: file}
}

func (be *ByteEncoder) WriteInt32(num int32) {
	intByteSlice := make([]byte, 4)
	for i := 3; i >= 0; i-- {
		index := 3 - i
		intByteSlice[index] = byte(num >> (8 * i))
	}
	_, err := be.file.Write(intByteSlice)
	check(err)
}

func (be *ByteEncoder) WriteEncoding(tokenEncoding TokenEncodingPair) {
	be.WriteInt32(int32(tokenEncoding.token))
	be.WriteInt32(int32(len(tokenEncoding.encoding)))
	for _, bitRune := range tokenEncoding.encoding {
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
