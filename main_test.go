package main

import (
	"container/heap"
	"os"
	"testing"
)

func TestCountCharacters(t *testing.T) {
	file, _ := os.Open("135-0.txt")
	got, err := countCharacters(file)
	if err != nil {
		t.Error("got error")
	}
	if got["X"] != 333 {
		t.Errorf("got %v, want 333", got["X"])
	}

}

func createCharacterMap() map[string]int {
	charCountMap := make(map[string]int)
	charCountMap["C"] = 32
	charCountMap["D"] = 42
	charCountMap["E"] = 120
	charCountMap["K"] = 7
	charCountMap["M"] = 24
	charCountMap["U"] = 37
	charCountMap["Z"] = 2
	return charCountMap
}

func TestCreateHuffmanPartialTreeQueue(t *testing.T) {
	charMap := createCharacterMap()
	got := createHuffmanPartialTreeQueue(charMap)
	want := make(PriorityQueue, len(charMap))
	i := 0
	for character, weight := range charMap {
		want[i] = &HuffmanNode{weight: weight, element: character, isLeaf: true, index: i}
		i++
	}
	heap.Init(&want)

	for want.Len() > 0 {
		wantedElem := heap.Pop(&want).(*HuffmanNode)
		gotElem := heap.Pop(&got).(*HuffmanNode)

		if wantedElem.element != gotElem.element {
			t.Error("error")
		}
	}
}
