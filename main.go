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
	outputFileName := flag.Args()[1]

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
	lookupMap := make(map[string]string, len(charCountMap))
	populateLookupMap(heap.Pop(&pq).(*HuffmanNode), "", lookupMap)
	encodedData, err := getEncodedData(lookupMap, file)
	if err != nil {
		fmt.Print("File not found")
		return
	}
	addDataToOutputFile(outputFileName, lookupMap, encodedData)
}

func getEncodedData(lookupMap map[string]string, inputFile *os.File) ([]byte, error) {
	reader := bufio.NewReader(inputFile)
	var encoded []byte
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		encoded = append(encoded, []byte(lookupMap[string(r)])...)
	}
	return encoded, nil
}

func addDataToOutputFile(fileName string, lookupMap map[string]string, encodedData []byte) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString("Header Section start \n")
	for character, encoding := range lookupMap {
		_, err = file.WriteString(character + ":" + encoding + "\n")
		if err != nil {
			return err
		}
	}
	file.WriteString("Header Section end \n")
	file.WriteString("Encoded text section start \n")
	file.Write(encodedData)
	return nil

}

func populateLookupMap(item *HuffmanNode, code string, lookupMap map[string]string) {
	if item.isLeaf {
		lookupMap[item.element] = code
	}
	if item.left != nil {
		populateLookupMap(item.left, code+"0", lookupMap)
	}
	if item.right != nil {
		populateLookupMap(item.right, code+"1", lookupMap)
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
