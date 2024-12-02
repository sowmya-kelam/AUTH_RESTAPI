package utils

import (
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

// Encrypting Password 
func Hashpassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Validating Password While Login
func Checkpasswordhash(hashedpassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hashedpassword))
	fmt.Println(err)
	return err == nil
}