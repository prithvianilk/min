package main

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var h Huffman
	h.Compress("min.go", "min.zip")
	h.Decompress("min.zip", "mim2.go")
}
