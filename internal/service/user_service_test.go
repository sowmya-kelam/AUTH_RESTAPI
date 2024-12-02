package service

import (
	"errors"
	"restapi/models"
	"database/sql"
	"testing"
	"time"
	"github.com/golang-jwt/jwt/v5"

	"github.com/stretchr/testify/assert"
)

type mockUserService struct {
	saveErr     error
	validateErr error
}

func (m *mockUserService) Save(db *sql.DB, u *models.User) error {
	return m.saveErr
}

func (m *mockUserService) Validate(db *sql.DB, u *models.User) error {
	return m.validateErr
}

func TestUserService_Save_Success(t *testing.T) {
	mockService := &mockUserService{saveErr: nil}
	user := &models.User{Email: "test@example.com", Password: "password123"}

	err := mockService.Save(nil, user)
	assert.NoError(t, err, "Expected no error when saving a valid user")
}

func TestUserService_Save_Failure(t *testing.T) {
	mockService := &mockUserService{saveErr: errors.New("save error")}
	user := &models.User{Email: "test@example.com", Password: "password123"}

	err := mockService.Save(nil, user)
	assert.EqualError(t, err, "save error", "Expected 'save error' when save fails")
}

func TestUserService_Validate_Success(t *testing.T) {
	mockService := &mockUserService{validateErr: nil}
	user := &models.User{Email: "test@example.com", Password: "password123"}

	err := mockService.Validate(nil, user)
	assert.NoError(t, err, "Expected no error when validating a valid user")
}

func TestUserService_Validate_InvalidCredentials(t *testing.T) {
	mockService := &mockUserService{validateErr: errors.New("invalid credentials")}
	user := &models.User{Email: "test@example.com", Password: "wrongpassword"}

	err := mockService.Validate(nil, user)
	assert.EqualError(t, err, "invalid credentials", "Expected 'invalid credentials' when validation fails")
}

func TestGetNewAccessToken_ValidToken(t *testing.T) {
	us := &Userservice{}

	// Create a valid refresh token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  "test@example.com",
		"userId": float64(1),
		"exp":    time.Now().Add(time.Hour).Unix(),
	})
	refreshToken, _ := token.SignedString([]byte("supersecret"))

	// Act
	newAccessToken, err := us.GetNewAccessToken(refreshToken)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)
}


