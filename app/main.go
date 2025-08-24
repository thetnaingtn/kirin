package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!!!")
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)

	log.Println("Starting server on :9999")
	log.Fatal(http.ListenAndServe(":9999", mux))
}
