package main

type HuffmanNode struct {
	weight  int
	element string
	isLeaf  bool
	left    *HuffmanNode
	right   *HuffmanNode
	index   int
}
