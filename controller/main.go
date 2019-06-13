package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cleitonmonteiro/frequency-words/messages"
	pb "github.com/cleitonmonteiro/frequency-words/service/proto"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	addrServices    chan string
	frequencyErrors chan error
)


const (
	BlockSize = 20
	ServerAddress = "localhost:4007"
)

func init() {
	addrServices = make(chan string, 128)
	frequencyErrors = make(chan error)
}

func main() {
	log.Printf("| starting server on '%v'", ServerAddress)

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/api/v1/", homeHandler).Methods("GET")
	r.HandleFunc("/api/v1/frequency", frequencyHandler).Methods("GET")
	r.HandleFunc("/api/v1/serve", serveHandler).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	log.Fatal(http.ListenAndServe(ServerAddress, handlers.CompressHandler(r)))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "<h1>index handler!</h1>")
}

func frequencyHandler(w http.ResponseWriter, r *http.Request) {
	freqRequest := &messages.FrequencyRequest{}
	_ = json.NewDecoder(r.Body).Decode(freqRequest)

	log.Printf("| new request '%v', content '%v'", r.RemoteAddr, freqRequest)

	textReader := io.Reader(strings.NewReader(freqRequest.Text))

	if freqRequest.Text == "" {
		log.Printf("| GET '%v'", freqRequest.Url)
		response, _ := http.Get(freqRequest.Url)
		textReader = response.Body
	}

	frequency := make(chan []*pb.FrequencyWord, 64)
	done := make(chan bool)

	go splitRequest(buildBlocksText(textReader), frequency, done)
	result := make(map[string]int64)

done:
	for {
		select {
		case serviceFrequency := <-frequency:
			for _, frequencyWord := range serviceFrequency {
				result[frequencyWord.Word] += frequencyWord.Frequency
			}
		case err := <-frequencyErrors:
			log.Printf("| requisition failed '%v'\n", r.RemoteAddr)
			_ = json.NewEncoder(w).Encode(messages.FrequencyResponse{ErrorCode: 1, Error: err.Error()})
			return
		case <-done:
			break done
		}
	}
	log.Printf("| requisition completed '%v'\n", r.RemoteAddr)

	_ = json.NewEncoder(w).Encode(messages.FrequencyResponse{Result: result})
}

func serveHandler(w http.ResponseWriter, r *http.Request) {
	addr := &messages.ServiceAddr{}
	_ = json.NewDecoder(r.Body).Decode(addr)
	addrServices <- addr.String()
	log.Printf("| new service available '%v'\n", addr)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "<h1>Not Found!</h1>")
}

func splitRequest(blocks []string, frequency chan<- []*pb.FrequencyWord, done chan bool) {
	wg := &sync.WaitGroup{}
	for _, block := range blocks {
		go serviceHandler(block, frequency, wg)
		wg.Add(1)
	}
	wg.Wait()
	done <- true
}

func serviceHandler(block string, frequency chan<- []*pb.FrequencyWord, wg *sync.WaitGroup) {
	defer wg.Done()
	serviceAddr, err := serviceAddress()
	if err != nil {
		time.Sleep(time.Second*2)
		serviceAddr, err = serviceAddress()
		if err != nil {
			frequencyErrors <- err
			return
		}
	}
	log.Printf("| created rpc connection '%v'\n", serviceAddr)
	conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := pb.NewFrequencyWordsClient(conn)
	ctx := context.Background()
	log.Printf("| rcp call to'%v'\n", serviceAddr)
	response, err := client.Frequency(ctx, &pb.Text{Body: block})
	if err != nil {
		log.Fatalln(err)
	}
	frequency <- response.FrequencyAll
}

func serviceAddress() (string, error) {
	select {
	case addr := <-addrServices:
		addrServices <- addr
		return addr, nil
	default:
		return "", errors.New("no service available")
	}
}

func buildBlocksText(body io.Reader) []string {
	scanner := bufio.NewScanner(body)
	scanner.Split(bufio.ScanLines)

	var blocks []string
	var strTemp string
	var counter = 1

	for scanner.Scan() {
		strTemp += scanner.Text()
		if counter % BlockSize == 0 {
			blocks = append(blocks, strings.ToUpper(strTemp))
			strTemp = ""
		} else {
			strTemp += "\n"
		}
		counter++
	}
	if strTemp != "" {
		blocks = append(blocks, strings.ToUpper(strTemp))
	}

	return blocks
}
