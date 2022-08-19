package main

import (
	"os"

	pq "github.com/emirpasic/gods/queues/priorityqueue"
)

const BUFFER_SIZE int = 4 * 1024

type Huffman struct {
	byteEncoder ByteEncoder
}

func (huffman *Huffman) Compress(fpath, zipfpath string) {
	freqMap := huffman.getFreqMap(fpath)
	root := huffman.getTree(freqMap)
	encodingMap := make(map[byte]string)
	populateEncodingMap(&root, "", encodingMap)
	zipFile, err := os.Create(zipfpath)
	check(err)
	huffman.byteEncoder = createNewByteEncoder(zipFile)
	huffman.writeEncodings(encodingMap)
	huffman.writeEncodedFile(fpath, encodingMap)
}

func (*Huffman) getTree(freqMap map[byte]int) Node {
	queue := pq.NewWith(byPriority)

	for token, freq := range freqMap {
		queue.Enqueue(createNewLeaf(byte(token), freq))
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

func (*Huffman) getFreqMap(fpath string) map[byte]int {
	file, err := os.Open(fpath)
	check(err)
	defer file.Close()
	freqMap := make(map[byte]int)
	for {
		buffer := make([]byte, BUFFER_SIZE)
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		for index := 0; index < n; index++ {
			token := buffer[index]
			freqMap[token] += 1
		}
	}
	return freqMap
}

func (huffman *Huffman) writeEncodings(encodingMap map[byte]string) {
	huffman.byteEncoder.WriteInt32(int32(len(encodingMap)))
	for token, encoding := range encodingMap {
		huffman.byteEncoder.WriteTokenEncodingPair(token, encoding)
	}
}

func (huffman *Huffman) writeEncodedFile(fpath string, encodingMap map[byte]string) {
	file, err := os.Open(fpath)
	check(err)
	defer file.Close()
	stat, err := os.Stat(fpath)
	check(err)
	huffman.byteEncoder.WriteInt32(int32(stat.Size()))
	for {
		buffer := make([]byte, BUFFER_SIZE)
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		for i := 0; i < n; i++ {
			huffman.byteEncoder.WriteEncoding(encodingMap[buffer[i]])
		}
	}
	huffman.byteEncoder.Flush()
}

func (h *Huffman) Decompress(zipfpath, fpath string) {
	buffer, err := os.ReadFile(zipfpath)
	check(err)
	byteDecoder := createNewByteDecoder(buffer)
	file, err := os.Create(fpath)
	check(err)
	byteDecoder.Decode(file)
}
