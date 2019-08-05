package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Home")
	if err != nil {
		log.Fatalf("Problem writing response with error %v", err)
	}
}

func main() {
	http.HandleFunc("/", HomeHandler)
	log.Print("Listening on https://localhost:8080/")
	file, e := os.Create("users.csv")
	if e != nil {
		log.Fatal(e)
	} else {
		log.Printf("Sucessfully created %v", file.Name())
	}
	// log the error
	log.Fatal(http.ListenAndServeTLS(":8080", "public.crt", "private.key", nil))
}
