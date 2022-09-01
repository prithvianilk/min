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
	return ByteDecoder{
		bufferSize:    0,
		contents:      contents,
		contentsIndex: 0,
		encodingMap:   make(map[string]byte),
	}
}

func (bd *ByteDecoder) decodeAndWrite(file *os.File) {
	bd.getEncodingMap()
	numberOfEncodings := int(bd.readInt32())
	for i := 0; i < numberOfEncodings; i++ {
		encoding := bd.readEncoding()
		token := string(bd.encodingMap[encoding])
		_, err := file.WriteString(token)
		check(err)
	}
}

func (bd *ByteDecoder) getEncodingMap() {
	encodingMapSize := int(bd.readInt32())
	for i := 0; i < encodingMapSize; i++ {
		token := bd.readToken()
		encoding := bd.readEncodingWithLength()
		bd.encodingMap[encoding] = token
	}
}

func (bd *ByteDecoder) readInt32() int32 {
	var num int32 = 0
	for i := 31; i >= 0; i-- {
		if bd.isCurrentBitSet() {
			num |= (1 << i)
		}
		bd.incrementAndResetIfFull()
	}
	return num
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

func (bd *ByteDecoder) readToken() byte {
	var token byte = 0
	for i := 7; i >= 0; i-- {
		if bd.isCurrentBitSet() {
			token |= (1 << i)
		}
		bd.incrementAndResetIfFull()
	}
	return token
}

func (bd *ByteDecoder) readEncodingWithLength() string {
	len := int(bd.readInt32())
	encoding := ""
	for i := 0; i < len; i++ {
		encoding = bd.updateEncoding(encoding)
	}
	return encoding
}

func (bd *ByteDecoder) updateEncoding(encoding string) string {
	if bd.isCurrentBitSet() {
		encoding += "1"
	} else {
		encoding += "0"
	}
	bd.incrementAndResetIfFull()
	return encoding
}

func (bd *ByteDecoder) readEncoding() string {
	encoding := ""
	for !bd.isValidEncoding(encoding) {
		encoding = bd.updateEncoding(encoding)
	}
	return encoding
}

func (bd *ByteDecoder) isValidEncoding(encoding string) bool {
	_, isInMap := bd.encodingMap[encoding]
	return isInMap
}
