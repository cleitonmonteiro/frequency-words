package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/cleitonmonteiro/frequency-words/service/proto"
	"google.golang.org/grpc"
)

const (
	serviceAddr = "localhost:5001"
)

func main() {

	conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := pb.NewFrequencyWordsClient(conn)
	ctx := context.Background()
	log.Printf("| rcp call to'%v'\n", serviceAddr)
	text := "TEST text client pb google word nice try"
	response, err := client.Frequency(ctx, &pb.Text{Body: text})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(response)
}
