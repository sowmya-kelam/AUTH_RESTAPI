package utils

import (
	"testing"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"restapi/constants"
)

func TestGenerateToken(t *testing.T) {
	email := "testuser@example.com"
	userid := int64(1)

	accessToken, refreshToken, err := GenerateToken(email, userid)

	// Assertions
	assert.NoError(t, err, "Expected no error while generating tokens")
	assert.NotEmpty(t, accessToken, "Expected access token to be non-empty")
	assert.NotEmpty(t, refreshToken, "Expected refresh token to be non-empty")
}

func TestVerifyToken_ValidToken(t *testing.T) {
	email := "testuser@example.com"
	userid := int64(1)

	accessToken, _, err := GenerateToken(email, userid)
	assert.NoError(t, err, "Expected no error while generating token")

	verifiedUserID, err := VerifyToken(accessToken)

	// Assertions for valid token
	assert.NoError(t, err, "Expected no error while verifying token")
	assert.Equal(t, userid, verifiedUserID, "Expected verified user ID to match the original")
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	invalidToken := "Bearer invalid.token.signature"

	verifiedUserID, err := VerifyToken(invalidToken)

	// Assertions for invalid token
	assert.Error(t, err, "Expected an error for invalid token")
	assert.Equal(t, int64(0), verifiedUserID, "Expected verified user ID to be 0 for invalid token")
}

func TestVerifyToken_ExpiredToken(t *testing.T) {
	email := "testuser@example.com"
	userid := int64(1)

	// Generate an expired token for testing
	expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userid,
		"exp":    time.Now().Add(-time.Hour).Unix(), // Token already expired
	}).SignedString([]byte(constants.Secretkey))
	assert.NoError(t, err, "Expected no error while generating an expired token")

	verifiedUserID, err := VerifyToken(expiredToken)

	// Assertions for expired token
	assert.Error(t, err, "Expected an error for expired token")
	assert.Equal(t, int64(0), verifiedUserID, "Expected verified user ID to be 0 for expired token")
}
