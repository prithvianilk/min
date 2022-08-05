package main

import (
	pq "github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/emirpasic/gods/utils"
)

type Node struct {
	token       rune
	freq        int
	left, right *Node
}

func createNewNode(token rune, freq int, left *Node, right *Node) Node {
	return Node{token: token, freq: freq, left: left, right: right}
}

func createNewLeaf(token rune, freq int) Node {
	return createNewNode(token, freq, nil, nil)
}

type TokenEncodingPair struct {
	token    rune
	encoding string
}

func createNewEncodingPair(token rune, encoding string) TokenEncodingPair {
	return TokenEncodingPair{token: token, encoding: encoding}
}

func populateEncodings(node *Node, encoding string, encodings *[]TokenEncodingPair) {
	if node == nil {
		return
	}
	if node.isLeaf() {
		*encodings = append(*encodings, createNewEncodingPair(node.token, encoding))
		return
	}
	populateEncodings(node.left, encoding+"0", encodings)
	populateEncodings(node.right, encoding+"1", encodings)
}

func (node *Node) isLeaf() bool {
	if node == nil {
		return false
	}
	return node.left == nil && node.right == nil
}

func byPriority(a, b interface{}) int {
	pa, pb := a.(Node).freq, b.(Node).freq
	return utils.IntComparator(pa, pb)
}

func dequeueAndGetNode(q *pq.Queue) Node {
	top, _ := q.Dequeue()
	return top.(Node)
}
