package controller

import (
	"restapi/models"
	"restapi/utils"
	"restapi/internal/service"

	_ "github.com/mattn/go-sqlite3"

	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Global Hashmap and  Mutex to handle revoked tokens safely
var revokedTokens = make(map[string]bool)
var mu sync.Mutex

// UserController handles HTTP requests related to user actions
type UserController struct{
 	Service service.UserService   // Dendency injection of user service 
}

// Signup User Handler godoc
// @Summary Signingup User
// @Description This endpoint allows user to sign up using email and password .
// @Tags Auth Rest Api's
// @Accept  json
// @Produce json
// @Param user body models.User  true "User Details"
// @Success 201 {object} map[string]string "message: user created"
// @Failure 400 {object} map[string]string "Error: Bad Request"
// @Failure 500 {object} map[string]string "Error: Internal Server Error"
// @Router /signup [post]
func(uc *UserController) Signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = uc.Service.Save(&user) // Saving The User Credentials to Database
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "user created "})
}

// Login User Handler godoc
// @Summary Login User
// @Description This endpoint allows user to Login using given email and password .
// @Tags Auth Rest Api's
// @Accept  json
// @Produce json
// @Param user body models.User  true "User Details"
// @Success 200 {object} map[string]string "message: Login successful"
// @Failure 400 {object} map[string]string "Error: Bad Request"
// @Failure 401 {object} map[string]string "Error: Unauthorized"
// @Failure 500 {object} map[string]string "Error: Internal Server Error"
// @Router /login [post]
func(uc *UserController) Login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = uc.Service.Validate(&user)  // Validating User Credentials

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	accesstoken,refreshtoken, err := utils.GenerateToken(user.Email, user.ID) // Generate Access Token, Refresh Token	

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful", "accesstoken": accesstoken,"refreshtoken":refreshtoken})
}

// Authorize Handler godoc
// @Summary Authorize Token
// @Description This endpoint allows user to Authorize access token .
// @Tags Auth Rest Api's
// @Accept  json
// @Produce json
// @Param token body models.TokenRequest  true "Access Token"
// @Success 200 {object} map[string]string "message: valid token"
// @Failure 400 {object} map[string]string "Error: Bad Request"
// @Failure 401 {object} map[string]string "Error: unauthorized"
// @Router /authorize-token [post]
func(uc *UserController) AuthorizeToken(context *gin.Context) {
	var token models.TokenRequest
	err := context.ShouldBindJSON(&token)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if token.Token == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "token is empty"})
		return
	}

	_, err = utils.VerifyToken(token.Token)

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	if IsRevokedToken(token.Token) {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Revoked Token"})
		return
	}
	
	

	context.JSON(http.StatusOK, gin.H{"message": " Vaild Token "})
}

// Revoke Token Handler godoc
// @Summary Revoke Token 
// @Description This endpoint allows user to revoke the access token .
// @Tags Auth Rest Api's
// @Accept  json
// @Produce json
// @Param token body models.TokenRequest  true "Access Token"
// @Success 200 {object} map[string]string "message: Token is Revoked "
// @Failure 400 {object} map[string]string "Error: Bad Request"
// @Router /revoke-token [post]
func(uc *UserController) RevokeToken(context *gin.Context) {
	var token models.TokenRequest
	err := context.ShouldBindJSON(&token)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if token.Token == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "token is empty"})
		return
	}

	if IsRevokedToken(token.Token) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "token is already revoked"})
		return
	}
	
	mu.Lock()
    revokedTokens[token.Token] = true 
    mu.Unlock()

	context.JSON(http.StatusOK, gin.H{"message": " Token is Revoked "})
}


func IsRevokedToken(token string) bool {
	mu.Lock()
    defer mu.Unlock()
    _, exists := revokedTokens[token]
    return exists
}

// Refresh Token Handler godoc
// @Summary  Refresh Token 
// @Description This endpoint allows user to Get New Accesstoken using refresh token .
// @Tags Auth Rest Api's
// @Accept  json
// @Produce json
// @Param refreshtoken body models.RefreshTokenRequest  true "Refresh token request"
// @Success 200 {object} map[string]string "message: New Token Generated"
// @Failure 400 {object} map[string]string "Error: Bad Request"
// @Failure 500 {object} map[string]string "Error: Internal Server Error"
// @Router /refresh-token [post]
func(uc *UserController) RefreshToken(context *gin.Context) {

	var refreshtoken models.RefreshTokenRequest
	err := context.ShouldBindJSON(&refreshtoken)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if refreshtoken.RefreshToken == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "token is empty"})
		return
	}

	if IsRevokedToken(refreshtoken.RefreshToken) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "token is revoked, login again to get new access token"})
		return
	}	
	
	newAccessToken, err := uc.Service.GetNewAccessToken(refreshtoken.RefreshToken)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message":"New Token Generated","Token": newAccessToken})
}
