package main

type ByteDecoder struct {
	bufferIndex int
	buffer      []byte
}

func createNewByteDecoder(buffer []byte) ByteDecoder {
	return ByteDecoder{bufferIndex: 0, buffer: buffer}
}

func (bd *ByteDecoder) Decode() {
	// encodingMap := bd.getEncodingMap()
	// fmt.Println(encodingMap)
}

func (bd *ByteDecoder) getEncodingMap() map[rune]string {
	len := int(bd.readInt32())
	encodingMap := make(map[rune]string)
	for i := 0; i < len; i++ {
		token := bd.readInt32()
		encoding := bd.readBitString()
		encodingMap[token] = encoding
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

func (bd *ByteDecoder) readBitString() string {
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
