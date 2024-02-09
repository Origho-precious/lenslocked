package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func hash(password string) {
	hashByte, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost,
	)

	if err != nil {
		panic(err)
	}

	hashedPassword := string(hashByte)

	// fmt.Println("hashByte", hashByte)

	fmt.Printf("Password hash: %q\n", hashedPassword)
}

func compare(password, hash string) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		fmt.Printf("Incorrect password: %q\n", password)
		return
	}

	fmt.Println("Password is correct!")
}

func main() {
	// fmt.Println(len(os.Args))
	if len(os.Args) < 3 {
		panic(`Need to specify a command e.g "hash" or "compare"`)
	}
	
	var command, password, hashedPassword string
	
	command = os.Args[1]
	password = os.Args[2]
	
	if len(os.Args) == 4 {
		hashedPassword = os.Args[3]
	}

	switch command {
	case "hash":
		hash(password)
	case "compare":
		compare(password, hashedPassword)
	default:
		fmt.Printf("Invalid command: %v\n", command)
	}
}
