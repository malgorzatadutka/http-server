package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var counter int
var mutex = &sync.Mutex{}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
	}
	if r.Method != "GET" {
		http.Error(w, "method is not supported", http.StatusNotFound)
	}
	fmt.Fprint(w, "hello!")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() error: %v", err)
	}
	fmt.Fprintf(w, "POST request successful!\n")
	name := r.FormValue("name")
	address := r.FormValue("address")
	fmt.Fprintf(w, "Name =%s\n", name)
	fmt.Fprintf(w, "Adress =%s\n", address)
}

func incrementCounterHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	fmt.Fprintf(w, strconv.Itoa(counter))
	mutex.Unlock()
}

func headersHandler(w http.ResponseWriter, r *http.Request) {
	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	//check static folder to index.html file
	fileServer := http.FileServer((http.Dir("/static")))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/counter", incrementCounterHandler)
	http.HandleFunc("/headers", headersHandler)

	fmt.Printf("Starting server at port 8080\n")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
