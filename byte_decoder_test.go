package main

import (
	"testing"
)

func TestReadBitString000(t *testing.T) {
	panicIfBitStringNotEqual([]byte{0x00, 0x00, 0x00, 0x03, 0x00}, "000", t)
}

func TestReadBitString10(t *testing.T) {
	panicIfBitStringNotEqual([]byte{0x00, 0x00, 0x00, 0x02, 0x80}, "10", t)
}

func TestReadBitString1010(t *testing.T) {
	panicIfBitStringNotEqual([]byte{0x00, 0x00, 0x00, 0x04, 0xA0}, "1010", t)
}

func TestReadInt1(t *testing.T) {
	panicIfIntEqual([]byte{0x00, 0x00, 0x00, 0x01}, 1, t)
}

func TestReadInt4(t *testing.T) {
	panicIfIntEqual([]byte{0x00, 0x00, 0x00, 0x04}, 4, t)
}

func TestReadInt32(t *testing.T) {
	panicIfIntEqual([]byte{0x00, 0x00, 0x00, 0x20}, 32, t)
}

func panicIfIntEqual(buffer []byte, expectedNum int32, t *testing.T) {
	bd := createNewByteDecoder(buffer)
	num := bd.readInt32()
	if num != expectedNum {
		t.Error("Expected", expectedNum, "got", num)
	}
}

func panicIfBitStringNotEqual(buffer []byte, expectedBitString string, t *testing.T) {
	bd := createNewByteDecoder(buffer)
	bitString := bd.ReadEncoding()
	if bitString != expectedBitString {
		t.Error("Expected", expectedBitString, "got", bitString)
	}
}
