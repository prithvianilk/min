package main

import (
	"os"

	pq "github.com/emirpasic/gods/queues/priorityqueue"
)

const BUFFER_SIZE int = 4 * 1024

type Huffman struct{}

func (huffman *Huffman) Compress(fpath, zipfpath string) {
	freqMap := huffman.getFreqMap(fpath)
	root := huffman.getTree(freqMap)
	encodingMap := make(map[rune]string)
	populateEncodingMap(&root, "", encodingMap)
	zipFile, err := os.Create(zipfpath)
	check(err)
	huffman.writeEncodings(zipFile, encodingMap)
	huffman.writeEncodedFile(fpath, zipFile, encodingMap)
}

func (*Huffman) getTree(freqMap map[rune]int) Node {
	queue := pq.NewWith(byPriority)

	for token, freq := range freqMap {
		queue.Enqueue(createNewLeaf(token, freq))
	}

	for queue.Size() > 1 {
		minFreqNode := dequeueAndGetNode(queue)
		secondMinFreqNode := dequeueAndGetNode(queue)
		newFreq := minFreqNode.freq + secondMinFreqNode.freq
		node := createNewNode('$', newFreq, &minFreqNode, &secondMinFreqNode)
		queue.Enqueue(node)
	}

	root, _ := queue.Dequeue()
	return root.(Node)
}

func (*Huffman) getFreqMap(fpath string) map[rune]int {
	file, err := os.Open(fpath)
	check(err)
	defer file.Close()
	freqMap := make(map[rune]int)
	for {
		buffer := make([]byte, BUFFER_SIZE)
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		for index := 0; index < n; index++ {
			token := rune(buffer[index])
			freqMap[token] += 1
		}
	}
	return freqMap
}

func (huffman *Huffman) writeEncodings(zipFile *os.File, encodingMap map[rune]string) {
	byteEncoder := createNewByteEncoder(zipFile)
	byteEncoder.WriteInt32(int32(len(encodingMap)))
	for token, encoding := range encodingMap {
		byteEncoder.WriteTokenEncodingPair(token, encoding)
	}
}

func (huffman *Huffman) writeEncodedFile(fpath string, zipFile *os.File, encodingMap map[rune]string) {
	file, err := os.Open(fpath)
	check(err)
	defer file.Close()
	byteEncoder := createNewByteEncoder(zipFile)
	stat, err := os.Stat(fpath)
	check(err)
	byteEncoder.WriteInt32(int32(stat.Size()))
	for {
		buffer := make([]byte, BUFFER_SIZE)
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		for i := 0; i < n; i++ {
			byteEncoder.WriteEncoding(encodingMap[rune(buffer[i])])
		}
	}
	byteEncoder.WriteInt32(int32(len(encodingMap)))
	for token, encoding := range encodingMap {
		byteEncoder.WriteTokenEncodingPair(token, encoding)
	}
}

func (h *Huffman) Decompress(zipfpath, fpath string) {
	buffer, err := os.ReadFile(zipfpath)
	check(err)
	byteDecoder := createNewByteDecoder(buffer)
	file, err := os.Create(fpath)
	check(err)
	byteDecoder.Decode(file)
}
