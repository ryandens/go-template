package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
)

type User struct {
	Name           string
	HashedPassword string
}

// calls GetUsersFromFile with the default users.csv file
func GetUsers() ([]User, error) {
	file, e := os.Open("users.csv")
	if e != nil {
		log.Fatalf("Problem opening users.csv %v", e)
	}
	return GetUsersFromFile(file)
}

// read existing users into memory
// adapted from example provided in https://golang.org/pkg/encoding/csv/
func GetUsersFromFile(usersCsv *os.File) ([]User, error) {
	reader := csv.NewReader(bufio.NewReader(usersCsv))
	var users []User

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Problem reading line with error %v", err)
		}
		if len(line) != 2 {
			log.Fatalf("Encountered CSV entry with invalid number of columns")
		}
		users = append(users, User{Name: line[0], HashedPassword: line[1]})
	}
	return users, nil
}

// calls GetUsersFromFile with the default users.csv file
func AddUser(newUser User) (bool, error) {
	file, e := os.Open("users.csv")
	if e != nil {
		log.Fatalf("Problem opening users.csv %v", e)
	}
	return AddUserToFile(file, newUser)
}

// Reads in all the users into memory. If userName is present, replace it with newName
// and rewrite the file to be updated. This could absolutely be optimized, I'm just not familiar
// enough with go's APIs
func UpdateUserWithName(userName string, newName string) bool {
	users, e := GetUsers()
	if e != nil {
		return false
	}

	var shouldRewrite = false
	for i, user := range users {
		if userName == user.Name {
			users[i] = User{newName, user.HashedPassword}
			shouldRewrite = true
			break
		}
	}
	if !shouldRewrite {
		return false
	}

	if err := os.Remove("users.csv"); err != nil {
		log.Print(err)
		return false
	}

	file, createErr := os.Create("users.csv")
	if createErr != nil {
		log.Print(createErr)
		return false
	}

	for _, user := range users {
		_, writeErr := file.WriteString(user.Name + "," + user.HashedPassword)
		if writeErr != nil {
			log.Print(writeErr)
		}
	}
	return true
}

// read existing users into memory, verify the user name is not in the file, then add the user to that file
// return true if successful, otherwise false
func AddUserToFile(usersCsv *os.File, newUser User) (bool, error) {
	users, err := GetUsersFromFile(usersCsv)
	if err != nil {
		return false, err
	}

	// check if there is already a user. If there is, don't overwrite user and return false
	for _, existingUser := range users {
		if existingUser.Name == newUser.Name {
			return false, nil
		}
	}

	writeableFile, wfErr := os.OpenFile(usersCsv.Name(), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if wfErr != nil {
		return false, wfErr
	}
	_, wrErr := writeableFile.WriteString(newUser.Name + "," + newUser.HashedPassword)
	if wrErr != nil {
		return false, wrErr
	}
	return true, nil
}
