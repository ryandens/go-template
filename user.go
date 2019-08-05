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

// read existing users into memory
// adapted from example provided in https://golang.org/pkg/encoding/csv/
func GetUsers(usersCsv *os.File) ([]User, error) {
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
