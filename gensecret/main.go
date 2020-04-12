package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	var inp []byte
	if _, err := fmt.Scanln(&inp); err != nil {
		panic(err)
	}
	hashed, err := bcrypt.GenerateFromPassword(inp, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(hashed))
}
