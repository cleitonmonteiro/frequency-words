package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	ServerAddress = "localhost:7005"
	basePath = "/home/monteiro/go/src/github.com/cleitonmonteiro/frequency-words/text-server/files/"
)

func main() {
	log.Printf("| starting text server on '%v'", ServerAddress)

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/txt/{file}", txtHandler).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	log.Fatal(http.ListenAndServe(ServerAddress, handlers.CompressHandler(r)))
}

func txtHandler(w http.ResponseWriter, r *http.Request) {
	file := mux.Vars(r)["file"]
	log.Printf("| request file '%v'\n", file)

	data, err := ioutil.ReadFile(basePath + file)
	if err != nil {
		log.Printf("| file not found | error: %v\n", err)
		_, _ = w.Write([]byte("<h1> File not found! </h1>"))
		return
	}
	log.Println("| success!")

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write(data)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "<h1>Page not Found!</h1> <h2> Use: /txt/{file} </h2>")
}
