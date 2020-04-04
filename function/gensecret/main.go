package main

import (
	"fmt"
	"io/ioutil"
	"os"

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
	fmt.Println()
	hashed, err := bcrypt.GenerateFromPassword(inp, bcrypt.DefaultCost)
	ioutil.WriteFile("bin/secret", hashed, 0444)
}
