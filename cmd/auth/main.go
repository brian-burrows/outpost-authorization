package main

import (
	"fmt"

	"github.com/brian-burrows/outpost-authorization/internal/auth"
	"github.com/brian-burrows/outpost-authorization/internal/db"
)

func main() {
	fmt.Println(auth.HandlerFileName)
	fmt.Println(auth.RepositoryFileName)
	fmt.Println(auth.ServiceFileName)
	fmt.Println(db.DbFileName)
	// use mongodb.Connect() to get client pointer to read instance
	// defer mongodb.Close()
	// use mongodb.correct() to get client pointer to write instance
	// create concrete instance of UserRepository type for read instance
	// create concrete instance of UserRepository type for write instance
	// create handlers for Users and Tokens, passing in UserRepository for each
	// register handlers with some type of MUX
}
