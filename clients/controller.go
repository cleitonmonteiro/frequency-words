package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cleitonmonteiro/frequency-words/messages"
	"log"
	"net/http"
)

func main() {
	url := "http://localhost:4007/api/v1/frequency"

	var jsonStr = []byte(`{"url": "http://localhost:7005/txt/arq02", "text": ""}`)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	response := &messages.FrequencyResponse{}
	err = json.NewDecoder(resp.Body).Decode(response)

	if err != nil {
		log.Fatal(err)
	}

	if response.ErrorCode != 0 {
		log.Fatal(response.Error)
	}

	fmt.Printf("| %30v | %v\n", "Word", "Frequency")
	fmt.Println("|=================================================================")
	for word, freq := range response.Result {
		fmt.Printf("| %30v | %v\n", word, freq)
	}
}
