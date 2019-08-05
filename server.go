package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
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

// Accepts POST requests. If username and password are valid, it creates a user and stores in CSV file
// TODO there is definitely some opportunity for code re-use here, but as I am the only reviewer, I wanted
// to make it crystal clear where user data (potentially attacker data) is being used.
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username, password, ok := r.BasicAuth()

		if !ok {
			log.Print("Problem parsing Basic Auth header")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// only allow names and passwords that are less than 32 characters in length to prevent the
		// user from writing an arbitrary number of bytes to disk
		// The number 32 is insignificant here, there just needs to be a sensible ceiling
		if len(username) > 32 || len(password) > 32 {
			message := "Username or password did not comply with length guidelines"
			log.Print(message)
			_, writeErr := w.Write([]byte(message))
			if writeErr != nil {
				log.Print(writeErr)
			}
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// we need to be extra careful about the user name because it gets written to the response
		// to prevent XSS, make sure it at least one English letter, followed by any amount of english letters or spaces.
		// As support is added for other dialects, this regex can be expanded
		userNameRegex := regexp.MustCompile("[A-Za-z][A-Za-z ]*")
		if !userNameRegex.MatchString(username) {
			message := "Username must contain only letters A-Z or a-z"
			log.Print(message)
			_, writeErr := w.Write([]byte(message))
			if writeErr != nil {
				log.Print(writeErr)
			}
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		hashedPassword, hashErr := HashPassword(password)
		if hashErr != nil {
			message := "Problem creating user"
			log.Print(message)
			_, writeErr := w.Write([]byte(message))
			if writeErr != nil {
				log.Print(writeErr)
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		success, addUserErr := AddUser(User{username, hashedPassword})

		if addUserErr != nil || !success {
			message := "Problem creating user"
			log.Print(message)
			_, writeErr := w.Write([]byte(message))
			if writeErr != nil {
				log.Print(writeErr)
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// it is safe to print the username here because we know it passes our regex of safe values
		_, writeErr := fmt.Fprintf(w, "Successfully created user with name %v", username)
		if writeErr != nil {
			log.Print("problem writing to response")
		}
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
	http.HandleFunc("/signup", SignUpHandler)
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
