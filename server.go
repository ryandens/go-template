package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Print("Unsupported HTTP method")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
	_, err := fmt.Fprintf(w, "Home")
	if err != nil {
		log.Fatalf("Problem writing response with error %v", err)
	}
}

// wraps HTTP handlers in another HTTP handler which first does authorization checks and then executes the original handler
func AuthWrapperHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			log.Print("Basic Auth header returned false for ok")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}

		users, userError := GetUsers()
		if userError != nil {
			log.Print("Problem getting users from users.csv")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}

		if !VerifyBasicAuth(username, password, ok, users) {
			log.Print("Invalid login attempt")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}

		// execute the request if the user is authorized
		handler(w, r)
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
