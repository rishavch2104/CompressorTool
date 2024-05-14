package main

import (
	"bufio"
	"container/heap"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {

	flag.Parse()
	if len(flag.Args()) < 3 {
		fmt.Print("invalid arguments")
		return
	}
	fileName := flag.Args()[0]
	outputFileName := flag.Args()[1]
	operation := flag.Args()[2]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Print("File not found")
		return
	}

	if operation == "encode" {

		charCountMap, err := countCharacters(file)
		if err != nil {
			fmt.Print("File not found")
			return
		}

		pq := createHuffmanPartialTreeQueue(charCountMap)

		createHuffManTreeFromPq(&pq)
		lookupMap := make(map[rune]string, len(charCountMap))
		populateLookupMap(heap.Pop(&pq).(*HuffmanNode), "", lookupMap)
		encodedData, err := getEncodedData(lookupMap, file)
		if err != nil {
			fmt.Print("File not found")
			return
		}
		addDataToOutputFile(outputFileName, lookupMap, encodedData)
	} else if operation == "decode" {
		lookupMap := createLookupMapFromFile(file)
		decodedData := decodeDataUsingLookup(lookupMap, fileName)

		writeDecodedDataToFile(decodedData, outputFileName)
	} else {
		fmt.Print("invalid operation!")
	}

}

func writeDecodedDataToFile(data string, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	err = file.Truncate(0)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(data)
	return nil
}

func decodeDataUsingLookup(lookupMap map[string]rune, fileName string) string {
	var decodedString string
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return ""
	}

	fileContent := string(content)

	headerEndIndex := strings.Index(fileContent, "Header Section end")
	if headerEndIndex == -1 {
		return ""
	}

	encodedString := fileContent[headerEndIndex+len("Header Section end"):]
	encodedString = strings.TrimSpace(encodedString)

	fmt.Println("Encoded string:", len(encodedString))
	var currentCharacter string = ""
	for i := 0; i < len(encodedString); i++ {
		fmt.Print(i)
		currentCharacter = currentCharacter + string(encodedString[i])
		if decoding, ok := lookupMap[currentCharacter]; ok {
			decodedString = decodedString + string(decoding)
			currentCharacter = ""
		}
	}
	return decodedString
}

func createLookupMapFromFile(file *os.File) map[string]rune {
	scanner := bufio.NewScanner(file)
	headerStarted := false
	headerEnded := false

	data := make(map[rune]string)
	var jsonDataBuilder strings.Builder
	for scanner.Scan() {
		line := scanner.Text()

		if !headerStarted {
			if strings.HasPrefix(line, "Header Section start") {
				headerStarted = true
			}
			continue
		}

		if headerStarted && !headerEnded {
			if strings.HasPrefix(line, "Header Section end") {
				headerEnded = true
				break
			}

			jsonDataBuilder.WriteString(line)
		}
	}
	if err := json.Unmarshal([]byte(jsonDataBuilder.String()), &data); err != nil {
		fmt.Println("Error parsing JSON data:", err)
		return nil
	}
	reversedData := make(map[string]rune)
	for key, value := range data {
		reversedData[value] = key
	}

	return reversedData

}

func getEncodedData(lookupMap map[rune]string, inputFile *os.File) ([]byte, error) {
	if _, err := inputFile.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	reader := bufio.NewReader(inputFile)
	var encoded []byte
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		encoded = append(encoded, []byte(lookupMap[r])...)
	}
	return encoded, nil
}

func addDataToOutputFile(fileName string, lookupMap map[rune]string, encodedData []byte) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	err = file.Truncate(0)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString("Header Section start \n")
	fmt.Print(len(lookupMap))
	jsonData, _ := json.Marshal(lookupMap)

	file.Write(jsonData)
	file.WriteString("\n")
	file.WriteString("Header Section end \n")
	file.WriteString(string(encodedData))
	return nil

}

func populateLookupMap(item *HuffmanNode, code string, lookupMap map[rune]string) {
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

func createHuffmanPartialTreeQueue(charCountMap map[rune]int) PriorityQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	for character, occurence := range charCountMap {
		node := &HuffmanNode{weight: occurence, element: character, isLeaf: true}
		heap.Push(&pq, node)
	}
	return pq
}

func countCharacters(file *os.File) (map[rune]int, error) {

	charCountMap := make(map[rune]int)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			charCountMap['\n']++
		}
		for _, char := range line {
			charCountMap[char]++
		}
	}

	return charCountMap, nil

}
