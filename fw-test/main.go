package main

import (
	"bufio"
	"io"
	"log"
)

const (
	BlockSize = 2
)


func buildBlocksText(body io.Reader) []string {
	scanner := bufio.NewScanner(body)
	scanner.Split(bufio.ScanLines)

	var blocks []string
	var strTemp string
	var counter = 1

	for scanner.Scan() {
		strTemp += scanner.Text()
		if counter % BlockSize == 0 {
			blocks = append(blocks, strTemp)
			strTemp = ""
		} else {
			strTemp += "\n"
		}
		counter++
	}
	blocks = append(blocks, strTemp)
	return blocks
}

const (
	url = "http://www.gutenberg.org/cache/epub/55752/pg55752.txt"
)

func main() {
	log.Printf("| ---------------------------------------- Received ----------------------------------------\n")
	log.Printf("| _-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_ Received _-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_\n")
}
