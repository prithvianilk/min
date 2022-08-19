package main

import (
	"os"
)

const MAX_DECODER_BUFFER_SIZE int = 8

type ByteDecoder struct {
	bufferSize    int
	contents      []byte
	contentsIndex int
	encodingMap   map[string]byte
}

func createNewByteDecoder(contents []byte) ByteDecoder {
	return ByteDecoder{bufferSize: 0, contents: contents, contentsIndex: 0, encodingMap: nil}
}

func (bd *ByteDecoder) Decode(file *os.File) {
	bd.encodingMap = bd.getEncodingMap()
	size := int(bd.ReadInt32())
	for i := 0; i < size; i++ {
		encoding := bd.ReadEncoding()
		token := string(bd.encodingMap[encoding])
		file.WriteString(token)
	}
}

func (bd *ByteDecoder) getEncodingMap() map[string]byte {
	len := int(bd.ReadInt32())
	encodingMap := make(map[string]byte)
	for i := 0; i < len; i++ {
		token := bd.ReadToken()
		encoding := bd.ReadEncodingWithLength()
		encodingMap[encoding] = token
	}
	return encodingMap
}

func (bd *ByteDecoder) ReadToken() byte {
	var token byte = 0
	for i := 7; i >= 0; i-- {
		if bd.isCurrentBitSet() {
			token |= (1 << i)
		}
		bd.incrementAndResetIfFull()
	}
	return token
}

func (bd *ByteDecoder) ReadInt32() int32 {
	var num int32 = 0
	for i := 31; i >= 0; i-- {
		if bd.isCurrentBitSet() {
			num |= (1 << i)
		}
		bd.incrementAndResetIfFull()
	}
	return num
}

func (bd *ByteDecoder) ReadEncodingWithLength() string {
	len := int(bd.ReadInt32())
	encoding := ""
	for i := 0; i < len; i++ {
		if bd.isCurrentBitSet() {
			encoding += "1"
		} else {
			encoding += "0"
		}
		bd.incrementAndResetIfFull()
	}
	return encoding
}

func (bd *ByteDecoder) ReadEncoding() string {
	encoding := ""
	for {
		if bd.isCurrentBitSet() {
			encoding += "1"
		} else {
			encoding += "0"
		}
		bd.incrementAndResetIfFull()
		_, isInMap := bd.encodingMap[encoding]
		if isInMap {
			break
		}
	}
	return encoding
}

func (bd *ByteDecoder) isCurrentBitSet() bool {
	mask := byte(1 << (7 - bd.bufferSize))
	return (bd.getBuffer() & mask) > 0
}

func (bd *ByteDecoder) getBuffer() byte {
	return bd.contents[bd.contentsIndex]
}

func (bd *ByteDecoder) incrementAndResetIfFull() {
	bd.bufferSize++
	if bd.bufferSize == MAX_DECODER_BUFFER_SIZE {
		bd.reset()
	}
}

func (bd *ByteDecoder) reset() {
	bd.bufferSize = 0
	bd.contentsIndex++
}
