package main

import (
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	command := os.Args[1]
	h := Huffman{}
	if command == "zip" {
		filepath, zipPath := os.Args[2], os.Args[3]
		h.Compress(filepath, zipPath)
	} else if command == "unzip" {
		zipPath, filepath := os.Args[2], os.Args[3]
		h.Decompress(zipPath, filepath)
	} else {
		panic("Invalid command: " + command)
	}
}
