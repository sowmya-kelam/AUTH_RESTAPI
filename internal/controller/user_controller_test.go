package controller

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"net/http"
	"net/http/httptest"
	"restapi/models"
	"restapi/internal/service"

	"restapi/utils"
	//"database/sql"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Save(u *models.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserService) Validate(u *models.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserService) GetNewAccessToken(refreshtoken string) (string,error) {
	args := m.Called(refreshtoken)
	return args.String(0),args.Error(1)
}

func setUpRouter(service service.UserService) *gin.Engine {
	router := gin.Default()
	Controller := &UserController{Service: service}
	//mockdb := &sql.DB{}
	router.POST("/signup", func(c *gin.Context) { Controller.Signup(c) })
	router.POST("/login", func(c *gin.Context) { Controller.Login(c) })
	router.POST("/authorize-token", func(c *gin.Context) { Controller.AuthorizeToken(c) })
	router.POST("revoke-token", func(c *gin.Context) { Controller.RevokeToken(c) })
	router.POST("/refresh-token", func(c *gin.Context) { Controller.RefreshToken(c) })
	return router
}

func TestSignup(t *testing.T) {
    gin.SetMode(gin.TestMode)
    mockService := new(MockUserService)
    router := setUpRouter(mockService)
    user := models.User{Email: "test@example.com", Password: "password123"}

    mockService.On("Save", &user).Return(nil)

    
    body, _ := json.Marshal(user)
    req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

   
    router.ServeHTTP(w, req)

   
    assert.Equal(t, http.StatusCreated, w.Code)
    assert.JSONEq(t, `{"message": "user created "}`, w.Body.String())
    mockService.AssertExpectations(t)
}


func TestLogin(t *testing.T) {	
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	router := setUpRouter(mockService)
	user := models.User{
		Email:    "test@example.com",
		Password: "password123",
	}	
	mockService.On("Validate", &user).Return(nil)	
	body, _ := json.Marshal(user)	
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")	
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestAuthorizeToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	router := setUpRouter(mockService)
	token, _, _ := utils.GenerateToken("test@example.com", 1)

	tokenReq := models.TokenRequest{Token:token}
	body, _ := json.Marshal(tokenReq)
	
	req, _ := http.NewRequest(http.MethodPost, "/authorize-token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}


func TestRevokeToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	router := setUpRouter(mockService)

	token, _, _ := utils.GenerateToken("test@example.com", 1)
	tokenReq := models.TokenRequest{Token: "Bearer " + token}

	body, _ := json.Marshal(tokenReq)
	req, _ := http.NewRequest(http.MethodPost, "/revoke-token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Token is Revoked")
}

func TestRefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	router := setUpRouter(mockService)

	_, refreshToken, _ := utils.GenerateToken("test@example.com", 1)
	refreshTokenReq := models.RefreshTokenRequest{RefreshToken: refreshToken}

	mockService.On("GetNewAccessToken", refreshToken).Return("",nil)	

	body, _ := json.Marshal(refreshTokenReq)
	req, _ := http.NewRequest(http.MethodPost, "/refresh-token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "New Token Generated")
}
