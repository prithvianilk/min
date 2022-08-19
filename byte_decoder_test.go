package main

import (
	"testing"
)

func TestReadBitString000WithLength(t *testing.T) {
	panicIfBitStringNotEqual([]byte{0x00, 0x00, 0x00, 0x03, 0x00}, "000", t)
}

func TestReadBitString10WithLength(t *testing.T) {
	panicIfBitStringNotEqual([]byte{0x00, 0x00, 0x00, 0x02, 0x80}, "10", t)
}

func TestReadBitString1010WithLength(t *testing.T) {
	panicIfBitStringNotEqual([]byte{0x00, 0x00, 0x00, 0x04, 0xA0}, "1010", t)
}

func TestReadInt1(t *testing.T) {
	panicIfIntNotEqual([]byte{0x00, 0x00, 0x00, 0x01}, 1, t)
}

func TestReadInt4(t *testing.T) {
	panicIfIntNotEqual([]byte{0x00, 0x00, 0x00, 0x04}, 4, t)
}

func TestReadInt32(t *testing.T) {
	panicIfIntNotEqual([]byte{0x00, 0x00, 0x00, 0x20}, 32, t)
}

func panicIfBitStringNotEqual(buffer []byte, expectedBitString string, t *testing.T) {
	bd := createNewByteDecoder(buffer)
	bitString := bd.ReadEncodingWithLength()
	if bitString != expectedBitString {
		t.Error("Expected", expectedBitString, "got", bitString)
	}
}

func panicIfIntNotEqual(buffer []byte, expectedNum int32, t *testing.T) {
	bd := createNewByteDecoder(buffer)
	num := bd.ReadInt32()
	if num != expectedNum {
		t.Error("Expected", expectedNum, "got", num)
	}
}
