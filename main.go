package main

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	flag.Parse()

	fileName := flag.Args()[0]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Print("File not found")
		return
	}

	charCountMap, err := countCharacters(file)
	if err != nil {
		fmt.Print("File not found")
		return
	}

	pq := createHuffmanPartialTreeQueue(charCountMap)

	createHuffManTreeFromPq(&pq)
	printTree(heap.Pop(&pq).(*HuffmanNode))

}

func printTree(item *HuffmanNode) {
	fmt.Printf("Element: %s, Weight: %d, is leaf %t \n", item.element, item.weight, item.isLeaf)
	if item.left != nil {
		printTree(item.left)
	}
	if item.right != nil {
		printTree(item.right)
	}
}

func createHuffManTreeFromPq(pq *PriorityQueue) {
	for pq.Len() > 1 {
		item1 := heap.Pop(pq).(*HuffmanNode)
		item2 := heap.Pop(pq).(*HuffmanNode)
		newNode := &HuffmanNode{weight: item1.weight + item2.weight, isLeaf: false, left: item1, right: item2}
		heap.Push(pq, newNode)
	}
}

func createHuffmanPartialTreeQueue(charCountMap map[string]int) PriorityQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	for character, occurence := range charCountMap {
		node := &HuffmanNode{weight: occurence, element: character, isLeaf: true}
		heap.Push(&pq, node)
	}
	return pq
}

func countCharacters(file *os.File) (map[string]int, error) {
	charCountMap := make(map[string]int)

	reader := bufio.NewReader(file)
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		charCountMap[string(r)]++
	}
	return charCountMap, nil
}
