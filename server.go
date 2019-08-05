package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Print("Unsupported HTTP method")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	_, err := fmt.Fprintf(w, "Home")
	if err != nil {
		log.Printf("Problem writing response with error %v", err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	username, _, _ := r.BasicAuth()
	// the username was already validated before reaching this point, as a result of being validated in SignupHandler
	// before getting written to the CSV and in AuthWrapperHandler which only gets values from the CSV, but lets make
	// sure its safe to write to response
	_, err := fmt.Fprintf(w, "Hello, %v", template.HTMLEscapeString(username))
	if err != nil {
		log.Print(err)
	}
}

func ChangeNameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Print("Unsupported HTTP method")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	username, _, _ := r.BasicAuth()
	bodyBytes, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Print(e)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	newUserName := string(bodyBytes)

	if !ValidateUserName(newUserName, w) {
		return
	} else {
		UpdateUserWithName(username, newUserName)
	}
	_, writeErr := fmt.Fprintf(w, "Hello, %v", template.HTMLEscapeString(newUserName))
	if writeErr != nil {
		log.Print(writeErr)
	}
}

// return true if valid, otherwise false
func ValidatePassword(password string, w http.ResponseWriter) bool {
	if len(password) > 32 {
		message := "Username or password did not comply with length guidelines"
		log.Print(message)
		_, writeErr := w.Write([]byte(message))
		if writeErr != nil {
			log.Print(writeErr)
		}
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}
	return true
}

// return true if valid, otherwise false
func ValidateUserName(userName string, w http.ResponseWriter) bool {

	// only allow names and passwords that are less than 32 characters in length to prevent the
	// user from writing an arbitrary number of bytes to disk
	// The number 32 is insignificant here, there just needs to be a sensible ceiling
	if len(userName) > 32 {
		message := "Username or password did not comply with length guidelines"
		log.Print(message)
		_, writeErr := w.Write([]byte(message))
		if writeErr != nil {
			log.Print(writeErr)
		}
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}

	// we need to be extra careful about the user name because it gets written to the response
	// to prevent XSS, make sure it at least one English letter, followed by any amount of english letters or spaces.
	// As support is added for other dialects, this regex can be expanded
	userNameRegex := regexp.MustCompile("[A-Za-z][A-Za-z ]*")
	if !userNameRegex.MatchString(userName) {
		message := "Username must contain only letters A-Z or a-z"
		log.Print(message)
		_, writeErr := w.Write([]byte(message))
		if writeErr != nil {
			log.Print(writeErr)
		}
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}

	return true
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

		// if either username or password aren't valid, return
		if !ValidateUserName(username, w) || !ValidatePassword(password, w) {
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
		HelloHandler(w, r)
	} else {
		log.Print("Unsupported HTTP method")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

// wraps HTTP handlers in another HTTP handler which first does authorization checks and then executes the original handler
func AuthWrapperHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			log.Print("Basic Auth header returned false for ok")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		users, userError := GetUsers()
		if userError != nil {
			log.Print("Problem getting users from users.csv")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if !VerifyBasicAuth(username, password, ok, users) {
			log.Print("Invalid login attempt")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// execute the request if the user is authorized
		handler(w, r)
	}
}

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/signup", SignUpHandler)
	http.HandleFunc("/hello", AuthWrapperHandler(HelloHandler))
	http.HandleFunc("/update-name", AuthWrapperHandler(ChangeNameHandler))
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
