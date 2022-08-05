package main

import "os"

type ByteDecoder struct {
	bufferIndex int
	buffer      []byte
}

func createNewByteDecoder(buffer []byte) ByteDecoder {
	return ByteDecoder{bufferIndex: 0, buffer: buffer}
}

func (bd *ByteDecoder) Decode(file *os.File) {
	encodingMap := bd.getEncodingMap()
	size := int(bd.readInt32())
	for i := 0; i < size; i++ {
		encoding := bd.ReadEncoding()
		token := string(encodingMap[encoding])
		file.WriteString(token)
	}
}

func (bd *ByteDecoder) getEncodingMap() map[string]rune {
	len := int(bd.readInt32())
	encodingMap := make(map[string]rune)
	for i := 0; i < len; i++ {
		token := bd.readInt32()
		encoding := bd.ReadEncoding()
		encodingMap[encoding] = token
	}
	return encodingMap
}

func (bd *ByteDecoder) readInt32() int32 {
	var num int32 = 0
	for i := 3; i >= 0; i-- {
		currByte := bd.buffer[bd.bufferIndex]
		num |= (int32(currByte) << (i * 8))
		bd.bufferIndex++
	}
	return num
}

func (bd *ByteDecoder) ReadEncoding() string {
	len := int(bd.readInt32())
	encoding := ""
	for i := 0; i < len; i++ {
		currIndex := bd.bufferIndex + (i / 8)
		if (bd.buffer[currIndex] & (1 << (7 - (i % 8)))) > 0 {
			encoding += "1"
		} else {
			encoding += "0"
		}
	}
	bd.bufferIndex += (len / 8)
	if (len % 8) > 0 {
		bd.bufferIndex++
	}
	return encoding
}
