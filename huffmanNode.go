package main

type HuffmanNode struct {
	weight  int
	element rune
	isLeaf  bool
	left    *HuffmanNode
	right   *HuffmanNode
	index   int
}
