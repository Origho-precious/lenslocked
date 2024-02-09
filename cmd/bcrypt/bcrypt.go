package main

import (
	"fmt"
	"os"
)

func hash(password string) {

}

func compare(password, hash string) {

}

func main() {
	command := os.Args[1]
	password := os.Args[2]
	hashedPassword := os.Args[3]

	switch command {
	case "hash":
		hash(password)
	case "compare":
		compare(password, hashedPassword)
	default:
		fmt.Printf("Invalid command: %v\n", command)
	}
}
