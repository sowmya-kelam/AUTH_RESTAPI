package utils


import (
	"errors"
	"time"
	"strings"
	"restapi/constants"

	"github.com/golang-jwt/jwt/v5"
)
//Generate AccessToken and RefreshToken 
func GenerateToken(email string, userid int64) (string,string, error) {
	accessToken,err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userid,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	}).SignedString([]byte(constants.Secretkey))
	if err != nil {
        return "", "", err
    }

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email":  email,
		"userId": userid,
        "exp":    time.Now().Add(7 * 24 * time.Hour).Unix(), // Long-lived token (e.g., 7 days)
    }).SignedString([]byte(constants.Secretkey))
    if err != nil {
        return "", "", err
    }

	return accessToken,refreshToken,err
}

//Verifying Token 
func VerifyToken(token string) (int64, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(constants.Secretkey), nil
	})

	if err != nil {
		return 0, errors.New("could not parse token")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return 0, errors.New("invalid Token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("invalid token claims")
	}
	// email := claims["email"].(string)
	userid := int64(claims["userId"].(float64))

	return userid, nil
}