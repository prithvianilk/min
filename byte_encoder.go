package main

import (
	"os"
)

const MAX_ENCODER_BUFFER_SIZE int = 8

type ByteEncoder struct {
	bufferSize int
	buffer     byte
	file       *os.File
}

func createNewByteEncoder(file *os.File) ByteEncoder {
	return ByteEncoder{bufferSize: 0, buffer: 0, file: file}
}

func (be *ByteEncoder) WriteToken(token byte) {
	for i := 7; i >= 0; i-- {
		if (token & (1 << i)) > 0 {
			be.setCurrentBit()
		}
		be.incrementAndFlushIfFull()
	}
}

func (be *ByteEncoder) setCurrentBit() {
	be.buffer |= (1 << (7 - be.bufferSize))
}

func (be *ByteEncoder) incrementAndFlushIfFull() {
	be.bufferSize++
	if be.bufferSize == MAX_ENCODER_BUFFER_SIZE {
		be.Flush()
	}
}

func (be *ByteEncoder) WriteInt32(num int32) {
	for i := 31; i >= 0; i-- {
		if (num & (1 << i)) > 0 {
			be.setCurrentBit()
		}
		be.incrementAndFlushIfFull()
	}
}

func (be *ByteEncoder) WriteTokenEncodingPair(token byte, encoding string) {
	be.WriteToken(token)
	be.WriteInt32(int32(len(encoding)))
	be.WriteEncoding(encoding)
}

func (be *ByteEncoder) WriteEncoding(encoding string) {
	for _, bitRune := range encoding {
		if bitRune == '1' {
			be.setCurrentBit()
		}
		be.incrementAndFlushIfFull()
	}
}

func (be *ByteEncoder) Flush() {
	if be.bufferSize == 0 {
		return
	}
	_, err := be.file.Write([]byte{be.buffer})
	check(err)
	be.buffer = 0
	be.bufferSize = 0
}
