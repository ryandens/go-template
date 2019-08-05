package main

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// Takes an HTTP request and verifies that the supplied username and password matches an entry in the
// supplied array of Users
func VerifyBasicAuth(suppliedUsername string, suppliedPassword string, ok bool, users []User) bool {
	suppliedPasswordByteArray := []byte(suppliedPassword)

	if !ok {
		log.Fatal("Problem parsing Basic Auth header")
		return false
	}

	for _, user := range users {
		if user.Name == suppliedUsername {
			return nil == bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), suppliedPasswordByteArray)
		}
	}
	return false
}

// hashes password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedPassword), err
}
