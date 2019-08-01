package main

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home")
}

func main() {
	http.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":8080", nil)
}
