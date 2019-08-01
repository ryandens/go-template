package main

import (
	"fmt"
	"log"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Home")
	if err != nil {
		log.Fatalf("Problem writing response with error %v", err)
	}
}

func main() {
	http.HandleFunc("/", HomeHandler)
	log.Print("Listening on localhost")
	// log the error
	log.Fatal(http.ListenAndServe(":8080", nil))
}
