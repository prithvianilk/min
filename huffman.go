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
	encodings := make([]TokenEncodingPair, 0)
	populateEncodings(&root, "", &encodings)
	huffman.writeEncodings(zipfpath, encodings)
}

func (huffman *Huffman) writeEncodings(zipfpath string, encodings []TokenEncodingPair) {
	byteEncoder := createNewByteEncoder(zipfpath)
	byteEncoder.WriteInt32(int32(len(encodings)))
	for _, encoding := range encodings {
		byteEncoder.WriteEncoding(encoding)
	}
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
	file, _ := os.Open(fpath)
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

func (h *Huffman) Decompress(zipfpath, fpath string) {
	buffer, err := os.ReadFile(zipfpath)
	byteDecoder := createNewByteDecoder(buffer)
	byteDecoder.Decode()
	check(err)
}
