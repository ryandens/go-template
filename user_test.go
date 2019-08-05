package main

import (
	"os"
	"testing"
)

func TestGetUsers(t *testing.T) {
	// GIVEN an array of expected users
	var expectedUsers []User
	expectedUsers = append(expectedUsers, User{
		"ryandens",
		"abcd",
	}, User{"josesolis", "defg"})

	// WHEN we call GetUsersFromFile on the test_users.csv artifact
	file, fileErr := os.Open("test_users.csv")
	if fileErr != nil {
		t.Fatalf("Problem opening test artifact test_users.csv")
	}
	users, e := GetUsersFromFile(file)
	if e != nil {
		t.Fatal(e)
	}

	// VERIFY the lengths match
	if len(users) != len(expectedUsers) {
		t.Fatalf("Unexpected number of users parsed from CSV")
	}

	// VERIFY the contents are the same
	for i, expectedUser := range expectedUsers {
		if expectedUser != users[i] {
			t.Fatalf("user parsed from test_users.csv did not match. Expected: %v, Actual :%v", expectedUser, users[i])
		}
	}
}

func TestAddUserToFile(t *testing.T) {
	file, fileErr := os.Open("test_users.csv")
	if fileErr != nil {
		t.Fatalf("Problem opening test artifact test_users.csv")
	}

	duplicateSuccess, duplicateErr := AddUserToFile(file, User{
		"ryandens",
		"FakeHash",
	})
	if duplicateErr != nil {
		t.Fatal(duplicateErr)
	}
	if duplicateSuccess {
		t.Fatal("Duplicate user name should not be successfully added")
	}

	newSuccess, newErr := AddUserToFile(file, User{"newuser", "FakeHash"})
	if newErr != nil {
		t.Fatal(newErr)
	}
	if !newSuccess {
		t.Fatal("New user should have been  added to file")
	}
}
