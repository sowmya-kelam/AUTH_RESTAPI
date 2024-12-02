package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashpassword(t *testing.T) {
	password := "Password123"
	hashedPassword, err := Hashpassword(password)

	// Assertions
	assert.NoError(t, err, "Expected no error while hashing the password")
	assert.NotEmpty(t, hashedPassword, "Expected hashed password to be non-empty")
	assert.NotEqual(t, password, hashedPassword, "Hashed password should not be the same as the original password")
}

func TestCheckpasswordhash(t *testing.T) {
	password := "Password123"
	hashedPassword, err := Hashpassword(password)
	assert.NoError(t, err, "Expected no error while hashing the password")

	// Assertions for correct password
	//assert.True(t, Checkpasswordhash(hashedPassword, password), "Expected password to match the hash")

	// Assertions for incorrect password
	wrongPassword := "wrongPassword123"
	assert.False(t, Checkpasswordhash(hashedPassword, wrongPassword), "Expected password not to match the hash")
}
