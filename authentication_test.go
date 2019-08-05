package main

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestVerifyBasicAuth(t *testing.T) {
	// Given an array of users
	ryanPasswordBytes, ryanHashError := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	if ryanHashError != nil {
		t.Fatal(ryanHashError)
	}
	josePasswordBytes, joseHashError := bcrypt.GenerateFromPassword([]byte("password456"), 10)
	if joseHashError != nil {
		t.Fatal(joseHashError)
	}
	users := []User{User{"ryan", string(ryanPasswordBytes)}, User{"jose", string(josePasswordBytes)}}

	// WHEN we attempt to verify a user with a correct username/password combination
	if !VerifyBasicAuth("ryan", "password123", true, users) {
		// VERIFY that the bool returned is true
		t.Fatal("Expected user to be verified")
	}

	// WHEN we attempt to verify a user with a valid user name and a valid password, but not matching
	if VerifyBasicAuth("ryan", "password456", true, users) {
		// VERIFY that the  bool returned is false
		t.Fatal("Expected user not to be verified")
	}

	// WHEN we attempt to verify a user with a valid user name and an invalid password
	if VerifyBasicAuth("jose", "fail", true, users) {
		// VERIFY that the  bool returned is false
		t.Fatal("Expected user not to be verified")
	}

	// WHEN we attempt to verify a user with an invalid user name and an invalid password
	if VerifyBasicAuth("guest", "password123", true, users) {
		// VERIFY that the  bool returned is false
		t.Fatal("Expected user not to be verified")
	}

	// WHEN we attempt to verify a user with an invalid user name and an invalid password
	if VerifyBasicAuth("guest", "guest", true, users) {
		// VERIFY that the  bool returned is false
		t.Fatal("Expected user not to be verified")
	}
}
