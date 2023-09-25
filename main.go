package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var counter struct {
	Total int
	sync.Mutex
}

func get(writer http.ResponseWriter, _ *http.Request) {
	log.Printf("GET counter request: %v", counter.Total)
	_, _ = fmt.Fprintf(writer, "Counter is at: %d\n", counter.Total)
}

func set(writer http.ResponseWriter, req *http.Request) {
	log.Printf("SET counter request: %v", req.RequestURI)
	value := req.URL.Query().Get("value")
	intval, err := strconv.Atoi(value)

	if err != nil {
		log.Println("SET handler: non-integer parameter value.")
	}

	counter.Total = intval
	log.Printf("counter set to: %v", counter.Total)
	_, _ = fmt.Fprintf(writer, "Counter set to: %d\n", counter.Total)
}

func inc(_ http.ResponseWriter, _ *http.Request) {
	counter.Lock()
	defer counter.Unlock()
	counter.Total++
	log.Printf("counter incremented to: %v", counter.Total)
}

func dec(_ http.ResponseWriter, _ *http.Request) {
	counter.Lock()
	defer counter.Unlock()
	counter.Total--
	log.Printf("counter incremented to: %v", counter.Total)
}

func main() {
	http.HandleFunc("/counter", get)
	http.HandleFunc("/counter/set", set)
	http.HandleFunc("/increment", inc)
	http.HandleFunc("/decrement", dec)

	port := 9095
	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1])
	}
	log.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(port), nil))
}
