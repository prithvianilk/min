package main

import (
	"os"

	pq "github.com/emirpasic/gods/queues/priorityqueue"
)

const BUFFER_SIZE int = 4 * 1024

type HuffmanEncoder struct {
	fpath       string
	byteEncoder ByteEncoder
	encodingMap map[byte]string
}

func CreateNewHuffmanEncoder(fpath, zipPath string) HuffmanEncoder {
	zipFile, err := os.Create(zipPath)
	check(err)
	return HuffmanEncoder{
		fpath:       fpath,
		byteEncoder: createNewByteEncoder(zipFile),
		encodingMap: make(map[byte]string),
	}
}

func (encoder *HuffmanEncoder) Compress() {
	freqMap := encoder.getFreqMap()
	root := encoder.getTree(freqMap)
	populateEncodingMap(&root, "", encoder.encodingMap)
	encoder.writeEncodingMap()
	encoder.writeEncodedFile()
}

func (*HuffmanEncoder) getTree(freqMap map[byte]int) Node {
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

func (encoder *HuffmanEncoder) getFreqMap() map[byte]int {
	file, err := os.Open(encoder.fpath)
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

func (encoder *HuffmanEncoder) writeEncodingMap() {
	encoder.byteEncoder.writeInt32(int32(len(encoder.encodingMap)))
	for token, encoding := range encoder.encodingMap {
		encoder.byteEncoder.WriteTokenEncodingPair(token, encoding)
	}
}

func (encoder *HuffmanEncoder) writeEncodedFile() {
	file, err := os.Open(encoder.fpath)
	check(err)
	defer file.Close()
	stat, err := os.Stat(encoder.fpath)
	check(err)
	encoder.byteEncoder.writeInt32(int32(stat.Size()))
	for {
		buffer := make([]byte, BUFFER_SIZE)
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		for i := 0; i < n; i++ {
			encoder.byteEncoder.writeEncoding(encoder.encodingMap[buffer[i]])
		}
	}
	encoder.byteEncoder.flush()
}

func Decompress(zipfpath, fpath string) {
	buffer, err := os.ReadFile(zipfpath)
	check(err)
	byteDecoder := createNewByteDecoder(buffer)
	file, err := os.Create(fpath)
	check(err)
	byteDecoder.decodeAndWrite(file)
}
