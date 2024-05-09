package main

import (
	"bufio"
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
	fmt.Print(charCountMap["t"])
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
