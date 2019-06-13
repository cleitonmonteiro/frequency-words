package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cleitonmonteiro/frequency-words/messages"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	pb "github.com/cleitonmonteiro/frequency-words/service/proto"
	"google.golang.org/grpc"
)

type server struct{}

var serverAddress string
const (
	urlController = "http://localhost:4007/api/v1/serve"
	defaultServerPort = "40051"
	host = "localhost"
)

func init() {
	port := defaultServerPort
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	serverAddress = fmt.Sprintf("%v:%v", host, port)
	addr := messages.ServiceAddr{Host: "localhost", Port: port}
	data, _ := json.Marshal(addr)
	buff := bytes.NewBuffer(data)

	log.Printf("| POST '%v'", urlController)
	_, _ = http.Post(urlController, "application/json", buff)
}

func main() {
	log.Printf("| starting service on '%v'", serverAddress)
	conn, err := net.Listen("tcp", serverAddress)

	if err != nil {
		log.Fatalln(err)
	}

	s := grpc.NewServer()
	pb.RegisterFrequencyWordsServer(s, &server{})

	if err := s.Serve(conn); err != nil {
		log.Fatalln(err)
	}
}

func (s *server) Frequency(c context.Context, r *pb.Text) (*pb.Response, error) {
	log.Printf("| _-_-_-_-_-_ Received _-_-_-_-_-_\n%v\n\n", r.Body)
	db := make(chan map[string]int64, 64) // cap
	done := make(chan bool)

	// producer
	go handler(r.Body, db, done)

	// consumer
	result := make(map[string]int64)
done:
	for {
		select {
		case fw := <-db:
			for word, f := range fw {
				result[word] += f
			}
		case <-done:
			break done
		}
	}
	//log.Println("| END CALL!")
	return pbResult(result), nil
}

func handler(text string, db chan<- map[string]int64, done chan<- bool) {
	lines := strings.Split(text, "\n")
	wg := &sync.WaitGroup{}
	for _, line := range lines {
		go frequencyWords(line, db, wg)
		wg.Add(1)
	}
	wg.Wait()
	done <- true
}

func frequencyWords(text string, db chan<- map[string]int64, wg *sync.WaitGroup) {
	defer wg.Done()

	frequency := make(map[string]int64)
	words := strings.Split(text, " ")

	for _, word := range words {
		frequency[word]++
	}
	db <- frequency
}

func pbResult(r map[string]int64) *pb.Response {
	response := &pb.Response{}
	var result []*pb.FrequencyWord
	for word, f := range r {
		result = append(result, &pb.FrequencyWord{Word: word, Frequency: f})
	}
	response.FrequencyAll = result
	return response
}
