package service


import (
	"database/sql"
	"errors"

	"restapi/utils"
	"restapi/models"
	"restapi/constants"
	"github.com/golang-jwt/jwt/v5"


	_ "github.com/mattn/go-sqlite3"
)


type UserService interface {

	Save(u *models.User) error
	Validate(u *models.User) error
	GetNewAccessToken(refreshtoken string) (string, error)

}

type Userservice struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) UserService {
	return &Userservice{db}
}

func (us *Userservice) Save(u *models.User) error {
	if us.db == nil {
		return errors.New("database is nill")
	}

	query := `INSERT INTO users(email,password)
	 VALUES (?,?)`

	stmt, err := us.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedpassword, err := utils.Hashpassword(u.Password) // encrypting password 
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.Email, hashedpassword)

	if err != nil {
		return err
	}

	return nil
}

func (us *Userservice) Validate(u *models.User) error {
	if us.db == nil {
		return errors.New("database is nill")
	}

	query := `select id,password from users where email=?`

	row := us.db.QueryRow(query, u.Email)
	var retrievedpassword string
	err := row.Scan(&u.ID, &retrievedpassword)
	if err != nil {
		return err
	}

	passwordvalid := utils.Checkpasswordhash(u.Password, retrievedpassword) // validating given password with hashed password

	if !passwordvalid {
		return errors.New("invalid credentials")
	}

	return nil
}

func (us *Userservice) GetNewAccessToken(refreshtoken string) (string, error) {

	//parsing the refreshtoken
	parsedToken, err := jwt.Parse(refreshtoken, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, nil
        }
        return []byte(constants.Secretkey), nil
    })

    if err != nil || !parsedToken.Valid {
		return "", err
    }
	// extracting Claims That has User info
    claims, ok := parsedToken.Claims.(jwt.MapClaims)
    if !ok || !parsedToken.Valid {
		return "", errors.New("invalid claims")
    }

    userId, ok := claims["userId"].(float64)
    if !ok {
		
		return "", errors.New("invalid userid")
    }

    email, ok := claims["email"].(string)
    if !ok {
		return "", errors.New("invalid email")
    }
	
	newAccessToken,_, err := utils.GenerateToken(string(email),int64(userId)) // generating new token
	if err != nil {
		return "", err
	}

	return newAccessToken,nil
}