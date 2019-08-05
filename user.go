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

// calls GetUsrsFromFile with the default users.csv file
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
