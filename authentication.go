package main

import "net/http"

// Takes an HTTP request and verifies that the supplied username and password matches an entry in the
// supplied array of Users
func VerifyBasicAuth(r *http.Request, users []User) bool {
	return false
}
