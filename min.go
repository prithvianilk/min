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
	if command == "zip" {
		filepath, zipPath := os.Args[2], os.Args[3]
		encoder := CreateNewHuffmanEncoder(filepath, zipPath)
		encoder.Compress()
	} else if command == "unzip" {
		zipPath, filepath := os.Args[2], os.Args[3]
		Decompress(zipPath, filepath)
	} else {
		panic("Invalid command: " + command)
	}
}
