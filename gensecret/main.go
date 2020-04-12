package main

import (
	"fmt"
	"os"
	"io/ioutil"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	fmt.Print("enter new secret: ")
	stdin := int(os.Stdin.Fd())
	inp, err := terminal.ReadPassword(stdin)
	if err != nil {
		panic(err)
	}
	hashed, err := bcrypt.GenerateFromPassword(inp, bcrypt.DefaultCost)
	ioutil.WriteFile("bin/handler/secret", hashed, 0644)
}
